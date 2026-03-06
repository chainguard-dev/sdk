/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"
	"log"

	v1 "chainguard.dev/sdk/proto/platform/libraries/v1"
	"chainguard.dev/sdk/proto/platform/libraries/v1/test"
)

// ExampleMockLibrariesClients demonstrates using the mock Libraries clients
// for testing.
func ExampleMockLibrariesClients() {
	ctx := context.Background()

	// Create a mock with configured responses
	mock := test.MockLibrariesClients{
		EntitlementsClient: test.MockEntitlementsClient{
			OnList: []test.EntitlementsOnList{{
				Given: &v1.EntitlementFilter{ParentId: "org/example"},
				List: &v1.EntitlementList{
					Items: []*v1.Entitlement{{
						Id:        "ent-123",
						Ecosystem: v1.Ecosystem_JAVA,
					}},
				},
			}},
		},
	}

	// Use the mock like a real client
	list, err := mock.Entitlements().List(ctx, &v1.EntitlementFilter{ParentId: "org/example"})
	if err != nil {
		log.Fatalf("failed to list entitlements: %v", err)
	}

	fmt.Printf("Found %d entitlements\n", len(list.Items))
	// Output: Found 1 entitlements
}

// ExampleMockArtifactsClient demonstrates using the mock Artifacts client.
func ExampleMockArtifactsClient() {
	ctx := context.Background()

	// Configure the mock with expected responses
	mock := test.MockArtifactsClient{
		OnList: []test.ArtifactsOnList{{
			Given: &v1.ArtifactFilter{Ecosystems: []v1.Ecosystem{v1.Ecosystem_JAVA}},
			List: &v1.ArtifactList{
				Items: []*v1.Artifact{{
					Id:   "artifact-123",
					Name: "example-artifact",
				}},
			},
		}},
		OnListVersions: []test.ArtifactsOnListVersions{{
			Given: &v1.ArtifactVersionFilter{Id: "artifact-123"},
			List: &v1.ArtifactVersionList{
				Items: []*v1.ArtifactVersion{{
					Id:      "version-456",
					Version: "1.0.0",
				}},
			},
		}},
	}

	// List artifacts
	artifacts, err := mock.List(ctx, &v1.ArtifactFilter{Ecosystems: []v1.Ecosystem{v1.Ecosystem_JAVA}})
	if err != nil {
		log.Fatalf("failed to list artifacts: %v", err)
	}

	// List versions for an artifact
	versions, err := mock.ListVersions(ctx, &v1.ArtifactVersionFilter{Id: "artifact-123"})
	if err != nil {
		log.Fatalf("failed to list versions: %v", err)
	}

	fmt.Printf("Found %d artifacts and %d versions\n", len(artifacts.Items), len(versions.Items))
	// Output: Found 1 artifacts and 1 versions
}

// ExampleMockEntitlementsClient demonstrates using the mock Entitlements client.
func ExampleMockEntitlementsClient() {
	ctx := context.Background()

	// Configure the mock with expected responses
	mock := test.MockEntitlementsClient{
		OnCreate: []test.EntitlementsOnCreate{{
			Given: &v1.CreateEntitlementRequest{
				ParentId:  "org/example",
				Ecosystem: v1.Ecosystem_PYTHON,
			},
			Created: &v1.Entitlement{
				Id:        "ent-789",
				Ecosystem: v1.Ecosystem_PYTHON,
			},
		}},
		OnDelete: []test.EntitlementsOnDelete{{
			Given: &v1.DeleteEntitlementRequest{
				Id: "ent-789",
			},
		}},
	}

	// Create an entitlement
	created, err := mock.Create(ctx, &v1.CreateEntitlementRequest{
		ParentId:  "org/example",
		Ecosystem: v1.Ecosystem_PYTHON,
	})
	if err != nil {
		log.Fatalf("failed to create entitlement: %v", err)
	}

	// Delete an entitlement
	_, err = mock.Delete(ctx, &v1.DeleteEntitlementRequest{Id: "ent-789"})
	if err != nil {
		log.Fatalf("failed to delete entitlement: %v", err)
	}

	fmt.Printf("Created entitlement: %s\n", created.Id)
	// Output: Created entitlement: ent-789
}

// ExampleMockNpmPackagesClient demonstrates using the mock NPM Packages client.
func ExampleMockNpmPackagesClient() {
	ctx := context.Background()

	// Configure the mock with expected responses
	mock := test.MockNpmPackagesClient{
		OnList: []test.NpmPackagesOnList{{
			Given: &v1.NpmPackageFilter{Query: "example"},
			List: &v1.NpmPackageList{
				Items: []*v1.NpmPackage{{
					PackageName: "example-package",
				}},
			},
		}},
		OnListVersions: []test.NpmPackagesOnListVersions{{
			Given: &v1.NpmPackageVersionFilter{PackageName: "example-package"},
			List: &v1.NpmPackageVersionList{
				Items: []*v1.NpmPackageVersion{{
					PackageName: "example-package",
					Version:     "2.0.0",
				}},
			},
		}},
	}

	// List NPM packages
	packages, err := mock.List(ctx, &v1.NpmPackageFilter{Query: "example"})
	if err != nil {
		log.Fatalf("failed to list packages: %v", err)
	}

	// List versions for a package
	versions, err := mock.ListVersions(ctx, &v1.NpmPackageVersionFilter{PackageName: "example-package"})
	if err != nil {
		log.Fatalf("failed to list versions: %v", err)
	}

	fmt.Printf("Found %d packages and %d versions\n", len(packages.Items), len(versions.Items))
	// Output: Found 1 packages and 1 versions
}
