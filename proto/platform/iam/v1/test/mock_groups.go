/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ iam.GroupsClient = (*MockGroupsClient)(nil)

type MockGroupsClient struct {
	OnCreate []GroupOnCreate
	OnUpdate []GroupOnUpdate
	OnDelete []GroupOnDelete
	OnList   []GroupOnList
}

type GroupOnCreate struct {
	Given   *iam.CreateGroupRequest
	Created *iam.Group
	Error   error
}

type GroupOnUpdate struct {
	Given   *iam.Group
	Updated *iam.Group
	Error   error
}

type GroupOnDelete struct {
	Given *iam.DeleteGroupRequest
	Error error
}

type GroupOnList struct {
	Given *iam.GroupFilter
	List  *iam.GroupList
	Error error
}

func (m MockGroupsClient) Create(_ context.Context, given *iam.CreateGroupRequest, _ ...grpc.CallOption) (*iam.Group, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupsClient) Update(_ context.Context, given *iam.Group, _ ...grpc.CallOption) (*iam.Group, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupsClient) Delete(_ context.Context, given *iam.DeleteGroupRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupsClient) List(_ context.Context, given *iam.GroupFilter, _ ...grpc.CallOption) (*iam.GroupList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

// -- Server --

type MockGroupsServer struct {
	iam.UnimplementedGroupsServer
	Client MockGroupsClient
}

func (m MockGroupsServer) Create(ctx context.Context, req *iam.CreateGroupRequest) (*iam.Group, error) {
	return m.Client.Create(ctx, req)
}
func (m MockGroupsServer) Update(ctx context.Context, req *iam.Group) (*iam.Group, error) {
	return m.Client.Update(ctx, req)
}
func (m MockGroupsServer) List(ctx context.Context, req *iam.GroupFilter) (*iam.GroupList, error) {
	return m.Client.List(ctx, req)
}
func (m MockGroupsServer) Delete(ctx context.Context, req *iam.DeleteGroupRequest) (*emptypb.Empty, error) {
	return m.Client.Delete(ctx, req)
}
