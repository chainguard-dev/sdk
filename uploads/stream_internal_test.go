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
	"io"
	"math"
	"strings"
	"testing"
)

// capTestKey returns an RSA key for the chunk-cap tests. 2048 bits
// because these tests only need a key that wraps and unwraps, not the
// 3072-bit modulus the format is deployed with.
func capTestKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	return priv
}

// TestStreamChunkCapIsEnforcedOnBothSides drives the 2^32-chunk cap from
// inside the package. Reaching it over the wire would take 4 PiB of
// payload, so both counters are positioned directly: a black-box test
// here would only be able to assert that the cap is unreachable, which
// is not the same as asserting that it holds.
//
// Both sides must agree on where the cap falls, or a stream one side
// seals is a stream the other refuses. The last index is a legal home
// for the final chunk and for nothing else.
func TestStreamChunkCapIsEnforcedOnBothSides(t *testing.T) {
	priv := capTestKey(t)
	unwrap := func(_ string, wrapped []byte) ([]byte, error) {
		return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, wrapped, nil)
	}

	t.Run("writer refuses a non-final chunk at the last index", func(t *testing.T) {
		w, err := SealStream(io.Discard, &priv.PublicKey, EncryptionAlgorithm, "1")
		if err != nil {
			t.Fatalf("seal stream: %v", err)
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
		var header bytes.Buffer
		w, err := SealStream(&header, &priv.PublicKey, EncryptionAlgorithm, "1")
		if err != nil {
			t.Fatalf("seal stream: %v", err)
		}
		sw := w.(*streamWriter)
		chunk := sw.sc.aead.Seal(nil, sw.sc.nonceFor(math.MaxUint32, false), make([]byte, StreamChunkSize), sw.sc.aad)

		r, err := OpenStream(bytes.NewReader(bytes.Join([][]byte{header.Bytes(), chunk, chunk}, nil)), unwrap)
		if err != nil {
			t.Fatalf("open stream: %v", err)
		}
		sr := r.(*streamReader)
		sr.index = math.MaxUint32

		if _, err := io.ReadAll(sr); err == nil || !strings.Contains(err.Error(), "maximum") {
			t.Fatalf("reading a non-final chunk at the last index: got err = %v, want the chunk cap", err)
		}
	})

	t.Run("reader accepts the final chunk at the last index", func(t *testing.T) {
		var header bytes.Buffer
		w, err := SealStream(&header, &priv.PublicKey, EncryptionAlgorithm, "1")
		if err != nil {
			t.Fatalf("seal stream: %v", err)
		}
		sw := w.(*streamWriter)
		want := []byte("the last chunk a stream may carry")
		chunk := sw.sc.aead.Seal(nil, sw.sc.nonceFor(math.MaxUint32, true), want, sw.sc.aad)

		r, err := OpenStream(bytes.NewReader(bytes.Join([][]byte{header.Bytes(), chunk}, nil)), unwrap)
		if err != nil {
			t.Fatalf("open stream: %v", err)
		}
		sr := r.(*streamReader)
		sr.index = math.MaxUint32

		got, err := io.ReadAll(sr)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("final chunk at the last index: got = %q, want = %q", got, want)
		}
	})
}
