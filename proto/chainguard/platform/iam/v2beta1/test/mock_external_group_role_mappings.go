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

var _ iam.ExternalGroupRoleMappingsServiceClient = (*MockExternalGroupRoleMappingsServiceClient)(nil)

type MockExternalGroupRoleMappingsServiceClient struct {
	iam.ExternalGroupRoleMappingsServiceClient
	T *testing.T

	OnGetExternalGroupRoleMapping    []test.On[*iam.GetExternalGroupRoleMappingRequest, *iam.ExternalGroupRoleMapping]
	OnCreateExternalGroupRoleMapping []test.On[*iam.CreateExternalGroupRoleMappingRequest, *iam.ExternalGroupRoleMapping]
	OnDeleteExternalGroupRoleMapping []test.On[*iam.DeleteExternalGroupRoleMappingRequest, *emptypb.Empty]
	OnListExternalGroupRoleMappings  []test.On[*iam.ListExternalGroupRoleMappingsRequest, *iam.ListExternalGroupRoleMappingsResponse]
}

func (m MockExternalGroupRoleMappingsServiceClient) GetExternalGroupRoleMapping(_ context.Context, given *iam.GetExternalGroupRoleMappingRequest, _ ...grpc.CallOption) (*iam.ExternalGroupRoleMapping, error) {
	return test.Match(m.T, m.OnGetExternalGroupRoleMapping, given, "get-external-group-role-mapping", protocmp.Transform())
}

func (m MockExternalGroupRoleMappingsServiceClient) CreateExternalGroupRoleMapping(_ context.Context, given *iam.CreateExternalGroupRoleMappingRequest, _ ...grpc.CallOption) (*iam.ExternalGroupRoleMapping, error) {
	return test.Match(m.T, m.OnCreateExternalGroupRoleMapping, given, "create-external-group-role-mapping", protocmp.Transform())
}

func (m MockExternalGroupRoleMappingsServiceClient) DeleteExternalGroupRoleMapping(_ context.Context, given *iam.DeleteExternalGroupRoleMappingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteExternalGroupRoleMapping, given, "delete-external-group-role-mapping", protocmp.Transform())
}

func (m MockExternalGroupRoleMappingsServiceClient) ListExternalGroupRoleMappings(_ context.Context, given *iam.ListExternalGroupRoleMappingsRequest, _ ...grpc.CallOption) (*iam.ListExternalGroupRoleMappingsResponse, error) {
	return test.Match(m.T, m.OnListExternalGroupRoleMappings, given, "list-external-group-role-mappings", protocmp.Transform())
}
