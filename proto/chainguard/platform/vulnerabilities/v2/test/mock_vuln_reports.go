/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	"chainguard.dev/sdk/proto/chainguard/platform/test"
	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2"
)

var _ vuln.VulnReportsServiceClient = (*MockVulnReportsServiceClient)(nil)

type MockVulnReportsServiceClient struct {
	vuln.VulnReportsServiceClient
	T *testing.T

	OnListVulnCountReports []test.On[*vuln.ListVulnCountReportsRequest, *vuln.ListVulnCountReportsResponse]
}

func (m MockVulnReportsServiceClient) ListVulnCountReports(_ context.Context, given *vuln.ListVulnCountReportsRequest, _ ...grpc.CallOption) (*vuln.ListVulnCountReportsResponse, error) {
	return test.Match(m.T, m.OnListVulnCountReports, given, "list-vuln-count-reports", protocmp.Transform())
}
