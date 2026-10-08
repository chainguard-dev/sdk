/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	billing "chainguard.dev/sdk/proto/chainguard/platform/billing/v2beta1"
)

func ExampleNewClientsFromConnection() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := billing.NewClientsFromConnection(conn)
	defer clients.Close()

	fmt.Println(clients.BillingService() != nil)
	// Output: true
}
