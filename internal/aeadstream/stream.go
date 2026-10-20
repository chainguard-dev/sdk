/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aeadstream

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
)

const (
	// KeySize is the AES-256 DEK length in bytes.
	KeySize = 32

	// NoncePrefixSize is the per-stream random nonce prefix length in
	// bytes. The remaining 5 bytes of the 12-byte GCM nonce are the chunk
	// index and the final-chunk flag.
	NoncePrefixSize = 7

	// MaxChunkSize bounds the chunk size NewReader honors. A format reads
	// its chunk size from a header that is unauthenticated until a chunk
	// opens under it, so an unbounded value would let a hostile stream
	// name a multi-gigabyte frame. This bounds the frame only; the
	// reader's buffer tracks the bytes that actually arrive (see
	// chunkBufGrowth).
	MaxChunkSize = 16 << 20

	// nonceSize is the AES-GCM standard nonce length in bytes.
	nonceSize = 12

	// chunkBufGrowth is the reader's first ciphertext buffer allocation,
	// doubling from there until a chunk is fully read. Sizing the buffer
	// from the declared chunk size instead would let a tiny stream
	// cost 32 MiB, because nothing in the header authenticates until a
	// chunk opens under it.
	chunkBufGrowth = 64 << 10
)

// errStreamClosed is returned by writes that arrive after Close has
// sealed the final chunk; accepting them would silently drop data.
var errStreamClosed = errors.New("write on a closed stream")

// streamCipher is the per-stream AEAD state shared by the writer and the
// reader: the DEK-derived AEAD, the additional data bound into every
// chunk, and the nonce whose leading NoncePrefixSize bytes are the fixed
// per-stream prefix.
type streamCipher struct {
	aead      cipher.AEAD
	aad       []byte
	chunkSize int
	nonce     []byte
}

// newStreamCipher validates the arguments both constructors share and
// binds them into the STREAM state.
func newStreamCipher(dek, aad, noncePrefix []byte, chunkSize int) (*streamCipher, error) {
	if l := len(dek); l != KeySize {
		return nil, fmt.Errorf("dek has wrong length: got %d, want %d", l, KeySize)
	}
	if l := len(noncePrefix); l != NoncePrefixSize {
		return nil, fmt.Errorf("nonce prefix has wrong length: got %d, want %d", l, NoncePrefixSize)
	}
	if chunkSize < 1 || chunkSize > MaxChunkSize {
		return nil, fmt.Errorf("chunk size %d, want 1..%d", chunkSize, MaxChunkSize)
	}
	block, err := aes.NewCipher(dek)
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("aes-gcm: %w", err)
	}
	sc := &streamCipher{
		aead:      aead,
		aad:       aad,
		chunkSize: chunkSize,
		nonce:     make([]byte, nonceSize),
	}
	copy(sc.nonce, noncePrefix)
	return sc, nil
}

// nonceFor returns the STREAM nonce for a chunk: the per-stream prefix,
// the big-endian chunk index, and a trailing flag that is 1 only on the
// final chunk. Index and flag are what make a chunk unmovable — a
// reordered chunk authenticates under the wrong index, and a truncated
// stream leaves a chunk sealed with flag 0 where the reader demands 1.
func (s *streamCipher) nonceFor(index uint32, final bool) []byte {
	binary.BigEndian.PutUint32(s.nonce[NoncePrefixSize:], index)
	s.nonce[nonceSize-1] = 0
	if final {
		s.nonce[nonceSize-1] = 1
	}
	return s.nonce
}

// NewWriter returns a writer that seals everything written to it into w
// as AES-256-GCM chunks of chunkSize plaintext bytes, with aad as every
// chunk's additional data. It writes nothing itself, so the caller
// writes its header to w first.
//
// The caller must Close the returned writer to emit the final chunk; a
// stream whose final chunk is missing does not open. An empty stream
// still emits one. Close does not close w. Writes after Close fail, as
// does a stream that would need more than 2^32 chunks.
//
// noncePrefix is copied, but aad is retained and must not be modified
// while the writer is in use. The returned writer is not safe for
// concurrent use.
func NewWriter(w io.Writer, dek, aad, noncePrefix []byte, chunkSize int) (io.WriteCloser, error) {
	if w == nil {
		return nil, errors.New("nil writer")
	}
	sc, err := newStreamCipher(dek, aad, noncePrefix, chunkSize)
	if err != nil {
		return nil, err
	}
	return &streamWriter{
		w:   w,
		sc:  sc,
		buf: make([]byte, 0, chunkSize),
		ct:  make([]byte, 0, chunkSize+sc.aead.Overhead()),
	}, nil
}

// streamWriter seals successive plaintext chunks into the underlying
// writer. A full buffer is held back rather than flushed, because
// whether a chunk is final is only known once more bytes arrive or Close
// is called.
type streamWriter struct {
	w      io.Writer
	sc     *streamCipher
	buf    []byte
	ct     []byte
	index  uint32
	closed bool
	err    error
}

var _ io.WriteCloser = (*streamWriter)(nil)

func (s *streamWriter) Write(p []byte) (int, error) {
	if s.err != nil {
		return 0, s.err
	}
	if s.closed {
		return 0, errStreamClosed
	}
	written := 0
	for len(p) > 0 {
		n := min(s.sc.chunkSize-len(s.buf), len(p))
		s.buf = append(s.buf, p[:n]...)
		p = p[n:]
		written += n
		if len(s.buf) == s.sc.chunkSize && len(p) > 0 {
			if err := s.seal(false); err != nil {
				s.err = err
				return written, err
			}
		}
	}
	return written, nil
}

