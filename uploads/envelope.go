/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package uploads provides helpers for client-side-encrypted upload
// storage: the wire envelope shape and Seal / Parse / Open functions
// that mirror the browser's encryption.ts. Consumers:
//
//   - Client-side, on upload: Seal a fresh AES-256 key around the
//     plaintext, wrap the key under the org's RSA-OAEP public key, and
//     ship the JSON envelope.
//   - Server-side, on upload: Parse the envelope to validate its shape
//     without ever touching the plaintext.
//   - Client-side, on break-glass read: Open the envelope, which
//     verifies the JSON shell with Parse and AES-GCM-decrypts the
//     ciphertext using an injected unwrap callback (typically a KMS
//     AsymmetricDecrypt against the version recorded in the envelope).
package uploads

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"time"
)

// EncryptionAlgorithm is the KMS algorithm string the client is hardcoded
// against. The org public key endpoint advertises this; SealEnvelope
// refuses to run against any other value rather than silently producing
// payloads no reader can decrypt.
const EncryptionAlgorithm = "RSA_DECRYPT_OAEP_3072_SHA256"

const (
	// aesKeySize is the AES-256-GCM key length in bytes. SealEnvelope
	// generates a key of this size and OpenEnvelope refuses to proceed
	// with a wrapped-then-unwrapped key of any other length (defense
	// against a buggy unwrap callback returning truncated bytes).
	aesKeySize = 32

	// gcmIVSize is the AES-GCM standard nonce length in bytes (96 bits).
	// OpenEnvelope validates the decoded IV against this so a malformed
	// envelope fails fast rather than further down inside cipher.NewGCM.
	gcmIVSize = 12
)

// Envelope is the wire shape the client encrypts and the server stores
// as-is. All binary fields are base64-encoded; the server treats the
// JSON-encoded envelope as opaque and never inspects its contents.
type Envelope struct {
	Ciphertext   string `json:"ciphertext"`
	EncryptedKey string `json:"encryptedKey"`
	IV           string `json:"iv"`
	Timestamp    string `json:"timestamp"`

	// KeyVersion is the integer CryptoKeyVersion the client used for
	// RSA-OAEP wrapping (matches Cloud KMS' version numbering, e.g. "1").
	// Carried in the envelope so decrypters can pick the right version
	// regardless of what's currently the published primary — required
	// for any post-rotation read of pre-rotation blobs. Empty means "1"
	// for backward compatibility with envelopes uploaded before this
	// field existed.
	KeyVersion string `json:"keyVersion,omitempty"`
}

// ParseEnvelope decodes the JSON envelope without performing the
// RSA-OAEP/AES-GCM decryption. Useful on the upload path when the server
// wants to read non-sensitive metadata (Timestamp) without touching
// the plaintext.
func ParseEnvelope(payload string) (*Envelope, error) {
	var env Envelope
	if err := json.Unmarshal([]byte(payload), &env); err != nil {
		return nil, fmt.Errorf("payload is not a valid JSON envelope: %w", err)
	}
	return &env, nil
}

