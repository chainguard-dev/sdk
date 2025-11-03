/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"fmt"
	"net/url"

	"google.golang.org/grpc"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
)

type Clients interface {
	STS() SecurityTokenServiceClient

	Close() error
}

type options struct {
	userAgent string
}

type ClientOption func(*options)

func WithUserAgent(agent string) ClientOption {
	return func(o *options) {
		o.userAgent = agent
	}
}
func NewClients(addr string, opts ...ClientOption) (Clients, error) {
	conf := new(options)
	for _, opt := range opts {
		opt(conf)
	}

	uri, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse oidc service address, must be a url: %w", err)
	}

	target, rpcOpts := delegate.GRPCOptions(*uri)

	if conf.userAgent != "" {
		rpcOpts = append(rpcOpts, grpc.WithUserAgent(conf.userAgent))
	}

	conn, err := grpc.NewClient(target, rpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("oidc.NewClients: failed to connect to the oidc server: %w", err)
	}

	return &clients{
		sts: NewSecurityTokenServiceClient(conn),

		conn: conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		sts: NewSecurityTokenServiceClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	sts SecurityTokenServiceClient

	conn *grpc.ClientConn
}

func (c *clients) STS() SecurityTokenServiceClient {
	return c.sts
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
