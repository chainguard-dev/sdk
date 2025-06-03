/*
Copyright 2025 Chainguard, Inc.
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

var _ libraries.EntitlementsClient = (*MockEntitlementsClient)(nil)

type MockEntitlementsClient struct {
	libraries.EntitlementsClient

	OnCreate []EntitlementsOnCreate
	OnDelete []EntitlementsOnDelete
	OnList   []EntitlementsOnList
}

type EntitlementsOnCreate struct {
	Given   *libraries.CreateEntitlementRequest
	Created *libraries.Entitlement
	Error   error
}

type EntitlementsOnDelete struct {
	Given *libraries.DeleteEntitlementRequest
	Error error
}

type EntitlementsOnList struct {
	Given *libraries.EntitlementFilter
	List  *libraries.EntitlementList
	Error error
}

func (m MockEntitlementsClient) Create(_ context.Context, given *libraries.CreateEntitlementRequest, _ ...grpc.CallOption) (*libraries.Entitlement, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockEntitlementsClient) Delete(_ context.Context, given *libraries.DeleteEntitlementRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockEntitlementsClient) List(_ context.Context, given *libraries.EntitlementFilter, _ ...grpc.CallOption) (*libraries.EntitlementList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
