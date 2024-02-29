/*
Copyright 2023 Chainguard, Inc.
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

	api "chainguard.dev/sdk/proto/platform/iam/v1"
)

var _ api.IdentityProvidersClient = (*MockIdentityProvidersClient)(nil)

type MockIdentityProvidersClient struct {
	OnCreate []IdentityProvidersOnCreate
	OnUpdate []IdentityProvidersOnUpdate
	OnDelete []IdentityProvidersOnDelete
	OnList   []IdentityProvidersOnList
}

type IdentityProvidersOnCreate struct {
	Given   *api.CreateIdentityProviderRequest
	Created *api.IdentityProvider
	Error   error
}

type IdentityProvidersOnUpdate struct {
	Given   *api.IdentityProvider
	Updated *api.IdentityProvider
	Error   error
}

type IdentityProvidersOnDelete struct {
	Given *api.DeleteIdentityProviderRequest
	Error error
}

type IdentityProvidersOnList struct {
	Given *api.IdentityProviderFilter
	List  *api.IdentityProviderList
	Error error
}

func (m MockIdentityProvidersClient) Create(_ context.Context, given *api.CreateIdentityProviderRequest, _ ...grpc.CallOption) (*api.IdentityProvider, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentityProvidersClient) Update(_ context.Context, given *api.IdentityProvider, _ ...grpc.CallOption) (*api.IdentityProvider, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentityProvidersClient) Delete(_ context.Context, given *api.DeleteIdentityProviderRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockIdentityProvidersClient) List(_ context.Context, given *api.IdentityProviderFilter, _ ...grpc.CallOption) (*api.IdentityProviderList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
