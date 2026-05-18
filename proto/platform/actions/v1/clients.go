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
	ActionsAuthorization() ActionsAuthorizationClient
	Actions() ActionsClient

	Close() error
}

func NewClients(ctx context.Context, addr string, token string) (Clients, error) {
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse actions service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*uri)

	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.WarnContext(ctx, "No authentication provided, this may end badly.")
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("actions.NewClients: failed to connect to the actions server: %w", err)
	}

	return &clients{
		authz:   NewActionsAuthorizationClient(conn),
		actions: NewActionsClient(conn),
		conn:    conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		authz:   NewActionsAuthorizationClient(conn),
		actions: NewActionsClient(conn),
	}
}

type clients struct {
	authz   ActionsAuthorizationClient
	actions ActionsClient
	conn    *grpc.ClientConn
}

func (c *clients) ActionsAuthorization() ActionsAuthorizationClient {
	return c.authz
}

func (c *clients) Actions() ActionsClient {
	return c.actions
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
