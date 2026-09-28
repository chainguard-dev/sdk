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

// Clients provides access to v2beta1 versions service clients.
type Clients interface {
	VersionsService() VersionsServiceClient

	ListProjectsIter(ctx context.Context, req *ListProjectsRequest) iter.Seq2[*Project, error]
	ListProjectsAll(ctx context.Context, req *ListProjectsRequest) ([]*Project, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 versions clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		versionsService: NewVersionsServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	versionsService VersionsServiceClient

	conn *grpc.ClientConn
}

func (c *clients) VersionsService() VersionsServiceClient {
	return c.versionsService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListProjectsIter returns an iterator over the projects matching the request.
func (c *clients) ListProjectsIter(ctx context.Context, req *ListProjectsRequest) iter.Seq2[*Project, error] {
	return v2iter.Paginate(ctx, req, "projects", func(ctx context.Context, r *ListProjectsRequest) ([]*Project, string, error) {
		resp, err := c.VersionsService().ListProjects(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetProjects(), resp.GetNextPageToken(), nil
	})
}

// ListProjectsAll fetches all projects matching the request by automatically handling pagination.
// For large result sets, consider using ListProjectsIter directly to process items incrementally.
func (c *clients) ListProjectsAll(ctx context.Context, req *ListProjectsRequest) ([]*Project, error) {
	return v2iter.All(c.ListProjectsIter(ctx, req))
}
