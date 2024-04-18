/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
)

type MockTenantClient struct {
	OnClose error

	ClustersClient       MockClustersClient
	RecordsClient        MockRecordsClient
	RecordContextsClient MockRecordContextsClient
	SbomsClient          MockSbomsClient
	SignaturesClient     MockSignaturesClient
	PolicyResultsClient  MockPolicyResultsClient
	RisksClient          MockRisksClient
	VulnReportsClient    MockVulnReportsClient
	AttestationClient    MockAttestationsClientt
}

var _ tenant.Clients = (*MockTenantClient)(nil)

func (m MockTenantClient) Clusters() tenant.ClustersClient {
	return &m.ClustersClient
}

func (m MockTenantClient) Records() tenant.RecordsClient {
	return &m.RecordsClient
}

func (m MockTenantClient) RecordContexts() tenant.RecordContextsClient {
	return &m.RecordContextsClient
}

func (m MockTenantClient) Sboms() tenant.SbomsClient {
	return &m.SbomsClient
}

func (m MockTenantClient) Risks() tenant.RisksClient {
	return &m.RisksClient
}

func (m MockTenantClient) Signatures() tenant.SignaturesClient {
	return &m.SignaturesClient
}

func (m MockTenantClient) PolicyResults() tenant.PolicyResultsClient {
	return &m.PolicyResultsClient
}

func (m MockTenantClient) VulnReports() tenant.VulnReportsClient {
	return &m.VulnReportsClient
}

func (m MockTenantClient) Attestations() tenant.AttestationsClient {
	return &m.AttestationClient
}

func (m MockTenantClient) Close() error {
	return m.OnClose
}
