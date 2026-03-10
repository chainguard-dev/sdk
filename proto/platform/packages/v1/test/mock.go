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

	packages "chainguard.dev/sdk/proto/platform/packages/v1"
)

var _ packages.Clients = (*MockPackagesClients)(nil)

type MockPackagesClients struct {
	EntitlementsOnClient MockEntitlementsClient

	OnClose error
}

func (m MockPackagesClients) Entitlements() packages.EntitlementsClient {
	return &m.EntitlementsOnClient
}

func (m MockPackagesClients) Close() error {
	return m.OnClose
}

// MockEntitlementsClient mocks the Entitlements service.
var _ packages.EntitlementsClient = (*MockEntitlementsClient)(nil)

type MockEntitlementsClient struct {
	packages.EntitlementsClient

	OnList []OnListEntitlements
}

type OnListEntitlements struct {
	Given *packages.PackageEntitlementFilter
	List  *packages.PackageEntitlementList
	Error error
}

func (m *MockEntitlementsClient) List(_ context.Context, given *packages.PackageEntitlementFilter, _ ...grpc.CallOption) (*packages.PackageEntitlementList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
