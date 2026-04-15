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
	vuln "chainguard.dev/sdk/proto/chainguard/platform/vulnerabilities/v2beta1"
)

var _ vuln.AdvisoriesServiceClient = (*MockAdvisoriesServiceClient)(nil)

type MockAdvisoriesServiceClient struct {
	vuln.AdvisoriesServiceClient
	T *testing.T

	OnGetAdvisory    []test.On[*vuln.GetAdvisoryRequest, *vuln.Advisory]
	OnListAdvisories []test.On[*vuln.ListAdvisoriesRequest, *vuln.ListAdvisoriesResponse]
}

func (m MockAdvisoriesServiceClient) GetAdvisory(_ context.Context, given *vuln.GetAdvisoryRequest, _ ...grpc.CallOption) (*vuln.Advisory, error) {
	return test.Match(m.T, m.OnGetAdvisory, given, "get-advisory", protocmp.Transform())
}

func (m MockAdvisoriesServiceClient) ListAdvisories(_ context.Context, given *vuln.ListAdvisoriesRequest, _ ...grpc.CallOption) (*vuln.ListAdvisoriesResponse, error) {
	return test.Match(m.T, m.OnListAdvisories, given, "list-advisories", protocmp.Transform())
}
