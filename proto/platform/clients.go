/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package platform

import (
	"context"
	"fmt"
	"net/url"
	"time"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	advisory "chainguard.dev/sdk/proto/platform/advisory/v1"
	platformauth "chainguard.dev/sdk/proto/platform/auth/v1"
	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	platformoidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	ping "chainguard.dev/sdk/proto/platform/ping/v1"
	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
	"github.com/chainguard-dev/clog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type userAgentString struct{}

type Clients interface {
	IAM() iam.Clients
	Tenant() tenant.Clients
	Registry() registry.Clients
	Advisory() advisory.Clients
	Ping() ping.Clients

	Close() error
}

func NewPlatformClients(ctx context.Context, apiURL string, cred credentials.PerRPCCredentials, addlOpts ...grpc.DialOption) (Clients, error) {
	apiURI, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse iam service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*apiURI)

	// TODO: we may want to require transport security at some future point.
	if cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}
	if ua := GetUserAgent(ctx); ua != "" {
		opts = append(opts, grpc.WithUserAgent(ua))
	}
	opts = append(opts, addlOpts...)

	var cancel context.CancelFunc
	if _, timeoutSet := ctx.Deadline(); !timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
	}
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the iam server: %w", err)
	}

	return &clients{
		iam:      iam.NewClientsFromConnection(conn),
		tenant:   tenant.NewClientsFromConnection(conn),
		registry: registry.NewClientsFromConnection(conn),
		advisory: advisory.NewClientsFromConnection(conn),
		ping:     ping.NewClientsFromConnection(conn),
		conn:     conn,
	}, nil
}

type clients struct {
	iam      iam.Clients
	tenant   tenant.Clients
	registry registry.Clients
	advisory advisory.Clients
	ping     ping.Clients

	conn *grpc.ClientConn
}

func (c *clients) IAM() iam.Clients {
	return c.iam
}

func (c *clients) Tenant() tenant.Clients {
	return c.tenant
}

func (c *clients) Registry() registry.Clients {
	return c.registry
}

func (c *clients) Advisory() advisory.Clients {
	return c.advisory
}

func (c *clients) Ping() ping.Clients {
	return c.ping
}

func (c *clients) Close() error {
	return c.conn.Close()
}

type OIDCClients interface {
	Auth() platformauth.AuthClient
	OIDC() platformoidc.Clients
	OIDCPing() ping.Clients

	Close() error
}

func NewOIDCClients(ctx context.Context, issuerURL string, cred credentials.PerRPCCredentials) (OIDCClients, error) {
	issuerURI, err := url.Parse(issuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse issuer URL: %w", err)
	}

	target, opts := delegate.GRPCOptions(*issuerURI)

	// TODO: we may want to require transport security at some future point.
	if cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}
	if ua := GetUserAgent(ctx); ua != "" {
		opts = append(opts, grpc.WithUserAgent(ua))
	}

	var cancel context.CancelFunc
	if _, timeoutSet := ctx.Deadline(); !timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
	}
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the OIDC issuer: %w", err)
	}

	return &oidcClients{
		auth: platformauth.NewAuthClient(conn),
		oidc: platformoidc.NewClientsFromConnection(conn),
		ping: ping.NewClientsFromConnection(conn),
		conn: conn,
	}, nil
}

type oidcClients struct {
	auth platformauth.AuthClient
	oidc platformoidc.Clients
	ping ping.Clients

	conn *grpc.ClientConn
}

func (c *oidcClients) Auth() platformauth.AuthClient {
	return c.auth
}

func (c *oidcClients) OIDC() platformoidc.Clients {
	return c.oidc
}

func (c *oidcClients) OIDCPing() ping.Clients {
	return c.ping
}

func (c *oidcClients) Close() error {
	return c.conn.Close()
}

// WithUserAgent adds a UserAgent string to the context
// passed to the GRPC client
func WithUserAgent(ctx context.Context, userAgent string) context.Context {
	if userAgent == "" {
		return ctx
	}
	return context.WithValue(ctx, userAgentString{}, userAgent)
}

// GetUserAgent extracts the user agent string from the context
func GetUserAgent(ctx context.Context) string {
	if ua := ctx.Value(userAgentString{}); ua != nil {
		return ua.(string)
	}
	return ""
}
