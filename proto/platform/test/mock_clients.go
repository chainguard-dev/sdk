/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"google.golang.org/grpc"

	"chainguard.dev/sdk/proto/platform"
	advisory "chainguard.dev/sdk/proto/platform/advisory/v1"
	advisorytest "chainguard.dev/sdk/proto/platform/advisory/v1/test"
	apk "chainguard.dev/sdk/proto/platform/apk/v1"
	apktest "chainguard.dev/sdk/proto/platform/apk/v1/test"
	auth "chainguard.dev/sdk/proto/platform/auth/v1"
	authtest "chainguard.dev/sdk/proto/platform/auth/v1/test"
	ecosystems "chainguard.dev/sdk/proto/platform/ecosystems/v1"
	ecosystemstest "chainguard.dev/sdk/proto/platform/ecosystems/v1/test"
	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	iamtest "chainguard.dev/sdk/proto/platform/iam/v1/test"
	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
	librariestest "chainguard.dev/sdk/proto/platform/libraries/v1/test"
	notifications "chainguard.dev/sdk/proto/platform/notifications/v1"
	notificationstest "chainguard.dev/sdk/proto/platform/notifications/v1/test"
	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	oidctest "chainguard.dev/sdk/proto/platform/oidc/v1/test"
	ping "chainguard.dev/sdk/proto/platform/ping/v1"
	pingtest "chainguard.dev/sdk/proto/platform/ping/v1/test"
	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	registrytest "chainguard.dev/sdk/proto/platform/registry/v1/test"
	vulnerabilities "chainguard.dev/sdk/proto/platform/vulnerabilities/v1"
	vulnerabilitiestest "chainguard.dev/sdk/proto/platform/vulnerabilities/v1/test"
)

var _ platform.Clients = (*MockPlatformClients)(nil)

type MockPlatformClients struct {
	OnError error
	Conn    *grpc.ClientConn

	IAMClient             iamtest.MockIAMClient
	RegistryClient        registrytest.MockRegistryClients
	AdvisoryClient        advisorytest.MockSecurityAdvisoryClients
	PingClient            pingtest.MockPingServiceClients
	NotificationsClient   notificationstest.MockNotificationsClients
	APKClient             apktest.MockAPKClients
	EcosystemsClient      ecosystemstest.MockEcosystemsClients
	LibrariesClient       librariestest.MockLibrariesClients
	VulnerabilitiesClient vulnerabilitiestest.MockVulnerabilitiesClients
}

func (m MockPlatformClients) Close() error {
	return m.OnError
}

func (m MockPlatformClients) IAM() iam.Clients {
	return m.IAMClient
}

func (m MockPlatformClients) Connection() *grpc.ClientConn {
	return m.Conn
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

func (m MockPlatformClients) Notifications() notifications.Clients {
	return m.NotificationsClient
}

func (m MockPlatformClients) APK() apk.Clients {
	return m.APKClient
}

func (m MockPlatformClients) Ecosystems() ecosystems.Clients {
	return m.EcosystemsClient
}

func (m MockPlatformClients) Libraries() libraries.Clients {
	return m.LibrariesClient
}

func (m MockPlatformClients) Vulnerabilities() vulnerabilities.Clients {
	return m.VulnerabilitiesClient
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
