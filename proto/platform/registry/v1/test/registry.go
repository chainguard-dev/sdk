/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ registry.Clients = (*MockRegistryClients)(nil)

type MockRegistryClients struct {
	OnClose error

	RegistryClient MockRegistryClient
}

func (m MockRegistryClients) Registry() registry.RegistryClient {
	return &m.RegistryClient
}

func (m MockRegistryClients) Close() error {
	return m.OnClose
}

var _ registry.RegistryClient = (*MockRegistryClient)(nil)

type MockRegistryClient struct {
	registry.RegistryClient

	OnCreateRepos    []ReposOnCreate
	OnDeleteRepos    []ReposOnDelete
	OnListRepos      []ReposOnList
	OnCreateTags     []TagsOnCreate
	OnDeleteTags     []TagsOnDelete
	OnUpdateTag      []TagOnUpdate
	OnListTags       []TagsOnList
	OnUpdateRepo     []RepoOnUpdate
	OnListTagHistory []TagHistoryOnList
	OnGetSbom        []SbomOnGet
	OnGetVulnReport  []VulnReportOnGet
	OnDiffImage      []DiffImage
}

type DiffImage struct {
	Given *registry.DiffImageRequest
	Diff  *registry.DiffImageResponse
	Error error
}

type ReposOnCreate struct {
	Given   *registry.CreateRepoRequest
	Created *registry.Repo
	Error   error
}

type ReposOnDelete struct {
	Given *registry.DeleteRepoRequest
	Error error
}

type ReposOnList struct {
	Given *registry.RepoFilter
	List  *registry.RepoList
	Error error
}

type TagsOnCreate struct {
	Given   *registry.CreateTagRequest
	Created *registry.Tag
	Error   error
}

type TagsOnDelete struct {
	Given *registry.DeleteTagRequest
	Error error
}

type TagOnUpdate struct {
	Given   *registry.Tag
	Updated *registry.Tag
	Error   error
}

type TagsOnList struct {
	Given *registry.TagFilter
	List  *registry.TagList
	Error error
}

type RepoOnUpdate struct {
	Given   *registry.Repo
	Updated *registry.Repo
	Error   error
}

type TagHistoryOnList struct {
	Given *registry.TagHistoryFilter
	List  *registry.TagHistoryList
	Error error
}

type SbomOnGet struct {
	Given *registry.SbomRequest
	Get   *tenant.Sbom2
	Error error
}

type VulnReportOnGet struct {
	Given *registry.VulnReportRequest
	Get   *tenant.VulnReport
	Error error
}

func (m MockRegistryClient) CreateRepo(_ context.Context, given *registry.CreateRepoRequest, _ ...grpc.CallOption) (*registry.Repo, error) {
	for _, o := range m.OnCreateRepos {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) DeleteRepo(_ context.Context, given *registry.DeleteRepoRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteRepos {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) ListRepos(_ context.Context, given *registry.RepoFilter, _ ...grpc.CallOption) (*registry.RepoList, error) { //nolint: revive
	for _, o := range m.OnListRepos {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) CreateTag(_ context.Context, given *registry.CreateTagRequest, _ ...grpc.CallOption) (*registry.Tag, error) {
	for _, o := range m.OnCreateTags {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) UpdateTag(_ context.Context, given *registry.Tag, _ ...grpc.CallOption) (*registry.Tag, error) {
	for _, o := range m.OnUpdateTag {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) DeleteTag(_ context.Context, given *registry.DeleteTagRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteTags {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) ListTags(_ context.Context, given *registry.TagFilter, _ ...grpc.CallOption) (*registry.TagList, error) { //nolint: revive
	for _, o := range m.OnListTags {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) UpdateRepo(_ context.Context, given *registry.Repo, _ ...grpc.CallOption) (*registry.Repo, error) { //nolint: revive
	for _, o := range m.OnUpdateRepo {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) ListTagHistory(_ context.Context, given *registry.TagHistoryFilter, _ ...grpc.CallOption) (*registry.TagHistoryList, error) { //nolint: revive
	for _, o := range m.OnListTagHistory {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetSbom(_ context.Context, given *registry.SbomRequest, _ ...grpc.CallOption) (*tenant.Sbom2, error) { //nolint: revive
	for _, o := range m.OnGetSbom {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetVulnReport(_ context.Context, given *registry.VulnReportRequest, _ ...grpc.CallOption) (*tenant.VulnReport, error) { //nolint: revive
	for _, o := range m.OnGetVulnReport {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) DiffImage(_ context.Context, given *registry.DiffImageRequest, _ ...grpc.CallOption) (*registry.DiffImageResponse, error) {
	for _, o := range m.OnDiffImage {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Diff, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
