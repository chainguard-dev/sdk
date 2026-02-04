/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"fmt"
	"net/url"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"github.com/chainguard-dev/clog"
	"google.golang.org/grpc"

	"chainguard.dev/sdk/auth"
)

type Clients interface {
	ImageMatcher() ImageMatcherClient

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
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("advisory.NewClients: failed to connect to the iam server: %w", err)
	}

	return NewClientsFromConnection(conn), nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		matcher: NewImageMatcherClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	matcher ImageMatcherClient

	conn *grpc.ClientConn
}

func (c *clients) ImageMatcher() ImageMatcherClient {
	return c.matcher
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
