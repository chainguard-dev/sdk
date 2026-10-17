/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package reportenvelope_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	"chainguard.dev/sdk/harness/reportenvelope"
	"chainguard.dev/sdk/uploads"
	"github.com/google/go-cmp/cmp"
)

// frameSize is a full chunk's ciphertext length: the plaintext plus the
// AES-GCM tag.
const frameSize = reportenvelope.ChunkSize + 16

func randomBytes(t *testing.T, n int) []byte {
	t.Helper()
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("random bytes: %v", err)
	}
	return b
}

func randomUint64(t *testing.T) uint64 {
	t.Helper()
	return binary.BigEndian.Uint64(randomBytes(t, 8))
}

// orgSlot and customerSlot carry random 384-byte wrapped values, the size of
// an RSA-3072 OAEP ciphertext. The package treats them as opaque.
func orgSlot(t *testing.T) reportenvelope.Slot {
	t.Helper()
	return reportenvelope.Slot{
		Kind:       reportenvelope.OrgKEK,
		KeyVersion: randomUint64(t)>>1 + 1,
		Wrapped:    randomBytes(t, 384),
	}
}

func customerSlot(t *testing.T) reportenvelope.Slot {
	t.Helper()
	return reportenvelope.Slot{
		Kind:    reportenvelope.CustomerSubmitKey,
		Wrapped: randomBytes(t, 384),
	}
}

func seal(t *testing.T, dek []byte, keys reportenvelope.KeyRegion, plaintext []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w, err := reportenvelope.Seal(&buf, dek, keys)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if _, err := w.Write(plaintext); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	return buf.Bytes()
}

// open reads the header and the whole body, returning the first error.
func open(sealed, dek []byte) (*reportenvelope.Header, []byte, error) {
	r := bytes.NewReader(sealed)
	h, err := reportenvelope.ReadHeader(r)
	if err != nil {
		return nil, nil, err
	}
	pr, err := h.Open(r, dek)
	if err != nil {
		return nil, nil, err
	}
	plaintext, err := io.ReadAll(pr)
	return h, plaintext, err
}

// split returns an envelope's header and body, found by ReadHeader.
func split(t *testing.T, sealed []byte) (header, body []byte) {
	t.Helper()
	r := bytes.NewReader(sealed)
	if _, err := reportenvelope.ReadHeader(r); err != nil {
		t.Fatalf("ReadHeader: %v", err)
	}
	at := len(sealed) - r.Len()
	return sealed[:at], sealed[at:]
}

func TestRoundTrip(t *testing.T) {
	for _, slots := range []int{1, 2} {
		for _, size := range []int{0, 1, 100, reportenvelope.ChunkSize - 1, reportenvelope.ChunkSize, reportenvelope.ChunkSize + 1, 2*reportenvelope.ChunkSize + 17} {
			t.Run(fmt.Sprintf("%d slots, %d bytes", slots, size), func(t *testing.T) {
				keys := reportenvelope.KeyRegion{Epoch: randomUint64(t), Slots: []reportenvelope.Slot{orgSlot(t)}}
				if slots == 2 {
					keys.Slots = append(keys.Slots, customerSlot(t))
				}
				dek := randomBytes(t, 32)
				plaintext := randomBytes(t, size)

				h, got, err := open(seal(t, dek, keys, plaintext), dek)
				if err != nil {
					t.Fatalf("open: %v", err)
				}
				if !bytes.Equal(got, plaintext) {
					t.Error("plaintext differs after round trip")
				}
				if diff := cmp.Diff(keys, h.Keys); diff != "" {
					t.Errorf("key region (-want, +got):\n%s", diff)
				}
			})
		}
	}
}

