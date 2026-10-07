/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"bytes"
	"crypto/rsa"
	"encoding/binary"
	"errors"
	"io"
	"math"
	"strings"
	"testing"

	"chainguard.dev/sdk/uploads"
	"github.com/google/go-cmp/cmp"
)

func TestOpenStreamSelectsStoredKeyVersion(t *testing.T) {
	keys := map[string]*rsa.PrivateKey{
		"1":                    streamKeypair(t, 0),
		"7":                    streamKeypair(t, 1),
		"18446744073709551615": streamKeypair(t, 1),
	}
	for _, tc := range []struct {
		version string
		wire    uint64
	}{
		{"1", 1},
		{"7", 7},
		{"18446744073709551615", math.MaxUint64},
	} {
		t.Run(tc.version, func(t *testing.T) {
			want := bytesOfLen(4096)
			var sealed bytes.Buffer
			w, err := uploads.SealStream(&sealed, &keys[tc.version].PublicKey, uploads.EncryptionAlgorithm, tc.version)
			if err != nil {
				t.Fatalf("seal: %v", err)
			}
			if _, err := w.Write(want); err != nil {
				t.Fatalf("write: %v", err)
			}
			if err := w.Close(); err != nil {
				t.Fatalf("close: %v", err)
			}
			if got := binary.BigEndian.Uint64(sealed.Bytes()[keyVersionAt:]); got != tc.wire {
				t.Errorf("wire key version: got = %d, want = %d", got, tc.wire)
			}
			if got := binary.BigEndian.Uint16(sealed.Bytes()[wrappedLenAt:]); got != 3072/8 {
				t.Errorf("wire wrapped key size: got = %d, want = %d", got, 3072/8)
			}

			var versions []string
			r, err := uploads.OpenStream(&sealed, func(version string, wrapped []byte) ([]byte, error) {
				versions = append(versions, version)
				key, ok := keys[version]
				if !ok {
					return nil, errors.New("unknown key version")
				}
				return rsaUnwrap(key)(wrapped)
			})
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			got, err := io.ReadAll(r)
			if err != nil {
				t.Fatalf("read: %v", err)
			}
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("plaintext (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff([]string{tc.version}, versions); diff != "" {
				t.Errorf("unwrap key versions (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestSealStreamRejectsInvalidKeyVersion(t *testing.T) {
	priv := streamKeypair(t, 0)
	for _, version := range []string{"", "0", "-1", "+1", "01", "unknown", "18446744073709551616"} {
		t.Run(version, func(t *testing.T) {
			var sealed bytes.Buffer
			w, err := uploads.SealStream(&sealed, &priv.PublicKey, uploads.EncryptionAlgorithm, version)
			if err == nil || !strings.Contains(err.Error(), "key version") || w != nil {
				t.Errorf("seal: got (%v, %v), want (nil, key version error)", w, err)
			}
			if sealed.Len() != 0 {
				t.Errorf("wrote %d bytes for an invalid key version, want 0", sealed.Len())
			}
		})
	}
}

func TestOpenStreamAuthenticatesKeyVersion(t *testing.T) {
	priv := streamKeypair(t, 0)
	sealed := sealStream(t, &priv.PublicKey, bytesOfLen(4096), 4096)
	binary.BigEndian.PutUint64(sealed[keyVersionAt:], 7)
	// Returning the original DEK isolates the header authentication from
	// failures due to selecting a different RSA key.
	r, err := uploads.OpenStream(bytes.NewReader(sealed), rsaStreamUnwrap(priv))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	got, err := io.ReadAll(r)
	if err == nil || !strings.Contains(err.Error(), "authenticate chunk") {
		t.Errorf("read: got err = %v, want an authentication error", err)
	}
	if len(got) != 0 {
		t.Errorf("released %d plaintext bytes after changing the key version, want 0", len(got))
	}
}

func TestOpenStreamRejectsFailedUnwrap(t *testing.T) {
	priv := streamKeypair(t, 0)
	sealed := sealStream(t, &priv.PublicKey, bytesOfLen(32), 32)
	denied := errors.New("kms permission denied")
	for _, tc := range []struct {
		name string
		key  []byte
		err  error
	}{
		{name: "unwrap error", err: denied},
		{name: "key and unwrap error", key: bytesOfLen(32), err: denied},
		{name: "empty key"},
		{name: "short key", key: bytesOfLen(31)},
		{name: "long key", key: bytesOfLen(33)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r, err := uploads.OpenStream(bytes.NewReader(sealed), func(string, []byte) ([]byte, error) {
				return tc.key, tc.err
			})
			if err == nil || r != nil {
				t.Fatalf("open: got (%v, %v), want (nil, error)", r, err)
			}
			if tc.err != nil {
				if !errors.Is(err, tc.err) {
					t.Errorf("open: got err = %v, want %v", err, tc.err)
				}
			} else if !strings.Contains(err.Error(), "wrong length") {
				t.Errorf("open: got err = %v, want a key length error", err)
			}
		})
	}
}
