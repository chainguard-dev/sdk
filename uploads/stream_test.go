/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"runtime"
	"strings"
	"sync"
	"testing"

	uploads "chainguard.dev/sdk/uploads"
)

const (
	// gcmTagSize is the AES-GCM authentication tag length the stream
	// format appends to every chunk. Spelled out here so the tests can
	// frame a sealed stream without importing crypto/cipher.
	gcmTagSize = 16

	// streamHeaderFixedSize is the stream header length up to the wrapped
	// DEK: magic(4), version(1), prefix(7), chunkSize(4), keyVersion(8),
	// wrappedLen(2). Restated rather than imported, so an accidental
	// layout change fails the tests that edit the header on the wire.
	streamHeaderFixedSize = 26
	noncePrefixAt         = 5
	chunkSizeAt           = 12
	keyVersionAt          = 16
	wrappedLenAt          = 24

	// headerLen is the full header for the RSA-3072 keys these tests use:
	// the fixed part plus the wrapped DEK, which is one modulus wide.
	headerLen = streamHeaderFixedSize + 3072/8
)

// streamKeys caches two RSA-3072 keypairs for the package. Key
// generation dominates the runtime of these tests and none of them
// depends on a fresh modulus.
var streamKeys = sync.OnceValues(func() ([2]*rsa.PrivateKey, error) {
	var keys [2]*rsa.PrivateKey
	for i := range keys {
		k, err := rsa.GenerateKey(rand.Reader, 3072)
		if err != nil {
			return keys, err
		}
		keys[i] = k
	}
	return keys, nil
})

// streamKeypair returns the i'th cached test keypair.
func streamKeypair(t *testing.T, i int) *rsa.PrivateKey {
	t.Helper()
	keys, err := streamKeys()
	if err != nil {
		t.Fatalf("generate rsa keys: %v", err)
	}
	return keys[i]
}

// sealStream seals plaintext under pub, writing it in pieces of
// writeSize so the writer's chunk buffering is exercised across Write
// boundaries, and returns the complete stream bytes.
func sealStream(t *testing.T, pub *rsa.PublicKey, plaintext []byte, writeSize int) []byte {
	t.Helper()
	var buf bytes.Buffer
	sw, err := uploads.SealStream(&buf, pub, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal stream: %v", err)
	}
	for rest := plaintext; len(rest) > 0; {
		n := min(writeSize, len(rest))
		if _, err := sw.Write(rest[:n]); err != nil {
			t.Fatalf("write: %v", err)
		}
		rest = rest[n:]
	}
	if err := sw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	return buf.Bytes()
}

// splitStream frames a sealed stream into its header and its chunks, so
// a test can drop, swap or replace one of them. It requires every chunk
// to be full length, which holds for a plaintext that is an exact
// multiple of StreamChunkSize.
func splitStream(t *testing.T, sealed []byte) (header []byte, out [][]byte) {
	t.Helper()
	chunkLen := uploads.StreamChunkSize + gcmTagSize
	if len(sealed) < headerLen {
		t.Fatalf("sealed stream is %d bytes, shorter than its %d-byte header", len(sealed), headerLen)
	}
	body := sealed[headerLen:]
	if len(body) == 0 || len(body)%chunkLen != 0 {
		t.Fatalf("stream body is %d bytes, want a positive multiple of %d", len(body), chunkLen)
	}
	out = make([][]byte, 0, len(body)/chunkLen)
	for at := 0; at < len(body); at += chunkLen {
		out = append(out, body[at:at+chunkLen])
	}
	return sealed[:headerLen], out
}

