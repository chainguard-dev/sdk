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

var _ iam.IdentitiesClient = (*MockIdentitiesClient)(nil)

type MockIdentitiesClient struct {
	OnCreate []IdentityOnCreate
	OnUpdate []IdentityOnUpdate
	OnDelete []IdentityOnDelete
	OnList   []IdentityOnList
	OnLooKup []IdentityOnLookup
}

type IdentityOnCreate struct {
	Given        *iam.CreateIdentityRequest
	Created      *iam.Identity
	Error        error
	IgnoreFields cmp.Option
}

type IdentityOnUpdate struct {
	Given   *iam.Identity
	Updated *iam.Identity
	Error   error
}

type IdentityOnDelete struct {
	Given *iam.DeleteIdentityRequest
	Error error
}

type IdentityOnList struct {
	Given *iam.IdentityFilter
	List  *iam.IdentityList
	Error error
}

type IdentityOnLookup struct {
	Given *iam.LookupRequest
	Found *iam.Identity
	Error error
}

func (m MockIdentitiesClient) Create(_ context.Context, given *iam.CreateIdentityRequest, _ ...grpc.CallOption) (*iam.Identity, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform(), o.IgnoreFields) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentitiesClient) Update(_ context.Context, given *iam.Identity, _ ...grpc.CallOption) (*iam.Identity, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentitiesClient) Delete(_ context.Context, given *iam.DeleteIdentityRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentitiesClient) List(_ context.Context, given *iam.IdentityFilter, _ ...grpc.CallOption) (*iam.IdentityList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentitiesClient) Lookup(_ context.Context, given *iam.LookupRequest, _ ...grpc.CallOption) (*iam.Identity, error) {
	for _, o := range m.OnLooKup {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Found, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
