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

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ registry.OverlayBindingsServiceClient = (*MockOverlayBindingsServiceClient)(nil)

type MockOverlayBindingsServiceClient struct {
	registry.OverlayBindingsServiceClient
	T *testing.T

	OnGetOverlayBinding    []test.On[*registry.GetOverlayBindingRequest, *registry.OverlayBinding]
	OnListOverlayBindings  []test.On[*registry.ListOverlayBindingsRequest, *registry.ListOverlayBindingsResponse]
	OnCreateOverlayBinding []test.On[*registry.CreateOverlayBindingRequest, *registry.OverlayBinding]
	OnDeleteOverlayBinding []test.On[*registry.DeleteOverlayBindingRequest, *emptypb.Empty]
}

func (m MockOverlayBindingsServiceClient) GetOverlayBinding(_ context.Context, given *registry.GetOverlayBindingRequest, _ ...grpc.CallOption) (*registry.OverlayBinding, error) {
	return test.Match(m.T, m.OnGetOverlayBinding, given, "get-overlay-binding", protocmp.Transform())
}

func (m MockOverlayBindingsServiceClient) ListOverlayBindings(_ context.Context, given *registry.ListOverlayBindingsRequest, _ ...grpc.CallOption) (*registry.ListOverlayBindingsResponse, error) {
	return test.Match(m.T, m.OnListOverlayBindings, given, "list-overlay-bindings", protocmp.Transform())
}

func (m MockOverlayBindingsServiceClient) CreateOverlayBinding(_ context.Context, given *registry.CreateOverlayBindingRequest, _ ...grpc.CallOption) (*registry.OverlayBinding, error) {
	return test.Match(m.T, m.OnCreateOverlayBinding, given, "create-overlay-binding", protocmp.Transform())
}

func (m MockOverlayBindingsServiceClient) DeleteOverlayBinding(_ context.Context, given *registry.DeleteOverlayBindingRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteOverlayBinding, given, "delete-overlay-binding", protocmp.Transform())
}
