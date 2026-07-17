/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"context"
	"fmt"

	advisory "chainguard.dev/sdk/proto/chainguard/platform/advisory/v2beta1"
)

func ExampleNewClientsFromConnection() {
	clients := advisory.NewClientsFromConnection(nil)
	fmt.Println(clients != nil)
	// Output: true
}

func ExampleClients_SecurityAdvisoryService() {
	clients := advisory.NewClientsFromConnection(nil)
	fmt.Println(clients.SecurityAdvisoryService() != nil)
	// Output: true
}

func ExampleClients_ListVulnerabilityMetadataIter() {
	clients := advisory.NewClientsFromConnection(nil)
	_ = clients.ListVulnerabilityMetadataIter(context.Background(), &advisory.ListVulnerabilityMetadataRequest{})
	fmt.Println("iterator created")
	// Output: iterator created
}

func ExampleClients_ListVulnerabilityMetadataAll() {
	clients := advisory.NewClientsFromConnection(nil)
	_ = clients.ListVulnerabilityMetadataAll
	fmt.Println("func available")
	// Output: func available
}

func ExampleClients_Close() {
	clients := advisory.NewClientsFromConnection(nil)
	_ = clients.Close()
	fmt.Println("closed")
	// Output: closed
}
