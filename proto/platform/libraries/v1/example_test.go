/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "chainguard.dev/sdk/proto/platform/libraries/v1"
)

// ExampleNewClients demonstrates creating a new Libraries client
// with authentication token.
func ExampleNewClients() {
	ctx := context.Background()

	// Create clients with the ecosystems service URL and auth token
	clients, err := v1.NewClients(ctx, "https://console-api.enforce.dev", "your-auth-token")
	if err != nil {
		log.Fatalf("failed to create clients: %v", err)
	}
	defer clients.Close()

	// Use the clients to interact with the Libraries API
	_ = clients.Artifacts()
	_ = clients.Entitlements()
	_ = clients.NpmPackages()

	fmt.Println("Clients created successfully")
	// Output: Clients created successfully
}

// ExampleNewClientsFromConnection demonstrates creating clients
// from an existing gRPC connection.
func ExampleNewClientsFromConnection() {
	// Create a gRPC connection (example uses insecure for demonstration)
	conn, err := grpc.NewClient(
		"console-api.enforce.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to create connection: %v", err)
	}
	defer conn.Close()

	// Create clients from the existing connection
	clients := v1.NewClientsFromConnection(conn)

	// Use the clients to interact with the Libraries API
	_ = clients.Artifacts()
	_ = clients.Entitlements()
	_ = clients.NpmPackages()

	// Note: When using NewClientsFromConnection, the caller is responsible
	// for closing the connection. Calling clients.Close() will not close
	// the connection.

	fmt.Println("Clients created from connection")
	// Output: Clients created from connection
}

// ExampleClients demonstrates using the Clients interface
// to interact with the Libraries API.
func ExampleClients() {
	ctx := context.Background()

	clients, err := v1.NewClients(ctx, "https://console-api.enforce.dev", "your-auth-token")
	if err != nil {
		log.Fatalf("failed to create clients: %v", err)
	}
	defer clients.Close()

	// Access individual service clients
	artifactsClient := clients.Artifacts()
	entitlementsClient := clients.Entitlements()
	npmPackagesClient := clients.NpmPackages()

	// Use the service clients to make API calls
	_ = artifactsClient
	_ = entitlementsClient
	_ = npmPackagesClient

	fmt.Println("Using Clients interface")
	// Output: Using Clients interface
}
