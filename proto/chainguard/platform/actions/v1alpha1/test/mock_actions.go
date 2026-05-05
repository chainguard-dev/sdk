/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package test provides mocks for the chainguard.platform.actions.v1alpha1
// Actions service, suitable for unit testing consumers of the SDK.
package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	actions "chainguard.dev/sdk/proto/chainguard/platform/actions/v1alpha1"
)

var _ actions.ActionsClient = (*MockActionsClient)(nil)

// MockActionsClient is a mock implementation of actions.ActionsClient for tests.
// Configure expected calls via the OnCreateEntitlement, OnGetEntitlement, and
// OnDeleteEntitlement slices. Calls without a matching configured expectation
// return an error.
type MockActionsClient struct {
	actions.ActionsClient

	OnCreateEntitlement []ActionsOnCreateEntitlement
	OnGetEntitlement    []ActionsOnGetEntitlement
	OnDeleteEntitlement []ActionsOnDeleteEntitlement
}

type ActionsOnCreateEntitlement struct {
	Given   *actions.CreateEntitlementRequest
	Created *actions.Entitlement
	Error   error
}

type ActionsOnGetEntitlement struct {
	Given *actions.GetEntitlementRequest
	Got   *actions.Entitlement
	Error error
}

type ActionsOnDeleteEntitlement struct {
	Given *actions.DeleteEntitlementRequest
	Error error
}

func (m MockActionsClient) CreateEntitlement(_ context.Context, given *actions.CreateEntitlementRequest, _ ...grpc.CallOption) (*actions.Entitlement, error) {
	for _, o := range m.OnCreateEntitlement {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockActionsClient) GetEntitlement(_ context.Context, given *actions.GetEntitlementRequest, _ ...grpc.CallOption) (*actions.Entitlement, error) {
	for _, o := range m.OnGetEntitlement {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Got, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockActionsClient) DeleteEntitlement(_ context.Context, given *actions.DeleteEntitlementRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteEntitlement {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}
