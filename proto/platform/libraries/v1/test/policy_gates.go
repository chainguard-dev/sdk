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
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
)

var _ libraries.LibraryPoliciesClient = (*MockLibraryPoliciesClient)(nil)

type MockLibraryPoliciesClient struct {
	libraries.LibraryPoliciesClient

	OnCreate []LibraryPoliciesOnCreate
	OnUpdate []LibraryPoliciesOnUpdate
	OnList   []LibraryPoliciesOnList
	OnDelete []LibraryPoliciesOnDelete
}

type LibraryPoliciesOnCreate struct {
	Given   *libraries.CreateLibraryPolicyRequest
	Created *libraries.LibraryPolicy
	Error   error
}

type LibraryPoliciesOnUpdate struct {
	Given   *libraries.LibraryPolicy
	Updated *libraries.LibraryPolicy
	Error   error
}

type LibraryPoliciesOnList struct {
	Given *libraries.LibraryPolicyFilter
	List  *libraries.LibraryPolicyList
	Error error
}

type LibraryPoliciesOnDelete struct {
	Given *libraries.DeleteLibraryPolicyRequest
	Error error
}

func (m MockLibraryPoliciesClient) CreatePolicy(_ context.Context, given *libraries.CreateLibraryPolicyRequest, _ ...grpc.CallOption) (*libraries.LibraryPolicy, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockLibraryPoliciesClient) UpdatePolicy(_ context.Context, given *libraries.LibraryPolicy, _ ...grpc.CallOption) (*libraries.LibraryPolicy, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockLibraryPoliciesClient) ListPolicies(_ context.Context, given *libraries.LibraryPolicyFilter, _ ...grpc.CallOption) (*libraries.LibraryPolicyList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockLibraryPoliciesClient) DeletePolicy(_ context.Context, given *libraries.DeleteLibraryPolicyRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

var _ libraries.LibraryPolicyBindingsClient = (*MockLibraryPolicyBindingsClient)(nil)

type MockLibraryPolicyBindingsClient struct {
	libraries.LibraryPolicyBindingsClient

	OnCreate []LibraryPolicyBindingsOnCreate
	OnUpdate []LibraryPolicyBindingsOnUpdate
	OnList   []LibraryPolicyBindingsOnList
	OnDelete []LibraryPolicyBindingsOnDelete
}

type LibraryPolicyBindingsOnCreate struct {
	Given   *libraries.CreateLibraryPolicyBindingRequest
	Created *libraries.LibraryPolicyBinding
	Error   error
}

type LibraryPolicyBindingsOnUpdate struct {
	Given   *libraries.LibraryPolicyBinding
	Updated *libraries.LibraryPolicyBinding
	Error   error
}

type LibraryPolicyBindingsOnList struct {
	Given *libraries.LibraryPolicyBindingFilter
	List  *libraries.LibraryPolicyBindingList
	Error error
}

type LibraryPolicyBindingsOnDelete struct {
	Given *libraries.DeleteLibraryPolicyBindingRequest
	Error error
}

func (m MockLibraryPolicyBindingsClient) CreateBinding(_ context.Context, given *libraries.CreateLibraryPolicyBindingRequest, _ ...grpc.CallOption) (*libraries.LibraryPolicyBinding, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockLibraryPolicyBindingsClient) UpdateBinding(_ context.Context, given *libraries.LibraryPolicyBinding, _ ...grpc.CallOption) (*libraries.LibraryPolicyBinding, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockLibraryPolicyBindingsClient) ListBindings(_ context.Context, given *libraries.LibraryPolicyBindingFilter, _ ...grpc.CallOption) (*libraries.LibraryPolicyBindingList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockLibraryPolicyBindingsClient) DeleteBinding(_ context.Context, given *libraries.DeleteLibraryPolicyBindingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

var _ libraries.LibraryPolicyBlockEventsClient = (*MockLibraryPolicyBlockEventsClient)(nil)

type MockLibraryPolicyBlockEventsClient struct {
	libraries.LibraryPolicyBlockEventsClient

	OnList []LibraryPolicyBlockEventsOnList
}

type LibraryPolicyBlockEventsOnList struct {
	Given *libraries.LibraryPolicyBlockEventFilter
	List  *libraries.LibraryPolicyBlockEventList
	Error error
}

func (m MockLibraryPolicyBlockEventsClient) ListBlockEvents(_ context.Context, given *libraries.LibraryPolicyBlockEventFilter, _ ...grpc.CallOption) (*libraries.LibraryPolicyBlockEventList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
