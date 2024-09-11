/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package headless

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/rand"
	"errors"
	"fmt"

	auth "chainguard.dev/sdk/proto/platform/auth/v1"
)

const (
	// X25519 keys are 32 bytes long
	ECDHKeyLength = 32
)

// ecdhCurve returns the ecdhCurve used for all ECDH operations in this library.
// It is important to use the same ecdhCurve for all operations in here, as
// the ECDH() calls only work when both sides use the same ecdhCurve.
//
// This is effectively a constant.
func ecdhCurve() ecdh.Curve {
	// X25519 gives really short public keys, which is nice as a URL-embedded headless code.
	return ecdh.X25519()
}

// GenerateKeyPair generates a new ECDSA key pair.
func GenerateKeyPair() (*ecdh.PrivateKey, error) {
	return ecdhCurve().GenerateKey(rand.Reader)
}

// DecryptIDToken decrypts the ID token using the private key.
func DecryptIDToken(sess *auth.HeadlessSession, pk *ecdh.PrivateKey) ([]byte, error) {
	serverPub, err := parsePublic(string(sess.EcdhPublicKey))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	symKey, err := pk.ECDH(serverPub)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt key: %w", err)
	}
	idtoken, err := symmetricDecrypt(sess.EncryptedIdtoken, symKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt id token: %w", err)
	}
	return idtoken, nil
}

// symmetricDecrypt decrypts the payload using the key, using AES-GCM.
// The payload is expected to be the output of #symmetricEncrypt, and
// having the nonce prepended.
//
// See https://pkg.go.dev/crypto/cipher#example-NewGCM-Decrypt
func symmetricDecrypt(ciphertext, key []byte) ([]byte, error) {
	if len(key) != ECDHKeyLength {
		return nil, fmt.Errorf("invalid key size %d != %d", len(key), ECDHKeyLength)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}
	// We know the ciphertext is actually nonce+ciphertext
	// and len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt using AES-GCM: %w", err)
	}
	return plain, nil
}
