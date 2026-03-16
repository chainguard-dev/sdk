/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package test provides mock implementations of Registry gRPC clients for testing.

This package offers mock clients that implement the Registry service interfaces,
enabling unit testing without requiring a live Registry service connection.

# Overview

The test package provides the following key features:

  - Mock implementations of all Registry gRPC client interfaces
  - Configurable responses for each RPC method
  - Request matching using protocol buffer comparison

# Basic Usage

Create a mock registry client and configure expected responses:

	mock := &test.MockRegistryClients{
		RegistryClient: test.MockRegistryClient{
			OnListRepos: []test.ReposOnList{{
				Given: &registry.RepoFilter{},
				List: &registry.RepoList{
					Items: []*registry.Repo{{
						Id:   "repo-id",
						Name: "test-repo",
					}},
				},
			}},
		},
	}

	// Use the mock in place of a real registry client
	repos, err := mock.Registry().ListRepos(ctx, &registry.RepoFilter{})

# Available Mock Clients

The package provides mock implementations for all Registry services:

  - MockRegistryClients: Aggregates all mock clients and implements registry.Clients
  - MockRegistryClient: Repository and tag management operations
  - MockVulnerabilitiesClient: Vulnerability report operations
  - MockApkoClient: Apko image build and config operations
  - MockEntitlementsClient: Entitlement management operations

# Request Matching

Mock responses are matched against incoming requests using protocol buffer
comparison. Configure multiple responses for different request patterns:

	mock := test.MockRegistryClient{
		OnListRepos: []test.ReposOnList{{
			Given: &registry.RepoFilter{},
			List:  &registry.RepoList{Items: repos},
		}},
	}

# Error Simulation

Simulate error conditions by setting the Error field:

	mock := test.MockRegistryClient{
		OnDeleteRepos: []test.ReposOnDelete{{
			Given: &registry.DeleteRepoRequest{Id: "protected-repo"},
			Error: status.Error(codes.PermissionDenied, "cannot delete protected repo"),
		}},
	}

# Thread Safety

Mock clients are safe for concurrent read access. Configure all expected
responses before using the mock in concurrent tests. Do not modify mock
configurations while tests are running.
*/
package test
