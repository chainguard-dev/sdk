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

// Clients provides access to v2beta1 advisory service clients.
type Clients interface {
	SecurityAdvisoryService() SecurityAdvisoryServiceClient

	ListVulnerabilityMetadataIter(ctx context.Context, req *ListVulnerabilityMetadataRequest) iter.Seq2[*VulnerabilityMetadata, error]
	ListVulnerabilityMetadataAll(ctx context.Context, req *ListVulnerabilityMetadataRequest) ([]*VulnerabilityMetadata, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 advisory clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		securityAdvisoryService: NewSecurityAdvisoryServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	securityAdvisoryService SecurityAdvisoryServiceClient

	conn *grpc.ClientConn
}

func (c *clients) SecurityAdvisoryService() SecurityAdvisoryServiceClient {
	return c.securityAdvisoryService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListVulnerabilityMetadataIter returns an iterator over vulnerability metadata matching the request.
func (c *clients) ListVulnerabilityMetadataIter(ctx context.Context, req *ListVulnerabilityMetadataRequest) iter.Seq2[*VulnerabilityMetadata, error] {
	return v2iter.Paginate(ctx, req, "vulnerability_metadata", func(ctx context.Context, r *ListVulnerabilityMetadataRequest) ([]*VulnerabilityMetadata, string, error) {
		resp, err := c.SecurityAdvisoryService().ListVulnerabilityMetadata(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetVulnerabilityMetadata(), resp.GetNextPageToken(), nil
	})
}

// ListVulnerabilityMetadataAll fetches all vulnerability metadata matching the request by automatically handling pagination.
// For large result sets, consider using ListVulnerabilityMetadataIter directly to process items incrementally.
func (c *clients) ListVulnerabilityMetadataAll(ctx context.Context, req *ListVulnerabilityMetadataRequest) ([]*VulnerabilityMetadata, error) {
	return v2iter.All(c.ListVulnerabilityMetadataIter(ctx, req))
}
