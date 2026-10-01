/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	sandbox "chainguard.dev/sdk/proto/chainguard/platform/sandbox/v1alpha1"
)

func ExampleNewInstallationsClient() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := sandbox.NewInstallationsClient(conn)

	fmt.Println(client != nil)
	// Output: true
}
