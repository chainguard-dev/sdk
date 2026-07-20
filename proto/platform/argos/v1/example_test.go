/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	argosv1 "chainguard.dev/sdk/proto/platform/argos/v1"
)

// ExampleNewClientsFromConnection demonstrates constructing an Argos client
// from an existing gRPC connection.
func ExampleNewClientsFromConnection() {
	cl := argosv1.NewClientsFromConnection(nil)
	fmt.Println(cl != nil)
	// Output:
	// true
}
