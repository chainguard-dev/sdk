/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"google.golang.org/grpc"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"chainguard.dev/sdk/auth"
	"github.com/chainguard-dev/clog"
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
func NewClients(ctx context.Context, addr string, token string, opts ...ClientOption) (Clients, error) {
	conf := new(options)
	for _, opt := range opts {
		opt(conf)
	}

	uri, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse oidc service address, must be a url: %w", err)
	}

	target, rpcOpts := delegate.GRPCOptions(*uri)

	// TODO: we may want to require transport security at some future point.
	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		rpcOpts = append(rpcOpts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}

	if conf.userAgent != "" {
		rpcOpts = append(rpcOpts, grpc.WithUserAgent(conf.userAgent))
	}

	var cancel context.CancelFunc
	if _, timeoutSet := ctx.Deadline(); !timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
	}
	// grpc.NewClient introduced a regression with respect to proxying requests.
	// Specifically, the target URI gets resolved to the IP and passed to the connection,
	// which causes issues for customers using proxies.
	// This issue is being tracked here https://github.com/grpc/grpc-go/issues/7556 and a fix
	// is expected by grpc-go 1.70
	//nolint:staticcheck // Revert back to grpc.NewClient once #7556 is resolved.
	conn, err := grpc.DialContext(ctx, target, rpcOpts...)
	if err != nil {
		return nil, fmt.Errorf("oidc.NewClients: failed to connect to the iam server: %w", err)
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
