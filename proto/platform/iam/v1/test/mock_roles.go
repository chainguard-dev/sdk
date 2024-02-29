/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
)

var _ iam.RolesClient = (*MockRolesClient)(nil)

type MockRolesClient struct {
	OnCreate []RoleOnCreate
	OnUpdate []RoleOnUpdate
	OnDelete []RoleOnDelete
	OnList   []RoleOnList
}

type RoleOnCreate struct {
	Given   *iam.CreateRoleRequest
	Created *iam.Role
	Error   error
}

type RoleOnUpdate struct {
	Given   *iam.Role
	Updated *iam.Role
	Error   error
}

type RoleOnDelete struct {
	Given *iam.DeleteRoleRequest
	Error error
}

type RoleOnList struct {
	Given *iam.RoleFilter
	List  *iam.RoleList
	Error error
}

func (m MockRolesClient) Create(_ context.Context, given *iam.CreateRoleRequest, _ ...grpc.CallOption) (*iam.Role, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRolesClient) Update(_ context.Context, given *iam.Role, _ ...grpc.CallOption) (*iam.Role, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRolesClient) Delete(_ context.Context, given *iam.DeleteRoleRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockRolesClient) List(_ context.Context, given *iam.RoleFilter, _ ...grpc.CallOption) (*iam.RoleList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
