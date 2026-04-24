/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	advisorytest "chainguard.dev/sdk/proto/platform/advisory/v1/test"
)

// ExampleMockSecurityAdvisoryClients demonstrates constructing a mock advisory client.
func ExampleMockSecurityAdvisoryClients() {
	mock := advisorytest.MockSecurityAdvisoryClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
