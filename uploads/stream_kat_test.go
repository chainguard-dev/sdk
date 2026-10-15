/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"strings"
	"testing"

	"chainguard.dev/sdk/uploads"
	"github.com/google/go-cmp/cmp"
)

// The known-answer vectors below were sealed with crypto/cipher directly,
// following the v2 layout rather than calling SealStream: AES-256-GCM
// under katDEK, nonce a0..a6 || uint32 BE index || final flag, and the
// whole header as every chunk's additional data. The wrapped DEK is an
// opaque blob because OpenStream hands it to the unwrap callback
// untouched, so no private key has to be committed to pin the format.
const (
	// magic "CGUS", version 2, prefix a0..a6, chunk size, key version 7,
	// wrapped length 16, then the wrapped blob.
	katHeaderAtStreamChunkSize = "4347555302" + "a0a1a2a3a4a5a6" + "00100000" + "0000000000000007" + "0010" + "777261707065642d64656b2d626c6f62"
	katHeaderAt16ByteChunks    = "4347555302" + "a0a1a2a3a4a5a6" + "00000010" + "0000000000000007" + "0010" + "777261707065642d64656b2d626c6f62"

	katWrapped    = "wrapped-dek-blob"
	katKeyVersion = "7"
)

var katDEK = func() []byte {
	dek := make([]byte, 32)
	for i := range dek {
		dek[i] = byte(i)
	}
	return dek
}()

var katVectors = []struct {
	name      string
	header    string
	chunkSize int
	plaintext string
	chunks    []string
}{{
	// An empty payload is one final chunk holding only its tag.
	name:      "empty",
	header:    katHeaderAtStreamChunkSize,
	chunkSize: uploads.StreamChunkSize,
	chunks:    []string{"7fcc8f4f74ad37a78dc471f63fda06c7"},
}, {
	name:      "single chunk",
	header:    katHeaderAtStreamChunkSize,
	chunkSize: uploads.StreamChunkSize,
	plaintext: "Chainguard envelope v2 known-answer vector",
	chunks:    []string{"34092877d66195a17101f84aa1f8de34687b4ec9b4d05c7ddde47722613e93c98e5f0f3e1ecaf71f1c473e512c46eb7f72d2310001631e86e7e3"},
}, {
	// The reader honors the header's chunk size, which lets a small
	// vector cover the multi-chunk framing: two full chunks sealed
	// non-final, then a short final one.
	name:      "three chunks",
	header:    katHeaderAt16ByteChunks,
	chunkSize: 16,
	plaintext: "forty bytes of plaintext over 16B chunk",
	chunks: []string{
		"6cbe6eddef7427441846b7a0f666761b8e82df8d0e7cc324be389876b1f6d062",
		"3d837ea9f9598efbc0e2c811bf029f066cf1a3a46ecd95c8ca0422da55eed1a5",
		"a38de28683807486739dfee90f1c0f1338a6cec8c58928",
	},
}, {
	// An exact multiple of the chunk size ends on a full chunk sealed
	// final, with no empty chunk after it.
	name:      "exact multiple of the chunk size",
	header:    katHeaderAt16ByteChunks,
	chunkSize: 16,
	plaintext: "thirty-two bytes, two full chunk",
	chunks: []string{
		"7eb975dbe22d68491b4ce4e2e07433182a91d57199405baa3239d4680edd83f7",
		"c526b49e6f0afd4ad833c1d43709c720ac62835408af53e327441e901b29f15b",
	},
}}

// TestOpenStreamKnownAnswer pins the v2 wire format from the reader's
// side. The vectors are fixed bytes, so a change to the header layout,
// the nonce construction, the additional data or the chunk framing fails
// here even if SealStream and OpenStream change in step.
func TestOpenStreamKnownAnswer(t *testing.T) {
	for _, tc := range katVectors {
		t.Run(tc.name, func(t *testing.T) {
			sealed, err := hex.DecodeString(tc.header + strings.Join(tc.chunks, ""))
			if err != nil {
				t.Fatalf("decode vector: %v", err)
			}
			r, err := uploads.OpenStream(bytes.NewReader(sealed), func(keyVersion string, wrapped []byte) ([]byte, error) {
				if keyVersion != katKeyVersion {
					t.Errorf("unwrap key version: got = %q, want = %q", keyVersion, katKeyVersion)
				}
				if string(wrapped) != katWrapped {
					t.Errorf("unwrap wrapped key: got = %q, want = %q", wrapped, katWrapped)
				}
				return bytes.Clone(katDEK), nil
			})
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			if string(got) != tc.plaintext {
				t.Errorf("plaintext: got = %q, want = %q", got, tc.plaintext)
			}
		})
	}
}