// openStream opens sealed with priv and returns the full plaintext.
func openStream(t *testing.T, sealed []byte, priv *rsa.PrivateKey) ([]byte, error) {
	t.Helper()
	r, err := uploads.OpenStream(bytes.NewReader(sealed), rsaStreamUnwrap(priv))
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

func rsaStreamUnwrap(priv *rsa.PrivateKey) func(string, []byte) ([]byte, error) {
	return func(_ string, wrapped []byte) ([]byte, error) {
		return rsaUnwrap(priv)(wrapped)
	}
}

func TestSealOpenStreamRoundTrip(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, tc := range []struct {
		name string
		size int
	}{
		{"empty", 0},
		{"one byte", 1},
		{"under one chunk", 4096},
		{"one byte under a chunk", uploads.StreamChunkSize - 1},
		{"exactly one chunk", uploads.StreamChunkSize},
		{"one byte over a chunk", uploads.StreamChunkSize + 1},
		{"exactly two chunks", 2 * uploads.StreamChunkSize},
		{"two chunks and a remainder", 2*uploads.StreamChunkSize + 7919},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := bytesOfLen(tc.size)
			// 7919 is prime, so writes land at every offset within a
			// chunk rather than only on chunk boundaries.
			sealed := sealStream(t, &priv.PublicKey, want, 7919)
			got, err := openStream(t, sealed, priv)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("round-trip mismatch: got %d bytes, want %d bytes", len(got), len(want))
			}
		})
	}
}

func TestSealStreamWritesCiphertext(t *testing.T) {
	priv := streamKeypair(t, 0)
	plaintext := bytesOfLen(64 * 1024)
	sealed := sealStream(t, &priv.PublicKey, plaintext, len(plaintext))
	if bytes.Contains(sealed, plaintext) {
		t.Fatal("sealed stream contains the plaintext verbatim")
	}
	if len(sealed) <= len(plaintext) {
		t.Fatalf("sealed stream is %d bytes, want more than the %d-byte plaintext", len(sealed), len(plaintext))
	}
}

// TestOpenStreamRejectsTamperedStream covers both directions of forging
// the final-chunk flag, which is not a byte on the wire an attacker can
// flip: the reader derives it from the framing and folds it into the
// nonce. "final chunk dropped" presents a chunk sealed non-final where
// the reader demands final, and "trailing bytes appended" presents the
// final chunk where the reader demands non-final. Both fail to open.
func TestOpenStreamRejectsTamperedStream(t *testing.T) {
	priv := streamKeypair(t, 0)
	other := streamKeypair(t, 1)
	const chunks = 3
	plaintext := bytesOfLen(chunks * uploads.StreamChunkSize)

	for _, tc := range []struct {
		name string
		// tamper rewrites a freshly sealed stream. It is handed the
		// header and the chunks already framed, and returns the bytes
		// OpenStream should reject.
		tamper func(t *testing.T, header []byte, chunk [][]byte) []byte
	}{{
		name: "final chunk dropped",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			return bytes.Join([][]byte{header, chunk[0], chunk[1]}, nil)
		},
	}, {
		name: "final chunk truncated mid-way",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			last := chunk[chunks-1]
			return bytes.Join([][]byte{header, chunk[0], chunk[1], last[:len(last)/2]}, nil)
		},
	}, {
		name: "all chunks dropped",
		tamper: func(_ *testing.T, header []byte, _ [][]byte) []byte {
			return header
		},
	}, {
		name: "chunks reordered",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			return bytes.Join([][]byte{header, chunk[1], chunk[0], chunk[2]}, nil)
		},
	}, {
		name: "final chunk moved first",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			return bytes.Join([][]byte{header, chunk[2], chunk[0], chunk[1]}, nil)
		},
	}, {
		name: "chunk duplicated",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			return bytes.Join([][]byte{header, chunk[0], chunk[0], chunk[1], chunk[2]}, nil)
		},
	}, {
		name: "chunk spliced in from another stream under the same key",
		tamper: func(t *testing.T, header []byte, chunk [][]byte) []byte {
			t.Helper()
			_, donor := splitStream(t, sealStream(t, &priv.PublicKey, plaintext, len(plaintext)))
			return bytes.Join([][]byte{header, chunk[0], donor[1], chunk[2]}, nil)
		},
	}, {
		name: "chunk spliced in from another stream under another key",
		tamper: func(t *testing.T, header []byte, chunk [][]byte) []byte {
			t.Helper()
			_, donor := splitStream(t, sealStream(t, &other.PublicKey, plaintext, len(plaintext)))
			return bytes.Join([][]byte{header, chunk[0], donor[1], chunk[2]}, nil)
		},
	}, {
		name: "header spliced in from another stream",
		tamper: func(t *testing.T, _ []byte, chunk [][]byte) []byte {
			t.Helper()
			donorHeader, _ := splitStream(t, sealStream(t, &priv.PublicKey, plaintext, len(plaintext)))
			return bytes.Join([][]byte{donorHeader, chunk[0], chunk[1], chunk[2]}, nil)
		},
	}, {
		name: "bit flipped in a chunk body",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			chunk[1][0] ^= 0x01
			return bytes.Join([][]byte{header, chunk[0], chunk[1], chunk[2]}, nil)
		},
	}, {
		name: "bit flipped in a chunk auth tag",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			last := chunk[chunks-1]
			last[len(last)-1] ^= 0x01
			return bytes.Join([][]byte{header, chunk[0], chunk[1], last}, nil)
		},
	}, {
		name: "trailing bytes appended",
		tamper: func(_ *testing.T, header []byte, chunk [][]byte) []byte {
			return bytes.Join([][]byte{header, chunk[0], chunk[1], chunk[2], bytesOfLen(gcmTagSize)}, nil)
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			// Seal fresh per row: the tamper functions mutate in place.
			header, chunk := splitStream(t, sealStream(t, &priv.PublicKey, plaintext, len(plaintext)))
			if _, err := openStream(t, tc.tamper(t, header, chunk), priv); err == nil {
				t.Fatal("got nil error, want the tampered stream to fail to open")
			}
		})
	}
}

