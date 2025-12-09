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

var _ iam.RoleBindingsClient = (*MockRoleBindingsClient)(nil)

type MockRoleBindingsClient struct {
	OnCreate      []RoleBindingOnCreate
	OnCreateBatch []RoleBindingOnCreateBatch
	OnUpdate      []RoleBindingOnUpdate
	OnDelete      []RoleBindingOnDelete
	OnList        []RoleBindingOnList
}

type RoleBindingOnCreate struct {
	Given   *iam.CreateRoleBindingRequest
	Created *iam.RoleBinding
	Error   error
}

type RoleBindingOnCreateBatch struct {
	Given   *iam.CreateRoleBindingBatchRequest
	Created *iam.RoleBindingBatch
	Error   error
}

type RoleBindingOnUpdate struct {
	Given   *iam.RoleBinding
	Updated *iam.RoleBinding
	Error   error
}

type RoleBindingOnDelete struct {
	Given *iam.DeleteRoleBindingRequest
	Error error
}

type RoleBindingOnList struct {
	Given *iam.RoleBindingFilter
	List  *iam.RoleBindingList
	Error error
}

func (m MockRoleBindingsClient) Create(_ context.Context, given *iam.CreateRoleBindingRequest, _ ...grpc.CallOption) (*iam.RoleBinding, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRoleBindingsClient) CreateBatch(_ context.Context, given *iam.CreateRoleBindingBatchRequest, _ ...grpc.CallOption) (*iam.RoleBindingBatch, error) {
	for _, o := range m.OnCreateBatch {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRoleBindingsClient) Update(_ context.Context, given *iam.RoleBinding, _ ...grpc.CallOption) (*iam.RoleBinding, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRoleBindingsClient) Delete(_ context.Context, given *iam.DeleteRoleBindingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockRoleBindingsClient) List(_ context.Context, given *iam.RoleBindingFilter, _ ...grpc.CallOption) (*iam.RoleBindingList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