func TestBodyTamperFailsToOpen(t *testing.T) {
	dek := randomBytes(t, 32)
	keys := reportenvelope.KeyRegion{Epoch: 1, Slots: []reportenvelope.Slot{orgSlot(t), customerSlot(t)}}
	plaintext := randomBytes(t, 2*reportenvelope.ChunkSize+100)
	header, body := split(t, seal(t, dek, keys, plaintext))
	if want := 2*frameSize + 100 + 16; len(body) != want {
		t.Fatalf("body is %d bytes, want %d", len(body), want)
	}
	chunk := func(b []byte, i int) []byte { return b[i*frameSize : min((i+1)*frameSize, len(b))] }
	// other is the same plaintext sealed under the same DEK and key region;
	// only its nonce prefix differs.
	_, other := split(t, seal(t, dek, keys, plaintext))

	for _, tc := range []struct {
		name string
		body []byte
	}{{
		name: "body removed",
	}, {
		name: "final chunk removed",
		body: body[:2*frameSize],
	}, {
		name: "cut inside the final chunk",
		body: body[:len(body)-1],
	}, {
		name: "cut inside a full chunk",
		body: body[:frameSize+10],
	}, {
		name: "chunks reordered",
		body: slices.Concat(chunk(body, 1), chunk(body, 0), chunk(body, 2)),
	}, {
		name: "chunk repeated",
		body: slices.Concat(chunk(body, 0), chunk(body, 0), chunk(body, 1), chunk(body, 2)),
	}, {
		name: "chunk spliced from another envelope",
		body: slices.Concat(chunk(body, 0), chunk(other, 1), chunk(body, 2)),
	}, {
		name: "bytes appended",
		body: slices.Concat(body, []byte{0}),
	}, {
		name: "ciphertext byte flipped",
		body: func() []byte {
			b := slices.Clone(body)
			b[frameSize+5] ^= 1
			return b
		}(),
	}, {
		name: "tag byte flipped",
		body: func() []byte {
			b := slices.Clone(body)
			b[len(b)-1] ^= 1
			return b
		}(),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			if _, _, err := open(slices.Concat(header, tc.body), dek); err == nil {
				t.Error("open: got nil error, want a failure")
			}
		})
	}

	t.Run("wrong DEK", func(t *testing.T) {
		if _, _, err := open(slices.Concat(header, body), randomBytes(t, 32)); err == nil {
			t.Error("open: got nil error, want a failure")
		}
	})
}

func TestPrefixByteChangeFailsToOpen(t *testing.T) {
	dek := randomBytes(t, 32)
	keys := reportenvelope.KeyRegion{Epoch: 1, Slots: []reportenvelope.Slot{orgSlot(t)}}
	sealed := seal(t, dek, keys, randomBytes(t, reportenvelope.ChunkSize+1))
	if _, _, err := open(sealed, dek); err != nil {
		t.Fatalf("unmodified envelope: %v", err)
	}
	// The prefix is the first 49 bytes; envelope_internal_test.go shows each
	// byte is in every chunk's additional data.
	for i := range 49 {
		modified := slices.Clone(sealed)
		modified[i] ^= 1
		if _, _, err := open(modified, dek); err == nil {
			t.Errorf("prefix byte %d flipped: open got nil error, want a failure", i)
		}
	}
}

// TestMisalignedKeyRegionFailsToOpen covers a key region whose slot count
// disagrees with the slots that follow it, but which still parses: the body
// is then read from the wrong offset and chunk 0 fails.
func TestMisalignedKeyRegionFailsToOpen(t *testing.T) {
	dek := randomBytes(t, 32)
	plaintext := randomBytes(t, 100)
	org := orgSlot(t)
	cust := customerSlot(t)
	twoSlots := seal(t, dek, reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{org, cust}}, plaintext)
	oneSlot := seal(t, dek, reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{org}}, plaintext)
	// The slot count follows the 49-byte prefix and the 8-byte epoch.
	const countAt = 49 + 8

	for _, tc := range []struct {
		name   string
		sealed []byte
	}{{
		name: "count lowered, customer slot left in place",
		sealed: func() []byte {
			b := slices.Clone(twoSlots)
			b[countAt] = 1
			return b
		}(),
	}, {
		name: "customer slot inserted, count unchanged",
		sealed: func() []byte {
			header, body := split(t, oneSlot)
			slot := binary.BigEndian.AppendUint64([]byte{byte(reportenvelope.CustomerSubmitKey)}, 0)
			slot = binary.BigEndian.AppendUint16(slot, uint16(len(cust.Wrapped)))
			return slices.Concat(header, slot, cust.Wrapped, body)
		}(),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			r := bytes.NewReader(tc.sealed)
			h, err := reportenvelope.ReadHeader(r)
			if err != nil {
				t.Fatalf("ReadHeader: %v", err)
			}
			pr, err := h.Open(r, dek)
			if err != nil {
				t.Fatalf("Open: %v", err)
			}
			if _, err := io.ReadAll(pr); err == nil {
				t.Error("read body: got nil error, want a failure")
			}
		})
	}
}

