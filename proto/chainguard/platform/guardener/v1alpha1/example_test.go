/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	guardener "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1"
)

func ExampleNewClientsFromConnection() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := guardener.NewClientsFromConnection(conn)

	fmt.Println(clients.Guardener() != nil, clients.GitHubAssociations() != nil)
	// Output: true true
}

func ExampleNewGuardenerClient() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := guardener.NewGuardenerClient(conn)

	fmt.Println(client != nil)
	// Output: true
}
