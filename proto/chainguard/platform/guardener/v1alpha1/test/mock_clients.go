/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import guardener "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1"

var _ guardener.Clients = (*MockGuardenerClients)(nil)

type MockGuardenerClients struct {
	GuardenerClient MockGuardenerClient

	OnClose error
}

func (m MockGuardenerClients) Guardener() guardener.GuardenerClient {
	return &m.GuardenerClient
}

// GitHubAssociations is not mocked; calling methods on the returned nil client
// panics, matching the leave-unconfigured convention of the other mocks.
func (m MockGuardenerClients) GitHubAssociations() guardener.GitHubAssociationsClient {
	return nil
}

func (m MockGuardenerClients) Close() error {
	return m.OnClose
}
