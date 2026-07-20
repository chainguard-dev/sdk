/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	actions "chainguard.dev/sdk/proto/chainguard/platform/actions/v1alpha1"
)

func ExampleNewActionsClient() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := actions.NewActionsClient(conn)

	fmt.Println(client != nil)
	// Output: true
}
