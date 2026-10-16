/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
)

var _ iam.ScimGroupsClient = (*MockScimGroupsClient)(nil)

type MockScimGroupsClient struct {
	OnList []ScimGroupsOnList
}

type ScimGroupsOnList struct {
	Given *iam.ScimGroupFilter
	List  *iam.ScimGroupList
	Error error
}

func (m MockScimGroupsClient) List(_ context.Context, given *iam.ScimGroupFilter, _ ...grpc.CallOption) (*iam.ScimGroupList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return proto.Clone(o.List).(*iam.ScimGroupList), o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
