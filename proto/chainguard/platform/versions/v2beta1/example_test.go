/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"context"
	"fmt"

	versions "chainguard.dev/sdk/proto/chainguard/platform/versions/v2beta1"
)

func ExampleNewClientsFromConnection() {
	clients := versions.NewClientsFromConnection(nil)
	fmt.Println(clients != nil)
	// Output: true
}

func ExampleClients_VersionsService() {
	clients := versions.NewClientsFromConnection(nil)
	fmt.Println(clients.VersionsService() != nil)
	// Output: true
}

func ExampleClients_ListProjectsIter() {
	clients := versions.NewClientsFromConnection(nil)
	_ = clients.ListProjectsIter(context.Background(), &versions.ListProjectsRequest{})
	fmt.Println("iterator created")
	// Output: iterator created
}

func ExampleClients_ListProjectsAll() {
	clients := versions.NewClientsFromConnection(nil)
	_ = clients.ListProjectsAll
	fmt.Println("func available")
	// Output: func available
}

func ExampleClients_Close() {
	clients := versions.NewClientsFromConnection(nil)
	_ = clients.Close()
	fmt.Println("closed")
	// Output: closed
}
