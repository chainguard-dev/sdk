/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package test provides mock implementations of IAM gRPC clients for testing.

This package offers mock clients that implement the IAM service interfaces,
enabling unit testing without requiring a live IAM service connection.

# Overview

The test package provides the following key features:

  - Mock implementations of all IAM gRPC client interfaces
  - Configurable responses for each RPC method
  - Request matching using protocol buffer comparison
  - Both client and server mock implementations
  - Thread-safe mock configurations

# Basic Usage

Create a mock IAM client and configure expected responses:

	mock := &test.MockIAMClient{
		GroupsClient: test.MockGroupsClient{
			OnList: []test.GroupOnList{{
				Given: &iam.GroupFilter{},
				List: &iam.GroupList{
					Items: []*iam.Group{{
						Id:   "group-id",
						Name: "test-group",
					}},
				},
			}},
		},
	}

	// Use the mock in place of a real IAM client
	groups, err := mock.Groups().List(ctx, &iam.GroupFilter{})

# Available Mock Clients

The package provides mock implementations for all IAM services:

  - MockGroupsClient: Groups management operations
  - MockGroupInvitesClient: Group invitation operations
  - MockRolesClient: Role management operations
  - MockRoleBindingsClient: Role binding operations
  - MockIdentitiesClient: Identity management operations
  - MockIdentityProvidersClient: Identity provider operations
  - MockGroupAccountAssociationsClient: Account association operations
  - MockSubscriptionsClient: Event subscription operations

# Request Matching

Mock responses are matched against incoming requests using protocol buffer
comparison. Configure multiple responses for different request patterns:

	mock := test.MockGroupsClient{
		OnCreate: []test.GroupOnCreate{{
			Given: &iam.CreateGroupRequest{
				Parent: "parent-id",
				Group:  &iam.Group{Name: "group-a"},
			},
			Created: &iam.Group{Id: "id-a", Name: "group-a"},
		}, {
			Given: &iam.CreateGroupRequest{
				Parent: "parent-id",
				Group:  &iam.Group{Name: "group-b"},
			},
			Created: &iam.Group{Id: "id-b", Name: "group-b"},
		}},
	}

# Error Simulation

Simulate error conditions by setting the Error field:

	mock := test.MockGroupsClient{
		OnDelete: []test.GroupOnDelete{{
			Given: &iam.DeleteGroupRequest{Id: "protected-group"},
			Error: status.Error(codes.PermissionDenied, "cannot delete protected group"),
		}},
	}

# Server Mocks

For integration tests requiring gRPC servers, use the server mock implementations:

	server := test.MockGroupsServer{
		Client: test.MockGroupsClient{
			OnList: []test.GroupOnList{{
				Given: &iam.GroupFilter{},
				List:  &iam.GroupList{Items: groups},
			}},
		},
	}

	// Register with a gRPC test server
	iam.RegisterGroupsServer(grpcServer, server)

# Integration with MockIAMClient

The MockIAMClient aggregates all mock clients and implements the iam.Clients
interface for complete IAM client mocking:

	mock := &test.MockIAMClient{
		GroupsClient:       groupsMock,
		RolesClient:        rolesMock,
		RoleBindingsClient: roleBindingsMock,
		OnClose:            nil, // No error on close
	}

	// Pass to code expecting iam.Clients
	service := NewService(mock)

# Thread Safety

Mock clients are safe for concurrent read access. Configure all expected
responses before using the mock in concurrent tests. Do not modify mock
configurations while tests are running.
*/
package test