// region encodes a key region by hand, so that the reader's refusals are
// tested against bytes the writer would never produce.
type rawSlot struct {
	kind       byte
	keyVersion uint64
	wrappedLen uint16
	wrapped    []byte
}

func region(count byte, slots ...rawSlot) []byte {
	b := binary.BigEndian.AppendUint64(nil, 9)
	b = append(b, count)
	for _, s := range slots {
		b = append(b, s.kind)
		b = binary.BigEndian.AppendUint64(b, s.keyVersion)
		b = binary.BigEndian.AppendUint16(b, s.wrappedLen)
		b = append(b, s.wrapped...)
	}
	return b
}

func TestReadHeaderRefusesKeyRegion(t *testing.T) {
	dek := randomBytes(t, 32)
	keys := reportenvelope.KeyRegion{Epoch: 1, Slots: []reportenvelope.Slot{orgSlot(t)}}
	header, body := split(t, seal(t, dek, keys, randomBytes(t, 100)))
	prefix := header[:49]

	wrapped := randomBytes(t, 384)
	org := rawSlot{kind: 1, keyVersion: 7, wrappedLen: 384, wrapped: wrapped}
	cust := rawSlot{kind: 2, wrappedLen: 384, wrapped: wrapped}

	for _, tc := range []struct {
		name    string
		region  []byte
		body    []byte
		wantErr string
	}{{
		// A control row: the hand encoding matches the writer's.
		name:   "org-kek and customer-submit-key",
		region: region(2, org, cust),
		body:   body,
	}, {
		name:    "no slots declared",
		region:  region(0),
		body:    body,
		wantErr: "0 key slots",
	}, {
		name:    "three slots declared",
		region:  region(3, org, cust, cust),
		body:    body,
		wantErr: "3 key slots",
	}, {
		name:    "declared customer slot missing before the body",
		region:  region(2, org),
		body:    body,
		wantErr: "slot 1",
	}, {
		name:    "declared customer slot missing at the end",
		region:  region(2, org),
		wantErr: "read slot 1 of 2",
	}, {
		name:    "org-kek slot missing",
		region:  region(1, cust),
		body:    body,
		wantErr: "slot 0 is customer-submit-key, want org-kek",
	}, {
		name:    "customer-submit-key slot first",
		region:  region(2, cust, org),
		body:    body,
		wantErr: "slot 0 is customer-submit-key, want org-kek",
	}, {
		name:    "org-kek slot twice",
		region:  region(2, org, org),
		body:    body,
		wantErr: "slot 1 is org-kek, want customer-submit-key",
	}, {
		name:    "unknown slot kind",
		region:  region(2, org, rawSlot{kind: 3, wrappedLen: 384, wrapped: wrapped}),
		body:    body,
		wantErr: "slot 1 is SlotKind(3)",
	}, {
		name:    "org-kek slot with key version zero",
		region:  region(1, rawSlot{kind: 1, wrappedLen: 384, wrapped: wrapped}),
		body:    body,
		wantErr: "key version zero",
	}, {
		name:    "customer-submit-key slot with a key version",
		region:  region(2, org, rawSlot{kind: 2, keyVersion: 1, wrappedLen: 384, wrapped: wrapped}),
		body:    body,
		wantErr: "key version 1, want 0",
	}, {
		name:    "empty wrapped key",
		region:  region(1, rawSlot{kind: 1, keyVersion: 7}),
		body:    body,
		wantErr: "0-byte wrapped key",
	}, {
		name:    "oversized wrapped key",
		region:  region(1, rawSlot{kind: 1, keyVersion: 7, wrappedLen: 2049, wrapped: randomBytes(t, 2049)}),
		body:    body,
		wantErr: "2049-byte wrapped key",
	}, {
		name:    "wrapped key cut short",
		region:  region(1, rawSlot{kind: 1, keyVersion: 7, wrappedLen: 384, wrapped: wrapped[:100]}),
		wantErr: "read slot 0 of 1",
	}, {
		name:    "key region cut short",
		region:  region(1)[:5],
		wantErr: "read key region",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := reportenvelope.ReadHeader(bytes.NewReader(slices.Concat(prefix, tc.region, tc.body)))
			switch {
			case tc.wantErr == "" && err != nil:
				t.Fatalf("ReadHeader: %v", err)
			case tc.wantErr != "" && err == nil:
				t.Fatalf("ReadHeader: got nil error, want one containing %q", tc.wantErr)
			case tc.wantErr != "" && !strings.Contains(err.Error(), tc.wantErr):
				t.Fatalf("ReadHeader: got %v, want an error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestReadHeaderRefusesOtherFormats(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	var upload bytes.Buffer
	w, err := uploads.SealStream(&upload, &key.PublicKey, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("SealStream: %v", err)
	}
	if _, err := w.Write(randomBytes(t, 100)); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	for _, tc := range []struct {
		name    string
		input   []byte
		wantErr string
	}{{
		name:    "v2 upload stream",
		input:   upload.Bytes(),
		wantErr: "not a report envelope",
	}, {
		name:    "empty",
		wantErr: "read magic",
	}, {
		name:    "magic only",
		input:   []byte("CGRE"),
		wantErr: "read prefix",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := reportenvelope.ReadHeader(bytes.NewReader(tc.input))
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("ReadHeader: got %v, want an error containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestWritersRefuseInvalidKeyRegion(t *testing.T) {
	org := orgSlot(t)
	cust := customerSlot(t)
	for _, tc := range []struct {
		name    string
		keys    reportenvelope.KeyRegion
		wantErr string
	}{{
		name:    "no slots",
		wantErr: "0 slots",
	}, {
		name:    "three slots",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{org, cust, cust}},
		wantErr: "3 slots",
	}, {
		name:    "customer-submit-key slot only",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{cust}},
		wantErr: "slot 0 is customer-submit-key",
	}, {
		name:    "slots out of order",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{cust, org}},
		wantErr: "slot 0 is customer-submit-key",
	}, {
		name:    "org-kek slot with key version zero",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{{Kind: reportenvelope.OrgKEK, Wrapped: org.Wrapped}}},
		wantErr: "key version zero",
	}, {
		name:    "customer-submit-key slot with a key version",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{org, {Kind: reportenvelope.CustomerSubmitKey, KeyVersion: 2, Wrapped: cust.Wrapped}}},
		wantErr: "key version 2, want 0",
	}, {
		name:    "empty wrapped key",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{{Kind: reportenvelope.OrgKEK, KeyVersion: 1}}},
		wantErr: "0-byte wrapped key",
	}, {
		name:    "oversized wrapped key",
		keys:    reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{{Kind: reportenvelope.OrgKEK, KeyVersion: 1, Wrapped: randomBytes(t, 2049)}}},
		wantErr: "2049-byte wrapped key",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			if _, err := reportenvelope.Seal(&buf, randomBytes(t, 32), tc.keys); err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("Seal: got %v, want an error containing %q", err, tc.wantErr)
			}
			if buf.Len() != 0 {
				t.Errorf("Seal wrote %d bytes, want 0", buf.Len())
			}

			valid := seal(t, randomBytes(t, 32), reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{org}}, nil)
			if err := reportenvelope.Rewrite(&buf, bytes.NewReader(valid), tc.keys); err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("Rewrite: got %v, want an error containing %q", err, tc.wantErr)
			}
			if buf.Len() != 0 {
				t.Errorf("Rewrite wrote %d bytes, want 0", buf.Len())
			}
		})
	}
}

