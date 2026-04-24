/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package platform_test

import (
	"context"
	"fmt"

	"chainguard.dev/sdk/proto/platform"
)

// ExampleWithUserAgent demonstrates storing a user agent string in context.
func ExampleWithUserAgent() {
	ctx := context.Background()
	ctx = platform.WithUserAgent(ctx, "my-client/1.0")
	fmt.Println(platform.GetUserAgent(ctx))
	// Output:
	// my-client/1.0
}

// ExampleGetUserAgent demonstrates retrieving a user agent from an empty context.
func ExampleGetUserAgent() {
	ctx := context.Background()
	fmt.Println(platform.GetUserAgent(ctx) == "")
	// Output:
	// true
}
