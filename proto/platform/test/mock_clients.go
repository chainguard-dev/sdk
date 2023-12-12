/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"chainguard.dev/sdk/proto/platform"
	advisory "chainguard.dev/sdk/proto/platform/advisory/v1"
	advisorytest "chainguard.dev/sdk/proto/platform/advisory/v1/test"
	auth "chainguard.dev/sdk/proto/platform/auth/v1"
	authtest "chainguard.dev/sdk/proto/platform/auth/v1/test"
	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	iamtest "chainguard.dev/sdk/proto/platform/iam/v1/test"
	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	oidctest "chainguard.dev/sdk/proto/platform/oidc/v1/test"
	ping "chainguard.dev/sdk/proto/platform/ping/v1"
	pingtest "chainguard.dev/sdk/proto/platform/ping/v1/test"
	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	registrytest "chainguard.dev/sdk/proto/platform/registry/v1/test"
	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
	tenanttest "chainguard.dev/sdk/proto/platform/tenant/v1/test"
)

var _ platform.Clients = (*MockPlatformClients)(nil)

type MockPlatformClients struct {
	OnError error

	IAMClient      iamtest.MockIAMClient
	TenantClient   tenanttest.MockTenantClient
	RegistryClient registrytest.MockRegistryClients
	AdvisoryClient advisorytest.MockSecurityAdvisoryClients
	PingClient     pingtest.MockPingServiceClients
}

func (m MockPlatformClients) Close() error {
	return m.OnError
}

func (m MockPlatformClients) IAM() iam.Clients {
	return m.IAMClient
}

func (m MockPlatformClients) Tenant() tenant.Clients {
	return m.TenantClient
}

func (m MockPlatformClients) Registry() registry.Clients {
	return m.RegistryClient
}

func (m MockPlatformClients) Advisory() advisory.Clients {
	return m.AdvisoryClient
}

func (m MockPlatformClients) Ping() ping.Clients {
	return m.PingClient
}

var _ platform.OIDCClients = (*MockOIDCClients)(nil)

type MockOIDCClients struct {
	OnError error

	AuthClient authtest.MockAuthClient
	OIDCClient oidctest.MockOIDCClient
	PingClient pingtest.MockPingServiceClients
}

func (m MockOIDCClients) Close() error {
	return m.OnError
}

func (m MockOIDCClients) Auth() auth.AuthClient {
	return m.AuthClient
}

func (m MockOIDCClients) OIDC() oidc.Clients {
	return m.OIDCClient
}

func (m MockOIDCClients) OIDCPing() ping.Clients {
	return m.PingClient
}
