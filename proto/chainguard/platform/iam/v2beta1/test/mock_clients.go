/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"iter"

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

type MockClients struct {
	OnClose error

	AccountAssociationsServiceClient MockAccountAssociationsServiceClient
	GroupInvitesServiceClient        MockGroupInvitesServiceClient
	GroupsServiceClient              MockGroupsServiceClient
	IdentitiesServiceClient          MockIdentitiesServiceClient
	IdentityProvidersServiceClient   MockIdentityProvidersServiceClient
	RolesServiceClient               MockRolesServiceClient
	RoleBindingsServiceClient        MockRoleBindingsServiceClient
}

// Close implements [v2beta1.Clients].
func (m *MockClients) Close() error {
	return m.OnClose
}

// AccountAssociationsService implements [v2beta1.Clients].
func (m *MockClients) AccountAssociationsService() iam.AccountAssociationsServiceClient {
	return &m.AccountAssociationsServiceClient
}

// GroupInvitesService implements [v2beta1.Clients].
func (m *MockClients) GroupInvitesService() iam.GroupInvitesServiceClient {
	return &m.GroupInvitesServiceClient
}

// GroupsService implements [v2beta1.Clients].
func (m *MockClients) GroupsService() iam.GroupsServiceClient {
	return &m.GroupsServiceClient
}

// IdentitiesService implements [v2beta1.Clients].
func (m *MockClients) IdentitiesService() iam.IdentitiesServiceClient {
	return &m.IdentitiesServiceClient
}

// IdentityProvidersService implements [v2beta1.Clients].
func (m *MockClients) IdentityProvidersService() iam.IdentityProvidersServiceClient {
	return &m.IdentityProvidersServiceClient
}

// RolesService implements [v2beta1.Clients].
func (m *MockClients) RolesService() iam.RolesServiceClient {
	return &m.RolesServiceClient
}

// RoleBindingsService implements [v2beta1.Clients].
func (m *MockClients) RoleBindingsService() iam.RoleBindingsServiceClient {
	return &m.RoleBindingsServiceClient
}

// ListAccountAssociationsAll implements [v2beta1.Clients].
func (m *MockClients) ListAccountAssociationsAll(ctx context.Context, req *iam.ListAccountAssociationsRequest) ([]*iam.AccountAssociation, error) {
	resp, err := m.AccountAssociationsServiceClient.ListAccountAssociations(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetAccountAssociations(), nil
}

// ListAccountAssociationsIter implements [v2beta1.Clients].
func (m *MockClients) ListAccountAssociationsIter(ctx context.Context, req *iam.ListAccountAssociationsRequest) iter.Seq2[*iam.AccountAssociation, error] {
	return test.MockIter(m.ListAccountAssociationsAll(ctx, req))
}

// ListGroupInvitesAll implements [v2beta1.Clients].
func (m *MockClients) ListGroupInvitesAll(ctx context.Context, req *iam.ListGroupInvitesRequest) ([]*iam.GroupInvite, error) {
	resp, err := m.GroupInvitesServiceClient.ListGroupInvites(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetGroupInvites(), nil
}

// ListGroupInvitesIter implements [v2beta1.Clients].
func (m *MockClients) ListGroupInvitesIter(ctx context.Context, req *iam.ListGroupInvitesRequest) iter.Seq2[*iam.GroupInvite, error] {
	return test.MockIter(m.ListGroupInvitesAll(ctx, req))
}

// ListGroupsAll implements [v2beta1.Clients].
func (m *MockClients) ListGroupsAll(ctx context.Context, req *iam.ListGroupsRequest) ([]*iam.Group, error) {
	resp, err := m.GroupsServiceClient.ListGroups(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetGroups(), nil
}

// ListGroupsIter implements [v2beta1.Clients].
func (m *MockClients) ListGroupsIter(ctx context.Context, req *iam.ListGroupsRequest) iter.Seq2[*iam.Group, error] {
	return test.MockIter(m.ListGroupsAll(ctx, req))
}

// ListIdentitiesAll implements [v2beta1.Clients].
func (m *MockClients) ListIdentitiesAll(ctx context.Context, req *iam.ListIdentitiesRequest) ([]*iam.Identity, error) {
	resp, err := m.IdentitiesServiceClient.ListIdentities(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetIdentities(), nil
}

// ListIdentitiesIter implements [v2beta1.Clients].
func (m *MockClients) ListIdentitiesIter(ctx context.Context, req *iam.ListIdentitiesRequest) iter.Seq2[*iam.Identity, error] {
	return test.MockIter(m.ListIdentitiesAll(ctx, req))
}

// ListIdentityProvidersAll implements [v2beta1.Clients].
func (m *MockClients) ListIdentityProvidersAll(ctx context.Context, req *iam.ListIdentityProvidersRequest) ([]*iam.IdentityProvider, error) {
	resp, err := m.IdentityProvidersServiceClient.ListIdentityProviders(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetIdentityProviders(), nil
}

// ListIdentityProvidersIter implements [v2beta1.Clients].
func (m *MockClients) ListIdentityProvidersIter(ctx context.Context, req *iam.ListIdentityProvidersRequest) iter.Seq2[*iam.IdentityProvider, error] {
	return test.MockIter(m.ListIdentityProvidersAll(ctx, req))
}

// ListRoleBindingsAll implements [v2beta1.Clients].
func (m *MockClients) ListRoleBindingsAll(ctx context.Context, req *iam.ListRoleBindingsRequest) ([]*iam.RoleBinding, error) {
	resp, err := m.RoleBindingsServiceClient.ListRoleBindings(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetRoleBindings(), nil
}

// ListRoleBindingsIter implements [v2beta1.Clients].
func (m *MockClients) ListRoleBindingsIter(ctx context.Context, req *iam.ListRoleBindingsRequest) iter.Seq2[*iam.RoleBinding, error] {
	return test.MockIter(m.ListRoleBindingsAll(ctx, req))
}

// ListRolesAll implements [v2beta1.Clients].
func (m *MockClients) ListRolesAll(ctx context.Context, req *iam.ListRolesRequest) ([]*iam.Role, error) {
	resp, err := m.RolesServiceClient.ListRoles(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetRoles(), nil
}

// ListRolesIter implements [v2beta1.Clients].
func (m *MockClients) ListRolesIter(ctx context.Context, req *iam.ListRolesRequest) iter.Seq2[*iam.Role, error] {
	return test.MockIter(m.ListRolesAll(ctx, req))
}

var _ iam.Clients = (*MockClients)(nil)
