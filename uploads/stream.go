/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"strconv"

	"chainguard.dev/sdk/internal/aeadstream"
)

// StreamChunkSize is the plaintext chunk size SealStream seals under a
// single AES-GCM nonce. 1 MiB keeps the in-memory working set small
// enough to stream an arbitrarily large upload while staying far below
// the per-key, per-nonce birthday bounds for AES-GCM.
const StreamChunkSize = 1 << 20

const (
	// streamMagic prefixes every stream envelope so a reader handed a v1
	// JSON envelope (or arbitrary bytes) fails with a clear error rather
	// than inside the cipher.
	streamMagic = "CGUS"

	// streamVersion is the format version in the header. v1 is the JSON
	// envelope in envelope.go, which carries no version byte of its own;
	// the stream format starts at 2 so the two never collide.
	streamVersion = 2

	// noncePrefixSize is the per-stream random prefix length, in bytes.
	// The prefix is not what makes nonces unique: a fresh DEK per stream
	// already does that. It is carried so that a DEK shared across streams
	// would collide on a 56-bit birthday bound rather than repeat with
	// certainty. See doc.go.
	noncePrefixSize = aeadstream.NoncePrefixSize

	// The wrapped length is explicit because the reader's unwrap callback
	// can use a KMS key whose modulus is unavailable to this package.
	// Layout: magic(4), version(1), prefix(7), chunkSize(4), keyVersion(8),
	// wrappedLen(2), followed by wrappedLen bytes of wrapped DEK.
	// Integer fields are big-endian.
	streamHeaderFixedSize = len(streamMagic) + 1 + noncePrefixSize + 4 + 8 + 2

	// maxWrappedKeySize bounds the wrapped DEK, so an absurdly large RSA
	// key cannot turn into an absurdly large header. 2048 bytes is an
	// RSA-16384 modulus, well above any key Chainguard KMS issues.
	maxWrappedKeySize = 2048
)

// crc32cTable is the Castagnoli polynomial, the CRC32C variant Google
// Cloud Storage uses for object checksums.
var crc32cTable = crc32.MakeTable(crc32.Castagnoli)

// CRC32CWriter passes bytes through to an underlying writer while
// accumulating their CRC32C (Castagnoli) checksum. It exists so a caller
// streaming a sealed upload straight to object storage can hand the
// checksum over at the end without ever holding the ciphertext:
//
//	cw := uploads.NewCRC32CWriter(objectWriter)
//	sw, err := uploads.SealStream(cw, orgPub, algorithm, keyVersion)
//	// ... io.Copy(sw, src); sw.Close() ...
//	checksum := cw.Sum32()
//
// Sum32 is meaningful at any point and covers every byte written so far,
// header included. A CRC32CWriter is not safe for concurrent use.
type CRC32CWriter struct {
	w io.Writer
	h hash.Hash32
}

var _ io.Writer = (*CRC32CWriter)(nil)

// NewCRC32CWriter returns a CRC32CWriter that forwards to w.
func NewCRC32CWriter(w io.Writer) *CRC32CWriter {
	return &CRC32CWriter{w: w, h: crc32.New(crc32cTable)}
}

// Write forwards p to the underlying writer and folds the bytes that
// reached it into the running checksum. A short write contributes only
// the accepted prefix, so the checksum always describes what the
// underlying writer actually holds.
func (c *CRC32CWriter) Write(p []byte) (int, error) {
	n, err := c.w.Write(p)
	if n < 0 || n > len(p) {
		// Reject rather than slice with it: a count outside the io.Writer
		// contract would otherwise panic inside the hash.
		return 0, fmt.Errorf("underlying writer reported %d bytes written for a %d-byte write", n, len(p))
	}
	// hash.Hash32.Write never returns an error.
	_, _ = c.h.Write(p[:n])
	return n, err
}

// Sum32 returns the CRC32C of everything written so far.
func (c *CRC32CWriter) Sum32() uint32 { return c.h.Sum32() }