func TestSealRefusesBadArguments(t *testing.T) {
	keys := reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{orgSlot(t)}}
	if _, err := reportenvelope.Seal(nil, randomBytes(t, 32), keys); err == nil {
		t.Error("Seal with a nil writer: got nil error")
	}
	var buf bytes.Buffer
	if _, err := reportenvelope.Seal(&buf, randomBytes(t, 16), keys); err == nil {
		t.Error("Seal with a 16-byte DEK: got nil error")
	}
	if buf.Len() != 0 {
		t.Errorf("Seal with a 16-byte DEK wrote %d bytes, want 0", buf.Len())
	}
}

func TestZeroHeaderRefused(t *testing.T) {
	var h reportenvelope.Header
	if _, err := h.Open(bytes.NewReader(nil), randomBytes(t, 32)); err == nil {
		t.Error("Open on a zero Header: got nil error")
	}
	var out bytes.Buffer
	if err := h.Rewrite(&out, bytes.NewReader(nil), reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{orgSlot(t)}}); err == nil {
		t.Error("Rewrite on a zero Header: got nil error")
	}
}

func TestOpenTwice(t *testing.T) {
	dek := randomBytes(t, 32)
	plaintext := randomBytes(t, 100)
	sealed := seal(t, dek, reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{orgSlot(t)}}, plaintext)
	h, err := reportenvelope.ReadHeader(bytes.NewReader(sealed))
	if err != nil {
		t.Fatalf("ReadHeader: %v", err)
	}
	_, body := split(t, sealed)
	for i := range 2 {
		pr, err := h.Open(bytes.NewReader(body), dek)
		if err != nil {
			t.Fatalf("Open %d: %v", i, err)
		}
		got, err := io.ReadAll(pr)
		if err != nil {
			t.Fatalf("read body %d: %v", i, err)
		}
		if !bytes.Equal(got, plaintext) {
			t.Errorf("Open %d: plaintext differs", i)
		}
	}
}

