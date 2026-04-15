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

var _ iam.GroupInvitesServiceClient = (*MockGroupInvitesServiceClient)(nil)

type MockGroupInvitesServiceClient struct {
	iam.GroupInvitesServiceClient
	T *testing.T

	OnCreateGroupInvite []test.On[*iam.CreateGroupInviteRequest, *iam.GroupInvite]
	OnGetGroupInvite    []test.On[*iam.GetGroupInviteRequest, *iam.GroupInvite]
	OnDeleteGroupInvite []test.On[*iam.DeleteGroupInviteRequest, *emptypb.Empty]
	OnListGroupInvites  []test.On[*iam.ListGroupInvitesRequest, *iam.ListGroupInvitesResponse]
}

func (m MockGroupInvitesServiceClient) CreateGroupInvite(_ context.Context, given *iam.CreateGroupInviteRequest, _ ...grpc.CallOption) (*iam.GroupInvite, error) {
	return test.Match(m.T, m.OnCreateGroupInvite, given, "create-group-invite", protocmp.Transform())
}

func (m MockGroupInvitesServiceClient) GetGroupInvite(_ context.Context, given *iam.GetGroupInviteRequest, _ ...grpc.CallOption) (*iam.GroupInvite, error) {
	return test.Match(m.T, m.OnGetGroupInvite, given, "get-group-invite", protocmp.Transform())
}

func (m MockGroupInvitesServiceClient) DeleteGroupInvite(_ context.Context, given *iam.DeleteGroupInviteRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteGroupInvite, given, "delete-group-invite", protocmp.Transform())
}

func (m MockGroupInvitesServiceClient) ListGroupInvites(_ context.Context, given *iam.ListGroupInvitesRequest, _ ...grpc.CallOption) (*iam.ListGroupInvitesResponse, error) {
	return test.Match(m.T, m.OnListGroupInvites, given, "list-group-invites", protocmp.Transform())
}
