/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	matchertest "chainguard.dev/sdk/proto/platform/matcher/v1/test"
)

// ExampleMockImageMatcherClients demonstrates constructing a mock image matcher client.
func ExampleMockImageMatcherClients() {
	mock := matchertest.MockImageMatcherClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
