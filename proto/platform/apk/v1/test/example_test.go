/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	apktest "chainguard.dev/sdk/proto/platform/apk/v1/test"
)

// ExampleMockAPKClients demonstrates constructing a mock APK client.
func ExampleMockAPKClients() {
	mock := apktest.MockAPKClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
