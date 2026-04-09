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

var _ iam.RoleBindingsServiceClient = (*MockRoleBindingsServiceClient)(nil)

type MockRoleBindingsServiceClient struct {
	iam.RoleBindingsServiceClient
	T *testing.T

	OnGetRoleBinding    []test.On[*iam.GetRoleBindingRequest, *iam.RoleBinding]
	OnCreateRoleBinding []test.On[*iam.CreateRoleBindingRequest, *iam.RoleBinding]
	OnUpdateRoleBinding []test.On[*iam.UpdateRoleBindingRequest, *iam.RoleBinding]
	OnDeleteRoleBinding []test.On[*iam.DeleteRoleBindingRequest, *emptypb.Empty]
	OnListRoleBindings  []test.On[*iam.ListRoleBindingsRequest, *iam.ListRoleBindingsResponse]
}

func (m MockRoleBindingsServiceClient) GetRoleBinding(_ context.Context, given *iam.GetRoleBindingRequest, _ ...grpc.CallOption) (*iam.RoleBinding, error) {
	return test.Match(m.T, m.OnGetRoleBinding, given, "get-role-binding", protocmp.Transform())
}

func (m MockRoleBindingsServiceClient) CreateRoleBinding(_ context.Context, given *iam.CreateRoleBindingRequest, _ ...grpc.CallOption) (*iam.RoleBinding, error) {
	return test.Match(m.T, m.OnCreateRoleBinding, given, "create-role-binding", protocmp.Transform())
}

func (m MockRoleBindingsServiceClient) UpdateRoleBinding(_ context.Context, given *iam.UpdateRoleBindingRequest, _ ...grpc.CallOption) (*iam.RoleBinding, error) {
	return test.Match(m.T, m.OnUpdateRoleBinding, given, "update-role-binding", protocmp.Transform())
}

func (m MockRoleBindingsServiceClient) DeleteRoleBinding(_ context.Context, given *iam.DeleteRoleBindingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteRoleBinding, given, "delete-role-binding", protocmp.Transform())
}

func (m MockRoleBindingsServiceClient) ListRoleBindings(_ context.Context, given *iam.ListRoleBindingsRequest, _ ...grpc.CallOption) (*iam.ListRoleBindingsResponse, error) {
	return test.Match(m.T, m.OnListRoleBindings, given, "list-role-bindings", protocmp.Transform())
}
