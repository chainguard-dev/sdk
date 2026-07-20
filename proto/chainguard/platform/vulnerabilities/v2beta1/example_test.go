/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2beta1"
)

func ExampleNewClientsFromConnection() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := vuln.NewClientsFromConnection(conn)
	defer c.Close()

	fmt.Println(c != nil)
	// Output: true
}

func ExampleClients_AdvisoriesService() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := vuln.NewClientsFromConnection(conn)
	defer c.Close()

	fmt.Println(c.AdvisoriesService() != nil)
	// Output: true
}

func ExampleClients_Close() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := vuln.NewClientsFromConnection(conn)
	_ = c.Close()
	fmt.Println("closed")
	// Output: closed
}
