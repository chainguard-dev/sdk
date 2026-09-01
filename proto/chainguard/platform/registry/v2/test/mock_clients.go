/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"iter"

	registry "chainguard.dev/sdk/proto/chainguard/platform/registry/v2"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

type MockClients struct {
	OnClose error

	ReposServiceClient       MockReposServiceClient
	TagsServiceClient        MockTagsServiceClient
	ImagesServiceClient      MockImagesServiceClient
	StigReportsServiceClient MockStigReportsServiceClient
}

// Close implements [v2.Clients].
func (m *MockClients) Close() error {
	return m.OnClose
}

// ReposService implements [v2.Clients].
func (m *MockClients) ReposService() registry.ReposServiceClient {
	return &m.ReposServiceClient
}

// TagsService implements [v2.Clients].
func (m *MockClients) TagsService() registry.TagsServiceClient {
	return &m.TagsServiceClient
}

// ImagesService implements [v2.Clients].
func (m *MockClients) ImagesService() registry.ImagesServiceClient {
	return &m.ImagesServiceClient
}

// StigReportsService implements [v2.Clients].
func (m *MockClients) StigReportsService() registry.StigReportsServiceClient {
	return &m.StigReportsServiceClient
}

// ListReposAll implements [v2.Clients].
func (m *MockClients) ListReposAll(ctx context.Context, req *registry.ListReposRequest) ([]*registry.Repo, error) {
	resp, err := m.ReposServiceClient.ListRepos(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetRepos(), nil
}

// ListReposIter implements [v2.Clients].
func (m *MockClients) ListReposIter(ctx context.Context, req *registry.ListReposRequest) iter.Seq2[*registry.Repo, error] {
	return test.MockIter(m.ListReposAll(ctx, req))
}

// ListCatalogImagesAll implements [v2.Clients].
func (m *MockClients) ListCatalogImagesAll(ctx context.Context, req *registry.ListCatalogImagesRequest) ([]*registry.Repo, error) {
	resp, err := m.ReposServiceClient.ListCatalogImages(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetRepos(), nil
}

// ListCatalogImagesIter implements [v2.Clients].
func (m *MockClients) ListCatalogImagesIter(ctx context.Context, req *registry.ListCatalogImagesRequest) iter.Seq2[*registry.Repo, error] {
	return test.MockIter(m.ListCatalogImagesAll(ctx, req))
}

// ListTagsAll implements [v2.Clients].
func (m *MockClients) ListTagsAll(ctx context.Context, req *registry.ListTagsRequest) ([]*registry.Tag, error) {
	resp, err := m.TagsServiceClient.ListTags(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetTags(), nil
}

// ListTagsIter implements [v2.Clients].
func (m *MockClients) ListTagsIter(ctx context.Context, req *registry.ListTagsRequest) iter.Seq2[*registry.Tag, error] {
	return test.MockIter(m.ListTagsAll(ctx, req))
}

var _ registry.Clients = (*MockClients)(nil)
