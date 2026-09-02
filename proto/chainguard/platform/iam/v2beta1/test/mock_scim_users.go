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

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ iam.ScimUsersServiceClient = (*MockScimUsersServiceClient)(nil)

type MockScimUsersServiceClient struct {
	iam.ScimUsersServiceClient
	T *testing.T

	OnListScimUsers []test.On[*iam.ListScimUsersRequest, *iam.ListScimUsersResponse]
}

func (m MockScimUsersServiceClient) ListScimUsers(_ context.Context, given *iam.ListScimUsersRequest, _ ...grpc.CallOption) (*iam.ListScimUsersResponse, error) {
	return test.Match(m.T, m.OnListScimUsers, given, "list-scim-users", protocmp.Transform())
}