func TestOpenStreamRejectsBadHeader(t *testing.T) {
	priv := streamKeypair(t, 0)
	sealed := sealStream(t, &priv.PublicKey, bytesOfLen(1024), 1024)

	for _, tc := range []struct {
		name    string
		input   func() []byte
		wantErr string
	}{{
		name:    "empty input",
		input:   func() []byte { return nil },
		wantErr: "read stream header",
	}, {
		name:    "truncated header",
		input:   func() []byte { return sealed[:4] },
		wantErr: "read stream header",
	}, {
		name: "wrong magic",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			bad[0] = 'X'
			return bad
		},
		wantErr: "not a stream envelope",
	}, {
		name: "unsupported version",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			bad[4] = 99
			return bad
		},
		wantErr: "unsupported stream envelope version",
	}, {
		name: "chunk size zero",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			copy(bad[chunkSizeAt:chunkSizeAt+4], []byte{0, 0, 0, 0})
			return bad
		},
		wantErr: "chunk size",
	}, {
		name: "chunk size beyond the cap",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			copy(bad[chunkSizeAt:chunkSizeAt+4], []byte{0xff, 0xff, 0xff, 0xff})
			return bad
		},
		wantErr: "chunk size",
	}, {
		name: "key version zero",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			binary.BigEndian.PutUint64(bad[keyVersionAt:], 0)
			return bad
		},
		wantErr: "key version",
	}, {
		name: "wrapped key size zero",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			binary.BigEndian.PutUint16(bad[wrappedLenAt:], 0)
			return bad
		},
		wantErr: "wrapped key size",
	}, {
		name: "wrapped key size beyond the cap",
		input: func() []byte {
			bad := bytes.Clone(sealed)
			binary.BigEndian.PutUint16(bad[wrappedLenAt:], 2049)
			return bad
		},
		wantErr: "wrapped key size",
	}, {
		name:    "header truncated before the wrapped dek",
		input:   func() []byte { return sealed[:streamHeaderFixedSize+8] },
		wantErr: "read wrapped dek",
	}, {
		name: "a v1 json envelope",
		input: func() []byte {
			_, pubPEM := keypair(t)
			payload, err := uploads.SealEnvelope([]byte("v1"), pubPEM, uploads.EncryptionAlgorithm, "1")
			if err != nil {
				t.Fatalf("seal v1 envelope: %v", err)
			}
			return []byte(payload)
		},
		wantErr: "not a stream envelope",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := uploads.OpenStream(bytes.NewReader(tc.input()), func(string, []byte) ([]byte, error) {
				t.Error("unwrap called for an invalid header")
				return nil, errors.New("unexpected unwrap")
			})
			if err == nil {
				t.Fatalf("got nil error, want error containing %q", tc.wantErr)
			}
			if !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("got error = %v, want error containing %q", err, tc.wantErr)
			}
		})
	}
}

