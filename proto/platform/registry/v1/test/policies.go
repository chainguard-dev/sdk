/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ registry.PoliciesClient = (*MockPoliciesClient)(nil)

type MockPoliciesClient struct {
	registry.PoliciesClient

	OnCheckPolicies []PoliciesOnCheck
}

type PoliciesOnCheck struct {
	Given    *registry.CheckPoliciesRequest
	Response *registry.CheckPoliciesResponse
	Error    error
}

func (m *MockPoliciesClient) CheckPolicies(_ context.Context, given *registry.CheckPoliciesRequest, _ ...grpc.CallOption) (*registry.CheckPoliciesResponse, error) {
	for _, o := range m.OnCheckPolicies {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
