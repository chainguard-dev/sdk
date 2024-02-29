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

	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
)

var _ tenant.VulnReportsClient = (*MockVulnReportsClient)(nil)

type MockVulnReportsClient struct {
	OnList []VulnReportsOnList
}

type VulnReportsOnList struct {
	Given *tenant.VulnReportFilter
	List  *tenant.VulnReportList
	Error error
}

func (m MockVulnReportsClient) List(_ context.Context, given *tenant.VulnReportFilter, _ ...grpc.CallOption) (*tenant.VulnReportList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