// A short final chunk fits several declared sizes. The header's AAD binding
// must reject an edit even when the ciphertext framing remains valid.
func TestOpenStreamRejectsEditedChunkSize(t *testing.T) {
	priv := streamKeypair(t, 0)
	sealed := sealStream(t, &priv.PublicKey, bytesOfLen(4096), 4096)

	for _, tc := range []struct {
		name      string
		chunkSize uint32
	}{
		// Only the first re-frames the chunk and so fails with or without
		// the binding. A short final chunk is read short at any declared
		// size it fits under, so the other two leave the framing intact
		// and are the rows the binding is for.
		{"lowered below the payload", 2048},
		{"lowered to a size the payload still fits under", 8192},
		{"raised to the cap", 16 << 20},
	} {
		t.Run(tc.name, func(t *testing.T) {
			edited := bytes.Clone(sealed)
			binary.BigEndian.PutUint32(edited[chunkSizeAt:], tc.chunkSize)
			if _, err := openStream(t, edited, priv); err == nil {
				t.Fatal("got nil error, want the edited chunk size to fail to open")
			}
		})
	}

	if _, err := openStream(t, sealed, priv); err != nil {
		t.Fatalf("open the unedited stream: %v", err)
	}
}

// TestOpenStreamBoundsAllocationFromUnauthenticatedHeader pins the
// reader's buffers to the bytes that arrive rather than to the chunk
// size the header asks for. No header field authenticates until a chunk
// opens under it, and the AAD binding does not help here: it rejects the
// stream only on the first Read, long after the buffers are sized. Doing
// it the other way made a tiny stream cost 32 MiB.
func TestOpenStreamBoundsAllocationFromUnauthenticatedHeader(t *testing.T) {
	priv := streamKeypair(t, 0)
	hostile := bytes.Clone(sealStream(t, &priv.PublicKey, bytesOfLen(16), 16))
	binary.BigEndian.PutUint32(hostile[chunkSizeAt:], 16<<20)

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	r, err := uploads.OpenStream(bytes.NewReader(hostile), rsaStreamUnwrap(priv))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	_, readErr := io.ReadAll(r)
	runtime.ReadMemStats(&after)

	if readErr == nil {
		t.Fatal("got nil error, want the edited chunk size to fail to open")
	}
	// Generous against the 33 MiB the preallocating form cost, and far
	// above the few hundred KiB this path actually needs, so incidental
	// allocation elsewhere in the binary cannot tip it.
	const budget = 4 << 20
	if grew := after.TotalAlloc - before.TotalAlloc; grew > budget {
		t.Errorf("opening a %d-byte stream declaring a 16 MiB chunk size allocated %d bytes, want at most %d", len(hostile), grew, budget)
	}
}

// TestSealStreamUsesAFreshDEKPerStream pins the precondition the whole
// nonce construction rests on. A chunk index is unique only within a
// key, so "a chunk cannot be moved between streams" holds because no two
// streams share a DEK, not because of anything in the nonce. If this
// ever fails, the index stops separating streams and the 7-byte random
// prefix becomes the only thing standing between them.
func TestSealStreamUsesAFreshDEKPerStream(t *testing.T) {
	priv := streamKeypair(t, 0)
	plaintext := bytesOfLen(4096)
	unwrap := func(sealed []byte) []byte {
		t.Helper()
		dek, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, sealed[streamHeaderFixedSize:headerLen], nil)
		if err != nil {
			t.Fatalf("unwrap dek: %v", err)
		}
		return dek
	}

	first := sealStream(t, &priv.PublicKey, plaintext, len(plaintext))
	second := sealStream(t, &priv.PublicKey, plaintext, len(plaintext))
	if bytes.Equal(unwrap(first), unwrap(second)) {
		t.Fatal("two streams sealed under the same public key share a dek; chunk indices no longer separate them")
	}
	if bytes.Equal(first[noncePrefixAt:chunkSizeAt], second[noncePrefixAt:chunkSizeAt]) {
		t.Fatal("two streams share a nonce prefix")
	}
}