// TestKnownAnswerVectorsAreSealedChunks ties the fixed vectors to
// sealChunks, the recomputation TestSealStreamKnownAnswer checks
// SealStream against, so the reader-side and writer-side tests pin one
// format rather than two that merely agree with themselves.
func TestKnownAnswerVectorsAreSealedChunks(t *testing.T) {
	for _, tc := range katVectors {
		t.Run(tc.name, func(t *testing.T) {
			header, err := hex.DecodeString(tc.header)
			if err != nil {
				t.Fatalf("decode header: %v", err)
			}
			chunks := sealChunks(t, katDEK, header[noncePrefixAt:chunkSizeAt], header, []byte(tc.plaintext), tc.chunkSize)
			got := make([]string, 0, len(chunks))
			for _, chunk := range chunks {
				got = append(got, hex.EncodeToString(chunk))
			}
			if diff := cmp.Diff(tc.chunks, got); diff != "" {
				t.Errorf("recomputed chunks (-want, +got):\n%s", diff)
			}
		})
	}
}

// TestSealStreamKnownAnswer pins the v2 wire format from the writer's
// side. SealStream draws its DEK and nonce prefix at random, so its output
// cannot be a fixed vector; instead the test unwraps the DEK with the
// private key, reads the prefix from the header, and recomputes every
// chunk with crypto/cipher. The sealed bytes must equal the header
// followed by exactly those chunks.
func TestSealStreamKnownAnswer(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, tc := range []struct {
		name string
		size int
	}{
		{"empty", 0},
		{"exactly one chunk", uploads.StreamChunkSize},
		{"exactly two chunks", 2 * uploads.StreamChunkSize},
		{"two chunks and a remainder", 2*uploads.StreamChunkSize + 7919},
	} {
		t.Run(tc.name, func(t *testing.T) {
			plaintext := bytesOfLen(tc.size)
			sealed := sealStream(t, &priv.PublicKey, plaintext, 7919)
			if len(sealed) < headerLen {
				t.Fatalf("sealed stream is %d bytes, shorter than its %d-byte header", len(sealed), headerLen)
			}
			header := sealed[:headerLen]

			if got := string(header[:noncePrefixAt]); got != "CGUS\x02" {
				t.Errorf("magic and version: got = %q, want = %q", got, "CGUS\x02")
			}
			if got := binary.BigEndian.Uint32(header[chunkSizeAt:]); got != uploads.StreamChunkSize {
				t.Errorf("chunk size: got = %d, want = %d", got, uploads.StreamChunkSize)
			}
			if got := binary.BigEndian.Uint64(header[keyVersionAt:]); got != 1 {
				t.Errorf("key version: got = %d, want = 1", got)
			}
			if got := binary.BigEndian.Uint16(header[wrappedLenAt:]); int(got) != headerLen-streamHeaderFixedSize {
				t.Errorf("wrapped length: got = %d, want = %d", got, headerLen-streamHeaderFixedSize)
			}
			dek, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, header[streamHeaderFixedSize:], nil)
			if err != nil {
				t.Fatalf("unwrap dek: %v", err)
			}

			want := bytes.Clone(header)
			for _, chunk := range sealChunks(t, dek, header[noncePrefixAt:chunkSizeAt], header, plaintext, uploads.StreamChunkSize) {
				want = append(want, chunk...)
			}
			if !bytes.Equal(sealed, want) {
				t.Errorf("sealed stream differs from the recomputed one: got %d bytes, want %d bytes", len(sealed), len(want))
			}
		})
	}
}

// sealChunks is the STREAM construction written out independently of the
// package. Only an empty payload has an empty final chunk; a payload that
// is an exact multiple of chunkSize ends on a full one.
func sealChunks(t *testing.T, dek, noncePrefix, aad, plaintext []byte, chunkSize int) [][]byte {
	t.Helper()
	block, err := aes.NewCipher(dek)
	if err != nil {
		t.Fatalf("aes cipher: %v", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("aes-gcm: %v", err)
	}
	chunks := make([][]byte, 0, len(plaintext)/chunkSize+1)
	for index := uint32(0); ; index++ {
		n := min(chunkSize, len(plaintext))
		final := n == len(plaintext)
		nonce := binary.BigEndian.AppendUint32(bytes.Clone(noncePrefix), index)
		if final {
			nonce = append(nonce, 1)
		} else {
			nonce = append(nonce, 0)
		}
		chunks = append(chunks, aead.Seal(nil, nonce, plaintext[:n], aad))
		plaintext = plaintext[n:]
		if final {
			return chunks
		}
	}
}