// TestHeaderRewrite reads the header, decides on it, and rewrites from the
// same reader, as a re-key does with one read of the object.
func TestHeaderRewrite(t *testing.T) {
	dek := randomBytes(t, 32)
	plaintext := randomBytes(t, 100)
	sealed := seal(t, dek, reportenvelope.KeyRegion{Epoch: 2, Slots: []reportenvelope.Slot{orgSlot(t), customerSlot(t)}}, plaintext)

	r := bytes.NewReader(sealed)
	h, err := reportenvelope.ReadHeader(r)
	if err != nil {
		t.Fatalf("ReadHeader: %v", err)
	}
	org, ok := h.Keys.Slot(reportenvelope.OrgKEK)
	if !ok {
		t.Fatal("no org-kek slot")
	}
	rekeyed := reportenvelope.KeyRegion{Epoch: h.Keys.Epoch + 1, Slots: []reportenvelope.Slot{org, customerSlot(t)}}
	var out bytes.Buffer
	if err := h.Rewrite(&out, r, rekeyed); err != nil {
		t.Fatalf("Rewrite: %v", err)
	}

	got, opened, err := open(out.Bytes(), dek)
	if err != nil {
		t.Fatalf("open the rewritten envelope: %v", err)
	}
	if !bytes.Equal(opened, plaintext) {
		t.Error("rewritten envelope opens to a different plaintext")
	}
	if diff := cmp.Diff(rekeyed, got.Keys); diff != "" {
		t.Errorf("key region (-want, +got):\n%s", diff)
	}
	_, body := split(t, sealed)
	pr, err := h.Open(bytes.NewReader(body), dek)
	if err != nil {
		t.Fatalf("Open after Rewrite: %v", err)
	}
	if _, err := io.ReadAll(pr); err != nil {
		t.Errorf("read body after Rewrite: %v", err)
	}
}

