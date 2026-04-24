/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	platformtest "chainguard.dev/sdk/proto/platform/test"
)

// ExampleMockPlatformClients demonstrates constructing a mock platform client.
func ExampleMockPlatformClients() {
	mock := platformtest.MockPlatformClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
