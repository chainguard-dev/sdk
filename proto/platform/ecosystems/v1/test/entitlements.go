/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

//nolint:staticcheck
package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	ecosystems "chainguard.dev/sdk/proto/platform/ecosystems/v1"
)

var _ ecosystems.EntitlementsClient = (*MockEntitlementsClient)(nil)

type MockEntitlementsClient struct {
	ecosystems.EntitlementsClient

	OnCreate []EntitlementsOnCreate
	OnDelete []EntitlementsOnDelete
	OnList   []EntitlementsOnList
}

type EntitlementsOnCreate struct {
	Given   *ecosystems.CreateEntitlementRequest
	Created *ecosystems.Entitlement
	Error   error
}

type EntitlementsOnDelete struct {
	Given *ecosystems.DeleteEntitlementRequest
	Error error
}

type EntitlementsOnList struct {
	Given *ecosystems.EntitlementFilter
	List  *ecosystems.EntitlementList
	Error error
}

func (m MockEntitlementsClient) Create(_ context.Context, given *ecosystems.CreateEntitlementRequest, _ ...grpc.CallOption) (*ecosystems.Entitlement, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockEntitlementsClient) Delete(_ context.Context, given *ecosystems.DeleteEntitlementRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockEntitlementsClient) List(_ context.Context, given *ecosystems.EntitlementFilter, _ ...grpc.CallOption) (*ecosystems.EntitlementList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
