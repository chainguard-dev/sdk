/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"context"
	"fmt"
	"net/url"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	iamv2 "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2beta1"
	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2beta1"
	"github.com/chainguard-dev/clog"
)

// Clients provides access to v2beta1 API clients.
type Clients interface {
	IAM() iamv2.Clients
	Registry() registry.Clients
	Vulnerabilities() vuln.Clients
	Close() error
}

type clients struct {
	iam      iamv2.Clients
	registry registry.Clients
	vuln     vuln.Clients
	conn     *grpc.ClientConn
}

// NewClients creates a v2beta1 API gRPC client. The caller is responsible for closing the connection.
func NewClients(ctx context.Context, apiURL, userAgent string, cred credentials.PerRPCCredentials, addlOpts ...grpc.DialOption) (Clients, error) {
	if userAgent == "" {
		return nil, fmt.Errorf("userAgent cannot be empty")
	}
	uri, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse api service address, must be a url: %w", err)
	}

	// Parse the target API URL and get default call options, including min version TLS1.2.
	target, opts := delegate.GRPCOptions(*uri)

	if cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.WarnContext(ctx, "No authentication provided, this may end badly.")
	}

	opts = append(append(opts, addlOpts...), grpc.WithUserAgent(userAgent))

	// Create a new client connection. No I/O is performed until an RPC is made with the connection.
	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the api server %s: %w", target, err)
	}

	return &clients{
		iam:      iamv2.NewClientsFromConnection(conn),
		registry: registry.NewClientsFromConnection(conn),
		vuln:     vuln.NewClientsFromConnection(conn),
		conn:     conn,
	}, nil
}

func (c *clients) IAM() iamv2.Clients {
	return c.iam
}

func (c *clients) Registry() registry.Clients {
	return c.registry
}

func (c *clients) Vulnerabilities() vuln.Clients {
	return c.vuln
}

func (c *clients) Close() error {
	return c.conn.Close()
}
