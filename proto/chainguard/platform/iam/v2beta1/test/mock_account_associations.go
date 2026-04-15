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

var _ iam.AccountAssociationsServiceClient = (*MockAccountAssociationsServiceClient)(nil)

type MockAccountAssociationsServiceClient struct {
	iam.AccountAssociationsServiceClient
	T *testing.T

	OnGetAccountAssociation    []test.On[*iam.GetAccountAssociationRequest, *iam.AccountAssociation]
	OnCreateAccountAssociation []test.On[*iam.CreateAccountAssociationRequest, *iam.AccountAssociation]
	OnUpdateAccountAssociation []test.On[*iam.UpdateAccountAssociationRequest, *iam.AccountAssociation]
	OnDeleteAccountAssociation []test.On[*iam.DeleteAccountAssociationRequest, *emptypb.Empty]
	OnListAccountAssociations  []test.On[*iam.ListAccountAssociationsRequest, *iam.ListAccountAssociationsResponse]
}

func (m MockAccountAssociationsServiceClient) GetAccountAssociation(_ context.Context, given *iam.GetAccountAssociationRequest, _ ...grpc.CallOption) (*iam.AccountAssociation, error) {
	return test.Match(m.T, m.OnGetAccountAssociation, given, "get-account-association", protocmp.Transform())
}

func (m MockAccountAssociationsServiceClient) CreateAccountAssociation(_ context.Context, given *iam.CreateAccountAssociationRequest, _ ...grpc.CallOption) (*iam.AccountAssociation, error) {
	return test.Match(m.T, m.OnCreateAccountAssociation, given, "create-account-association", protocmp.Transform())
}

func (m MockAccountAssociationsServiceClient) UpdateAccountAssociation(_ context.Context, given *iam.UpdateAccountAssociationRequest, _ ...grpc.CallOption) (*iam.AccountAssociation, error) {
	return test.Match(m.T, m.OnUpdateAccountAssociation, given, "update-account-association", protocmp.Transform())
}

func (m MockAccountAssociationsServiceClient) DeleteAccountAssociation(_ context.Context, given *iam.DeleteAccountAssociationRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteAccountAssociation, given, "delete-account-association", protocmp.Transform())
}

func (m MockAccountAssociationsServiceClient) ListAccountAssociations(_ context.Context, given *iam.ListAccountAssociationsRequest, _ ...grpc.CallOption) (*iam.ListAccountAssociationsResponse, error) {
	return test.Match(m.T, m.OnListAccountAssociations, given, "list-account-associations", protocmp.Transform())
}
