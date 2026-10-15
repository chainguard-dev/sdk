/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aeadstream_test

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"
	"testing"

	"chainguard.dev/sdk/internal/aeadstream"
)

func randomBytes(t *testing.T, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("random bytes: %v", err)
	}
	return b
}

func TestNewWriterAndNewReaderRejectBadArguments(t *testing.T) {
	dek := randomBytes(t, aeadstream.KeySize)
	prefix := randomBytes(t, aeadstream.NoncePrefixSize)
	for _, tc := range []struct {
		name      string
		nilStream bool
		dek       []byte
		prefix    []byte
		chunkSize int
		wantErr   string
	}{
		{name: "nil stream", nilStream: true, dek: dek, prefix: prefix, chunkSize: 16, wantErr: "nil"},
		{name: "empty dek", dek: nil, prefix: prefix, chunkSize: 16, wantErr: "dek"},
		{name: "aes-128 dek", dek: dek[:16], prefix: prefix, chunkSize: 16, wantErr: "dek"},
		{name: "dek one byte too long", dek: append(bytes.Clone(dek), 0), prefix: prefix, chunkSize: 16, wantErr: "dek"},
		{name: "prefix one byte too short", dek: dek, prefix: prefix[:aeadstream.NoncePrefixSize-1], chunkSize: 16, wantErr: "nonce prefix"},
		{name: "prefix one byte too long", dek: dek, prefix: append(bytes.Clone(prefix), 0), chunkSize: 16, wantErr: "nonce prefix"},
		{name: "chunk size zero", dek: dek, prefix: prefix, chunkSize: 0, wantErr: "chunk size"},
		{name: "chunk size negative", dek: dek, prefix: prefix, chunkSize: -1, wantErr: "chunk size"},
		{name: "chunk size beyond the cap", dek: dek, prefix: prefix, chunkSize: aeadstream.MaxChunkSize + 1, wantErr: "chunk size"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w := io.Discard
			var r io.Reader = bytes.NewReader(nil)
			if tc.nilStream {
				w, r = nil, nil
			}
			if sw, err := aeadstream.NewWriter(w, tc.dek, nil, tc.prefix, tc.chunkSize); err == nil || !strings.Contains(err.Error(), tc.wantErr) || sw != nil {
				t.Errorf("NewWriter: got (%v, %v), want (nil, error containing %q)", sw, err, tc.wantErr)
			}
			if sr, err := aeadstream.NewReader(r, tc.dek, nil, tc.prefix, tc.chunkSize); err == nil || !strings.Contains(err.Error(), tc.wantErr) || sr != nil {
				t.Errorf("NewReader: got (%v, %v), want (nil, error containing %q)", sr, err, tc.wantErr)
			}
		})
	}

	// The bounds are inclusive.
	for _, chunkSize := range []int{1, aeadstream.MaxChunkSize} {
		if _, err := aeadstream.NewWriter(io.Discard, dek, nil, prefix, chunkSize); err != nil {
			t.Errorf("NewWriter at chunk size %d: %v", chunkSize, err)
		}
		if _, err := aeadstream.NewReader(bytes.NewReader(nil), dek, nil, prefix, chunkSize); err != nil {
			t.Errorf("NewReader at chunk size %d: %v", chunkSize, err)
		}
	}
}

// TestRoundTrip also pins the constructors' handling of their inputs:
// the stream must not depend on the caller's noncePrefix slice after
// construction, and neither side may write to dek, aad or noncePrefix.
func TestRoundTrip(t *testing.T) {
	const chunkSize = 16
	dek := randomBytes(t, aeadstream.KeySize)
	aad := randomBytes(t, 40)
	prefix := randomBytes(t, aeadstream.NoncePrefixSize)
	dekWant, aadWant := bytes.Clone(dek), bytes.Clone(aad)
	// Three full chunks and a short final one, written across chunk
	// boundaries.
	want := randomBytes(t, 3*chunkSize+5)

	var sealed bytes.Buffer
	writerPrefix := bytes.Clone(prefix)
	w, err := aeadstream.NewWriter(&sealed, dek, aad, writerPrefix, chunkSize)
	if err != nil {
		t.Fatalf("new writer: %v", err)
	}
	if !bytes.Equal(writerPrefix, prefix) {
		t.Errorf("writer nonce prefix modified: got %x, want %x", writerPrefix, prefix)
	}
	clear(writerPrefix)
	for rest := want; len(rest) > 0; {
		n := min(7, len(rest))
		if _, err := w.Write(rest[:n]); err != nil {
			t.Fatalf("write: %v", err)
		}
		rest = rest[n:]
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if got, wantLen := sealed.Len(), len(want)+4*16; got != wantLen {
		t.Errorf("sealed length: got = %d, want = %d (plaintext plus one tag per chunk)", got, wantLen)
	}

	readerPrefix := bytes.Clone(prefix)
	r, err := aeadstream.NewReader(&sealed, dek, aad, readerPrefix, chunkSize)
	if err != nil {
		t.Fatalf("new reader: %v", err)
	}
	if !bytes.Equal(readerPrefix, prefix) {
		t.Errorf("reader nonce prefix modified: got %x, want %x", readerPrefix, prefix)
	}
	clear(readerPrefix)
	got, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("round trip: got %x, want %x", got, want)
	}

	for _, c := range []struct {
		name      string
		got, want []byte
	}{{"dek", dek, dekWant}, {"aad", aad, aadWant}} {
		if !bytes.Equal(c.got, c.want) {
			t.Errorf("%s modified: got %x, want %x", c.name, c.got, c.want)
		}
	}
}
