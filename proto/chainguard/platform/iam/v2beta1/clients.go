/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package v2beta1 provides clients for IAM v2beta1 services.
package v2beta1

import (
	"context"
	"iter"

	"google.golang.org/grpc"

	v2iter "chainguard.dev/sdk/proto/chainguard/platform/iter"
)

// Clients provides access to v2beta1 IAM service clients.
type Clients interface {
	AccountAssociationsService() AccountAssociationsServiceClient
	GroupInvitesService() GroupInvitesServiceClient
	GroupsService() GroupsServiceClient
	IdentitiesService() IdentitiesServiceClient
	IdentityProvidersService() IdentityProvidersServiceClient
	RolesService() RolesServiceClient
	RoleBindingsService() RoleBindingsServiceClient

	// Iterator methods for pagination - GroupInvites
	ListGroupInvitesIter(ctx context.Context, req *ListGroupInvitesRequest) iter.Seq2[*GroupInvite, error]
	ListGroupInvitesAll(ctx context.Context, req *ListGroupInvitesRequest) ([]*GroupInvite, error)

	// Iterator methods for pagination - Groups
	ListGroupsIter(ctx context.Context, req *ListGroupsRequest) iter.Seq2[*Group, error]
	ListGroupsAll(ctx context.Context, req *ListGroupsRequest) ([]*Group, error)

	// Iterator methods for pagination - Identities
	ListIdentitiesIter(ctx context.Context, req *ListIdentitiesRequest) iter.Seq2[*Identity, error]
	ListIdentitiesAll(ctx context.Context, req *ListIdentitiesRequest) ([]*Identity, error)

	// Iterator methods for pagination - Identity Providers
	ListIdentityProvidersIter(ctx context.Context, req *ListIdentityProvidersRequest) iter.Seq2[*IdentityProvider, error]
	ListIdentityProvidersAll(ctx context.Context, req *ListIdentityProvidersRequest) ([]*IdentityProvider, error)

	// Iterator methods for pagination - Roles
	ListRolesIter(ctx context.Context, req *ListRolesRequest) iter.Seq2[*Role, error]
	ListRolesAll(ctx context.Context, req *ListRolesRequest) ([]*Role, error)

	// Iterator methods for pagination - RoleBindings
	ListRoleBindingsIter(ctx context.Context, req *ListRoleBindingsRequest) iter.Seq2[*RoleBinding, error]
	ListRoleBindingsAll(ctx context.Context, req *ListRoleBindingsRequest) ([]*RoleBinding, error)

	// Iterator methods for pagination - AccountAssociations
	ListAccountAssociationsIter(ctx context.Context, req *ListAccountAssociationsRequest) iter.Seq2[*AccountAssociation, error]
	ListAccountAssociationsAll(ctx context.Context, req *ListAccountAssociationsRequest) ([]*AccountAssociation, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 IAM clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		accountAssociationsService: NewAccountAssociationsServiceClient(conn),
		groupInvitesService:        NewGroupInvitesServiceClient(conn),
		groupsService:              NewGroupsServiceClient(conn),
		identitiesService:          NewIdentitiesServiceClient(conn),
		identityProviderService:    NewIdentityProvidersServiceClient(conn),
		rolesService:               NewRolesServiceClient(conn),
		roleBindingsService:        NewRoleBindingsServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	accountAssociationsService AccountAssociationsServiceClient
	groupInvitesService        GroupInvitesServiceClient
	groupsService              GroupsServiceClient
	identitiesService          IdentitiesServiceClient
	identityProviderService    IdentityProvidersServiceClient
	rolesService               RolesServiceClient
	roleBindingsService        RoleBindingsServiceClient

	conn *grpc.ClientConn
}

var _ Clients = (*clients)(nil)

func (c *clients) AccountAssociationsService() AccountAssociationsServiceClient {
	return c.accountAssociationsService
}

func (c *clients) GroupInvitesService() GroupInvitesServiceClient {
	return c.groupInvitesService
}

func (c *clients) GroupsService() GroupsServiceClient {
	return c.groupsService
}

func (c *clients) IdentitiesService() IdentitiesServiceClient {
	return c.identitiesService
}

func (c *clients) IdentityProvidersService() IdentityProvidersServiceClient {
	return c.identityProviderService
}

func (c *clients) RolesService() RolesServiceClient {
	return c.rolesService
}

func (c *clients) RoleBindingsService() RoleBindingsServiceClient {
	return c.roleBindingsService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListGroupInvitesIter returns an iterator over group invites matching the request.
func (c *clients) ListGroupInvitesIter(ctx context.Context, req *ListGroupInvitesRequest) iter.Seq2[*GroupInvite, error] {
	return v2iter.Paginate(ctx, req, "group_invites", func(ctx context.Context, r *ListGroupInvitesRequest) ([]*GroupInvite, string, error) {
		resp, err := c.GroupInvitesService().ListGroupInvites(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetGroupInvites(), resp.GetNextPageToken(), nil
	})
}

// ListGroupInvitesAll fetches all group invites matching the request by automatically handling pagination.
// For large result sets, consider using ListGroupInvitesIter directly to process items incrementally.
func (c *clients) ListGroupInvitesAll(ctx context.Context, req *ListGroupInvitesRequest) ([]*GroupInvite, error) {
	return v2iter.All(c.ListGroupInvitesIter(ctx, req))
}

// ListGroupsIter returns an iterator over groups matching the request.
func (c *clients) ListGroupsIter(ctx context.Context, req *ListGroupsRequest) iter.Seq2[*Group, error] {
	return v2iter.Paginate(ctx, req, "groups", func(ctx context.Context, r *ListGroupsRequest) ([]*Group, string, error) {
		resp, err := c.GroupsService().ListGroups(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetGroups(), resp.GetNextPageToken(), nil
	})
}

// ListGroupsAll fetches all groups matching the request by automatically handling pagination.
// For large result sets, consider using ListGroupsIter directly to process items incrementally.
func (c *clients) ListGroupsAll(ctx context.Context, req *ListGroupsRequest) ([]*Group, error) {
	return v2iter.All(c.ListGroupsIter(ctx, req))
}

// ListIdentitiesIter returns an iterator over identities matching the request.
func (c *clients) ListIdentitiesIter(ctx context.Context, req *ListIdentitiesRequest) iter.Seq2[*Identity, error] {
	return v2iter.Paginate(ctx, req, "identities", func(ctx context.Context, r *ListIdentitiesRequest) ([]*Identity, string, error) {
		resp, err := c.IdentitiesService().ListIdentities(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetIdentities(), resp.GetNextPageToken(), nil
	})
}

// ListIdentitiesAll fetches all identities matching the request by automatically handling pagination.
// For large result sets, consider using ListIdentitiesIter directly to process items incrementally.
func (c *clients) ListIdentitiesAll(ctx context.Context, req *ListIdentitiesRequest) ([]*Identity, error) {
	return v2iter.All(c.ListIdentitiesIter(ctx, req))
}

// ListIdentityProvidersIter returns an iterator over identity providers matching the request.
func (c *clients) ListIdentityProvidersIter(ctx context.Context, req *ListIdentityProvidersRequest) iter.Seq2[*IdentityProvider, error] {
	return v2iter.Paginate(ctx, req, "identity_providers", func(ctx context.Context, r *ListIdentityProvidersRequest) ([]*IdentityProvider, string, error) {
		resp, err := c.IdentityProvidersService().ListIdentityProviders(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetIdentityProviders(), resp.GetNextPageToken(), nil
	})
}

// ListIdentityProvidersAll fetches all identity providers matching the request by automatically handling pagination.
// For large result sets, consider using ListIdentityProvidersIter directly to process items incrementally.
func (c *clients) ListIdentityProvidersAll(ctx context.Context, req *ListIdentityProvidersRequest) ([]*IdentityProvider, error) {
	return v2iter.All(c.ListIdentityProvidersIter(ctx, req))
}

// ListRolesIter returns an iterator over roles matching the request.
func (c *clients) ListRolesIter(ctx context.Context, req *ListRolesRequest) iter.Seq2[*Role, error] {
	return v2iter.Paginate(ctx, req, "roles", func(ctx context.Context, r *ListRolesRequest) ([]*Role, string, error) {
		resp, err := c.RolesService().ListRoles(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetRoles(), resp.GetNextPageToken(), nil
	})
}

// ListRolesAll fetches all roles matching the request by automatically handling pagination.
// For large result sets, consider using ListRolesIter directly to process items incrementally.
func (c *clients) ListRolesAll(ctx context.Context, req *ListRolesRequest) ([]*Role, error) {
	return v2iter.All(c.ListRolesIter(ctx, req))
}

// ListRoleBindingsIter returns an iterator over role bindings matching the request.
func (c *clients) ListRoleBindingsIter(ctx context.Context, req *ListRoleBindingsRequest) iter.Seq2[*RoleBinding, error] {
	return v2iter.Paginate(ctx, req, "role_bindings", func(ctx context.Context, r *ListRoleBindingsRequest) ([]*RoleBinding, string, error) {
		resp, err := c.RoleBindingsService().ListRoleBindings(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetRoleBindings(), resp.GetNextPageToken(), nil
	})
}

// ListRoleBindingsAll fetches all role bindings matching the request by automatically handling pagination.
// For large result sets, consider using ListRoleBindingsIter directly to process items incrementally.
func (c *clients) ListRoleBindingsAll(ctx context.Context, req *ListRoleBindingsRequest) ([]*RoleBinding, error) {
	return v2iter.All(c.ListRoleBindingsIter(ctx, req))
}

// ListAccountAssociationsIter returns an iterator over account associations matching the request.
func (c *clients) ListAccountAssociationsIter(ctx context.Context, req *ListAccountAssociationsRequest) iter.Seq2[*AccountAssociation, error] {
	return v2iter.Paginate(ctx, req, "account_associations", func(ctx context.Context, r *ListAccountAssociationsRequest) ([]*AccountAssociation, string, error) {
		resp, err := c.AccountAssociationsService().ListAccountAssociations(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetAccountAssociations(), resp.GetNextPageToken(), nil
	})
}

// ListAccountAssociationsAll fetches all account associations matching the request by automatically handling pagination.
// For large result sets, consider using ListAccountAssociationsIter directly to process items incrementally.
func (c *clients) ListAccountAssociationsAll(ctx context.Context, req *ListAccountAssociationsRequest) ([]*AccountAssociation, error) {
	return v2iter.All(c.ListAccountAssociationsIter(ctx, req))
}
