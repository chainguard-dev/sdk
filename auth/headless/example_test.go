/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package headless_test

import (
	"fmt"

	"chainguard.dev/sdk/auth/headless"
)

// ExampleGenerateKeyPair demonstrates generating an ECDH key pair.
func ExampleGenerateKeyPair() {
	pk, err := headless.GenerateKeyPair()
	fmt.Println(err)
	fmt.Println(pk != nil)
	// Output:
	// <nil>
	// true
}

// ExampleNewCode demonstrates creating a headless code from a public key.
func ExampleNewCode() {
	pk, _ := headless.GenerateKeyPair()
	code := headless.NewCode(pk.PublicKey())
	fmt.Println(len(code) > 0)
	// Output:
	// true
}

// ExampleVerifyCode demonstrates that a valid code passes verification.
func ExampleVerifyCode() {
	pk, _ := headless.GenerateKeyPair()
	code := headless.NewCode(pk.PublicKey())
	err := headless.VerifyCode(string(code))
	fmt.Println(err)
	// Output:
	// <nil>
}

// ExampleCode_ConfirmationCode demonstrates deriving and checking a confirmation code.
func ExampleCode_ConfirmationCode() {
	code := headless.Code("uFK0TyUaZCMB6QBGC-HVKGJKqI4-U9q5rV9k8I_YUQ4")
	confirmationCode, err := code.ConfirmationCode()
	fmt.Println(confirmationCode, err)
	fmt.Println(code.MatchesConfirmationCode("xrpm-cqyw"))
	fmt.Println(code.MatchesConfirmationCode("XRPMCQYX"))
	// Output:
	// XRPMCQYW <nil>
	// true
	// false
}
