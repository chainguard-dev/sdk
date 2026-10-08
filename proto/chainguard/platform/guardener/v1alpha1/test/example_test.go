/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	guardenertest "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1/test"
)

// ExampleMockGuardenerClients demonstrates constructing a mock guardener
// client.
func ExampleMockGuardenerClients() {
	mock := guardenertest.MockGuardenerClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
