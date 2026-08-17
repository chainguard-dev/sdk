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

	ReposServiceClient  MockReposServiceClient
	TagsServiceClient   MockTagsServiceClient
	ImagesServiceClient MockImagesServiceClient

	OverlaysServiceClient        MockOverlaysServiceClient
	OverlayBindingsServiceClient MockOverlayBindingsServiceClient
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

// ImagesService implements [v2beta1.Clients].
func (m *MockClients) ImagesService() registry.ImagesServiceClient {
	return &m.ImagesServiceClient
}

// OverlaysService implements [v2beta1.Clients].
func (m *MockClients) OverlaysService() registry.OverlaysServiceClient {
	return &m.OverlaysServiceClient
}

// OverlayBindingsService implements [v2beta1.Clients].
func (m *MockClients) OverlayBindingsService() registry.OverlayBindingsServiceClient {
	return &m.OverlayBindingsServiceClient
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

// ListOverlaysAll implements [v2beta1.Clients].
func (m *MockClients) ListOverlaysAll(ctx context.Context, req *registry.ListOverlaysRequest) ([]*registry.Overlay, error) {
	resp, err := m.OverlaysServiceClient.ListOverlays(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetOverlays(), nil
}

// ListOverlaysIter implements [v2beta1.Clients].
func (m *MockClients) ListOverlaysIter(ctx context.Context, req *registry.ListOverlaysRequest) iter.Seq2[*registry.Overlay, error] {
	return test.MockIter(m.ListOverlaysAll(ctx, req))
}

// ListOverlayBindingsAll implements [v2beta1.Clients].
func (m *MockClients) ListOverlayBindingsAll(ctx context.Context, req *registry.ListOverlayBindingsRequest) ([]*registry.OverlayBinding, error) {
	resp, err := m.OverlayBindingsServiceClient.ListOverlayBindings(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetOverlayBindings(), nil
}

// ListOverlayBindingsIter implements [v2beta1.Clients].
func (m *MockClients) ListOverlayBindingsIter(ctx context.Context, req *registry.ListOverlayBindingsRequest) iter.Seq2[*registry.OverlayBinding, error] {
	return test.MockIter(m.ListOverlayBindingsAll(ctx, req))
}

var _ registry.Clients = (*MockClients)(nil)
