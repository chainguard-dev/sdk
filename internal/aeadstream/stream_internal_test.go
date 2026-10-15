/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aeadstream

import (
	"bytes"
	"crypto/rand"
	"io"
	"math"
	"strings"
	"testing"
)

// TestStreamChunkCapIsEnforcedOnBothSides drives the 2^32-chunk cap from
// inside the package. Reaching it over the wire would take 2^32 chunks,
// so both counters are positioned directly: a black-box test could only
// assert that the cap is unreachable, which is not the same as asserting
// that it holds.
//
// Both sides must agree on where the cap falls, or a stream one side
// seals is a stream the other refuses. The last index is a legal home
// for the final chunk and for nothing else.
func TestStreamChunkCapIsEnforcedOnBothSides(t *testing.T) {
	const chunkSize = 16
	dek := make([]byte, KeySize)
	_, _ = rand.Read(dek)
	aad := []byte("header")
	prefix := []byte("prefix!")

	// sc seals the chunks the reader rows forge, under the same state the
	// reader opens with.
	sc, err := newStreamCipher(dek, aad, prefix, chunkSize)
	if err != nil {
		t.Fatalf("stream cipher: %v", err)
	}
	newReader := func(t *testing.T, chunks ...[]byte) *streamReader {
		t.Helper()
		r, err := NewReader(bytes.NewReader(bytes.Join(chunks, nil)), dek, aad, prefix, chunkSize)
		if err != nil {
			t.Fatalf("new reader: %v", err)
		}
		sr := r.(*streamReader)
		sr.index = math.MaxUint32
		return sr
	}

	t.Run("writer refuses a non-final chunk at the last index", func(t *testing.T) {
		w, err := NewWriter(io.Discard, dek, aad, prefix, chunkSize)
		if err != nil {
			t.Fatalf("new writer: %v", err)
		}
		sw := w.(*streamWriter)
		sw.index = math.MaxUint32

		if err := sw.seal(false); err == nil || !strings.Contains(err.Error(), "maximum") {
			t.Fatalf("sealing a non-final chunk at the last index: got err = %v, want the chunk cap", err)
		}
		if err := sw.seal(true); err != nil {
			t.Fatalf("sealing the final chunk at the last index: %v", err)
		}
	})

	t.Run("reader refuses a non-final chunk at the last index", func(t *testing.T) {
		// Forge what the writer refuses to emit. The chunk has to be a
		// full frame with more bytes behind it, or the reader reads it as
		// the final chunk and never reaches the counter.
		chunk := sc.aead.Seal(nil, sc.nonceFor(math.MaxUint32, false), make([]byte, chunkSize), aad)

		sr := newReader(t, chunk, chunk)
		if _, err := io.ReadAll(sr); err == nil || !strings.Contains(err.Error(), "maximum") {
			t.Fatalf("reading a non-final chunk at the last index: got err = %v, want the chunk cap", err)
		}
		// The error is sticky: the refused chunk authenticated, but it must
		// not surface on a later Read.
		if n, err := sr.Read(make([]byte, chunkSize)); n != 0 || err == nil || !strings.Contains(err.Error(), "maximum") {
			t.Fatalf("read after the chunk cap: got (%d, %v), want (0, the chunk cap)", n, err)
		}
	})

	t.Run("reader accepts the final chunk at the last index", func(t *testing.T) {
		want := []byte("the last chunk")
		chunk := sc.aead.Seal(nil, sc.nonceFor(math.MaxUint32, true), want, aad)

		got, err := io.ReadAll(newReader(t, chunk))
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("final chunk at the last index: got = %q, want = %q", got, want)
		}
	})
}
