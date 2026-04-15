/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v2beta1 "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
)

func ExampleNewClientsFromConnection() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	// Access individual service clients
	_ = clients.GroupsService()
	_ = clients.IdentitiesService()
	_ = clients.RolesService()
	_ = clients.RoleBindingsService()
}

func ExampleClients_GroupsService() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Get a specific group
	group, err := clients.GroupsService().GetGroup(ctx, &v2beta1.GetGroupRequest{
		Uid: "example-group-uid",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(group.GetName())
}

func ExampleClients_IdentitiesService() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Get a specific identity
	identity, err := clients.IdentitiesService().GetIdentity(ctx, &v2beta1.GetIdentityRequest{
		Uid: "example-identity-uid",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(identity.GetName())
}

func ExampleClients_RolesService() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Get a specific role
	role, err := clients.RolesService().GetRole(ctx, &v2beta1.GetRoleRequest{
		Uid: "example-role-uid",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(role.GetName())
}

func ExampleClients_RoleBindingsService() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Get a specific role binding
	roleBinding, err := clients.RoleBindingsService().GetRoleBinding(ctx, &v2beta1.GetRoleBindingRequest{
		Uid: "example-role-binding-uid",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(roleBinding.GetUid())
}

func ExampleClients_ListGroupsIter() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Iterate over groups with automatic pagination
	for group, err := range clients.ListGroupsIter(ctx, &v2beta1.ListGroupsRequest{
		PageSize: 50,
	}) {
		if err != nil {
			panic(err)
		}
		fmt.Println(group.GetName())
	}
}

func ExampleClients_ListGroupsAll() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Fetch all groups into a slice
	groups, err := clients.ListGroupsAll(ctx, &v2beta1.ListGroupsRequest{})
	if err != nil {
		panic(err)
	}
	for _, group := range groups {
		fmt.Println(group.GetName())
	}
}

func ExampleClients_ListIdentitiesIter() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Iterate over identities with automatic pagination
	for identity, err := range clients.ListIdentitiesIter(ctx, &v2beta1.ListIdentitiesRequest{
		PageSize: 50,
	}) {
		if err != nil {
			panic(err)
		}
		fmt.Println(identity.GetName())
	}
}

func ExampleClients_ListIdentitiesAll() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Fetch all identities into a slice
	identities, err := clients.ListIdentitiesAll(ctx, &v2beta1.ListIdentitiesRequest{})
	if err != nil {
		panic(err)
	}
	for _, identity := range identities {
		fmt.Println(identity.GetName())
	}
}

func ExampleClients_ListRolesIter() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Iterate over roles with automatic pagination
	for role, err := range clients.ListRolesIter(ctx, &v2beta1.ListRolesRequest{
		PageSize: 50,
	}) {
		if err != nil {
			panic(err)
		}
		fmt.Println(role.GetName())
	}
}

func ExampleClients_ListRolesAll() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Fetch all roles into a slice
	roles, err := clients.ListRolesAll(ctx, &v2beta1.ListRolesRequest{})
	if err != nil {
		panic(err)
	}
	for _, role := range roles {
		fmt.Println(role.GetName())
	}
}

func ExampleClients_ListRoleBindingsIter() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Iterate over role bindings with automatic pagination
	for roleBinding, err := range clients.ListRoleBindingsIter(ctx, &v2beta1.ListRoleBindingsRequest{
		PageSize: 50,
	}) {
		if err != nil {
			panic(err)
		}
		fmt.Println(roleBinding.GetUid())
	}
}

func ExampleClients_ListRoleBindingsAll() {
	conn, err := grpc.NewClient("api.chainguard.dev:443",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	clients := v2beta1.NewClientsFromConnection(conn)
	defer clients.Close()

	ctx := context.Background()

	// Fetch all role bindings into a slice
	roleBindings, err := clients.ListRoleBindingsAll(ctx, &v2beta1.ListRoleBindingsRequest{})
	if err != nil {
		panic(err)
	}
	for _, roleBinding := range roleBindings {
		fmt.Println(roleBinding.GetUid())
	}
}
