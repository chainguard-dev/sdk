/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"context"
	"fmt"

	packages "chainguard.dev/sdk/proto/platform/packages/v1"
	"chainguard.dev/sdk/proto/platform/packages/v1/test"
)

// ExampleClients demonstrates using the Clients interface with a mock.
func ExampleClients() {
	// In production, use packages.NewClients() with a real API URL and token.
	// Here we use a mock for demonstration.
	var clients packages.Clients = test.MockPackagesClients{
		EntitlementsOnClient: test.MockEntitlementsClient{
			OnList: []test.OnListEntitlements{{
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

	// Access the entitlements client.
	entitlements := clients.Entitlements()
	fmt.Printf("Entitlements client: %T\n", entitlements)

	// Always close the client when done.
	if err := clients.Close(); err != nil {
		fmt.Printf("Close error: %v\n", err)
	}

	// Output:
	// Entitlements client: *test.MockEntitlementsClient
}

// ExamplePackageTier demonstrates the PackageTier enum values.
func ExamplePackageTier() {
	tiers := []packages.PackageTier{
		packages.PackageTier_PACKAGE_TIER_UNSPECIFIED,
		packages.PackageTier_PACKAGE_TIER_BASE,
		packages.PackageTier_PACKAGE_TIER_FIPS,
	}

	for _, tier := range tiers {
		fmt.Println(tier.String())
	}

	// Output:
	// PACKAGE_TIER_UNSPECIFIED
	// PACKAGE_TIER_BASE
	// PACKAGE_TIER_FIPS
}

// ExamplePackageEntitlement demonstrates constructing a PackageEntitlement.
func ExamplePackageEntitlement() {
	e := &packages.PackageEntitlement{
		Id:   "entitlement-abc",
		Tier: packages.PackageTier_PACKAGE_TIER_FIPS,
	}

	fmt.Printf("ID: %s\n", e.GetId())
	fmt.Printf("Tier: %s\n", e.GetTier())

	// Output:
	// ID: entitlement-abc
	// Tier: PACKAGE_TIER_FIPS
}

// ExamplePackageEntitlementList demonstrates constructing a PackageEntitlementList.
func ExamplePackageEntitlementList() {
	list := &packages.PackageEntitlementList{
		Items: []*packages.PackageEntitlement{{
			Id:   "entitlement-1",
			Tier: packages.PackageTier_PACKAGE_TIER_BASE,
		}, {
			Id:   "entitlement-2",
			Tier: packages.PackageTier_PACKAGE_TIER_FIPS,
		}},
	}

	fmt.Printf("Count: %d\n", len(list.GetItems()))
	for _, item := range list.GetItems() {
		fmt.Printf("  %s: %s\n", item.GetId(), item.GetTier())
	}

	// Output:
	// Count: 2
	//   entitlement-1: PACKAGE_TIER_BASE
	//   entitlement-2: PACKAGE_TIER_FIPS
}

// ExamplePackageEntitlementFilter demonstrates constructing a filter.
func ExamplePackageEntitlementFilter() {
	filter := &packages.PackageEntitlementFilter{
		ParentId: "group-456",
	}

	fmt.Printf("ParentId: %s\n", filter.GetParentId())

	// Output:
	// ParentId: group-456
}

// ExampleNewClientsFromConnection demonstrates creating a Clients from an
// existing gRPC connection using NewClientsFromConnection.
func ExampleNewClientsFromConnection() {
	// NewClientsFromConnection wraps an existing *grpc.ClientConn.
	// Passing nil here is only for demonstration; in production use a real conn.
	clients := packages.NewClientsFromConnection(nil)

	// The returned Clients provides access to the Entitlements service.
	_ = clients.Entitlements()

	fmt.Println("clients created from connection")

	// Output:
	// clients created from connection
}

// Example_listEntitlements demonstrates listing entitlements via the mock.
func Example_listEntitlements() {
	mock := test.MockPackagesClients{
		EntitlementsOnClient: test.MockEntitlementsClient{
			OnList: []test.OnListEntitlements{{
				Given: &packages.PackageEntitlementFilter{
					ParentId: "group-789",
				},
				List: &packages.PackageEntitlementList{
					Items: []*packages.PackageEntitlement{{
						Id:   "ent-1",
						Tier: packages.PackageTier_PACKAGE_TIER_BASE,
					}},
				},
			}},
		},
	}

	list, err := mock.Entitlements().List(context.Background(), &packages.PackageEntitlementFilter{
		ParentId: "group-789",
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
