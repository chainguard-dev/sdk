/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	events "chainguard.dev/sdk/proto/platform/events/v1"
	iam "chainguard.dev/sdk/proto/platform/iam/v1"
)

type MockIAMClient struct {
	OnClose error

	GroupsClient                   MockGroupsClient
	GroupInvitesClient             MockGroupInvitesClient
	RolesClient                    MockRolesClient
	RoleBindingsClient             MockRoleBindingsClient
	IdentitiesClient               MockIdentitiesClient
	DeprecatedIdentitiesClient     MockDeprecatedIdentitiesClient
	IdentityProvidersClient        MockIdentityProvidersClient
	GroupAccountAssociationsClient MockGroupAccountAssociationsClient
	SubscriptionsClient            MockSubscriptionsClient
}

var _ iam.Clients = (*MockIAMClient)(nil)

func (m MockIAMClient) Groups() iam.GroupsClient {
	return &m.GroupsClient
}

func (m MockIAMClient) GroupInvites() iam.GroupInvitesClient {
	return &m.GroupInvitesClient
}

func (m MockIAMClient) Roles() iam.RolesClient {
	return &m.RolesClient
}

func (m MockIAMClient) RoleBindings() iam.RoleBindingsClient {
	return &m.RoleBindingsClient
}

func (m MockIAMClient) Identities() iam.IdentitiesClient {
	return &m.IdentitiesClient
}

func (m MockIAMClient) DeprecatedIdentities() events.IdentitiesClient {
	return &m.DeprecatedIdentitiesClient
}

func (m MockIAMClient) IdentityProviders() iam.IdentityProvidersClient {
	return &m.IdentityProvidersClient
}

func (m MockIAMClient) AccountAssociations() iam.GroupAccountAssociationsClient {
	return &m.GroupAccountAssociationsClient
}

func (m MockIAMClient) Subscriptions() events.SubscriptionsClient {
	return &m.SubscriptionsClient
}

func (m MockIAMClient) Close() error {
	return m.OnClose
}
