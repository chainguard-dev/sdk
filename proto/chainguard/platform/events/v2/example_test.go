/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2_test

import (
	"fmt"

	events "chainguard.dev/sdk/proto/chainguard/platform/events/v2"
)

func ExampleNewClientsFromConnection() {
	clients := events.NewClientsFromConnection(nil)
	fmt.Println(clients != nil)
	// Output: true
}

func ExampleClients_SubscriptionsService() {
	clients := events.NewClientsFromConnection(nil)
	fmt.Println(clients.SubscriptionsService() != nil)
	// Output: true
}

func ExampleClients_Close() {
	clients := events.NewClientsFromConnection(nil)
	_ = clients.Close()
	fmt.Println("closed")
	// Output: closed
}
