/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package reportenvelope

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"io"
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// The known-answer vectors below were built with crypto/cipher directly,
// following the layout in doc.go rather than calling seal: AES-256-GCM under
// katDEK, nonce a0..a6 || uint32 BE index || final flag, and the 49-byte
// prefix as every chunk's additional data. The wrapped values are opaque
// strings, so no private key has to be committed to pin the format.
const (
	// "CGRE", version length 32, "harness-report-envelope-v1alpha1", AEAD 1,
	// nonce prefix a0..a6, chunk size 1 MiB.
	katPrefix = "43475245" + "20" + "6861726e6573732d7265706f72742d656e76656c6f70652d7631616c70686131" + "01" + "a0a1a2a3a4a5a6" + "00100000"

	// Epoch 3, one slot: org-kek, key version 7, "org-kek-wrapped-dek".
	katRegionOrg = "0000000000000003" + "01" +
		"01" + "0000000000000007" + "0013" + "6f72672d6b656b2d777261707065642d64656b"

	// Epoch 3, two slots: the org-kek slot above, then customer-submit-key,
	// key version 0, "customer-wrapped-dek".
	katRegionBoth = "0000000000000003" + "02" +
		"01" + "0000000000000007" + "0013" + "6f72672d6b656b2d777261707065642d64656b" +
		"02" + "0000000000000000" + "0014" + "637573746f6d65722d777261707065642d64656b"

	// katLargeSize spans two full chunks and a short final one, too large to
	// commit as hex, so its vector is the SHA-256 of the sealed object.
	katLargeSize = 2*ChunkSize + 1234
)

var katDEK = func() []byte {
	dek := make([]byte, 32)
	for i := range dek {
		dek[i] = byte(i)
	}
	return dek
}()

var katNoncePrefix = []byte{0xa0, 0xa1, 0xa2, 0xa3, 0xa4, 0xa5, 0xa6}

var katOrg = Slot{Kind: OrgKEK, KeyVersion: 7, Wrapped: []byte("org-kek-wrapped-dek")}

var katCustomer = Slot{Kind: CustomerSubmitKey, Wrapped: []byte("customer-wrapped-dek")}

// katLargePlaintext is byte i%251, a period that never lines up with the
// chunk size.
func katLargePlaintext() []byte {
	b := make([]byte, katLargeSize)
	for i := range b {
		b[i] = byte(i % 251)
	}
	return b
}

func TestKnownAnswer(t *testing.T) {
	for _, tc := range []struct {
		name      string
		keys      KeyRegion
		plaintext []byte
		// Exactly one of sealed and sealedSHA256 is set.
		sealed       string
		sealedSHA256 string
	}{{
		// An empty report is one final chunk holding only its tag.
		name:   "empty, org-kek slot only",
		keys:   KeyRegion{Epoch: 3, Slots: []Slot{katOrg}},
		sealed: katPrefix + katRegionOrg + "1ad8bcd87434e00484a36d8f9ec69a41",
	}, {
		name:      "one chunk, both slots",
		keys:      KeyRegion{Epoch: 3, Slots: []Slot{katOrg, katCustomer}},
		plaintext: []byte("harness-report-envelope-v1alpha1 known-answer vector"),
		sealed: katPrefix + katRegionBoth +
			"1f003b70dd7593ed7100a840bdfa963d697d4e85ad92193bc5ba61203c379c8bd95113711fc1b90a1d467ec079ac5089d660b3b54531da6e0ee4eb3385a9248418999ec9",
	}, {
		name:         "three chunks, both slots",
		keys:         KeyRegion{Epoch: 3, Slots: []Slot{katOrg, katCustomer}},
		plaintext:    katLargePlaintext(),
		sealedSHA256: "4b9b02dfa5c924a2b0705b7ca7a4d68d8fec509b7d32cefa022ad1da27d10d62",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			w, err := seal(&buf, katDEK, tc.keys, katNoncePrefix)
			if err != nil {
				t.Fatalf("seal: %v", err)
			}
			if _, err := w.Write(tc.plaintext); err != nil {
				t.Fatalf("Write: %v", err)
			}
			if err := w.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}

			if tc.sealed != "" {
				if got := hex.EncodeToString(buf.Bytes()); got != tc.sealed {
					t.Fatalf("sealed bytes:\ngot  %s\nwant %s", got, tc.sealed)
				}
			} else {
				if got := sha256.Sum256(buf.Bytes()); hex.EncodeToString(got[:]) != tc.sealedSHA256 {
					t.Fatalf("sealed SHA-256: got %x, want %s", got, tc.sealedSHA256)
				}
			}

			r := bytes.NewReader(buf.Bytes())
			h, err := ReadHeader(r)
			if err != nil {
				t.Fatalf("ReadHeader: %v", err)
			}
			if diff := cmp.Diff(tc.keys, h.Keys); diff != "" {
				t.Errorf("key region (-want, +got):\n%s", diff)
			}
			pr, err := h.Open(r, katDEK)
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			got, err := io.ReadAll(pr)
			if err != nil {
				t.Fatalf("read body: %v", err)
			}
			if !bytes.Equal(got, tc.plaintext) {
				t.Error("known-answer vector opens to a different plaintext")
			}
		})
	}
}

