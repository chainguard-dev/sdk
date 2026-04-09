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

var _ registry.TagsServiceClient = (*MockTagsServiceClient)(nil)

type MockTagsServiceClient struct {
	registry.TagsServiceClient
	T *testing.T

	OnGetTag   []test.On[*registry.GetTagRequest, *registry.Tag]
	OnListTags []test.On[*registry.ListTagsRequest, *registry.ListTagsResponse]
}

func (m MockTagsServiceClient) GetTag(_ context.Context, given *registry.GetTagRequest, _ ...grpc.CallOption) (*registry.Tag, error) {
	return test.Match(m.T, m.OnGetTag, given, "get-tag", protocmp.Transform())
}

func (m MockTagsServiceClient) ListTags(_ context.Context, given *registry.ListTagsRequest, _ ...grpc.CallOption) (*registry.ListTagsResponse, error) {
	return test.Match(m.T, m.OnListTags, given, "list-tags", protocmp.Transform())
}
