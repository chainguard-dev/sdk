/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	"chainguard.dev/sdk/proto/platform/iam/v1/test"
)

// ExampleClients demonstrates the Clients interface methods.
func ExampleClients() {
	// Create a mock IAM client for demonstration
	// In production, use iam.NewClients() with a real IAM URL and token
	var clients iam.Clients = &test.MockIAMClient{
		GroupsClient: test.MockGroupsClient{
			OnList: []test.GroupOnList{{
				Given: &iam.GroupFilter{},
				List: &iam.GroupList{
					Items: []*iam.Group{{
						Id:   "group-123",
						Name: "example-group",
					}},
				},
			}},
		},
	}

	// Access individual service clients
	groupsClient := clients.Groups()
	rolesClient := clients.Roles()
	identitiesClient := clients.Identities()

	fmt.Printf("Groups client: %T\n", groupsClient)
	fmt.Printf("Roles client: %T\n", rolesClient)
	fmt.Printf("Identities client: %T\n", identitiesClient)

	// Always close the client when done
	if err := clients.Close(); err != nil {
		fmt.Printf("Close error: %v\n", err)
	}

	// Output:
	// Groups client: *test.MockGroupsClient
	// Roles client: *test.MockRolesClient
	// Identities client: *test.MockIdentitiesClient
}

// Example_groupOperations demonstrates working with groups.
func Example_groupOperations() {
	// Create a mock client with group data
	mock := &test.MockIAMClient{
		GroupsClient: test.MockGroupsClient{
			OnList: []test.GroupOnList{{
				Given: &iam.GroupFilter{},
				List: &iam.GroupList{
					Items: []*iam.Group{{
						Id:          "group-1",
						Name:        "engineering",
						Description: "Engineering team",
					}, {
						Id:          "group-2",
						Name:        "security",
						Description: "Security team",
					}},
				},
			}},
		},
	}

	// List all groups
	fmt.Printf("Available service clients:\n")
	fmt.Printf("  - Groups: %v\n", mock.Groups() != nil)
	fmt.Printf("  - Roles: %v\n", mock.Roles() != nil)
	fmt.Printf("  - RoleBindings: %v\n", mock.RoleBindings() != nil)
	fmt.Printf("  - Identities: %v\n", mock.Identities() != nil)
	fmt.Printf("  - IdentityProviders: %v\n", mock.IdentityProviders() != nil)
	fmt.Printf("  - GroupInvites: %v\n", mock.GroupInvites() != nil)
	fmt.Printf("  - AccountAssociations: %v\n", mock.AccountAssociations() != nil)
	fmt.Printf("  - Subscriptions: %v\n", mock.Subscriptions() != nil)

	// Output:
	// Available service clients:
	//   - Groups: true
	//   - Roles: true
	//   - RoleBindings: true
	//   - Identities: true
	//   - IdentityProviders: true
	//   - GroupInvites: true
	//   - AccountAssociations: true
	//   - Subscriptions: true
}

// Example_roleBindings demonstrates the role binding client accessor.
func Example_roleBindings() {
	// Create a mock client
	mock := &test.MockIAMClient{
		RoleBindingsClient: test.MockRoleBindingsClient{
			OnList: []test.RoleBindingOnList{{
				Given: &iam.RoleBindingFilter{},
				List: &iam.RoleBindingList{
					Items: []*iam.RoleBindingList_Binding{{
						Id:       "binding-1",
						Identity: "user@example.com",
						Role:     &iam.Role{Id: "roles/viewer"},
					}},
				},
			}},
		},
	}

	// Access the role bindings client
	rbClient := mock.RoleBindings()
	fmt.Printf("RoleBindings client type: %T\n", rbClient)

	// Output:
	// RoleBindings client type: *test.MockRoleBindingsClient
}