// TestOpenStreamDistinguishesEmptyFromTruncated covers the degenerate
// lengths. Every stream carries a final chunk, including one over an
// empty payload, so truncating a stream to its header is detectable
// rather than indistinguishable from a sender who wrote nothing.
func TestOpenStreamDistinguishesEmptyFromTruncated(t *testing.T) {
	priv := streamKeypair(t, 0)

	for _, tc := range []struct {
		name string
		size int
	}{
		{"zero payload bytes", 0},
		{"one short chunk", 4096},
		{"exactly one full chunk", uploads.StreamChunkSize},
	} {
		t.Run(tc.name, func(t *testing.T) {
			want := bytesOfLen(tc.size)
			sealed := sealStream(t, &priv.PublicKey, want, 4096)
			got, err := openStream(t, sealed, priv)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("round-trip mismatch: got %d bytes, want %d bytes", len(got), len(want))
			}
			// The header alone must not read as an empty payload, which is
			// exactly what the zero-byte row would decrypt to.
			if _, err := openStream(t, sealed[:headerLen], priv); err == nil {
				t.Fatal("got nil error, want a stream truncated to its header to fail to open")
			}
		})
	}

	// An empty payload is one tag on the wire, so it is a different
	// length from a header-only stream as well as a different outcome.
	if empty := sealStream(t, &priv.PublicKey, nil, 4096); len(empty) != headerLen+gcmTagSize {
		t.Errorf("empty stream: got %d bytes, want %d (header plus one tag)", len(empty), headerLen+gcmTagSize)
	}
}

func TestOpenStreamRejectsWrongKey(t *testing.T) {
	priv := streamKeypair(t, 0)
	other := streamKeypair(t, 1)
	sealed := sealStream(t, &priv.PublicKey, bytesOfLen(3*uploads.StreamChunkSize), 64*1024)

	if _, err := openStream(t, sealed, other); err == nil {
		t.Fatal("got nil error, want the wrong private key to fail to unwrap the dek")
	}
	// And the right key still opens the same bytes, so the row above is
	// rejecting the key rather than a corrupted stream.
	if _, err := openStream(t, sealed, priv); err != nil {
		t.Fatalf("open with the sealing key: %v", err)
	}
}

func TestSealStreamRejectsWriteAfterClose(t *testing.T) {
	priv := streamKeypair(t, 0)
	sw, err := uploads.SealStream(io.Discard, &priv.PublicKey, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal stream: %v", err)
	}
	if err := sw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if _, err := sw.Write([]byte("late")); err == nil {
		t.Fatal("got nil error, want error writing to a closed stream")
	}
	// Close is idempotent, so a deferred Close after an explicit one does
	// not append a second final chunk.
	if err := sw.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
}

func TestSealStreamRejectsNilArguments(t *testing.T) {
	priv := streamKeypair(t, 0)
	if _, err := uploads.SealStream(nil, &priv.PublicKey, uploads.EncryptionAlgorithm, "1"); err == nil {
		t.Fatal("got nil error, want error on a nil writer")
	}
	if _, err := uploads.SealStream(io.Discard, nil, uploads.EncryptionAlgorithm, "1"); err == nil {
		t.Fatal("got nil error, want error on a nil public key")
	}
	if _, err := uploads.OpenStream(nil, rsaStreamUnwrap(priv)); err == nil {
		t.Fatal("got nil error, want error on a nil reader")
	}
	if _, err := uploads.OpenStream(bytes.NewReader(nil), nil); err == nil {
		t.Fatal("got nil error, want error on a nil private key")
	}
}

