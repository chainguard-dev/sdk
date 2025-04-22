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

	RegistryClient        MockRegistryClient
	VulnerabilitiesClient MockVulnerabilitiesClient
	ApkoClient            MockApkoClient
	EntitlementsClient    MockEntitlementsClient
}

func (m MockRegistryClients) Registry() registry.RegistryClient {
	return &m.RegistryClient
}

func (m MockRegistryClients) Vulnerabilities() registry.VulnerabilitiesClient {
	return &m.VulnerabilitiesClient
}

func (m MockRegistryClients) Apko() registry.ApkoClient {
	return &m.ApkoClient
}

func (m MockRegistryClients) Entitlements() registry.EntitlementsClient {
	return &m.EntitlementsClient
}

func (m MockRegistryClients) Close() error {
	return m.OnClose
}

var _ registry.RegistryClient = (*MockRegistryClient)(nil)

type MockRegistryClient struct {
	registry.RegistryClient

	OnCreateRepos               []ReposOnCreate
	OnDeleteRepos               []ReposOnDelete
	OnListRepos                 []ReposOnList
	OnCreateTags                []TagsOnCreate
	OnDeleteTags                []TagsOnDelete
	OnUpdateTag                 []TagOnUpdate
	OnListTags                  []TagsOnList
	OnUpdateRepo                []RepoOnUpdate
	OnListTagHistory            []TagHistoryOnList
	OnGetImageConfig            []ImageConfigOnGet
	OnGetSbom                   []SbomOnGet
	OnGetVulnReport             []VulnReportOnGet
	OnListManifestMetadata      []ManifestMetadataOnList
	OnGetRawSbom                []RawSbomOnGet
	OnGetPackageVersionMetadata []PackageVersionMetadataOnGet
	OnListBuildReports          []BuildReportsOnList
	OnGetBuildStatus            []BuildStatusOnGet
	OnGetUpdateStatus           []UpdateStatusOnGet
	OnGetRepoCountBySource      []RepoCountBySourceOnGet
	OnGetArchs                  []ArchsOnGet
	OnGetSize                   []SizeOnGet
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

type ImageConfigOnGet struct {
	Given *registry.ImageConfigRequest
	Get   *registry.ImageConfig
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

type ManifestMetadataOnList struct {
	Given *registry.ManifestMetadataFilter
	List  *registry.ManifestMetadataList
	Error error
}

type RawSbomOnGet struct {
	Given *registry.RawSbomRequest
	Get   *registry.RawSbom
	Error error
}

type PackageVersionMetadataOnGet struct {
	Given *registry.PackageVersionMetadataRequest
	Get   *registry.PackageVersionMetadata
	Error error
}

type BuildReportsOnList struct {
	Given *registry.BuildReportFilter
	List  *registry.BuildReportList
	Error error
}

type BuildStatusOnGet struct {
	Given *registry.BuildReportFilter
	Get   *registry.BuildStatus
	Error error
}

type UpdateStatusOnGet struct {
	Given *registry.UpdateStatusRequest
	Get   *registry.UpdateStatus
	Error error
}

type RepoCountBySourceOnGet struct {
	Given *registry.GetRepoCountBySourceRequest
	Get   *registry.RepoCount
	Error error
}

type ArchsOnGet struct {
	Given *registry.ArchRequest
	Get   *registry.Archs
	Error error
}

type SizeOnGet struct {
	Given *registry.SizeRequest
	Get   *registry.Size
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

func (m MockRegistryClient) GetImageConfig(_ context.Context, given *registry.ImageConfigRequest, _ ...grpc.CallOption) (*registry.ImageConfig, error) { //nolint: revive
	for _, o := range m.OnGetImageConfig {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
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

func (m MockRegistryClient) ListManifestMetadata(_ context.Context, given *registry.ManifestMetadataFilter, _ ...grpc.CallOption) (*registry.ManifestMetadataList, error) {
	for _, o := range m.OnListManifestMetadata {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetRawSbom(_ context.Context, given *registry.RawSbomRequest, _ ...grpc.CallOption) (*registry.RawSbom, error) { //nolint: revive
	for _, o := range m.OnGetRawSbom {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetPackageVersionMetadata(_ context.Context, given *registry.PackageVersionMetadataRequest, _ ...grpc.CallOption) (*registry.PackageVersionMetadata, error) { //nolint: revive
	for _, o := range m.OnGetPackageVersionMetadata {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) ListBuildReports(_ context.Context, given *registry.BuildReportFilter, _ ...grpc.CallOption) (*registry.BuildReportList, error) {
	for _, o := range m.OnListBuildReports {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetBuildStatus(_ context.Context, given *registry.BuildReportFilter, _ ...grpc.CallOption) (*registry.BuildStatus, error) {
	for _, o := range m.OnGetBuildStatus {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetUpdateStatus(_ context.Context, given *registry.UpdateStatusRequest, _ ...grpc.CallOption) (*registry.UpdateStatus, error) {
	for _, o := range m.OnGetUpdateStatus {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetRepoCountBySource(_ context.Context, given *registry.GetRepoCountBySourceRequest, _ ...grpc.CallOption) (*registry.RepoCount, error) {
	for _, o := range m.OnGetRepoCountBySource {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetArchs(_ context.Context, given *registry.ArchRequest, _ ...grpc.CallOption) (*registry.Archs, error) {
	for _, o := range m.OnGetArchs {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRegistryClient) GetSize(_ context.Context, given *registry.SizeRequest, _ ...grpc.CallOption) (*registry.Size, error) {
	for _, o := range m.OnGetSize {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Get, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
