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

var _ registry.OverlaysServiceClient = (*MockOverlaysServiceClient)(nil)

type MockOverlaysServiceClient struct {
	registry.OverlaysServiceClient
	T *testing.T

	OnGetOverlay    []test.On[*registry.GetOverlayRequest, *registry.Overlay]
	OnListOverlays  []test.On[*registry.ListOverlaysRequest, *registry.ListOverlaysResponse]
	OnCreateOverlay []test.On[*registry.CreateOverlayRequest, *registry.Overlay]
	OnDeleteOverlay []test.On[*registry.DeleteOverlayRequest, *emptypb.Empty]
}

func (m MockOverlaysServiceClient) GetOverlay(_ context.Context, given *registry.GetOverlayRequest, _ ...grpc.CallOption) (*registry.Overlay, error) {
	return test.Match(m.T, m.OnGetOverlay, given, "get-overlay", protocmp.Transform())
}

func (m MockOverlaysServiceClient) ListOverlays(_ context.Context, given *registry.ListOverlaysRequest, _ ...grpc.CallOption) (*registry.ListOverlaysResponse, error) {
	return test.Match(m.T, m.OnListOverlays, given, "list-overlays", protocmp.Transform())
}

func (m MockOverlaysServiceClient) CreateOverlay(_ context.Context, given *registry.CreateOverlayRequest, _ ...grpc.CallOption) (*registry.Overlay, error) {
	return test.Match(m.T, m.OnCreateOverlay, given, "create-overlay", protocmp.Transform())
}

func (m MockOverlaysServiceClient) DeleteOverlay(_ context.Context, given *registry.DeleteOverlayRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteOverlay, given, "delete-overlay", protocmp.Transform())
}
