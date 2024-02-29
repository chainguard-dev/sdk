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

var _ iam.GroupInvitesClient = (*MockGroupInvitesClient)(nil)

type MockGroupInvitesClient struct {
	OnCreate          []GroupInviteOnCreate
	OnCreateWithGroup []GroupInviteOnCreateWithGroup
	OnDelete          []GroupInviteOnDelete
	OnList            []GroupInviteOnList
}

type GroupInviteOnCreate struct {
	Given   *iam.GroupInviteRequest
	Created *iam.GroupInvite
	Error   error
}

type GroupInviteOnCreateWithGroup struct {
	Given   *iam.GroupInviteRequest
	Created *iam.GroupInvite
	Error   error
}

type GroupInviteOnDelete struct {
	Given *iam.DeleteGroupInviteRequest
	Error error
}

type GroupInviteOnList struct {
	Given *iam.GroupInviteFilter
	List  *iam.GroupInviteList
	Error error
}

func (m MockGroupInvitesClient) Create(_ context.Context, given *iam.GroupInviteRequest, _ ...grpc.CallOption) (*iam.GroupInvite, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupInvitesClient) CreateWithGroup(_ context.Context, given *iam.GroupInviteRequest, _ ...grpc.CallOption) (*iam.GroupInvite, error) {
	for _, o := range m.OnCreateWithGroup {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupInvitesClient) Delete(_ context.Context, given *iam.DeleteGroupInviteRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupInvitesClient) List(_ context.Context, given *iam.GroupInviteFilter, _ ...grpc.CallOption) (*iam.GroupInviteList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
