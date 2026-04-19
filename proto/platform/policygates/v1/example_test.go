/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"context"
	"fmt"

	policygatev1 "chainguard.dev/sdk/proto/platform/policygates/v1"
)

// ExampleNewClients demonstrates that an invalid address returns an error.
func ExampleNewClients() {
	ctx := context.Background()
	_, err := policygatev1.NewClients(ctx, "http://%zz", "")
	fmt.Println(err != nil)
	// Output:
	// true
}
