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

var _ iam.IdentitiesServiceClient = (*MockIdentitiesServiceClient)(nil)

type MockIdentitiesServiceClient struct {
	iam.IdentitiesServiceClient
	T *testing.T

	OnCreateIdentity []test.On[*iam.CreateIdentityRequest, *iam.Identity]
	OnGetIdentity    []test.On[*iam.GetIdentityRequest, *iam.Identity]
	OnUpdateIdentity []test.On[*iam.UpdateIdentityRequest, *iam.Identity]
	OnDeleteIdentity []test.On[*iam.DeleteIdentityRequest, *emptypb.Empty]
	OnListIdentities []test.On[*iam.ListIdentitiesRequest, *iam.ListIdentitiesResponse]
}

func (m MockIdentitiesServiceClient) CreateIdentity(_ context.Context, given *iam.CreateIdentityRequest, _ ...grpc.CallOption) (*iam.Identity, error) {
	return test.Match(m.T, m.OnCreateIdentity, given, "create-identity", protocmp.Transform())
}

func (m MockIdentitiesServiceClient) GetIdentity(_ context.Context, given *iam.GetIdentityRequest, _ ...grpc.CallOption) (*iam.Identity, error) {
	return test.Match(m.T, m.OnGetIdentity, given, "get-identity", protocmp.Transform())
}

func (m MockIdentitiesServiceClient) UpdateIdentity(_ context.Context, given *iam.UpdateIdentityRequest, _ ...grpc.CallOption) (*iam.Identity, error) {
	return test.Match(m.T, m.OnUpdateIdentity, given, "update-identity", protocmp.Transform())
}

func (m MockIdentitiesServiceClient) DeleteIdentity(_ context.Context, given *iam.DeleteIdentityRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteIdentity, given, "delete-identity", protocmp.Transform())
}

func (m MockIdentitiesServiceClient) ListIdentities(_ context.Context, given *iam.ListIdentitiesRequest, _ ...grpc.CallOption) (*iam.ListIdentitiesResponse, error) {
	return test.Match(m.T, m.OnListIdentities, given, "list-identities", protocmp.Transform())
}
