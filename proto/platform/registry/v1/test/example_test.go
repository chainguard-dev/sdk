/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"chainguard.dev/sdk/proto/platform/registry/v1/test"
)

// ExampleMockRegistryClients demonstrates creating a complete mock registry client.
func ExampleMockRegistryClients() {
	// Create a mock registry client with configured responses.
	mock := &test.MockRegistryClients{
		RegistryClient: test.MockRegistryClient{
			OnListRepos: []test.ReposOnList{{
				Given: &registry.RepoFilter{},
				List: &registry.RepoList{
					Items: []*registry.Repo{{
						Id:   "repo-123",
						Name: "test-repo",
					}},
				},
			}},
		},
	}

	// Use the mock client.
	repos, err := mock.Registry().ListRepos(context.Background(), &registry.RepoFilter{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d repos\n", len(repos.Items))
	fmt.Printf("First repo: %s\n", repos.Items[0].Name)

	// Output:
	// Found 1 repos
	// First repo: test-repo
}

// ExampleMockRegistryClients_close demonstrates handling the Close method.
func ExampleMockRegistryClients_close() {
	mock := &test.MockRegistryClients{
		OnClose: nil,
	}

	err := mock.Close()
	if err != nil {
		fmt.Printf("Close error: %v\n", err)
		return
	}

	fmt.Println("Client closed successfully")

	// Output:
	// Client closed successfully
}

// ExampleMockVulnerabilitiesClient demonstrates mocking vulnerability operations.
func ExampleMockVulnerabilitiesClient() {
	mock := test.MockVulnerabilitiesClient{
		OnListVulnReports: []test.VulnReportsOnList{{
			Given: &registry.VulnReportFilter{},
			List:  &registry.VulnReportList{},
		}},
	}

	list, err := mock.ListVulnReports(context.Background(), &registry.VulnReportFilter{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Vuln reports listed: %v\n", list != nil)

	// Output:
	// Vuln reports listed: true
}

// ExampleMockApkoClient demonstrates mocking Apko operations.
func ExampleMockApkoClient() {
	mock := test.MockApkoClient{
		OnBuildImage: []test.OnBuildImage{{
			Given:  &registry.BuildImageRequest{},
			Result: &registry.BuildImageResponse{},
		}},
	}

	result, err := mock.BuildImage(context.Background(), &registry.BuildImageRequest{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Build image result: %v\n", result != nil)

	// Output:
	// Build image result: true
}
