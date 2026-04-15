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

var _ iam.RolesServiceClient = (*MockRolesServiceClient)(nil)

type MockRolesServiceClient struct {
	iam.RolesServiceClient
	T *testing.T

	OnGetRole    []test.On[*iam.GetRoleRequest, *iam.Role]
	OnCreateRole []test.On[*iam.CreateRoleRequest, *iam.Role]
	OnUpdateRole []test.On[*iam.UpdateRoleRequest, *iam.Role]
	OnDeleteRole []test.On[*iam.DeleteRoleRequest, *emptypb.Empty]
	OnListRoles  []test.On[*iam.ListRolesRequest, *iam.ListRolesResponse]
}

func (m MockRolesServiceClient) GetRole(_ context.Context, given *iam.GetRoleRequest, _ ...grpc.CallOption) (*iam.Role, error) {
	return test.Match(m.T, m.OnGetRole, given, "get-role", protocmp.Transform())
}

func (m MockRolesServiceClient) CreateRole(_ context.Context, given *iam.CreateRoleRequest, _ ...grpc.CallOption) (*iam.Role, error) {
	return test.Match(m.T, m.OnCreateRole, given, "create-role", protocmp.Transform())
}

func (m MockRolesServiceClient) UpdateRole(_ context.Context, given *iam.UpdateRoleRequest, _ ...grpc.CallOption) (*iam.Role, error) {
	return test.Match(m.T, m.OnUpdateRole, given, "update-role", protocmp.Transform())
}

func (m MockRolesServiceClient) DeleteRole(_ context.Context, given *iam.DeleteRoleRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteRole, given, "delete-role", protocmp.Transform())
}

func (m MockRolesServiceClient) ListRoles(_ context.Context, given *iam.ListRolesRequest, _ ...grpc.CallOption) (*iam.ListRolesResponse, error) {
	return test.Match(m.T, m.OnListRoles, given, "list-roles", protocmp.Transform())
}
