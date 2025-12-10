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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
)

var _ iam.GroupAccountAssociationsClient = (*MockGroupAccountAssociationsClient)(nil)

type MockGroupAccountAssociationsClient struct {
	OnCreate []AccountAssociationsOnCreate
	OnUpdate []AccountAssociationsOnUpdate
	OnDelete []AccountAssociationsOnDelete
	OnList   []AccountAssociationsOnList
	OnCheck  []AccountAssociationsOnCheck
}

type AccountAssociationsOnCreate struct {
	Given   *iam.AccountAssociations
	Created *iam.AccountAssociations
	Error   error
}

type AccountAssociationsOnUpdate struct {
	Given   *iam.AccountAssociations
	Updated *iam.AccountAssociations
	Error   error
}

type AccountAssociationsOnDelete struct {
	Given *iam.DeleteAccountAssociationsRequest
	Error error
}

type AccountAssociationsOnList struct {
	Given *iam.AccountAssociationsFilter
	List  *iam.AccountAssociationsList
	Error error
}

type AccountAssociationsOnCheck struct {
	Given  *iam.AccountAssociationsCheckRequest
	Status *iam.AccountAssociationsStatus
	Error  error
}

func (m MockGroupAccountAssociationsClient) Create(_ context.Context, given *iam.AccountAssociations, _ ...grpc.CallOption) (*iam.AccountAssociations, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupAccountAssociationsClient) Update(_ context.Context, given *iam.AccountAssociations, _ ...grpc.CallOption) (*iam.AccountAssociations, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupAccountAssociationsClient) Delete(_ context.Context, given *iam.DeleteAccountAssociationsRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupAccountAssociationsClient) List(_ context.Context, given *iam.AccountAssociationsFilter, _ ...grpc.CallOption) (*iam.AccountAssociationsList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return proto.Clone(o.List).(*iam.AccountAssociationsList), o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGroupAccountAssociationsClient) Check(_ context.Context, given *iam.AccountAssociationsCheckRequest, _ ...grpc.CallOption) (*iam.AccountAssociationsStatus, error) {
	for _, o := range m.OnCheck {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return proto.Clone(o.Status).(*iam.AccountAssociationsStatus), o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
