/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ registry.EntitlementsClient = (*MockEntitlementsClient)(nil)

type MockEntitlementsClient struct {
	registry.EntitlementsClient

	OnListEntitlements      []ListOnEntitlements
	OnListEntitlementImages []ListOnEntitlementImages
}

type ListOnEntitlements struct {
	Given *registry.EntitlementFilter
	List  *registry.EntitlementList
	Error error
}

type ListOnEntitlementImages struct {
	Given *registry.EntitlementImagesFilter
	List  *registry.EntitlementImagesList
	Error error
}

func (m *MockEntitlementsClient) ListEntitlements(_ context.Context, given *registry.EntitlementFilter, _ ...grpc.CallOption) (*registry.EntitlementList, error) {
	for _, o := range m.OnListEntitlements {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockEntitlementsClient) ListEntitlementImages(_ context.Context, given *registry.EntitlementImagesFilter, _ ...grpc.CallOption) (*registry.EntitlementImagesList, error) {
	for _, o := range m.OnListEntitlementImages {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
