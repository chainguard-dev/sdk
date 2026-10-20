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

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ iam.AuthPoliciesServiceClient = (*MockAuthPoliciesServiceClient)(nil)

type MockAuthPoliciesServiceClient struct {
	iam.AuthPoliciesServiceClient
	T *testing.T

	OnGetAuthPolicy    []test.On[*iam.GetAuthPolicyRequest, *iam.AuthPolicy]
	OnUpdateAuthPolicy []test.On[*iam.UpdateAuthPolicyRequest, *iam.AuthPolicy]
}

func (m MockAuthPoliciesServiceClient) GetAuthPolicy(_ context.Context, given *iam.GetAuthPolicyRequest, _ ...grpc.CallOption) (*iam.AuthPolicy, error) {
	return test.Match(m.T, m.OnGetAuthPolicy, given, "get-auth-policy", protocmp.Transform())
}

func (m MockAuthPoliciesServiceClient) UpdateAuthPolicy(_ context.Context, given *iam.UpdateAuthPolicyRequest, _ ...grpc.CallOption) (*iam.AuthPolicy, error) {
	return test.Match(m.T, m.OnUpdateAuthPolicy, given, "update-auth-policy", protocmp.Transform())
}
