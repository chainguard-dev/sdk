/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"iter"

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

type MockClients struct {
	OnClose error

	ReposServiceClient MockReposServiceClient
	TagsServiceClient  MockTagsServiceClient
}

// Close implements [v2beta1.Clients].
func (m *MockClients) Close() error {
	return m.OnClose
}

// ReposService implements [v2beta1.Clients].
func (m *MockClients) ReposService() registry.ReposServiceClient {
	return &m.ReposServiceClient
}

// TagsService implements [v2beta1.Clients].
func (m *MockClients) TagsService() registry.TagsServiceClient {
	return &m.TagsServiceClient
}

// ListReposAll implements [v2beta1.Clients].
func (m *MockClients) ListReposAll(ctx context.Context, req *registry.ListReposRequest) ([]*registry.Repo, error) {
	resp, err := m.ReposServiceClient.ListRepos(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetRepos(), nil
}

// ListReposIter implements [v2beta1.Clients].
func (m *MockClients) ListReposIter(ctx context.Context, req *registry.ListReposRequest) iter.Seq2[*registry.Repo, error] {
	return test.MockIter(m.ListReposAll(ctx, req))
}

// ListTagsAll implements [v2beta1.Clients].
func (m *MockClients) ListTagsAll(ctx context.Context, req *registry.ListTagsRequest) ([]*registry.Tag, error) {
	resp, err := m.TagsServiceClient.ListTags(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetTags(), nil
}

// ListTagsIter implements [v2beta1.Clients].
func (m *MockClients) ListTagsIter(ctx context.Context, req *registry.ListTagsRequest) iter.Seq2[*registry.Tag, error] {
	return test.MockIter(m.ListTagsAll(ctx, req))
}

var _ registry.Clients = (*MockClients)(nil)
