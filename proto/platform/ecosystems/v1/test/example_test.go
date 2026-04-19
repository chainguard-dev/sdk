/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	ecosystemstest "chainguard.dev/sdk/proto/platform/ecosystems/v1/test"
)

// ExampleMockEcosystemsClients demonstrates constructing a mock ecosystems client.
func ExampleMockEcosystemsClients() {
	mock := ecosystemstest.MockEcosystemsClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
