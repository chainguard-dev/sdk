/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth_test

import (
	"context"
	"fmt"

	"chainguard.dev/sdk/auth"
)

// ExampleGetToken demonstrates retrieving a token from a context that has none.
func ExampleGetToken() {
	ctx := context.Background()
	fmt.Println(auth.GetToken(ctx) == "")
	// Output:
	// true
}

// ExampleNewFromFile demonstrates creating credentials from a token file path.
func ExampleNewFromFile() {
	ctx := context.Background()
	// Returns nil when the file does not exist.
	cred := auth.NewFromFile(ctx, "/nonexistent/token", false)
	fmt.Println(cred == nil)
	// Output:
	// true
}