// SealEnvelope produces the JSON-encoded payload string the upload
// service expects in the request's payload field. It mirrors the
// browser's hybrid AES-GCM / RSA-OAEP encryption minus the network call:
//
//  1. Generate a fresh AES-256 key + 12-byte IV.
//  2. AES-256-GCM encrypt the plaintext.
//  3. Parse the SPKI PEM into an RSA public key.
//  4. RSA-OAEP-SHA256 wrap the AES key.
//  5. Base64-encode each binary field, marshal the envelope.
//
// The algorithm parameter is the value advertised by the server's
// GetOrgPublicKey response; mismatch against the client's hardcoded
// EncryptionAlgorithm is a hard failure rather than silently producing
// payloads no reader can decrypt.
//
// keyVersion is the CryptoKeyVersion (e.g. "1") that wrapped this
// envelope, recorded so post-rotation readers know which private key
// to call AsymmetricDecrypt against.
func SealEnvelope(plaintext []byte, publicKeyPEM, algorithm, keyVersion string) (string, error) {
	if algorithm != EncryptionAlgorithm {
		return "", fmt.Errorf("server advertises encryption algorithm %q, client expects %q", algorithm, EncryptionAlgorithm)
	}

	pub, err := parseSPKIPublicKey(publicKeyPEM)
	if err != nil {
		return "", fmt.Errorf("parse public key: %w", err)
	}

	aesKey := make([]byte, aesKeySize)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return "", fmt.Errorf("generate AES key: %w", err)
	}

	iv := make([]byte, gcmIVSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("generate IV: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("aes-gcm: %w", err)
	}
	// Seal appends ciphertext (including auth tag) to dst — using nil
	// destination matches the browser's encrypt() output shape.
	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	// RSA-OAEP-SHA256 wrap of the AES key. No label (matches the
	// browser's RSA-OAEP with default empty label).
	wrappedKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pub, aesKey, nil)
	if err != nil {
		return "", fmt.Errorf("rsa-oaep wrap aes key: %w", err)
	}

	env := Envelope{
		Ciphertext:   base64.StdEncoding.EncodeToString(ciphertext),
		EncryptedKey: base64.StdEncoding.EncodeToString(wrappedKey),
		IV:           base64.StdEncoding.EncodeToString(iv),
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		KeyVersion:   keyVersion,
	}
	out, err := json.Marshal(env)
	if err != nil {
		return "", fmt.Errorf("marshal envelope: %w", err)
	}
	return string(out), nil
}

// OpenEnvelope is the inverse of SealEnvelope. The caller supplies an
// unwrap function that decrypts the RSA-OAEP-wrapped AES key (typically
// a KMS AsymmetricDecrypt call against the version recorded in
// env.KeyVersion); this package stays GCP-free.
//
// The split lets break-glass / reader code reuse the same AES-GCM /
// base64 logic the upload side uses, without forcing every SDK consumer
// to pull in cloud.google.com/go/kms.
func OpenEnvelope(payload string, unwrapAESKey func(wrappedKey []byte) (aesKey []byte, err error)) ([]byte, error) {
	env, err := ParseEnvelope(payload)
	if err != nil {
		return nil, err
	}
	wrapped, err := base64.StdEncoding.DecodeString(env.EncryptedKey)
	if err != nil {
		return nil, fmt.Errorf("decode wrapped key: %w", err)
	}
	iv, err := base64.StdEncoding.DecodeString(env.IV)
	if err != nil {
		return nil, fmt.Errorf("decode iv: %w", err)
	}
	if len(iv) != gcmIVSize {
		return nil, fmt.Errorf("iv has wrong length: got %d, want %d", len(iv), gcmIVSize)
	}
	ct, err := base64.StdEncoding.DecodeString(env.Ciphertext)
	if err != nil {
		return nil, fmt.Errorf("decode ciphertext: %w", err)
	}

	aesKey, err := unwrapAESKey(wrapped)
	if err != nil {
		return nil, fmt.Errorf("unwrap aes key: %w", err)
	}
	if l := len(aesKey); l != aesKeySize {
		return nil, fmt.Errorf("unwrapped aes key has wrong length: got %d, want %d", l, aesKeySize)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("aes-gcm: %w", err)
	}
	plaintext, err := gcm.Open(nil, iv, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("aes-gcm open: %w", err)
	}
	return plaintext, nil
}

// parseSPKIPublicKey decodes a PEM-encoded SubjectPublicKeyInfo (the
// format returned by GetOrgPublicKey) into an *rsa.PublicKey.
func parseSPKIPublicKey(pemStr string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return nil, fmt.Errorf("no PEM block found in input")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse PKIX public key: %w", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not RSA (got %T)", pub)
	}
	return rsaPub, nil
}
