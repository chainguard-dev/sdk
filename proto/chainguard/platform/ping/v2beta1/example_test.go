/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	ping "chainguard.dev/sdk/proto/chainguard/platform/ping/v2beta1"
)

func ExampleNewClientsFromConnection() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := ping.NewClientsFromConnection(conn)
	defer clients.Close()

	fmt.Println(clients != nil)
	// Output: true
}

func ExampleClients_PingService() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := ping.NewClientsFromConnection(conn)
	defer clients.Close()

	fmt.Println(clients.PingService() != nil)
	// Output: true
}

func ExampleClients_Close() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := ping.NewClientsFromConnection(conn)
	_ = clients.Close()
	fmt.Println("closed")
	// Output: closed
}
