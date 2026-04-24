/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws_test

import (
	"context"
	"fmt"

	"chainguard.dev/sdk/auth/aws"
)

// ExampleVerifyToken demonstrates that an invalid token returns an error.
func ExampleVerifyToken() {
	ctx := context.Background()
	_, err := aws.VerifyToken(ctx, "not-a-valid-token")
	fmt.Println(err != nil)
	// Output:
	// true
}
