/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"fmt"
	"net/url"
	"time"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"google.golang.org/grpc"
	"knative.dev/pkg/logging"

	"chainguard.dev/sdk/auth"
	events "chainguard.dev/sdk/proto/platform/events/v1"
)

type Clients interface {
	Groups() GroupsClient
	GroupInvites() GroupInvitesClient
	Roles() RolesClient
	RoleBindings() RoleBindingsClient

	Identities() IdentitiesClient
	DeprecatedIdentities() events.IdentitiesClient
	IdentityProviders() IdentityProvidersClient

	AccountAssociations() GroupAccountAssociationsClient

	Subscriptions() events.SubscriptionsClient

	Policies() PoliciesClient

	Sigstore() SigstoreServiceClient

	Close() error
}

func NewClients(ctx context.Context, iamURL string, token string) (Clients, error) {
	iamURI, err := url.Parse(iamURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse iam service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*iamURI)

	// TODO: we may want to require transport security at some future point.
	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		logging.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}

	var cancel context.CancelFunc
	if _, timeoutSet := ctx.Deadline(); !timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
	}
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the iam server: %w", err)
	}

	return &clients{
		group:                NewGroupsClient(conn),
		groupInvite:          NewGroupInvitesClient(conn),
		role:                 NewRolesClient(conn),
		roleBinding:          NewRoleBindingsClient(conn),
		identities:           NewIdentitiesClient(conn),
		deprecatedIdentities: events.NewIdentitiesClient(conn),
		identityProviders:    NewIdentityProvidersClient(conn),

		accountAssociations: NewGroupAccountAssociationsClient(conn),

		subscription: events.NewSubscriptionsClient(conn),

		policy: NewPoliciesClient(conn),

		sigstore: NewSigstoreServiceClient(conn),

		conn: conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		group:                NewGroupsClient(conn),
		groupInvite:          NewGroupInvitesClient(conn),
		role:                 NewRolesClient(conn),
		roleBinding:          NewRoleBindingsClient(conn),
		identities:           NewIdentitiesClient(conn),
		deprecatedIdentities: events.NewIdentitiesClient(conn),
		identityProviders:    NewIdentityProvidersClient(conn),

		accountAssociations: NewGroupAccountAssociationsClient(conn),

		subscription: events.NewSubscriptionsClient(conn),

		policy:   NewPoliciesClient(conn),
		sigstore: NewSigstoreServiceClient(conn),

		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	group                GroupsClient
	groupInvite          GroupInvitesClient
	role                 RolesClient
	roleBinding          RoleBindingsClient
	identities           IdentitiesClient
	deprecatedIdentities events.IdentitiesClient
	identityProviders    IdentityProvidersClient

	accountAssociations GroupAccountAssociationsClient

	subscription events.SubscriptionsClient

	policy PoliciesClient

	sigstore SigstoreServiceClient

	conn *grpc.ClientConn
}

func (c *clients) Groups() GroupsClient {
	return c.group
}

func (c *clients) GroupInvites() GroupInvitesClient {
	return c.groupInvite
}

func (c *clients) Roles() RolesClient {
	return c.role
}

func (c *clients) RoleBindings() RoleBindingsClient {
	return c.roleBinding
}

func (c *clients) Identities() IdentitiesClient {
	return c.identities
}

func (c *clients) DeprecatedIdentities() events.IdentitiesClient {
	return c.deprecatedIdentities
}

func (c *clients) IdentityProviders() IdentityProvidersClient {
	return c.identityProviders
}

func (c *clients) AccountAssociations() GroupAccountAssociationsClient {
	return c.accountAssociations
}

func (c *clients) Subscriptions() events.SubscriptionsClient {
	return c.subscription
}

func (c *clients) Policies() PoliciesClient {
	return c.policy
}

func (c *clients) Sigstore() SigstoreServiceClient {
	return c.sigstore
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
