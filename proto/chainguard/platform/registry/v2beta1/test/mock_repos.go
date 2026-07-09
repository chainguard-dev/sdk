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

var _ registry.ReposServiceClient = (*MockReposServiceClient)(nil)

type MockReposServiceClient struct {
	registry.ReposServiceClient
	T *testing.T

	OnGetRepo          []test.On[*registry.GetRepoRequest, *registry.Repo]
	OnListRepos        []test.On[*registry.ListReposRequest, *registry.ListReposResponse]
	OnGetRepoReadme    []test.On[*registry.GetRepoReadmeRequest, *registry.RepoReadme]
	OnUpdateRepoReadme []test.On[*registry.UpdateRepoReadmeRequest, *registry.RepoReadme]
}

func (m MockReposServiceClient) GetRepo(_ context.Context, given *registry.GetRepoRequest, _ ...grpc.CallOption) (*registry.Repo, error) {
	return test.Match(m.T, m.OnGetRepo, given, "get-repo", protocmp.Transform())
}

func (m MockReposServiceClient) ListRepos(_ context.Context, given *registry.ListReposRequest, _ ...grpc.CallOption) (*registry.ListReposResponse, error) {
	return test.Match(m.T, m.OnListRepos, given, "list-repos", protocmp.Transform())
}

func (m MockReposServiceClient) GetRepoReadme(_ context.Context, given *registry.GetRepoReadmeRequest, _ ...grpc.CallOption) (*registry.RepoReadme, error) {
	return test.Match(m.T, m.OnGetRepoReadme, given, "get-repo-readme", protocmp.Transform())
}

func (m MockReposServiceClient) UpdateRepoReadme(_ context.Context, given *registry.UpdateRepoReadmeRequest, _ ...grpc.CallOption) (*registry.RepoReadme, error) {
	return test.Match(m.T, m.OnUpdateRepoReadme, given, "update-repo-readme", protocmp.Transform())
}
