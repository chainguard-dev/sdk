/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"errors"
	"fmt"

	packages "chainguard.dev/sdk/proto/platform/packages/v1"
	"chainguard.dev/sdk/proto/platform/packages/v1/test"
)

// ExampleMockPackagesClients demonstrates creating a mock Packages client.
func ExampleMockPackagesClients() {
	mock := test.MockPackagesClients{
		EntitlementsClient: test.MockEntitlementsClient{
			OnList: []test.EntitlementsOnList{{
				Given: &packages.PackageEntitlementFilter{
					ParentId: "group-123",
				},
				List: &packages.PackageEntitlementList{
					Items: []*packages.PackageEntitlement{{
						Id:   "entitlement-1",
						Tier: packages.PackageTier_PACKAGE_TIER_BASE,
					}},
				},
			}},
		},
	}

	list, err := mock.Entitlements().List(context.Background(), &packages.PackageEntitlementFilter{
		ParentId: "group-123",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Found %d entitlement(s)\n", len(list.GetItems()))
	fmt.Printf("Tier: %s\n", list.GetItems()[0].GetTier())

	// Output:
	// Found 1 entitlement(s)
	// Tier: PACKAGE_TIER_BASE
}

// ExampleMockPackagesClients_close demonstrates the Close method.
func ExampleMockPackagesClients_close() {
	mock := test.MockPackagesClients{
		OnClose: nil,
	}

	if err := mock.Close(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println("closed successfully")

	// Output:
	// closed successfully
}

// ExampleMockPackagesClients_closeError demonstrates Close returning an error.
func ExampleMockPackagesClients_closeError() {
	mock := test.MockPackagesClients{
		OnClose: errors.New("connection reset"),
	}

	if err := mock.Close(); err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// Output:
	// error: connection reset
}

// ExampleMockEntitlementsClient demonstrates mocking the Entitlements service.
func ExampleMockEntitlementsClient() {
	mock := test.MockEntitlementsClient{
		OnList: []test.EntitlementsOnList{{
			Given: &packages.PackageEntitlementFilter{
				ParentId: "group-456",
			},
			List: &packages.PackageEntitlementList{
				Items: []*packages.PackageEntitlement{{
					Id:   "ent-fips",
					Tier: packages.PackageTier_PACKAGE_TIER_FIPS,
				}, {
					Id:   "ent-base",
					Tier: packages.PackageTier_PACKAGE_TIER_BASE,
				}},
			},
		}},
	}

	list, err := mock.List(context.Background(), &packages.PackageEntitlementFilter{
		ParentId: "group-456",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("Found %d entitlement(s)\n", len(list.GetItems()))
	for _, item := range list.GetItems() {
		fmt.Printf("  %s: %s\n", item.GetId(), item.GetTier())
	}

	// Output:
	// Found 2 entitlement(s)
	//   ent-fips: PACKAGE_TIER_FIPS
	//   ent-base: PACKAGE_TIER_BASE
}

// ExampleMockEntitlementsClient_error demonstrates simulating a List error.
func ExampleMockEntitlementsClient_error() {
	mock := test.MockEntitlementsClient{
		OnList: []test.EntitlementsOnList{{
			Given: &packages.PackageEntitlementFilter{
				ParentId: "group-789",
			},
			Error: errors.New("permission denied"),
		}},
	}

	_, err := mock.List(context.Background(), &packages.PackageEntitlementFilter{
		ParentId: "group-789",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	// Output:
	// error: permission denied
}
