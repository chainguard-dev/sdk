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

var _ iam.PoliciesClient = (*MockPoliciesClient)(nil)

type MockPoliciesClient struct {
	OnCreate        []PoliciesOnCreate
	OnUpdate        []PoliciesOnUpdate
	OnDelete        []PoliciesOnDelete
	OnList          []PoliciesOnList
	OnListVersions  []PoliciesOnListVersions
	OnActiveVersion []PoliciesOnActivateVersion
}

type PoliciesOnCreate struct {
	Given   *iam.CreatePolicyRequest
	Created *iam.Policy
	Error   error
}

type PoliciesOnUpdate struct {
	Given   *iam.Policy
	Updated *iam.Policy
	Error   error
}

type PoliciesOnDelete struct {
	Given *iam.DeletePolicyRequest
	Error error
}

type PoliciesOnList struct {
	Given *iam.PolicyFilter
	List  *iam.PolicyList
	Error error
}

type PoliciesOnListVersions struct {
	Given *iam.ListVersionsRequest
	List  *iam.PolicyVersionList
	Error error
}

type PoliciesOnActivateVersion struct {
	Given  *iam.ActivateVersionRequest
	Active *iam.Policy
	Error  error
}

func (m MockPoliciesClient) Create(_ context.Context, given *iam.CreatePolicyRequest, _ ...grpc.CallOption) (*iam.Policy, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockPoliciesClient) Update(_ context.Context, given *iam.Policy, _ ...grpc.CallOption) (*iam.Policy, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockPoliciesClient) Delete(_ context.Context, given *iam.DeletePolicyRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockPoliciesClient) List(_ context.Context, given *iam.PolicyFilter, _ ...grpc.CallOption) (*iam.PolicyList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockPoliciesClient) ListVersions(_ context.Context, given *iam.ListVersionsRequest, _ ...grpc.CallOption) (*iam.PolicyVersionList, error) {
	for _, o := range m.OnListVersions {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockPoliciesClient) ActivateVersion(_ context.Context, given *iam.ActivateVersionRequest, _ ...grpc.CallOption) (*iam.Policy, error) {
	for _, o := range m.OnActiveVersion {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Active, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
