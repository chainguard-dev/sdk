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

// Clients provides access to v2beta1 vulnerabilities service clients.
type Clients interface {
	AdvisoriesService() AdvisoriesServiceClient

	ListAdvisoriesIter(ctx context.Context, req *ListAdvisoriesRequest) iter.Seq2[*Advisory, error]
	ListAdvisoriesAll(ctx context.Context, req *ListAdvisoriesRequest) ([]*Advisory, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 vulnerabilities clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		advisoriesService: NewAdvisoriesServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	advisoriesService AdvisoriesServiceClient

	conn *grpc.ClientConn
}

func (c *clients) AdvisoriesService() AdvisoriesServiceClient {
	return c.advisoriesService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListAdvisoriesIter returns an iterator over advisories matching the request.
func (c *clients) ListAdvisoriesIter(ctx context.Context, req *ListAdvisoriesRequest) iter.Seq2[*Advisory, error] {
	return v2iter.Paginate(ctx, req, "advisories", func(ctx context.Context, r *ListAdvisoriesRequest) ([]*Advisory, string, error) {
		resp, err := c.AdvisoriesService().ListAdvisories(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetAdvisories(), resp.GetNextPageToken(), nil
	})
}

// ListAdvisoriesAll fetches all advisories matching the request by automatically handling pagination.
// For large result sets, consider using ListAdvisoriesIter directly to process items incrementally.
func (c *clients) ListAdvisoriesAll(ctx context.Context, req *ListAdvisoriesRequest) ([]*Advisory, error) {
	return v2iter.All(c.ListAdvisoriesIter(ctx, req))
}