func TestKnownAnswerBodyDigest(t *testing.T) {
	var buf bytes.Buffer
	w, err := seal(&buf, katDEK, KeyRegion{Epoch: 3, Slots: []Slot{katOrg, katCustomer}}, katNoncePrefix)
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	if _, err := w.Write(katLargePlaintext()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	got, err := BodyDigest(&buf)
	if err != nil {
		t.Fatalf("BodyDigest: %v", err)
	}
	if want := "ee6284f8514d9f9a210686b5523ed5447b15e5c8406d948452ed77e575ee4746"; hex.EncodeToString(got[:]) != want {
		t.Errorf("BodyDigest: got %x, want %s", got, want)
	}
}

// TestPrefixIsEveryChunksAdditionalData opens each chunk of a three-chunk
// body with crypto/cipher directly. Every chunk opens under the prefix as
// sealed, and none opens once any one prefix byte is changed.
func TestPrefixIsEveryChunksAdditionalData(t *testing.T) {
	dek := make([]byte, 32)
	if _, err := rand.Read(dek); err != nil {
		t.Fatalf("random DEK: %v", err)
	}
	var buf bytes.Buffer
	w, err := Seal(&buf, dek, KeyRegion{Slots: []Slot{katOrg, katCustomer}})
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if _, err := w.Write(make([]byte, 2*ChunkSize+1)); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	r := bytes.NewReader(buf.Bytes())
	h, err := ReadHeader(r)
	if err != nil {
		t.Fatalf("ReadHeader: %v", err)
	}
	body := buf.Bytes()[buf.Len()-r.Len():]

	block, err := aes.NewCipher(dek)
	if err != nil {
		t.Fatalf("aes: %v", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		t.Fatalf("gcm: %v", err)
	}
	frame := ChunkSize + aead.Overhead()
	chunks := [][]byte{body[:frame], body[frame : 2*frame], body[2*frame:]}
	openChunk := func(i int, aad []byte) error {
		nonce := slices.Clone(h.prefix[noncePrefixAt:chunkSizeAt])
		nonce = binary.BigEndian.AppendUint32(nonce, uint32(i))
		final := byte(0)
		if i == len(chunks)-1 {
			final = 1
		}
		_, err := aead.Open(nil, append(nonce, final), chunks[i], aad)
		return err
	}

	for i := range chunks {
		if err := openChunk(i, h.prefix); err != nil {
			t.Fatalf("chunk %d under the sealed prefix: %v", i, err)
		}
	}
	for b := range prefixSize {
		aad := slices.Clone(h.prefix)
		aad[b] ^= 1
		for i := range chunks {
			if openChunk(i, aad) == nil {
				t.Errorf("prefix byte %d changed: chunk %d still opens", b, i)
			}
		}
	}
}
