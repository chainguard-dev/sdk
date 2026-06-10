/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
)

var _ iam.ExternalGroupRoleMappingsClient = (*MockExternalGroupRoleMappingsClient)(nil)

type MockExternalGroupRoleMappingsClient struct {
	OnCreate []EGRMOnCreate
	OnGet    []EGRMOnGet
	OnList   []EGRMOnList
	OnDelete []EGRMOnDelete
}

type EGRMOnCreate struct {
	Given   *iam.CreateExternalGroupRoleMappingRequest
	Created *iam.ExternalGroupRoleMapping
	Error   error
}

type EGRMOnGet struct {
	Given   *iam.GetExternalGroupRoleMappingRequest
	Mapping *iam.ExternalGroupRoleMapping
	Error   error
}

type EGRMOnList struct {
	Given *iam.ExternalGroupRoleMappingFilter
	List  *iam.ExternalGroupRoleMappingList
	Error error
}

type EGRMOnDelete struct {
	Given *iam.DeleteExternalGroupRoleMappingRequest
	Error error
}

func (m MockExternalGroupRoleMappingsClient) Create(_ context.Context, given *iam.CreateExternalGroupRoleMappingRequest, _ ...grpc.CallOption) (*iam.ExternalGroupRoleMapping, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockExternalGroupRoleMappingsClient) Get(_ context.Context, given *iam.GetExternalGroupRoleMappingRequest, _ ...grpc.CallOption) (*iam.ExternalGroupRoleMapping, error) {
	for _, o := range m.OnGet {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Mapping, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockExternalGroupRoleMappingsClient) List(_ context.Context, given *iam.ExternalGroupRoleMappingFilter, _ ...grpc.CallOption) (*iam.ExternalGroupRoleMappingList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return proto.Clone(o.List).(*iam.ExternalGroupRoleMappingList), o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockExternalGroupRoleMappingsClient) Delete(_ context.Context, given *iam.DeleteExternalGroupRoleMappingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}