// Close seals the buffered remainder as the final chunk. An empty stream
// still emits one final chunk, which is what lets the reader tell an
// empty payload from a stream truncated to its header.
func (s *streamWriter) Close() error {
	if s.err != nil {
		return s.err
	}
	if s.closed {
		return nil
	}
	s.closed = true
	if err := s.seal(true); err != nil {
		s.err = err
		return err
	}
	return nil
}

func (s *streamWriter) seal(final bool) error {
	// The 4-byte nonce index caps a stream at 2^32 chunks. Refusing a
	// further non-final chunk is what stops the index wrapping onto a
	// nonce this stream has already used; the last index is still a legal
	// home for the final chunk.
	if !final && s.index == math.MaxUint32 {
		return fmt.Errorf("stream exceeds the maximum of %d chunks", uint64(math.MaxUint32)+1)
	}
	s.ct = s.sc.aead.Seal(s.ct[:0], s.sc.nonceFor(s.index, final), s.buf, s.sc.aad)
	if _, err := s.w.Write(s.ct); err != nil {
		return fmt.Errorf("write chunk %d: %w", s.index, err)
	}
	s.buf = s.buf[:0]
	if !final {
		s.index++
	}
	return nil
}

// NewReader returns a reader of the plaintext of the chunks in r, sealed
// as NewWriter seals them under the same dek, aad, noncePrefix and
// chunkSize. r must be positioned at the first chunk, past any header.
//
// Read the returned reader to EOF to verify the complete stream. Each
// chunk authenticates before its plaintext is returned, and
// modification, deletion, reordering, splicing, a missing final chunk,
// or the wrong dek or aad causes a read error. Buffers grow with the
// bytes that arrive rather than with chunkSize, so a chunkSize taken from
// an unauthenticated header costs nothing until a peer sends that much.
//
// noncePrefix is copied, but aad is retained and must not be modified
// while the reader is in use. The returned reader is not safe for
// concurrent use.
func NewReader(r io.Reader, dek, aad, noncePrefix []byte, chunkSize int) (io.Reader, error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}
	sc, err := newStreamCipher(dek, aad, noncePrefix, chunkSize)
	if err != nil {
		return nil, err
	}
	return &streamReader{
		br:    bufio.NewReader(r),
		sc:    sc,
		frame: chunkSize + sc.aead.Overhead(),
	}, nil
}

// streamReader authenticates one chunk at a time. The bufio.Reader is
// there for Peek: a full-length chunk is the final one exactly when no
// bytes follow it, and that has to be known before choosing the nonce.
//
// ct and pt start empty and grow with the chunks that arrive, so the
// declared chunk size bounds them without preallocating them.
type streamReader struct {
	br    *bufio.Reader
	sc    *streamCipher
	frame int
	ct    []byte
	pt    []byte
	plain []byte
	index uint32
	done  bool
	err   error
}

var _ io.Reader = (*streamReader)(nil)

func (s *streamReader) Read(p []byte) (int, error) {
	for len(s.plain) == 0 {
		if s.err != nil {
			return 0, s.err
		}
		if s.done {
			return 0, io.EOF
		}
		if err := s.nextChunk(); err != nil {
			s.err = err
			return 0, err
		}
	}
	n := copy(p, s.plain)
	s.plain = s.plain[n:]
	return n, nil
}

// readFrame preserves the underlying error, including EOF alongside bytes.
// io.ReadFull would turn a clean short frame's EOF into ErrUnexpectedEOF,
// making it indistinguishable from a truncated transport.
func (s *streamReader) readFrame() (int, error) {
	n := 0
	for n < s.frame {
		if n == len(s.ct) {
			s.ct = append(s.ct, make([]byte, min(s.frame-n, max(len(s.ct), chunkBufGrowth)))...)
		}
		read, err := s.br.Read(s.ct[n:])
		n += read
		if err != nil {
			return n, err
		}
	}
	return n, nil
}

func (s *streamReader) nextChunk() error {
	final := false
	n, err := s.readFrame()
	switch {
	case err == nil:
		// A full-length chunk is the final one only if nothing follows it.
		if _, perr := s.br.Peek(1); perr != nil {
			if !errors.Is(perr, io.EOF) {
				return fmt.Errorf("read chunk %d: %w", s.index, perr)
			}
			final = true
		}
	case errors.Is(err, io.EOF):
		if n == 0 {
			return fmt.Errorf("stream ended before chunk %d: %w", s.index, io.ErrUnexpectedEOF)
		}
		final = true
	default:
		return fmt.Errorf("read chunk %d: %w", s.index, err)
	}

	plain, err := s.sc.aead.Open(s.pt[:0], s.sc.nonceFor(s.index, final), s.ct[:n], s.sc.aad)
	if err != nil {
		return fmt.Errorf("authenticate chunk %d: %w", s.index, err)
	}
	// Mirrors the writer's cap. The counter starts at zero and advances
	// only past an authenticated chunk, so a hostile stream can only
	// arrive here the long way, one authenticated chunk at a time. It is
	// checked before the plaintext is staged, or a Read after the error
	// would still return this chunk.
	if !final && s.index == math.MaxUint32 {
		return fmt.Errorf("stream exceeds the maximum of %d chunks", uint64(math.MaxUint32)+1)
	}
	// Keep whatever Open allocated as the next chunk's scratch; s.plain is
	// always drained before nextChunk runs again.
	s.pt, s.plain = plain, plain
	if final {
		s.done = true
		return nil
	}
	s.index++
	return nil
}