func TestRewrite(t *testing.T) {
	for _, tc := range []struct {
		name string
		from int
		to   int
	}{
		{name: "replace the customer-submit-key slot", from: 2, to: 2},
		{name: "attach a customer-submit-key slot", from: 1, to: 2},
		{name: "replace the org-kek slot alone", from: 1, to: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dek := randomBytes(t, 32)
			plaintext := randomBytes(t, reportenvelope.ChunkSize+500)
			slotsOf := func(n int) []reportenvelope.Slot {
				s := []reportenvelope.Slot{orgSlot(t)}
				if n == 2 {
					s = append(s, customerSlot(t))
				}
				return s
			}
			sealed := seal(t, dek, reportenvelope.KeyRegion{Epoch: 4, Slots: slotsOf(tc.from)}, plaintext)
			rekeyed := reportenvelope.KeyRegion{Epoch: 5, Slots: slotsOf(tc.to)}

			var out bytes.Buffer
			if err := reportenvelope.Rewrite(&out, bytes.NewReader(sealed), rekeyed); err != nil {
				t.Fatalf("Rewrite: %v", err)
			}

			header, body := split(t, sealed)
			gotHeader, gotBody := split(t, out.Bytes())
			if !bytes.Equal(gotBody, body) {
				t.Error("body changed across the rewrite")
			}
			if !bytes.Equal(gotHeader[:49], header[:49]) {
				t.Error("prefix changed across the rewrite")
			}
			h, got, err := open(out.Bytes(), dek)
			if err != nil {
				t.Fatalf("open the rewritten envelope: %v", err)
			}
			if !bytes.Equal(got, plaintext) {
				t.Error("rewritten envelope opens to a different plaintext")
			}
			if diff := cmp.Diff(rekeyed, h.Keys); diff != "" {
				t.Errorf("key region (-want, +got):\n%s", diff)
			}

			before, err := reportenvelope.BodyDigest(bytes.NewReader(sealed))
			if err != nil {
				t.Fatalf("BodyDigest before: %v", err)
			}
			after, err := reportenvelope.BodyDigest(bytes.NewReader(out.Bytes()))
			if err != nil {
				t.Fatalf("BodyDigest after: %v", err)
			}
			if before != after {
				t.Errorf("BodyDigest changed across the rewrite: got %x, want %x", after, before)
			}
		})
	}
}

func TestRewriteRefusesNonEnvelope(t *testing.T) {
	keys := reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{orgSlot(t)}}
	var out bytes.Buffer
	if err := reportenvelope.Rewrite(&out, bytes.NewReader([]byte("CGUS not an envelope")), keys); err == nil {
		t.Error("Rewrite: got nil error")
	}
	if out.Len() != 0 {
		t.Errorf("Rewrite wrote %d bytes, want 0", out.Len())
	}
	valid := seal(t, randomBytes(t, 32), keys, nil)
	if err := reportenvelope.Rewrite(nil, bytes.NewReader(valid), keys); err == nil || !strings.Contains(err.Error(), "nil writer") {
		t.Errorf("Rewrite with a nil writer: got %v, want a nil writer error", err)
	}
}

func TestBodyDigest(t *testing.T) {
	keys := reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{orgSlot(t), customerSlot(t)}}
	sealed := seal(t, randomBytes(t, 32), keys, randomBytes(t, 2*reportenvelope.ChunkSize))
	_, body := split(t, sealed)

	got, err := reportenvelope.BodyDigest(bytes.NewReader(sealed))
	if err != nil {
		t.Fatalf("BodyDigest: %v", err)
	}
	if want := sha256.Sum256(body); got != want {
		t.Errorf("BodyDigest: got %x, want %x", got, want)
	}

	if _, err := reportenvelope.BodyDigest(bytes.NewReader(sealed[:20])); err == nil {
		t.Error("BodyDigest of a truncated header: got nil error")
	}
}

func TestKeyRegionSlot(t *testing.T) {
	org := orgSlot(t)
	cust := customerSlot(t)
	keys := reportenvelope.KeyRegion{Slots: []reportenvelope.Slot{org, cust}}

	got, ok := keys.Slot(reportenvelope.CustomerSubmitKey)
	if !ok {
		t.Fatal("Slot(CustomerSubmitKey): not found")
	}
	if diff := cmp.Diff(cust, got); diff != "" {
		t.Errorf("Slot(CustomerSubmitKey) (-want, +got):\n%s", diff)
	}

	keys.Slots = keys.Slots[:1]
	if _, ok := keys.Slot(reportenvelope.CustomerSubmitKey); ok {
		t.Error("Slot(CustomerSubmitKey) on an org-kek-only region: found")
	}
}
