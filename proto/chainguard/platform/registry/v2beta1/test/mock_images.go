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

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ registry.ImagesServiceClient = (*MockImagesServiceClient)(nil)

type MockImagesServiceClient struct {
	registry.ImagesServiceClient
	T *testing.T

	OnGetArchitectures []test.On[*registry.GetArchitecturesRequest, *registry.Architectures]
	OnGetSize          []test.On[*registry.GetSizeRequest, *registry.Size]
}

func (m MockImagesServiceClient) GetArchitectures(_ context.Context, given *registry.GetArchitecturesRequest, _ ...grpc.CallOption) (*registry.Architectures, error) {
	return test.Match(m.T, m.OnGetArchitectures, given, "get-architectures", protocmp.Transform())
}

func (m MockImagesServiceClient) GetSize(_ context.Context, given *registry.GetSizeRequest, _ ...grpc.CallOption) (*registry.Size, error) {
	return test.Match(m.T, m.OnGetSize, given, "get-size", protocmp.Transform())
}
