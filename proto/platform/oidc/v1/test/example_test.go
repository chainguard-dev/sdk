/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	oidctest "chainguard.dev/sdk/proto/platform/oidc/v1/test"
)

// ExampleMockOIDCClient demonstrates constructing a mock OIDC client.
func ExampleMockOIDCClient() {
	mock := oidctest.MockOIDCClient{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
