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

	OnListEntitlements      []EntitlementsOnList
	OnListEntitlementImages []EntitlementImagesOnList
	OnGetEntitlementSummary []EntitlementSummaryOnGet
	OnGetFeatures           []FeaturesOnGet
}

type EntitlementsOnList struct {
	Given *registry.EntitlementFilter
	List  *registry.EntitlementList
	Error error
}

type EntitlementImagesOnList struct {
	Given *registry.EntitlementImagesFilter
	List  *registry.EntitlementImagesList
	Error error
}

type EntitlementSummaryOnGet struct {
	Given *registry.EntitlementSummaryRequest
	Get   *registry.EntitlementSummaryResponse
	Error error
}

type FeaturesOnGet struct {
	Given *registry.GetFeaturesRequest
	Get   *registry.GetFeaturesResponse
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

func (m *MockEntitlementsClient) Summary(_ context.Context, given *registry.EntitlementSummaryRequest, _ ...grpc.CallOption) (*registry.EntitlementSummaryResponse, error) {
	for _, o := range m.OnGetEntitlementSummary {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockEntitlementsClient) GetFeatures(_ context.Context, given *registry.GetFeaturesRequest, _ ...grpc.CallOption) (*registry.GetFeaturesResponse, error) {
	for _, o := range m.OnGetFeatures {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
