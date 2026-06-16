/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import libraries "chainguard.dev/sdk/proto/platform/libraries/v1"

var _ libraries.Clients = (*MockLibrariesClients)(nil)

type MockLibrariesClients struct {
	ArtifactsClient                MockArtifactsClient
	EntitlementsClient             MockEntitlementsClient
	NpmPackagesClient              MockNpmPackagesClient
	MalwareClient                  MockMalwareClient
	LibraryPoliciesClient          MockLibraryPoliciesClient
	LibraryPolicyBindingsClient    MockLibraryPolicyBindingsClient
	LibraryPolicyBlockEventsClient MockLibraryPolicyBlockEventsClient

	AWSMarketplaceSubscriptionsClient MockAWSMarketplaceSubscriptionsClient

	OnClose error
}

func (m MockLibrariesClients) Artifacts() libraries.ArtifactsClient {
	return &m.ArtifactsClient
}

func (m MockLibrariesClients) Entitlements() libraries.EntitlementsClient {
	return &m.EntitlementsClient
}

func (m MockLibrariesClients) NpmPackages() libraries.NpmPackagesClient {
	return &m.NpmPackagesClient
}

func (m MockLibrariesClients) Malware() libraries.MalwareClient {
	return &m.MalwareClient
}

func (m MockLibrariesClients) LibraryPolicies() libraries.LibraryPoliciesClient {
	return &m.LibraryPoliciesClient
}

func (m MockLibrariesClients) LibraryPolicyBindings() libraries.LibraryPolicyBindingsClient {
	return &m.LibraryPolicyBindingsClient
}

func (m MockLibrariesClients) LibraryPolicyBlockEvents() libraries.LibraryPolicyBlockEventsClient {
	return &m.LibraryPolicyBlockEventsClient
}

func (m MockLibrariesClients) AWSMarketplaceSubscriptions() libraries.AWSMarketplaceSubscriptionsClient {
	return &m.AWSMarketplaceSubscriptionsClient
}

func (m MockLibrariesClients) Close() error {
	return m.OnClose
}
