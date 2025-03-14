/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ registry.EntitlementsClient = (*MockEntitlementsClient)(nil)

type MockEntitlementsClient struct{}

func (*MockEntitlementsClient) ListEntitlements(context.Context, *registry.EntitlementFilter, ...grpc.CallOption) (*registry.EntitlementList, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (*MockEntitlementsClient) ListEntitlementImages(context.Context, *registry.EntitlementImagesFilter, ...grpc.CallOption) (*registry.EntitlementImagesList, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
