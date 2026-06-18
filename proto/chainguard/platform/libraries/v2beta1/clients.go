/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"context"
	"iter"

	"google.golang.org/grpc"

	v2iter "chainguard.dev/sdk/proto/chainguard/platform/iter"
)

// Clients provides access to v2beta1 Libraries service clients.
type Clients interface {
	ArtifactsService() ArtifactsServiceClient

	ListArtifactsIter(ctx context.Context, req *ListArtifactsRequest) iter.Seq2[*Artifact, error]
	ListArtifactsAll(ctx context.Context, req *ListArtifactsRequest) ([]*Artifact, error)

	ListArtifactVersionsIter(ctx context.Context, req *ListArtifactVersionsRequest) iter.Seq2[*ArtifactVersion, error]
	ListArtifactVersionsAll(ctx context.Context, req *ListArtifactVersionsRequest) ([]*ArtifactVersion, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 Libraries clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		artifactsService: NewArtifactsServiceClient(conn),
	}
}

type clients struct {
	artifactsService ArtifactsServiceClient

	conn *grpc.ClientConn
}

func (c *clients) ArtifactsService() ArtifactsServiceClient {
	return c.artifactsService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *clients) ListArtifactsIter(ctx context.Context, req *ListArtifactsRequest) iter.Seq2[*Artifact, error] {
	return v2iter.Paginate(ctx, req, "artifacts", func(ctx context.Context, r *ListArtifactsRequest) ([]*Artifact, string, error) {
		resp, err := c.ArtifactsService().ListArtifacts(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetArtifacts(), resp.GetNextPageToken(), nil
	})
}

func (c *clients) ListArtifactsAll(ctx context.Context, req *ListArtifactsRequest) ([]*Artifact, error) {
	return v2iter.All(c.ListArtifactsIter(ctx, req))
}

func (c *clients) ListArtifactVersionsIter(ctx context.Context, req *ListArtifactVersionsRequest) iter.Seq2[*ArtifactVersion, error] {
	return v2iter.Paginate(ctx, req, "artifact_versions", func(ctx context.Context, r *ListArtifactVersionsRequest) ([]*ArtifactVersion, string, error) {
		resp, err := c.ArtifactsService().ListArtifactVersions(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetArtifactVersions(), resp.GetNextPageToken(), nil
	})
}

func (c *clients) ListArtifactVersionsAll(ctx context.Context, req *ListArtifactVersionsRequest) ([]*ArtifactVersion, error) {
	return v2iter.All(c.ListArtifactVersionsIter(ctx, req))
}
