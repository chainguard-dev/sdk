/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import ecosystems "chainguard.dev/sdk/proto/platform/ecosystems/v1"

var _ ecosystems.Clients = (*MockEcosystemsClients)(nil)

type MockEcosystemsClients struct {
	EntitlementsClient MockEntitlementsClient

	OnClose error
}

func (m MockEcosystemsClients) Entitlements() ecosystems.EntitlementsClient {
	return &m.EntitlementsClient
}

func (m MockEcosystemsClients) Close() error {
	return m.OnClose
}
