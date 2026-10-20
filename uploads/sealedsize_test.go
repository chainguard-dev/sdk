/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
	"math"
	"math/big"
	"strings"
	"sync"
	"testing"

	uploads "chainguard.dev/sdk/uploads"
)

// sizedKeys holds one cached key per modulus size. The wrapped DEK is one
// modulus wide, so each size gives a different header length.
var sizedKeys = []struct {
	bits int
	key  func() (*rsa.PrivateKey, error)
}{
	{2048, generateKeyOnce(2048)},
	{3072, generateKeyOnce(3072)},
	{4096, generateKeyOnce(4096)},
}

func generateKeyOnce(bits int) func() (*rsa.PrivateKey, error) {
	return sync.OnceValues(func() (*rsa.PrivateKey, error) {
		return rsa.GenerateKey(rand.Reader, bits)
	})
}

// countingWriter counts the bytes written to it and keeps none of them.
type countingWriter struct{ n int64 }

func (c *countingWriter) Write(p []byte) (int, error) {
	c.n += int64(len(p))
	return len(p), nil
}

// zeroReader yields an endless run of zero bytes.
type zeroReader struct{}

func (zeroReader) Read(p []byte) (int, error) {
	clear(p)
	return len(p), nil
}

// sealedLen seals n zero bytes under pub and returns the stream's length.
func sealedLen(t *testing.T, pub *rsa.PublicKey, n int64) int64 {
	t.Helper()
	var cw countingWriter
	sw, err := uploads.SealStream(&cw, pub, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal stream: %v", err)
	}
	if _, err := io.CopyN(sw, zeroReader{}, n); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := sw.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	return cw.n
}

func TestSealedSizeMatchesSealStream(t *testing.T) {
	const c = uploads.StreamChunkSize
	for _, key := range sizedKeys {
		t.Run(fmt.Sprintf("rsa-%d", key.bits), func(t *testing.T) {
			priv, err := key.key()
			if err != nil {
				t.Fatalf("generate key: %v", err)
			}
			for _, tc := range []struct {
				name string
				n    int64
			}{
				{"empty", 0},
				{"one byte", 1},
				{"one byte under a chunk", c - 1},
				{"exactly one chunk", c},
				{"one byte over a chunk", c + 1},
				{"exactly two chunks", 2 * c},
				{"three chunks and a remainder", 3*c + 5},
			} {
				t.Run(tc.name, func(t *testing.T) {
					got, err := uploads.SealedSize(&priv.PublicKey, tc.n)
					if err != nil {
						t.Fatalf("sealed size: %v", err)
					}
					if want := sealedLen(t, &priv.PublicKey, tc.n); got != want {
						t.Errorf("sealed size: got = %d, want = %d (SealStream's length)", got, want)
					}
				})
			}
		})
	}
}

// TestSealedSizeBounds covers the inputs too large or malformed to seal,
// and the largest of each that still seals. None of them can be checked
// against SealStream: the cap rows would seal 4 PiB.
func TestSealedSizeBounds(t *testing.T) {
	// modulus returns a public key whose modulus is exactly size bytes.
	// SealedSize reads only its width, so it need not be a real key.
	modulus := func(size int) *rsa.PublicKey {
		return &rsa.PublicKey{N: new(big.Int).Lsh(big.NewInt(1), uint(size*8-1)), E: 65537}
	}
	const (
		maxChunks = 1 << 32
		maxN      = maxChunks * uploads.StreamChunkSize
	)
	for _, tc := range []struct {
		name    string
		pub     *rsa.PublicKey
		n       int64
		want    int64
		wantErr string
	}{{
		name:    "nil key",
		pub:     nil,
		wantErr: "nil public key",
	}, {
		name:    "nil modulus",
		pub:     &rsa.PublicKey{E: 65537},
		wantErr: "nil public key",
	}, {
		name:    "negative size",
		pub:     modulus(384),
		n:       -1,
		wantErr: "negative",
	}, {
		// 2048 bytes is the widest wrapped DEK the header admits.
		name: "widest modulus",
		pub:  modulus(2048),
		want: streamHeaderFixedSize + 2048 + gcmTagSize,
	}, {
		name:    "modulus one byte too wide",
		pub:     modulus(2049),
		wantErr: "wrapped dek",
	}, {
		// The last chunk index holds a full final chunk.
		name: "largest plaintext under the chunk cap",
		pub:  modulus(384),
		n:    maxN,
		want: streamHeaderFixedSize + 384 + maxN + maxChunks*gcmTagSize,
	}, {
		name:    "one byte past the chunk cap",
		pub:     modulus(384),
		n:       maxN + 1,
		wantErr: "chunks",
	}, {
		name:    "largest int64",
		pub:     modulus(384),
		n:       math.MaxInt64,
		wantErr: "chunks",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := uploads.SealedSize(tc.pub, tc.n)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("sealed size: got (%d, %v), want an error containing %q", got, err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("sealed size: %v", err)
			}
			if got != tc.want {
				t.Errorf("sealed size: got = %d, want = %d", got, tc.want)
			}
		})
	}
}
