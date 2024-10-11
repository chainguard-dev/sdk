/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
)

var _ registry.VulnerabilitiesClient = (*MockVulnerabilitiesClient)(nil)

type MockVulnerabilitiesClient struct {
	registry.VulnerabilitiesClient

	OnListVulnReports      []VulnReportsOnList
	OnGetRawVulnReport     []RawVulnReportOnGet
	OnListVulnCountReports []VulnCountReportsOnList
}

type VulnReportsOnList struct {
	Given *registry.VulnReportFilter
	List  *registry.VulnReportList
	Error error
}

type RawVulnReportOnGet struct {
	Given *registry.GetRawVulnReportRequest
	Get   *registry.RawVulnReport
	Error error
}

type VulnCountReportsOnList struct {
	Given *registry.VulnCountReportFilter
	List  *registry.VulnCountReportList
	Error error
}

func (m MockVulnerabilitiesClient) ListVulnReports(_ context.Context, given *registry.VulnReportFilter, _ ...grpc.CallOption) (*registry.VulnReportList, error) {
	for _, o := range m.OnListVulnReports {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockVulnerabilitiesClient) GetRawVulnReport(_ context.Context, given *registry.GetRawVulnReportRequest, _ ...grpc.CallOption) (*registry.RawVulnReport, error) {
	for _, o := range m.OnGetRawVulnReport {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockVulnerabilitiesClient) ListVulnCountReports(_ context.Context, given *registry.VulnCountReportFilter, _ ...grpc.CallOption) (*registry.VulnCountReportList, error) {
	for _, o := range m.OnListVulnCountReports {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
