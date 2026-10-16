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

var _ iam.ScimGroupsServiceClient = (*MockScimGroupsServiceClient)(nil)

type MockScimGroupsServiceClient struct {
	iam.ScimGroupsServiceClient
	T *testing.T

	OnListScimGroups []test.On[*iam.ListScimGroupsRequest, *iam.ListScimGroupsResponse]
}

func (m MockScimGroupsServiceClient) ListScimGroups(_ context.Context, given *iam.ListScimGroupsRequest, _ ...grpc.CallOption) (*iam.ListScimGroupsResponse, error) {
	return test.Match(m.T, m.OnListScimGroups, given, "list-scim-groups", protocmp.Transform())
}
