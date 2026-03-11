/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"chainguard.dev/sdk/proto/platform/registry/v1/test"
)

// ExampleClients demonstrates the Clients interface methods.
func ExampleClients() {
	// Create a mock registry client for demonstration.
	// In production, use registry.NewClients() with a real registry URL and token.
	var clients registry.Clients = &test.MockRegistryClients{}

	// Access individual service clients.
	registryClient := clients.Registry()
	vulnsClient := clients.Vulnerabilities()
	apkoClient := clients.Apko()
	entitlementsClient := clients.Entitlements()

	fmt.Printf("Registry client: %T\n", registryClient)
	fmt.Printf("Vulnerabilities client: %T\n", vulnsClient)
	fmt.Printf("Apko client: %T\n", apkoClient)
	fmt.Printf("Entitlements client: %T\n", entitlementsClient)

	// Always close the client when done.
	if err := clients.Close(); err != nil {
		fmt.Printf("Close error: %v\n", err)
	}

	// Output:
	// Registry client: *test.MockRegistryClient
	// Vulnerabilities client: *test.MockVulnerabilitiesClient
	// Apko client: *test.MockApkoClient
	// Entitlements client: *test.MockEntitlementsClient
}

// Example_mockClients demonstrates using mock clients for testing code that
// depends on the Clients interface.
func Example_mockClients() {
	// Use mock clients to test code that depends on registry.Clients.
	mock := &test.MockRegistryClients{
		RegistryClient: test.MockRegistryClient{
			OnListRepos: []test.ReposOnList{{
				Given: &registry.RepoFilter{},
				List: &registry.RepoList{
					Items: []*registry.Repo{{
						Id:   "repo-1",
						Name: "my-repo",
					}},
				},
			}},
		},
	}

	fmt.Printf("Registry client available: %v\n", mock.Registry() != nil)
	fmt.Printf("Vulnerabilities client available: %v\n", mock.Vulnerabilities() != nil)
	fmt.Printf("Apko client available: %v\n", mock.Apko() != nil)
	fmt.Printf("Entitlements client available: %v\n", mock.Entitlements() != nil)

	// Output:
	// Registry client available: true
	// Vulnerabilities client available: true
	// Apko client available: true
	// Entitlements client available: true
}
