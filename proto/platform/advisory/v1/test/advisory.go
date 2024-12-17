/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	advisory "chainguard.dev/sdk/proto/platform/advisory/v1"
)

var _ advisory.Clients = (*MockSecurityAdvisoryClients)(nil)

type MockSecurityAdvisoryClients struct {
	OnClose error

	SecurityAdvisoryClient MockSecurityAdvisoryClient
}

func (m MockSecurityAdvisoryClients) SecurityAdvisory() advisory.SecurityAdvisoryClient {
	return &m.SecurityAdvisoryClient
}

func (m MockSecurityAdvisoryClients) Close() error {
	return m.OnClose
}

var _ advisory.SecurityAdvisoryClient = (*MockSecurityAdvisoryClient)(nil)

type MockSecurityAdvisoryClient struct {
	OnListDocuments             []DocumentsOnList
	OnListVulnerabilityMetadata []VulnerabilityMetadataOnList
	OnListResolvedVulnsReports  []ResolvedVulnsReportsOnList
}

type DocumentsOnList struct {
	Given *advisory.DocumentFilter
	List  *advisory.DocumentList
	Error error
}

type VulnerabilityMetadataOnList struct {
	Given *advisory.VulnerabilityMetadataFilter
	List  *advisory.VulnerabilityMetadataList
	Error error
}

type ResolvedVulnsReportsOnList struct {
	Given *advisory.ResolvedVulnsReportFilter
	List  *advisory.ResolvedVulnsReportList
	Error error
}

func (m MockSecurityAdvisoryClient) ListDocuments(_ context.Context, given *advisory.DocumentFilter, _ ...grpc.CallOption) (*advisory.DocumentList, error) { //nolint: revive
	for _, o := range m.OnListDocuments {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSecurityAdvisoryClient) ListVulnerabilityMetadata(_ context.Context, given *advisory.VulnerabilityMetadataFilter, _ ...grpc.CallOption) (*advisory.VulnerabilityMetadataList, error) { //nolint: revive
	for _, o := range m.OnListVulnerabilityMetadata {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSecurityAdvisoryClient) ListResolvedVulnsReports(_ context.Context, given *advisory.ResolvedVulnsReportFilter, _ ...grpc.CallOption) (*advisory.ResolvedVulnsReportList, error) { //nolint: revive
	for _, o := range m.OnListResolvedVulnsReports {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
