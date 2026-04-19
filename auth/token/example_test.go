/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token_test

import (
	"fmt"

	"chainguard.dev/sdk/auth/token"
)

// ExampleKindAccess demonstrates the KindAccess constant value.
func ExampleKindAccess() {
	fmt.Println(token.KindAccess)
	// Output:
	// oidc-token
}

// ExampleKindRefresh demonstrates the KindRefresh constant value.
func ExampleKindRefresh() {
	fmt.Println(token.KindRefresh)
	// Output:
	// refresh-token
}

// ExampleWithAlias demonstrates creating an Option that sets a token alias.
func ExampleWithAlias() {
	opt := token.WithAlias("my-alias")
	fmt.Println(opt != nil)
	// Output:
	// true
}
