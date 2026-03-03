/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	"chainguard.dev/sdk/proto/platform/iam/v1/test"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ExampleMockIAMClient demonstrates creating a complete mock IAM client.
func ExampleMockIAMClient() {
	// Create a mock IAM client with configured responses
	mock := &test.MockIAMClient{
		GroupsClient: test.MockGroupsClient{
			OnList: []test.GroupOnList{{
				Given: &iam.GroupFilter{},
				List: &iam.GroupList{
					Items: []*iam.Group{{
						Id:   "group-123",
						Name: "test-group",
					}},
				},
			}},
		},
	}

	// Use the mock client
	groups, err := mock.Groups().List(context.Background(), &iam.GroupFilter{})
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d groups\n", len(groups.Items))
	fmt.Printf("First group: %s\n", groups.Items[0].Name)

	// Output:
	// Found 1 groups
	// First group: test-group
}

// ExampleMockGroupsClient demonstrates mocking group operations.
func ExampleMockGroupsClient() {
	// Configure mock responses for group operations
	mock := test.MockGroupsClient{
		OnCreate: []test.GroupOnCreate{{
			Given: &iam.CreateGroupRequest{
				Parent: "parent-id",
				Group:  &iam.Group{Name: "new-group"},
			},
			Created: &iam.Group{
				Id:   "created-id",
				Name: "new-group",
			},
		}},
		OnList: []test.GroupOnList{{
			Given: &iam.GroupFilter{},
			List: &iam.GroupList{
				Items: []*iam.Group{{
					Id:   "group-1",
					Name: "existing-group",
				}},
			},
		}},
	}

	ctx := context.Background()

	// Test create operation
	created, err := mock.Create(ctx, &iam.CreateGroupRequest{
		Parent: "parent-id",
		Group:  &iam.Group{Name: "new-group"},
	})
	if err != nil {
		fmt.Printf("Create error: %v\n", err)
		return
	}
	fmt.Printf("Created group: %s (ID: %s)\n", created.Name, created.Id)

	// Test list operation
	list, err := mock.List(ctx, &iam.GroupFilter{})
	if err != nil {
		fmt.Printf("List error: %v\n", err)
		return
	}
	fmt.Printf("Listed %d groups\n", len(list.Items))

	// Output:
	// Created group: new-group (ID: created-id)
	// Listed 1 groups
}

// ExampleMockGroupsClient_errorSimulation demonstrates simulating errors.
func ExampleMockGroupsClient_errorSimulation() {
	// Configure mock to return an error
	mock := test.MockGroupsClient{
		OnDelete: []test.GroupOnDelete{{
			Given: &iam.DeleteGroupRequest{Id: "protected-group"},
			Error: status.Error(codes.PermissionDenied, "cannot delete protected group"),
		}},
	}

	ctx := context.Background()

	// Attempt to delete a protected group
	_, err := mock.Delete(ctx, &iam.DeleteGroupRequest{Id: "protected-group"})
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			fmt.Printf("Error code: %s\n", st.Code())
			fmt.Printf("Error message: %s\n", st.Message())
		}
	}

	// Output:
	// Error code: PermissionDenied
	// Error message: cannot delete protected group
}

// ExampleMockIAMClient_close demonstrates handling the Close method.
func ExampleMockIAMClient_close() {
	// Configure mock with no close error
	mock := &test.MockIAMClient{
		OnClose: nil,
	}

	// Close the mock client
	err := mock.Close()
	if err != nil {
		fmt.Printf("Close error: %v\n", err)
		return
	}

	fmt.Println("Client closed successfully")

	// Output:
	// Client closed successfully
}