// TestSealStreamRejectsWrongAlgorithm mirrors the SealEnvelope check: an
// algorithm the client cannot produce must fail before any byte is
// written, not seal a stream no reader can decrypt.
func TestSealStreamRejectsWrongAlgorithm(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, algorithm := range []string{"", "RSA_DECRYPT_OAEP_2048_SHA256", "rsa_decrypt_oaep_3072_sha256"} {
		t.Run(algorithm, func(t *testing.T) {
			var sealed bytes.Buffer
			w, err := uploads.SealStream(&sealed, &priv.PublicKey, algorithm, "1")
			if err == nil || !strings.Contains(err.Error(), "algorithm") || w != nil {
				t.Errorf("seal: got (%v, %v), want (nil, algorithm error)", w, err)
			}
			if sealed.Len() != 0 {
				t.Errorf("wrote %d bytes for a rejected algorithm, want 0", sealed.Len())
			}
		})
	}
}

// TestCRC32CWriterUsesCastagnoli pins the polynomial to the one Google
// Cloud Storage checks object uploads against, using the standard CRC
// check vector. Swapping the table for crc32.IEEE yields 0xcbf43926 for
// the same input and fails here; GCS would reject every upload.
func TestCRC32CWriterUsesCastagnoli(t *testing.T) {
	const (
		input = "123456789"
		want  = uint32(0xe3069283)
	)
	cw := uploads.NewCRC32CWriter(io.Discard)
	if _, err := io.WriteString(cw, input); err != nil {
		t.Fatalf("write: %v", err)
	}
	if got := cw.Sum32(); got != want {
		t.Errorf("crc32c(%q): got = %#08x, want = %#08x", input, got, want)
	}
}

func TestCRC32CWriterMatchesChecksumOverCiphertext(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, tc := range []struct {
		name      string
		size      int
		writeSize int
	}{
		{"empty", 0, 1},
		{"under one chunk", 4096, 97},
		{"exactly one chunk", uploads.StreamChunkSize, 7919},
		{"several chunks", 3*uploads.StreamChunkSize + 11, 64 * 1024},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plaintext := bytesOfLen(tc.size)

			var buf bytes.Buffer
			cw := uploads.NewCRC32CWriter(&buf)
			sw, err := uploads.SealStream(cw, &priv.PublicKey, uploads.EncryptionAlgorithm, "1")
			if err != nil {
				t.Fatalf("seal stream: %v", err)
			}
			for rest := plaintext; len(rest) > 0; {
				n := min(tc.writeSize, len(rest))
				if _, err := sw.Write(rest[:n]); err != nil {
					t.Fatalf("write: %v", err)
				}
				rest = rest[n:]
			}
			if err := sw.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}

			want := crc32.Checksum(buf.Bytes(), crc32.MakeTable(crc32.Castagnoli))
			if got := cw.Sum32(); got != want {
				t.Errorf("incremental checksum: got = %#08x, want = %#08x", got, want)
			}

			// The checksum covers the header too, so the bytes it
			// describes are exactly the ones that open again.
			got, err := openStream(t, buf.Bytes(), priv)
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			if !bytes.Equal(got, plaintext) {
				t.Fatalf("round-trip mismatch: got %d bytes, want %d bytes", len(got), len(plaintext))
			}
		})
	}
}

func TestCRC32CWriterCountsOnlyAcceptedBytes(t *testing.T) {
	// A writer that accepts half of every write, as a short write would.
	halves := &halfWriter{}
	cw := uploads.NewCRC32CWriter(halves)
	payload := bytesOfLen(1024)
	n, err := cw.Write(payload)
	if err == nil {
		t.Fatal("got nil error, want the short write surfaced to the caller")
	}
	if n != len(payload)/2 {
		t.Fatalf("got n = %d, want %d", n, len(payload)/2)
	}
	want := crc32.Checksum(halves.written.Bytes(), crc32.MakeTable(crc32.Castagnoli))
	if got := cw.Sum32(); got != want {
		t.Errorf("checksum over a short write: got = %#08x, want = %#08x", got, want)
	}
}

// halfWriter accepts the first half of every write and reports the short
// write, the way a saturated network writer would.
type halfWriter struct {
	written bytes.Buffer
}

func (h *halfWriter) Write(p []byte) (int, error) {
	n, err := h.written.Write(p[:len(p)/2])
	if err != nil {
		return n, err
	}
	return n, io.ErrShortWrite
}
