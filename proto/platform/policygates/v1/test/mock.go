/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	policygates "chainguard.dev/sdk/proto/platform/policygates/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ policygates.Clients = (*MockPolicyGatesClients)(nil)

type MockPolicyGatesClients struct {
	PoliciesOnClient MockPoliciesClient
	BindingsOnClient MockBindingsClient

	OnClose error
}

func (m MockPolicyGatesClients) Policies() policygates.PoliciesClient {
	return &m.PoliciesOnClient
}

func (m MockPolicyGatesClients) Bindings() policygates.BindingsClient {
	return &m.BindingsOnClient
}

func (m MockPolicyGatesClients) Close() error {
	return m.OnClose
}

// MockPoliciesClient mocks the Policies service.
var _ policygates.PoliciesClient = (*MockPoliciesClient)(nil)

type MockPoliciesClient struct {
	policygates.PoliciesClient

	OnCreatePolicy []OnCreatePolicy
	OnUpdatePolicy []OnUpdatePolicy
	OnListPolicies []OnListPolicies
	OnDeletePolicy []OnDeletePolicy
}

type OnCreatePolicy struct {
	Given   *policygates.CreatePolicyRequest
	Created *policygates.Policy
	Error   error
}

type OnUpdatePolicy struct {
	Given   *policygates.Policy
	Updated *policygates.Policy
	Error   error
}

type OnListPolicies struct {
	Given *policygates.PolicyFilter
	List  *policygates.PolicyList
	Error error
}

type OnDeletePolicy struct {
	Given *policygates.DeletePolicyRequest
	Error error
}

func (m *MockPoliciesClient) CreatePolicy(_ context.Context, given *policygates.CreatePolicyRequest, _ ...grpc.CallOption) (*policygates.Policy, error) {
	for _, o := range m.OnCreatePolicy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockPoliciesClient) UpdatePolicy(_ context.Context, given *policygates.Policy, _ ...grpc.CallOption) (*policygates.Policy, error) {
	for _, o := range m.OnUpdatePolicy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockPoliciesClient) ListPolicies(_ context.Context, given *policygates.PolicyFilter, _ ...grpc.CallOption) (*policygates.PolicyList, error) {
	for _, o := range m.OnListPolicies {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockPoliciesClient) DeletePolicy(_ context.Context, given *policygates.DeletePolicyRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeletePolicy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

// MockBindingsClient mocks the Bindings service.
var _ policygates.BindingsClient = (*MockBindingsClient)(nil)

type MockBindingsClient struct {
	policygates.BindingsClient

	OnCreateBinding []OnCreateBinding
	OnUpdateBinding []OnUpdateBinding
	OnListBindings  []OnListBindings
	OnDeleteBinding []OnDeleteBinding
}

type OnCreateBinding struct {
	Given   *policygates.CreateBindingRequest
	Created *policygates.Binding
	Error   error
}

type OnUpdateBinding struct {
	Given   *policygates.Binding
	Updated *policygates.Binding
	Error   error
}

type OnListBindings struct {
	Given *policygates.BindingFilter
	List  *policygates.BindingList
	Error error
}

type OnDeleteBinding struct {
	Given *policygates.DeleteBindingRequest
	Error error
}

func (m *MockBindingsClient) CreateBinding(_ context.Context, given *policygates.CreateBindingRequest, _ ...grpc.CallOption) (*policygates.Binding, error) {
	for _, o := range m.OnCreateBinding {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockBindingsClient) UpdateBinding(_ context.Context, given *policygates.Binding, _ ...grpc.CallOption) (*policygates.Binding, error) {
	for _, o := range m.OnUpdateBinding {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockBindingsClient) ListBindings(_ context.Context, given *policygates.BindingFilter, _ ...grpc.CallOption) (*policygates.BindingList, error) {
	for _, o := range m.OnListBindings {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockBindingsClient) DeleteBinding(_ context.Context, given *policygates.DeleteBindingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteBinding {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
