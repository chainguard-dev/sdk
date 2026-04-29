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

var _ iam.GroupsServiceClient = (*MockGroupsServiceClient)(nil)

type MockGroupsServiceClient struct {
	iam.GroupsServiceClient
	T *testing.T

	OnCreateGroup        []test.On[*iam.CreateGroupRequest, *iam.Group]
	OnUpdateGroup        []test.On[*iam.UpdateGroupRequest, *iam.Group]
	OnDeleteGroup        []test.On[*iam.DeleteGroupRequest, *emptypb.Empty]
	OnListGroups         []test.On[*iam.ListGroupsRequest, *iam.ListGroupsResponse]
	OnGetGroup           []test.On[*iam.GetGroupRequest, *iam.Group]
	OnLookupGroup        []test.On[*iam.LookupGroupRequest, *iam.LookupGroupResponse]
	OnRequestGroupAccess []test.On[*iam.RequestGroupAccessRequest, *iam.RequestGroupAccessResponse]
}

func (m MockGroupsServiceClient) CreateGroup(_ context.Context, given *iam.CreateGroupRequest, _ ...grpc.CallOption) (*iam.Group, error) {
	return test.Match(m.T, m.OnCreateGroup, given, "create-group", protocmp.Transform())
}

func (m MockGroupsServiceClient) UpdateGroup(_ context.Context, given *iam.UpdateGroupRequest, _ ...grpc.CallOption) (*iam.Group, error) {
	return test.Match(m.T, m.OnUpdateGroup, given, "update-group", protocmp.Transform())
}

func (m MockGroupsServiceClient) DeleteGroup(_ context.Context, given *iam.DeleteGroupRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteGroup, given, "delete-group", protocmp.Transform())
}

func (m MockGroupsServiceClient) ListGroups(_ context.Context, given *iam.ListGroupsRequest, _ ...grpc.CallOption) (*iam.ListGroupsResponse, error) {
	return test.Match(m.T, m.OnListGroups, given, "list-groups", protocmp.Transform())
}

func (m MockGroupsServiceClient) GetGroup(_ context.Context, given *iam.GetGroupRequest, _ ...grpc.CallOption) (*iam.Group, error) {
	return test.Match(m.T, m.OnGetGroup, given, "get-group", protocmp.Transform())
}

func (m MockGroupsServiceClient) LookupGroup(_ context.Context, given *iam.LookupGroupRequest, _ ...grpc.CallOption) (*iam.LookupGroupResponse, error) {
	return test.Match(m.T, m.OnLookupGroup, given, "lookup-group", protocmp.Transform())
}

func (m MockGroupsServiceClient) RequestGroupAccess(_ context.Context, given *iam.RequestGroupAccessRequest, _ ...grpc.CallOption) (*iam.RequestGroupAccessResponse, error) {
	return test.Match(m.T, m.OnRequestGroupAccess, given, "request-group-access", protocmp.Transform())
}
