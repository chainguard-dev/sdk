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
	"encoding/base64"
	"fmt"
	"io"

	auth "chainguard.dev/sdk/proto/platform/auth/v1"
)

var URLSafeEncoding = base64.RawURLEncoding // Headless codes must be URL-safe

// headless.Code is a serialized public key that we use to exchange a shared symmetric key.
// This shared symmetric key is used to encrypt the ID token (see Code#NewSession).
//
// After obtaining the shared symmetric key, we throw away our own private key to guarantee
// that the content of the ID token can only be decrypted by the holder of this code's
// private key.
type Code string

// NewCode creates a code by serializing the public key in an url-safe format.
func NewCode(k *ecdh.PublicKey) Code {
	return Code(marshalPublic(k))
}

func marshalPublic(k *ecdh.PublicKey) string {
	return URLSafeEncoding.EncodeToString(k.Bytes())
}

func parsePublic(encoded string) (*ecdh.PublicKey, error) {
	decoded, err := URLSafeEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}
	return ecdhCurve().NewPublicKey(decoded)
}

// NewSession encrypts the idtoken using a shared symmetric key that is
// only available to us and the holder of the private key corresponding to
// this headless Code.
//
// It is important to recall how ECDH works:
//   - First, the user generates an EC keypair, and send us their public key
//     as the form of a headless login code.
//   - We generate a new ephemeral EC keypair for this session.
//   - With our private key and their public key, a shared symmetric key is
//     obtained by calling ourPriv.ECDH(theirPub).
//   - When we send our public key to the user, they can generate the same
//     shared symmetric key by calling theirPriv.ECDH(ourPub).
//
// The shared symmetric key obtained by ECDH in this function is used to encrypt
// the idtoken. After the idtoken is encrypted, we throw away our private key
// and the shared symmetric key, so that we ourselves cannot decrypt the idtoken
// ourselves.
//
// We then send the user our public key and the encrypted idtoken. As noted before,
// ECDH allows the user to generate the same shared symmetric key, which can be
// used to decrypt the idtoken.
func (h *Code) NewSession(idtoken []byte) (*auth.HeadlessSession, error) {
	theirPub, err := parsePublic(string(*h))
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}
	// First we generate a new ephemeral ECDH key pair for this session.
	ourPriv, err := GenerateKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}
	// With the user's public key, we can generate a shared symmetric key.
	// This shared symmetric key can only be produced by:
	// - ourPriv.ECDH(theirPub) OR
	// - theirPriv.ECDH(ourPub)
	symmetricKey, err := ourPriv.ECDH(theirPub)
	if err != nil {
		return nil, fmt.Errorf("failed to generate symmetric key: %w", err)
	}
	// Encrypt the idtoken with the symmetric key.
	encryptedIdtoken, err := symmetricEncrypt(idtoken, symmetricKey)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt token: %w", err)
	}
	// Throw away the private key and our copy of the symmetric key, so
	// that the symmetric key is only available to the other side by
	// calling theirPriv.ECDH(ourPub).
	ourPub := []byte(marshalPublic(ourPriv.Public().(*ecdh.PublicKey)))
	ourPriv, symmetricKey = nil, nil // nolint
	return &auth.HeadlessSession{
		EcdhPublicKey:    ourPub,
		EncryptedIdtoken: encryptedIdtoken,
	}, nil
}

// symmetricEncrypt uses AES-GCM to encrypt the payload with the provided key.
// The nonce is prepended to the ciphertext.
//
// See https://pkg.go.dev/crypto/cipher#example-NewGCM-Encrypt
func symmetricEncrypt(payload, key []byte) ([]byte, error) {
	if len(key) != ECDHKeyLength {
		return nil, fmt.Errorf("invalid key size %d != %d", len(key), ECDHKeyLength)
	}
	// First create a block cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}
	// The above is only usable for one block, not a stream.
	// Now create a GCM stream cipher based on that.
	//
	// The advantage of GCM vs CBC is that it provides both encryption and
	// authentication, so we can detect if the data has been tampered with,
	// instead of just decrypting garbage.
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM cipher: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	// gcm.Seal appends the result to the first parameter, which is our
	// nonce. So the result is nonce + ciphertext.
	result := gcm.Seal(nonce /* prefix */, nonce, payload, nil)
	return result, nil
}
