/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Generate the proto definitions
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true group.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true group_invites.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true role.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true role_binding.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true identity.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true account_associations.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true identity_providers.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true terms.platform.proto

/*
Package v1 provides gRPC client and server definitions for IAM interactions
with the Chainguard Console.

This package contains protocol buffer definitions and generated Go code for
managing identity and access management resources including groups, roles,
role bindings, identities, and identity providers.

# Overview

The IAM v1 package provides the following key features:

  - Unified client interface for all IAM services via the Clients interface
  - gRPC client implementations for each IAM resource type
  - Protocol buffer message types for requests, responses, and resources
  - gRPC-Gateway support for REST API access
  - Event types for IAM resource changes

# Basic Usage

Create an IAM client using NewClients with the IAM service URL and auth token:

	clients, err := v1.NewClients(ctx, "https://iam.example.com", token)
	if err != nil {
		return fmt.Errorf("failed to create IAM clients: %w", err)
	}
	defer clients.Close()

	// List groups
	groups, err := clients.Groups().List(ctx, &v1.GroupFilter{})
	if err != nil {
		return fmt.Errorf("failed to list groups: %w", err)
	}

	for _, group := range groups.Items {
		fmt.Printf("Group: %s (%s)\n", group.Name, group.Id)
	}

# Available Service Clients

The Clients interface provides access to the following service clients:

  - Groups: Manage organizational groups
  - GroupInvites: Handle group membership invitations
  - Roles: Define and manage roles with specific capabilities
  - RoleBindings: Assign roles to identities within groups
  - Identities: Manage user and service account identities
  - IdentityProviders: Configure external identity providers
  - AccountAssociations: Link groups to external accounts
  - Subscriptions: Subscribe to IAM event notifications

# Using an Existing Connection

If you already have a gRPC connection, use NewClientsFromConnection:

	conn, err := grpc.NewClient(target, opts...)
	if err != nil {
		return err
	}

	// Create IAM clients from existing connection
	// Note: The caller is responsible for closing the connection
	clients := v1.NewClientsFromConnection(conn)

	// Use the clients
	roles, err := clients.Roles().List(ctx, &v1.RoleFilter{})

# Working with Groups

Groups are the primary organizational unit for access control:

	// Create a new group
	group, err := clients.Groups().Create(ctx, &v1.CreateGroupRequest{
		Parent: parentGroupID,
		Group: &v1.Group{
			Name:        "engineering",
			Description: "Engineering team",
		},
	})

	// List groups with filtering
	groups, err := clients.Groups().List(ctx, &v1.GroupFilter{
		Parent: parentGroupID,
	})

# Managing Role Bindings

Role bindings associate identities with roles within a group:

	// Create a role binding
	binding, err := clients.RoleBindings().Create(ctx, &v1.CreateRoleBindingRequest{
		Parent: groupID,
		RoleBinding: &v1.RoleBinding{
			Identity: identityID,
			Role:     roleID,
		},
	})

	// List role bindings for a group
	bindings, err := clients.RoleBindings().List(ctx, &v1.RoleBindingFilter{
		Group: groupID,
	})

# Testing

For unit testing, use the mock implementations in the test subpackage:

	import "chainguard.dev/sdk/proto/platform/iam/v1/test"

	mock := &test.MockIAMClient{
		GroupsClient: test.MockGroupsClient{
			OnList: []test.GroupOnList{{
				Given: &v1.GroupFilter{},
				List: &v1.GroupList{
					Items: []*v1.Group{{
						Id:   "group-123",
						Name: "test-group",
					}},
				},
			}},
		},
	}

	// Use mock in place of real client
	service := NewService(mock)

# Thread Safety

All client methods are safe for concurrent use. The underlying gRPC connection
handles concurrent requests appropriately.
*/
package v1
