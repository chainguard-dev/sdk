/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	skillsv1 "chainguard.dev/sdk/proto/platform/skills/v1"
)

// ExampleNewClientsFromConnection demonstrates constructing a Skills client
// from an existing gRPC connection.
func ExampleNewClientsFromConnection() {
	cl := skillsv1.NewClientsFromConnection(nil)
	fmt.Println(cl != nil)
	// Output:
	// true
}
