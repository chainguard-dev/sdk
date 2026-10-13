/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package harness_test

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"math/big"
	"strings"
	"testing"

	"chainguard.dev/sdk/harness"
	"github.com/google/go-cmp/cmp"
)

func TestParseReportPublicKey(t *testing.T) {
	t.Parallel()

	rsaDER := func(t *testing.T, bits int) []byte {
		t.Helper()
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			t.Fatalf("rsa.GenerateKey(%d): %v", bits, err)
		}
		der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			t.Fatalf("x509.MarshalPKIXPublicKey: %v", err)
		}
		return der
	}

	ecdsaDER := func(t *testing.T) []byte {
		t.Helper()
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatalf("ecdsa.GenerateKey: %v", err)
		}
		der, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
		if err != nil {
			t.Fatalf("x509.MarshalPKIXPublicKey: %v", err)
		}
		return der
	}

	ed25519DER := func(t *testing.T) []byte {
		t.Helper()
		pub, _, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			t.Fatalf("ed25519.GenerateKey: %v", err)
		}
		der, err := x509.MarshalPKIXPublicKey(pub)
		if err != nil {
			t.Fatalf("x509.MarshalPKIXPublicKey: %v", err)
		}
		return der
	}

	for _, tc := range []struct {
		name    string
		der     func(*testing.T) []byte
		wantErr string
	}{{
		name: "rsa 3072 is accepted",
		der:  func(t *testing.T) []byte { return rsaDER(t, 3072) },
	}, {
		name:    "empty is not DER",
		der:     func(*testing.T) []byte { return nil },
		wantErr: "not a DER SubjectPublicKeyInfo",
	}, {
		name:    "truncated DER does not parse",
		der:     func(t *testing.T) []byte { return rsaDER(t, harness.ReportKeyBits)[:32] },
		wantErr: "not a DER SubjectPublicKeyInfo",
	}, {
		name:    "ecdsa is not RSA",
		der:     ecdsaDER,
		wantErr: "must be RSA",
	}, {
		name:    "ed25519 is not RSA",
		der:     ed25519DER,
		wantErr: "must be RSA",
	}, {
		name:    "rsa 2048 is too small",
		der:     func(t *testing.T) []byte { return rsaDER(t, 2048) },
		wantErr: "must be 3072 bits, got 2048",
	}, {
		name:    "rsa 4096 is too large",
		der:     func(t *testing.T) []byte { return rsaDER(t, 4096) },
		wantErr: "must be 3072 bits, got 4096",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			der := tc.der(t)
			got, err := harness.ParseReportPublicKey(der)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("ParseReportPublicKey: got nil error, want one containing %q", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Errorf("ParseReportPublicKey error: got = %q, want one containing %q", err, tc.wantErr)
				}
				if got != nil {
					t.Error("ParseReportPublicKey returned a key with an error")
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseReportPublicKey: %v", err)
			}
			// The recorded bytes must round-trip: what the client sent is what a
			// later reader seals the report under.
			back, err := x509.MarshalPKIXPublicKey(got)
			if err != nil {
				t.Fatalf("x509.MarshalPKIXPublicKey: %v", err)
			}
			if diff := cmp.Diff(der, back); diff != "" {
				t.Errorf("round-tripped key (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestParseReportPublicKeyRejectUnusableKeys(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name        string
		exponent    int64
		evenModulus bool
	}{
		{name: "exponent one", exponent: 1},
		{name: "even exponent", exponent: 2},
		{name: "exponent too large", exponent: 1<<32 | 1},
		{name: "even modulus", exponent: 65537, evenModulus: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			n := new(big.Int).Lsh(big.NewInt(1), 3071)
			if !tc.evenModulus {
				n.SetBit(n, 0, 1)
			}
			der, err := x509.MarshalPKIXPublicKey(&rsa.PublicKey{N: n, E: int(tc.exponent)})
			if err != nil {
				t.Fatalf("x509.MarshalPKIXPublicKey: %v", err)
			}
			got, err := harness.ParseReportPublicKey(der)
			if err == nil {
				t.Fatal("ParseReportPublicKey accepted a key that cannot encrypt with RSA-OAEP")
			}
			if !strings.Contains(err.Error(), "cannot encrypt with RSA-OAEP-SHA256") {
				t.Errorf("ParseReportPublicKey error: got = %q, want an RSA-OAEP validation error", err)
			}
			if got != nil {
				t.Error("ParseReportPublicKey returned an unusable key with an error")
			}
		})
	}
}
