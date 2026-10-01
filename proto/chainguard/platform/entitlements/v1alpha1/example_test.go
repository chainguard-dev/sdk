/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	entitlements "chainguard.dev/sdk/proto/chainguard/platform/entitlements/v1alpha1"
)

func ExampleNewFeaturesClient() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := entitlements.NewFeaturesClient(conn)

	fmt.Println(client != nil)
	// Output: true
}
