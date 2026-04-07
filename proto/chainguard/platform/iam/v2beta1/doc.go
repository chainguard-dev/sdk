/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package v2beta1 provides Go clients for the Chainguard IAM v2beta1 API.

# Overview

This package contains generated protobuf types and gRPC clients for managing
Chainguard IAM resources including groups, identities, roles, role bindings,
and account associations. It provides a unified Clients interface for accessing
all IAM services with built-in pagination support.

# Features

  - Groups: Manage organizations and folders in the IAM hierarchy
  - Identities: Manage IAM identities with OIDC, static keys, or AWS IAM
  - Roles: Manage custom roles and query available roles and their capabilities
  - RoleBindings: Associate identities with roles within group scopes
  - AccountAssociations: Manage cloud provider account associations for groups
  - GroupInvites: Manage group invitations with role assignment and email notification
  - Pagination: Iterator-based and slice-based pagination helpers

# Services

GroupsService provides operations for managing IAM groups:
  - GetGroup: Retrieve a single group by UID
  - ListGroups: List groups with filtering and pagination
  - UpdateGroup: Update group properties

IdentitiesService provides operations for managing IAM identities:
  - GetIdentity: Retrieve a single identity by UID
  - ListIdentities: List identities with filtering and pagination

RolesService provides operations for managing IAM roles:
  - CreateRole: Create a new custom role
  - GetRole: Retrieve a single role by UID
  - ListRoles: List roles with filtering and pagination
  - UpdateRole: Update a custom role's properties
  - DeleteRole: Delete a custom role

RoleBindingsService provides operations for managing role bindings:
  - GetRoleBinding: Retrieve a single role binding by UID
  - CreateRoleBinding: Create a new role binding (identity + role + group)
  - DeleteRoleBinding: Delete a role binding by UID (idempotent)
  - ListRoleBindings: List role bindings with filtering and pagination

AccountAssociationsService provides operations for managing cloud provider associations:
  - GetAccountAssociation: Retrieve a single account association by group UID
  - CreateAccountAssociation: Create a new account association for a group
  - DeleteAccountAssociation: Delete an account association by group UID (idempotent)
  - ListAccountAssociations: List account associations with filtering and pagination

GroupInvitesService provides operations for managing group invitations:
  - CreateGroupInvite: Create a new invite with role assignment and optional email
  - GetGroupInvite: Retrieve a single group invite by UID
  - DeleteGroupInvite: Delete a group invite by UID (idempotent)
  - ListGroupInvites: List group invites with filtering and pagination

# Usage

Create a client from an existing gRPC connection:

	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)

Access individual service clients:

	groupsClient := clients.GroupsService()
	identitiesClient := clients.IdentitiesService()
	rolesClient := clients.RolesService()
	roleBindingsClient := clients.RoleBindingsService()

# Pagination

The package provides two pagination patterns for list operations.

Iterator-based pagination processes items one at a time:

	for group, err := range clients.ListGroupsIter(ctx, &v2beta1.ListGroupsRequest{}) {
		if err != nil {
			return err
		}
		fmt.Println(group.GetName())
	}

Slice-based pagination collects all results:

	groups, err := clients.ListGroupsAll(ctx, &v2beta1.ListGroupsRequest{})
	if err != nil {
		return err
	}
	for _, group := range groups {
		fmt.Println(group.GetName())
	}

# Thread Safety

All service clients are safe for concurrent use. The Clients interface methods
can be called from multiple goroutines simultaneously.
*/
package v2beta1

//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/group_invites.proto
//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/groups.proto
//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/identities.proto
//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/roles.proto
//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/role_bindings.proto
//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/account_associations.proto
//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/iam/v2beta1/identity_providers.proto
