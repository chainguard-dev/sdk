/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import libraries "chainguard.dev/sdk/proto/platform/libraries/v1"

var _ libraries.Clients = (*MockLibrariesClients)(nil)

type MockLibrariesClients struct {
	ArtifactsClient    MockArtifactsClient
	EntitlementsClient MockEntitlementsClient

	OnClose error
}

func (m MockLibrariesClients) Artifacts() libraries.ArtifactsClient {
	return &m.ArtifactsClient
}

func (m MockLibrariesClients) Entitlements() libraries.EntitlementsClient {
	return &m.EntitlementsClient
}

func (m MockLibrariesClients) Close() error {
	return m.OnClose
}
