/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"

	uploads "chainguard.dev/sdk/uploads"
)

// keypair generates an RSA-3072 keypair and returns the private key
// plus the PEM-encoded SubjectPublicKeyInfo (the shape Chainguard
// KMS exports via GetOrgPublicKey).
func keypair(t *testing.T) (*rsa.PrivateKey, string) {
	t.Helper()
	priv, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		t.Fatalf("generate rsa key: %v", err)
	}
	pubDER, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		t.Fatalf("marshal pkix public key: %v", err)
	}
	pubPEM := string(pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubDER}))
	return priv, pubPEM
}

// rsaUnwrap mirrors what `cg uploads break-glass` does server-side
// (KMS AsymmetricDecrypt) but locally, against a generated keypair.
func rsaUnwrap(priv *rsa.PrivateKey) func(wrapped []byte) ([]byte, error) {
	return func(wrapped []byte) ([]byte, error) {
		return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, wrapped, nil)
	}
}

func TestSealOpenRoundTrip(t *testing.T) {
	priv, pubPEM := keypair(t)
	for _, tc := range []struct {
		name      string
		plaintext []byte
	}{
		{"empty", []byte{}},
		{"small text", []byte("hello argos")},
		{"newlines", []byte("line one\nline two\n")},
		{"binary", append([]byte{0x00, 0xff, 0xa5}, []byte("after-nul")...)},
		{"1KB", bytesOfLen(1024)},
		{"100KB", bytesOfLen(100 * 1024)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := uploads.SealEnvelope(tc.plaintext, pubPEM, uploads.EncryptionAlgorithm, "1")
			if err != nil {
				t.Fatalf("seal: %v", err)
			}
			got, err := uploads.OpenEnvelope(payload, rsaUnwrap(priv))
			if err != nil {
				t.Fatalf("open: %v", err)
			}
			if string(got) != string(tc.plaintext) {
				t.Fatalf("round-trip mismatch:\n got: %q (len=%d)\nwant: %q (len=%d)", got, len(got), tc.plaintext, len(tc.plaintext))
			}
		})
	}
}

func TestSealEnvelope_RejectsWrongAlgorithm(t *testing.T) {
	_, pubPEM := keypair(t)
	_, err := uploads.SealEnvelope([]byte("x"), pubPEM, "RSA_DECRYPT_OAEP_2048_SHA256", "1")
	if err == nil {
		t.Fatal("expected error on algorithm mismatch, got nil")
	}
	if !strings.Contains(err.Error(), "algorithm") {
		t.Fatalf("expected error to mention algorithm, got: %v", err)
	}
}

func TestSealEnvelope_RejectsMalformedPEM(t *testing.T) {
	_, err := uploads.SealEnvelope([]byte("x"), "not a pem", uploads.EncryptionAlgorithm, "1")
	if err == nil {
		t.Fatal("expected error on malformed PEM, got nil")
	}
}

func TestSealEnvelope_CarriesKeyVersion(t *testing.T) {
	_, pubPEM := keypair(t)
	payload, err := uploads.SealEnvelope([]byte("hi"), pubPEM, uploads.EncryptionAlgorithm, "7")
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	env, err := uploads.ParseEnvelope(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if env.KeyVersion != "7" {
		t.Fatalf("KeyVersion not preserved: got %q, want %q", env.KeyVersion, "7")
	}
}

func TestParseEnvelope_Shape(t *testing.T) {
	priv, pubPEM := keypair(t)
	payload, err := uploads.SealEnvelope([]byte("hi"), pubPEM, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	env, err := uploads.ParseEnvelope(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// Each binary field must be base64-decodable and non-empty.
	for name, b64 := range map[string]string{
		"Ciphertext":   env.Ciphertext,
		"EncryptedKey": env.EncryptedKey,
		"IV":           env.IV,
	} {
		if b64 == "" {
			t.Fatalf("%s is empty", name)
		}
		if _, err := base64.StdEncoding.DecodeString(b64); err != nil {
			t.Fatalf("%s is not valid base64: %v", name, err)
		}
	}
	if env.Timestamp == "" {
		t.Fatal("Timestamp is empty")
	}
	// And the wrapped key must actually be decryptable.
	if _, err := uploads.OpenEnvelope(payload, rsaUnwrap(priv)); err != nil {
		t.Fatalf("open against the parsed envelope failed: %v", err)
	}
}

func TestParseEnvelope_RejectsNonJSON(t *testing.T) {
	_, err := uploads.ParseEnvelope("not json {{{")
	if err == nil {
		t.Fatal("expected error on non-JSON payload, got nil")
	}
}

func TestOpenEnvelope_RejectsTamperedCiphertext(t *testing.T) {
	priv, pubPEM := keypair(t)
	payload, err := uploads.SealEnvelope([]byte("tamper me"), pubPEM, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal: %v", err)
	}

	// Flip one bit in the ciphertext to trigger the GCM auth-tag failure.
	env, err := uploads.ParseEnvelope(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	ctBytes, err := base64.StdEncoding.DecodeString(env.Ciphertext)
	if err != nil {
		t.Fatalf("decode ciphertext: %v", err)
	}
	if len(ctBytes) == 0 {
		t.Fatal("ciphertext is unexpectedly empty")
	}
	ctBytes[0] ^= 0x01
	env.Ciphertext = base64.StdEncoding.EncodeToString(ctBytes)

	tamperedJSON, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("remarshal: %v", err)
	}
	if _, err := uploads.OpenEnvelope(string(tamperedJSON), rsaUnwrap(priv)); err == nil {
		t.Fatal("expected GCM authentication failure on tampered ciphertext, got nil")
	}
}

func TestOpenEnvelope_RejectsBadIVLength(t *testing.T) {
	priv, pubPEM := keypair(t)
	payload, err := uploads.SealEnvelope([]byte("x"), pubPEM, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	env, err := uploads.ParseEnvelope(payload)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	// Replace the 12-byte IV with a 16-byte one — a real-world misuse
	// would be a sender that confused AES-CBC and AES-GCM nonce sizes.
	env.IV = base64.StdEncoding.EncodeToString(make([]byte, 16))
	bad, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("remarshal: %v", err)
	}
	_, err = uploads.OpenEnvelope(string(bad), rsaUnwrap(priv))
	if err == nil {
		t.Fatal("expected error on wrong-length IV, got nil")
	}
	if !strings.Contains(err.Error(), "iv has wrong length") {
		t.Fatalf("expected iv-length error, got: %v", err)
	}
}

func TestOpenEnvelope_RejectsBadKeyLength(t *testing.T) {
	_, pubPEM := keypair(t)
	payload, err := uploads.SealEnvelope([]byte("x"), pubPEM, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		t.Fatalf("seal: %v", err)
	}
	// Unwrap callback returns a 24-byte key (AES-192-like) instead of 32.
	_, err = uploads.OpenEnvelope(payload, func(_ []byte) ([]byte, error) {
		return make([]byte, 24), nil
	})
	if err == nil {
		t.Fatal("expected error on wrong-length AES key, got nil")
	}
	if !strings.Contains(err.Error(), "aes key has wrong length") {
		t.Fatalf("expected aes-key-length error, got: %v", err)
	}
}

func bytesOfLen(n int) []byte {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return b
}
