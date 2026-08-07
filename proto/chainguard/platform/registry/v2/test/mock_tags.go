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

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ registry.TagsServiceClient = (*MockTagsServiceClient)(nil)

type MockTagsServiceClient struct {
	registry.TagsServiceClient
	T *testing.T

	OnGetTag    []test.On[*registry.GetTagRequest, *registry.Tag]
	OnListTags  []test.On[*registry.ListTagsRequest, *registry.ListTagsResponse]
	OnCreateTag []test.On[*registry.CreateTagRequest, *registry.Tag]
	OnDeleteTag []test.On[*registry.DeleteTagRequest, *emptypb.Empty]
}

func (m MockTagsServiceClient) GetTag(_ context.Context, given *registry.GetTagRequest, _ ...grpc.CallOption) (*registry.Tag, error) {
	return test.Match(m.T, m.OnGetTag, given, "get-tag", protocmp.Transform())
}

func (m MockTagsServiceClient) ListTags(_ context.Context, given *registry.ListTagsRequest, _ ...grpc.CallOption) (*registry.ListTagsResponse, error) {
	return test.Match(m.T, m.OnListTags, given, "list-tags", protocmp.Transform())
}

func (m MockTagsServiceClient) CreateTag(_ context.Context, given *registry.CreateTagRequest, _ ...grpc.CallOption) (*registry.Tag, error) {
	return test.Match(m.T, m.OnCreateTag, given, "create-tag", protocmp.Transform())
}

func (m MockTagsServiceClient) DeleteTag(_ context.Context, given *registry.DeleteTagRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteTag, given, "delete-tag", protocmp.Transform())
}
