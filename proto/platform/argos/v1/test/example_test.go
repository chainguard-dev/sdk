/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	argostest "chainguard.dev/sdk/proto/platform/argos/v1/test"
)

// ExampleMockArgosClients demonstrates constructing a mock Argos client.
func ExampleMockArgosClients() {
	mock := argostest.MockArgosClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
