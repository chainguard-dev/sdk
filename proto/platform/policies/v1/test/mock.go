/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	policies "chainguard.dev/sdk/proto/platform/policies/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ policies.Clients = (*MockPoliciesClients)(nil)

type MockPoliciesClients struct {
	PoliciesOnClient MockPoliciesClient
	BindingsOnClient MockBindingsClient

	OnClose error
}

func (m MockPoliciesClients) Policies() policies.PoliciesClient {
	return &m.PoliciesOnClient
}

func (m MockPoliciesClients) Bindings() policies.BindingsClient {
	return &m.BindingsOnClient
}

func (m MockPoliciesClients) Close() error {
	return m.OnClose
}

// MockPoliciesClient mocks the Policies service.
var _ policies.PoliciesClient = (*MockPoliciesClient)(nil)

type MockPoliciesClient struct {
	policies.PoliciesClient

	OnCreatePolicy []OnCreatePolicy
	OnUpdatePolicy []OnUpdatePolicy
	OnListPolicies []OnListPolicies
	OnDeletePolicy []OnDeletePolicy
}

type OnCreatePolicy struct {
	Given   *policies.CreatePolicyRequest
	Created *policies.Policy
	Error   error
}

type OnUpdatePolicy struct {
	Given   *policies.Policy
	Updated *policies.Policy
	Error   error
}

type OnListPolicies struct {
	Given *policies.PolicyFilter
	List  *policies.PolicyList
	Error error
}

type OnDeletePolicy struct {
	Given *policies.DeletePolicyRequest
	Error error
}

func (m *MockPoliciesClient) CreatePolicy(_ context.Context, given *policies.CreatePolicyRequest, _ ...grpc.CallOption) (*policies.Policy, error) {
	for _, o := range m.OnCreatePolicy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockPoliciesClient) UpdatePolicy(_ context.Context, given *policies.Policy, _ ...grpc.CallOption) (*policies.Policy, error) {
	for _, o := range m.OnUpdatePolicy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockPoliciesClient) ListPolicies(_ context.Context, given *policies.PolicyFilter, _ ...grpc.CallOption) (*policies.PolicyList, error) {
	for _, o := range m.OnListPolicies {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockPoliciesClient) DeletePolicy(_ context.Context, given *policies.DeletePolicyRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeletePolicy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

// MockBindingsClient mocks the Bindings service.
var _ policies.BindingsClient = (*MockBindingsClient)(nil)

type MockBindingsClient struct {
	policies.BindingsClient

	OnCreateBinding []OnCreateBinding
	OnUpdateBinding []OnUpdateBinding
	OnListBindings  []OnListBindings
	OnDeleteBinding []OnDeleteBinding
}

type OnCreateBinding struct {
	Given   *policies.CreateBindingRequest
	Created *policies.Binding
	Error   error
}

type OnUpdateBinding struct {
	Given   *policies.Binding
	Updated *policies.Binding
	Error   error
}

type OnListBindings struct {
	Given *policies.BindingFilter
	List  *policies.BindingList
	Error error
}

type OnDeleteBinding struct {
	Given *policies.DeleteBindingRequest
	Error error
}

func (m *MockBindingsClient) CreateBinding(_ context.Context, given *policies.CreateBindingRequest, _ ...grpc.CallOption) (*policies.Binding, error) {
	for _, o := range m.OnCreateBinding {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockBindingsClient) UpdateBinding(_ context.Context, given *policies.Binding, _ ...grpc.CallOption) (*policies.Binding, error) {
	for _, o := range m.OnUpdateBinding {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockBindingsClient) ListBindings(_ context.Context, given *policies.BindingFilter, _ ...grpc.CallOption) (*policies.BindingList, error) {
	for _, o := range m.OnListBindings {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockBindingsClient) DeleteBinding(_ context.Context, given *policies.DeleteBindingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteBinding {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
