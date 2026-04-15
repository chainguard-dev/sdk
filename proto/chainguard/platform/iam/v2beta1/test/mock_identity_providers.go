/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ iam.IdentityProvidersServiceClient = (*MockIdentityProvidersServiceClient)(nil)

type MockIdentityProvidersServiceClient struct {
	iam.IdentityProvidersServiceClient
	T *testing.T

	OnGetIdentityProvider    []test.On[*iam.GetIdentityProviderRequest, *iam.IdentityProvider]
	OnCreateIdentityProvider []test.On[*iam.CreateIdentityProviderRequest, *iam.IdentityProvider]
	OnUpdateIdentityProvider []test.On[*iam.UpdateIdentityProviderRequest, *iam.IdentityProvider]
	OnDeleteIdentityProvider []test.On[*iam.DeleteIdentityProviderRequest, *emptypb.Empty]
	OnListIdentityProviders  []test.On[*iam.ListIdentityProvidersRequest, *iam.ListIdentityProvidersResponse]
}

func (m MockIdentityProvidersServiceClient) GetIdentityProvider(_ context.Context, given *iam.GetIdentityProviderRequest, _ ...grpc.CallOption) (*iam.IdentityProvider, error) {
	return test.Match(m.T, m.OnGetIdentityProvider, given, "get-identity-provider", protocmp.Transform())
}

func (m MockIdentityProvidersServiceClient) CreateIdentityProvider(_ context.Context, given *iam.CreateIdentityProviderRequest, _ ...grpc.CallOption) (*iam.IdentityProvider, error) {
	return test.Match(m.T, m.OnCreateIdentityProvider, given, "create-identity-provider", protocmp.Transform())
}

func (m MockIdentityProvidersServiceClient) UpdateIdentityProvider(_ context.Context, given *iam.UpdateIdentityProviderRequest, _ ...grpc.CallOption) (*iam.IdentityProvider, error) {
	return test.Match(m.T, m.OnUpdateIdentityProvider, given, "update-identity-provider", protocmp.Transform())
}

func (m MockIdentityProvidersServiceClient) DeleteIdentityProvider(_ context.Context, given *iam.DeleteIdentityProviderRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteIdentityProvider, given, "delete-identity-provider", protocmp.Transform())
}

func (m MockIdentityProvidersServiceClient) ListIdentityProviders(_ context.Context, given *iam.ListIdentityProvidersRequest, _ ...grpc.CallOption) (*iam.ListIdentityProvidersResponse, error) {
	return test.Match(m.T, m.OnListIdentityProviders, given, "list-identity-providers", protocmp.Transform())
}
