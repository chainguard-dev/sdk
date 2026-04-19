/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"context"
	"fmt"

	advisoryv1 "chainguard.dev/sdk/proto/platform/advisory/v1"
)

// ExampleNewClients demonstrates that an invalid address returns an error.
func ExampleNewClients() {
	ctx := context.Background()
	_, err := advisoryv1.NewClients(ctx, "http://%zz", "")
	fmt.Println(err != nil)
	// Output:
	// true
}
