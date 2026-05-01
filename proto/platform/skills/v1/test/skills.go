/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import skills "chainguard.dev/sdk/proto/platform/skills/v1"

var _ skills.Clients = (*MockSkillsClients)(nil)

type MockSkillsClients struct {
	EntitlementsClient MockEntitlementsClient

	OnClose error
}

func (m MockSkillsClients) Entitlements() skills.EntitlementsClient {
	return &m.EntitlementsClient
}

func (m MockSkillsClients) Close() error {
	return m.OnClose
}
