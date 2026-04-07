/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"chainguard.dev/sdk/auth"
	"github.com/chainguard-dev/clog"
)

type Clients interface {
	SecurityAdvisory() SecurityAdvisoryClient

	Close() error
}

func NewClients(ctx context.Context, addr string, token string) (Clients, error) {
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse advisory service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*uri)

	// TODO: we may want to require transport security at some future point.
	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.WarnContext(ctx, "No authentication provided, this may end badly.")
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("advisory.NewClients: failed to connect to the iam server: %w", err)
	}

	return &clients{
		advisory: NewSecurityAdvisoryClient(conn),

		conn: conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		advisory: NewSecurityAdvisoryClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	advisory SecurityAdvisoryClient

	conn *grpc.ClientConn
}

func (c *clients) SecurityAdvisory() SecurityAdvisoryClient {
	return c.advisory
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListVulnerabilityMetadataAll is a helper to fetch all vulnerability metadata that matches the given filter, handling
// pagination for the caller and returning the final slice of *VulnerabilityMetadata. Note that for large result
// sets this may incur a high memory cost, consider paginating manually in these cases.
func (c *clients) ListVulnerabilityMetadataAll(ctx context.Context, filter *VulnerabilityMetadataFilter) ([]*VulnerabilityMetadata, error) {
	if filter == nil {
		return nil, errors.New("at least one vulnerability ID is required")
	}
	f := proto.Clone(filter).(*VulnerabilityMetadataFilter)
	if f.GetPageSize() == 0 {
		f.PageSize = 50
	}

	all := make([]*VulnerabilityMetadata, 0, f.PageSize)

	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		res, err := c.advisory.ListVulnerabilityMetadata(ctx, f)
		if err != nil {
			return nil, err
		}
		all = append(all, res.Items...)

		if res.NextPageToken == "" {
			return all, nil
		}

		f.PageToken = res.NextPageToken
	}
}