// SealStream returns a writer that encrypts everything written to it
// into w as a stream envelope: a binary header followed by AES-GCM
// chunks in the STREAM construction.
//
// A fresh AES-256 data encryption key is generated per stream and
// RSA-OAEP-SHA256-wrapped under orgPub in the header, mirroring
// SealEnvelope's hybrid scheme. The header is written before SealStream
// returns, and is bound into every chunk as additional authenticated
// data.
//
// algorithm and keyVersion are the values advertised with orgPub by the
// server's GetOrgPublicKey response. As in SealEnvelope, a mismatch
// against the client's hardcoded EncryptionAlgorithm is a hard failure
// rather than silently producing payloads no reader can decrypt.
// keyVersion is the positive decimal CryptoKeyVersion, recorded so
// readers can select the same key after rotation.
//
// The caller must Close the returned writer to emit the final chunk; a
// stream whose final chunk is missing does not open. Close does not
// close w, which the caller still owns.
//
// The returned writer is not safe for concurrent use.
func SealStream(w io.Writer, orgPub *rsa.PublicKey, algorithm, keyVersion string) (io.WriteCloser, error) {
	if w == nil {
		return nil, errors.New("nil writer")
	}
	if orgPub == nil {
		return nil, errors.New("nil public key")
	}
	if algorithm != EncryptionAlgorithm {
		return nil, fmt.Errorf("server advertises encryption algorithm %q, client expects %q", algorithm, EncryptionAlgorithm)
	}
	version, err := strconv.ParseUint(keyVersion, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("parse key version %q: %w", keyVersion, err)
	}
	if version == 0 || strconv.FormatUint(version, 10) != keyVersion {
		return nil, fmt.Errorf("key version %q must be a positive decimal integer without leading zeros", keyVersion)
	}

	dek := make([]byte, aesKeySize)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, fmt.Errorf("generate dek: %w", err)
	}
	wrapped, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, orgPub, dek, nil)
	if err != nil {
		return nil, fmt.Errorf("rsa-oaep wrap dek: %w", err)
	}
	wrappedLen := len(wrapped)
	if wrappedLen > maxWrappedKeySize {
		return nil, fmt.Errorf("wrapped dek is %d bytes, max %d", wrappedLen, maxWrappedKeySize)
	}
	noncePrefix := make([]byte, noncePrefixSize)
	if _, err := io.ReadFull(rand.Reader, noncePrefix); err != nil {
		return nil, fmt.Errorf("generate nonce prefix: %w", err)
	}

	header := make([]byte, 0, streamHeaderFixedSize+wrappedLen)
	header = append(header, streamMagic...)
	header = append(header, streamVersion)
	header = append(header, noncePrefix...)
	header = binary.BigEndian.AppendUint32(header, StreamChunkSize)
	header = binary.BigEndian.AppendUint64(header, version)
	header = binary.BigEndian.AppendUint16(header, uint16(wrappedLen))
	header = append(header, wrapped...)

	sw, err := aeadstream.NewWriter(w, dek, header, noncePrefix, StreamChunkSize)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(header); err != nil {
		return nil, fmt.Errorf("write stream header: %w", err)
	}
	return sw, nil
}

// OpenStream parses the header and calls unwrapAESKey with the stored key
// version and RSA-OAEP-SHA256-wrapped DEK. As with OpenEnvelope, the callback
// can use KMS without exposing the RSA private key to this package.
//
// The key version is untrusted until a chunk authenticates; the callback
// must scope key lookup to the caller's authorized org. It must not modify
// wrappedKey, which is part of the authenticated header.
//
// Read the returned reader to EOF to verify the complete stream. Each chunk
// authenticates before its plaintext is returned, and modification, deletion,
// reordering, splicing or a missing final chunk causes a read error.
//
// The returned reader is not safe for concurrent use.
func OpenStream(r io.Reader, unwrapAESKey func(keyVersion string, wrappedKey []byte) ([]byte, error)) (io.Reader, error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}
	if unwrapAESKey == nil {
		return nil, errors.New("nil unwrap function")
	}

	var fixed [streamHeaderFixedSize]byte
	if _, err := io.ReadFull(r, fixed[:]); err != nil {
		return nil, fmt.Errorf("read stream header: %w", err)
	}
	if !bytes.HasPrefix(fixed[:], []byte(streamMagic)) {
		return nil, errors.New("input is not a stream envelope")
	}
	if v := fixed[len(streamMagic)]; v != streamVersion {
		return nil, fmt.Errorf("unsupported stream envelope version %d, want %d", v, streamVersion)
	}
	prefixAt := len(streamMagic) + 1
	chunkSize := binary.BigEndian.Uint32(fixed[prefixAt+noncePrefixSize:])
	// aeadstream bounds the chunk size too, but checking it here refuses a
	// hostile header before unwrapAESKey reaches KMS on its behalf.
	if chunkSize == 0 || chunkSize > aeadstream.MaxChunkSize {
		return nil, fmt.Errorf("stream header declares chunk size %d, want 1..%d", chunkSize, aeadstream.MaxChunkSize)
	}
	keyVersionAt := prefixAt + noncePrefixSize + 4
	keyVersion := binary.BigEndian.Uint64(fixed[keyVersionAt:])
	if keyVersion == 0 {
		return nil, errors.New("stream header declares key version zero")
	}
	wrappedLen := int(binary.BigEndian.Uint16(fixed[keyVersionAt+8:]))
	if wrappedLen == 0 || wrappedLen > maxWrappedKeySize {
		return nil, fmt.Errorf("stream header declares wrapped key size %d, want 1..%d", wrappedLen, maxWrappedKeySize)
	}

	header := make([]byte, streamHeaderFixedSize+wrappedLen)
	copy(header, fixed[:])
	noncePrefix := header[prefixAt : prefixAt+noncePrefixSize]
	wrapped := header[streamHeaderFixedSize:]
	if _, err := io.ReadFull(r, wrapped); err != nil {
		return nil, fmt.Errorf("read wrapped dek: %w", err)
	}

	dek, err := unwrapAESKey(strconv.FormatUint(keyVersion, 10), wrapped)
	if err != nil {
		return nil, fmt.Errorf("unwrap dek: %w", err)
	}
	if l := len(dek); l != aesKeySize {
		return nil, fmt.Errorf("unwrapped dek has wrong length: got %d, want %d", l, aesKeySize)
	}
	return aeadstream.NewReader(r, dek, header, noncePrefix, int(chunkSize))
}
