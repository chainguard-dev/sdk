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

var _ iam.ScimUsersClient = (*MockScimUsersClient)(nil)

type MockScimUsersClient struct {
	OnList []ScimUsersOnList
}

type ScimUsersOnList struct {
	Given *iam.ScimUserFilter
	List  *iam.ScimUserList
	Error error
}

func (m MockScimUsersClient) List(_ context.Context, given *iam.ScimUserFilter, _ ...grpc.CallOption) (*iam.ScimUserList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return proto.Clone(o.List).(*iam.ScimUserList), o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
