/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	events "chainguard.dev/sdk/proto/platform/events/v1"
)

var _ events.IdentitiesClient = (*MockDeprecatedIdentitiesClient)(nil)

type MockDeprecatedIdentitiesClient struct {
	OnCreate []DeprecatedIdentityOnCreate
}

type DeprecatedIdentityOnCreate struct {
	Given   *events.Identity
	Created *events.Identity
	Error   error
}

func (m MockDeprecatedIdentitiesClient) Create(_ context.Context, given *events.Identity, _ ...grpc.CallOption) (*events.Identity, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockDeprecatedIdentitiesClient) UpdateMetadata(ctx context.Context, in *events.IdentityMetadata, opts ...grpc.CallOption) (*events.IdentityMetadata, error) { //nolint: revive
	return in, nil
}
