/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ libraries.ArtifactsClient = (*MockArtifactsClient)(nil)

type MockArtifactsClient struct {
	libraries.EntitlementsClient

	OnList             []ArtifactsOnList
	OnListVersions     []ArtifactsOnListVersions
	OnGetArtifactCount []ArtifactsOnGetArtifactCount
}

type ArtifactsOnList struct {
	Given *libraries.ArtifactFilter
	List  *libraries.ArtifactList
	Error error
}

type ArtifactsOnListVersions struct {
	Given *libraries.ArtifactVersionFilter
	List  *libraries.ArtifactVersionList
	Error error
}

type ArtifactsOnGetArtifactCount struct {
	Given  *libraries.GetArtifactCountRequest
	Result *libraries.GetArtifactCountResponse
	Error  error
}

func (m MockArtifactsClient) List(_ context.Context, given *libraries.ArtifactFilter, _ ...grpc.CallOption) (*libraries.ArtifactList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArtifactsClient) ListVersions(_ context.Context, given *libraries.ArtifactVersionFilter, _ ...grpc.CallOption) (*libraries.ArtifactVersionList, error) {
	for _, o := range m.OnListVersions {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArtifactsClient) GetArtifactCount(_ context.Context, given *libraries.GetArtifactCountRequest, _ ...grpc.CallOption) (*libraries.GetArtifactCountResponse, error) {
	for _, o := range m.OnGetArtifactCount {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Result, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
