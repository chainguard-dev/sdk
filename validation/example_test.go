/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation_test

import (
	"fmt"

	"chainguard.dev/sdk/validation"
)

// ExampleValidateName demonstrates that a valid name passes validation.
func ExampleValidateName() {
	err := validation.ValidateName("my-repo")
	fmt.Println(err)
	// Output:
	// <nil>
}

// ExampleValidateName_invalid demonstrates that an invalid name returns an error.
func ExampleValidateName_invalid() {
	err := validation.ValidateName("INVALID NAME!")
	fmt.Println(err != nil)
	// Output:
	// true
}
