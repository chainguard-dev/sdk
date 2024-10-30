/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
)

var _ registry.ApkoClient = (*MockApkoClient)(nil)

type MockApkoClient struct {
	registry.ApkoClient

	OnResolveConfig []OnResolveConfig
	OnBuildImage    []OnBuildImage
}

type OnResolveConfig struct {
	Given  *registry.ResolveConfigRequest
	Result *registry.ApkoConfig
	Error  error
}

type OnBuildImage struct {
	Given  *registry.BuildImageRequest
	Result *registry.BuildImageResponse
	Error  error
}

func (m MockApkoClient) ResolveConfig(_ context.Context, given *registry.ResolveConfigRequest, _ ...grpc.CallOption) (*registry.ApkoConfig, error) {
	for _, o := range m.OnResolveConfig {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Result, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockApkoClient) BuildImage(_ context.Context, given *registry.BuildImageRequest, _ ...grpc.CallOption) (*registry.BuildImageResponse, error) {
	for _, o := range m.OnBuildImage {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Result, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
