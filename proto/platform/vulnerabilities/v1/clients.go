/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"fmt"
	"net/url"
	"time"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"github.com/chainguard-dev/clog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Clients holds the gRPC clients for the advisory v2 services.
type Clients interface {
	Advisories() AdvisoriesClient

	Close() error
}

type clients struct {
	conn       *grpc.ClientConn
	advisories AdvisoriesClient
}

// NewClientsFromConnection creates a new set of clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		conn:       conn,
		advisories: NewAdvisoriesClient(conn),
	}
}

func (c *clients) Advisories() AdvisoriesClient {
	return c.advisories
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// NewClients creates a new set of clients for the advisory v2 API
func NewClients(ctx context.Context, apiURL string, cred credentials.PerRPCCredentials, addlOpts ...grpc.DialOption) (Clients, error) {
	apiURI, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse api service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*apiURI)

	var cancel context.CancelFunc
	if _, timeoutSet := ctx.Deadline(); !timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
	}

	// TODO: we may want to require transport security at some future point.
	if cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}
	opts = append(opts, addlOpts...)

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("NewAdvisoryClients: failed to connect to the api server %s: %w", target, err)
	}

	return &clients{
		advisories: NewAdvisoriesClient(conn),

		conn: conn,
	}, nil
}
