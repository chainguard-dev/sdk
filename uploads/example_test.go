/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uploads_test

import (
	"fmt"

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
