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
	"google.golang.org/protobuf/types/known/emptypb"

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ iam.SessionsServiceClient = (*MockSessionsServiceClient)(nil)

type MockSessionsServiceClient struct {
	iam.SessionsServiceClient
	T *testing.T

	OnListSessions      []test.On[*iam.ListSessionsRequest, *iam.ListSessionsResponse]
	OnDeleteSession     []test.On[*iam.DeleteSessionRequest, *emptypb.Empty]
	OnRevokeAllSessions []test.On[*iam.RevokeAllSessionsRequest, *iam.RevokeAllSessionsResponse]
}

func (m MockSessionsServiceClient) ListSessions(_ context.Context, given *iam.ListSessionsRequest, _ ...grpc.CallOption) (*iam.ListSessionsResponse, error) {
	return test.Match(m.T, m.OnListSessions, given, "list-sessions", protocmp.Transform())
}

func (m MockSessionsServiceClient) DeleteSession(_ context.Context, given *iam.DeleteSessionRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteSession, given, "delete-session", protocmp.Transform())
}

func (m MockSessionsServiceClient) RevokeAllSessions(_ context.Context, given *iam.RevokeAllSessionsRequest, _ ...grpc.CallOption) (*iam.RevokeAllSessionsResponse, error) {
	return test.Match(m.T, m.OnRevokeAllSessions, given, "revoke-all-sessions", protocmp.Transform())
}
