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
	Policies() PoliciesClient
	Bindings() BindingsClient

	Close() error
}

func NewClients(ctx context.Context, addr string, token string) (Clients, error) {
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse policy gates service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*uri)

	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("policygates.NewClients: failed to connect to the server: %w", err)
	}

	return &clients{
		policies: NewPoliciesClient(conn),
		bindings: NewBindingsClient(conn),

		conn: conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		policies: NewPoliciesClient(conn),
		bindings: NewBindingsClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	policies PoliciesClient
	bindings BindingsClient

	conn *grpc.ClientConn
}

func (c *clients) Policies() PoliciesClient {
	return c.policies
}

func (c *clients) Bindings() BindingsClient {
	return c.bindings
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
