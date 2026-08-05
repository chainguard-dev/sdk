/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2_test

import (
	"context"
	"fmt"
	"log"

	v2 "chainguard.dev/sdk/proto/chainguard/platform/libraries/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ExampleClients_ListArtifactsIter() {
	ctx := context.Background()

	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	clients := v2.NewClientsFromConnection(conn)

	req := &v2.ListArtifactsRequest{
		Ecosystems: []v2.Ecosystem{v2.Ecosystem_ECOSYSTEM_PYTHON},
		Query:      "requests",
		PageSize:   10,
	}

	for artifact, err := range clients.ListArtifactsIter(ctx, req) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		fmt.Printf("Artifact: %s (v%s)\n", artifact.Name, artifact.LatestVersion)
	}
}

func ExampleClients_ListArtifactVersionsIter() {
	ctx := context.Background()

	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	clients := v2.NewClientsFromConnection(conn)

	req := &v2.ListArtifactVersionsRequest{
		ArtifactId: "pypi/requests",
		PageSize:   50,
		OrderBy:    "version desc",
	}

	for version, err := range clients.ListArtifactVersionsIter(ctx, req) {
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		fmt.Printf("Version: %s\n", version.Version)
	}
}
