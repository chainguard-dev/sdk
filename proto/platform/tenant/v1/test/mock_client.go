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

	SbomsClient       MockSbomsClient
	SignaturesClient  MockSignaturesClient
	VulnReportsClient MockVulnReportsClient
	AttestationClient MockAttestationsClientt
}

var _ tenant.Clients = (*MockTenantClient)(nil)

func (m MockTenantClient) Sboms() tenant.SbomsClient {
	return &m.SbomsClient
}

func (m MockTenantClient) Signatures() tenant.SignaturesClient {
	return &m.SignaturesClient
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
