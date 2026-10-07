/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"hash/crc32"
	"io"

	"chainguard.dev/sdk/uploads"
)

// ExampleParseEnvelope demonstrates decoding a JSON envelope without
// performing any cryptographic operations.
func ExampleParseEnvelope() {
	payload := `{"ciphertext":"abc=","encryptedKey":"def=","iv":"ghi=","timestamp":"2024-01-01T00:00:00Z","keyVersion":"1"}`
	env, err := uploads.ParseEnvelope(payload)
	fmt.Println(err)
	fmt.Println(env.KeyVersion)
	fmt.Println(env.Timestamp)
	// Output:
	// <nil>
	// 1
	// 2024-01-01T00:00:00Z
}

// ExampleParseEnvelope_invalid demonstrates the error returned when the
// payload is not a valid JSON envelope.
func ExampleParseEnvelope_invalid() {
	_, err := uploads.ParseEnvelope("not-json")
	fmt.Println(err != nil)
	// Output:
	// true
}

// ExampleEncryptionAlgorithm demonstrates the constant that identifies
// the algorithm the client is hardcoded against.
func ExampleEncryptionAlgorithm() {
	fmt.Println(uploads.EncryptionAlgorithm)
	// Output:
	// RSA_DECRYPT_OAEP_3072_SHA256
}

// ExampleSealStream demonstrates sealing a payload as a stream envelope
// while capturing the CRC32C of the ciphertext, then opening it again.
func ExampleSealStream() {
	priv, err := rsa.GenerateKey(rand.Reader, 3072)
	if err != nil {
		panic(err)
	}

	// In production the destination is an object-storage writer and the
	// ciphertext is never held in memory; a buffer stands in here so the
	// example can check the checksum against the finished bytes.
	var sealed bytes.Buffer
	cw := uploads.NewCRC32CWriter(&sealed)
	sw, err := uploads.SealStream(cw, &priv.PublicKey, uploads.EncryptionAlgorithm, "1")
	if err != nil {
		panic(err)
	}
	if _, err := io.WriteString(sw, "hello argos"); err != nil {
		panic(err)
	}
	// Close seals the final chunk. A stream missing it does not open.
	if err := sw.Close(); err != nil {
		panic(err)
	}
	// The incrementally accumulated checksum is ready as soon as the last
	// byte is written, which is what a caller that cannot buffer needs.
	fmt.Println(cw.Sum32() == crc32.Checksum(sealed.Bytes(), crc32.MakeTable(crc32.Castagnoli)))

	// A production callback selects the stored key version in Cloud KMS.
	r, err := uploads.OpenStream(&sealed, func(keyVersion string, wrapped []byte) ([]byte, error) {
		if keyVersion != "1" {
			return nil, fmt.Errorf("unknown key version %q", keyVersion)
		}
		return rsa.DecryptOAEP(sha256.New(), rand.Reader, priv, wrapped, nil)
	})
	if err != nil {
		panic(err)
	}
	plaintext, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", plaintext)
	// Output:
	// true
	// hello argos
}
