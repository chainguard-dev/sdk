/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	notificationstest "chainguard.dev/sdk/proto/platform/notifications/v1/test"
)

// ExampleMockNotificationsClients demonstrates constructing a mock notifications client.
func ExampleMockNotificationsClients() {
	mock := notificationstest.MockNotificationsClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
