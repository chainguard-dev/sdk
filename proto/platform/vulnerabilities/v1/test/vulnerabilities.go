/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	vulnerabilities "chainguard.dev/sdk/proto/platform/vulnerabilities/v1"
)

var _ vulnerabilities.Clients = (*MockVulnerabilitiesClients)(nil)

type MockVulnerabilitiesClients struct {
	AdvisoriesClients MockAdvisoriesClient

	OnClose error
}

func (m MockVulnerabilitiesClients) Advisories() vulnerabilities.AdvisoriesClient {
	return &m.AdvisoriesClients
}

func (m MockVulnerabilitiesClients) Close() error {
	return m.OnClose
}
