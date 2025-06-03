/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"fmt"
	"net/url"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"google.golang.org/grpc"

	"chainguard.dev/sdk/auth"
	"github.com/chainguard-dev/clog"
)

type Clients interface {
	Entitlements() EntitlementsClient

	Close() error
}

func NewClients(ctx context.Context, ecoURL string, token string) (Clients, error) {
	ecoURI, err := url.Parse(ecoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ecosystems service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*ecoURI)

	// TODO: we may want to require transport security at some future point.
	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("NewClients: failed to connect to the ecosystems server: %w", err)
	}

	return &clients{
		entitlements: NewEntitlementsClient(conn),

		conn: conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		entitlements: NewEntitlementsClient(conn),

		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	entitlements EntitlementsClient

	conn *grpc.ClientConn
}

func (c *clients) Entitlements() EntitlementsClient {
	return c.entitlements
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
