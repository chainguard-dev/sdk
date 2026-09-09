/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"testing"

	"google.golang.org/genproto/googleapis/api/httpbody"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ registry.StigReportsServiceClient = (*MockStigReportsServiceClient)(nil)

type MockStigReportsServiceClient struct {
	registry.StigReportsServiceClient
	T *testing.T

	OnGetStigReport      []test.On[*registry.GetStigReportRequest, *registry.StigReport]
	OnDownloadStigReport []test.On[*registry.DownloadStigReportRequest, *httpbody.HttpBody]
}

func (m MockStigReportsServiceClient) GetStigReport(_ context.Context, given *registry.GetStigReportRequest, _ ...grpc.CallOption) (*registry.StigReport, error) {
	return test.Match(m.T, m.OnGetStigReport, given, "get-stig-report", protocmp.Transform())
}

func (m MockStigReportsServiceClient) DownloadStigReport(_ context.Context, given *registry.DownloadStigReportRequest, _ ...grpc.CallOption) (*httpbody.HttpBody, error) {
	return test.Match(m.T, m.OnDownloadStigReport, given, "download-stig-report", protocmp.Transform())
}
