/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	vulntest "chainguard.dev/sdk/proto/platform/vulnerabilities/v1/test"
)

// ExampleMockVulnerabilitiesClients demonstrates constructing a mock vulnerabilities client.
func ExampleMockVulnerabilitiesClients() {
	mock := vulntest.MockVulnerabilitiesClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
