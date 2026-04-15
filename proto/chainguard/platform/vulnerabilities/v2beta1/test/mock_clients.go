/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"iter"

	"chainguard.dev/sdk/proto/chainguard/platform/test"
	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2beta1"
)

type MockClients struct {
	OnClose error

	AdvisoriesServiceClient MockAdvisoriesServiceClient
}

// Close implements [v2beta1.Clients].
func (m *MockClients) Close() error {
	return m.OnClose
}

// AdvisoriesService implements [v2beta1.Clients].
func (m *MockClients) AdvisoriesService() vuln.AdvisoriesServiceClient {
	return &m.AdvisoriesServiceClient
}

// ListAdvisoriesAll implements [v2beta1.Clients].
func (m *MockClients) ListAdvisoriesAll(ctx context.Context, req *vuln.ListAdvisoriesRequest) ([]*vuln.Advisory, error) {
	resp, err := m.AdvisoriesServiceClient.ListAdvisories(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetAdvisories(), nil
}

// ListAdvisoriesIter implements [v2beta1.Clients].
func (m *MockClients) ListAdvisoriesIter(ctx context.Context, req *vuln.ListAdvisoriesRequest) iter.Seq2[*vuln.Advisory, error] {
	return test.MockIter(m.ListAdvisoriesAll(ctx, req))
}

var _ vuln.Clients = (*MockClients)(nil)
