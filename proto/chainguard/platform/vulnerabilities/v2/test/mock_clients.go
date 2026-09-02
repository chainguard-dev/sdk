/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"iter"

	"chainguard.dev/sdk/proto/chainguard/platform/test"
	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2"
)

type MockClients struct {
	OnClose error

	AdvisoriesServiceClient  MockAdvisoriesServiceClient
	VulnReportsServiceClient MockVulnReportsServiceClient
}

// Close implements [v2.Clients].
func (m *MockClients) Close() error {
	return m.OnClose
}

// AdvisoriesService implements [v2.Clients].
func (m *MockClients) AdvisoriesService() vuln.AdvisoriesServiceClient {
	return &m.AdvisoriesServiceClient
}

// ListAdvisoriesAll implements [v2.Clients].
func (m *MockClients) ListAdvisoriesAll(ctx context.Context, req *vuln.ListAdvisoriesRequest) ([]*vuln.Advisory, error) {
	resp, err := m.AdvisoriesServiceClient.ListAdvisories(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetAdvisories(), nil
}

// ListAdvisoriesIter implements [v2.Clients].
func (m *MockClients) ListAdvisoriesIter(ctx context.Context, req *vuln.ListAdvisoriesRequest) iter.Seq2[*vuln.Advisory, error] {
	return test.MockIter(m.ListAdvisoriesAll(ctx, req))
}

// VulnReportsService implements [v2.Clients].
func (m *MockClients) VulnReportsService() vuln.VulnReportsServiceClient {
	return &m.VulnReportsServiceClient
}

// ListVulnCountReportsAll implements [v2.Clients].
func (m *MockClients) ListVulnCountReportsAll(ctx context.Context, req *vuln.ListVulnCountReportsRequest) ([]*vuln.VulnCountReport, error) {
	resp, err := m.VulnReportsServiceClient.ListVulnCountReports(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetVulnCountReports(), nil
}

// ListVulnCountReportsIter implements [v2.Clients].
func (m *MockClients) ListVulnCountReportsIter(ctx context.Context, req *vuln.ListVulnCountReportsRequest) iter.Seq2[*vuln.VulnCountReport, error] {
	return test.MockIter(m.ListVulnCountReportsAll(ctx, req))
}

var _ vuln.Clients = (*MockClients)(nil)
