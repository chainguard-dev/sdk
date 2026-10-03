/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"

var _ skills.Clients = (*MockSkillsClients)(nil)

type MockSkillsClients struct {
	CatalogClient          MockCatalogClient
	HardeningReportsClient MockHardeningReportsClient

	OnClose error
}

func (m MockSkillsClients) Skills() skills.SkillsClient {
	return &m.CatalogClient
}

func (m MockSkillsClients) HardeningReports() skills.SkillsHardeningReportsClient {
	return &m.HardeningReportsClient
}

func (m MockSkillsClients) Close() error {
	return m.OnClose
}
