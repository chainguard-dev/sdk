/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	telemetry "chainguard.dev/sdk/proto/chainguard/platform/telemetry/v1alpha1"
)

func ExampleNewClientsFromConnection() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := telemetry.NewClientsFromConnection(conn)

	fmt.Println(clients.Telemetry() != nil)
	// Output: true
}

func ExampleNewTelemetryServiceClient() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := telemetry.NewTelemetryServiceClient(conn)

	fmt.Println(client != nil)
	// Output: true
}
