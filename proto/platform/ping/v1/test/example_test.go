/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	pingtest "chainguard.dev/sdk/proto/platform/ping/v1/test"
)

// ExampleMockPingServiceClients demonstrates constructing a mock ping client.
func ExampleMockPingServiceClients() {
	mock := pingtest.MockPingServiceClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
