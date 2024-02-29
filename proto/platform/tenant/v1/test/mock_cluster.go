/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
)

var _ tenant.ClustersClient = (*MockClustersClient)(nil)

type MockClustersClient struct {
	OnCreate   []ClustersOnCreate
	OnDelete   []ClustersOnDelete
	OnList     []ClustersOnList
	OnUpdate   []ClustersOnUpdate
	OnConfig   []ClustersOnConfig
	OnProfiles []ClustersOnProfiles
	OnDiscover []ClustersOnDiscover
	OnCIDR     []ClustersOnCIDR
}

type ClustersOnCreate struct {
	Given   *tenant.CreateClusterRequest
	Created *tenant.Cluster
	Error   error
}

type ClustersOnDelete struct {
	Given *tenant.CreateClusterRequest
	Error error
}

type ClustersOnList struct {
	Given *tenant.ClusterFilter
	List  *tenant.ClusterList
	Error error
}

type ClustersOnUpdate struct {
	Given   *tenant.Cluster
	Updated *tenant.Cluster
	Error   error
}

type ClustersOnConfig struct {
	Given  *tenant.ClusterConfigRequest
	Config *tenant.ClusterConfigResponse
	Error  error
}

type ClustersOnProfiles struct {
	Given  *tenant.ClusterProfilesRequest
	Config *tenant.ClusterProfilesResponse
	Error  error
}

type ClustersOnDiscover struct {
	Given  *tenant.ClusterDiscoveryRequest
	Config *tenant.ClusterDiscoveryResponse
	Error  error
}

type ClustersOnCIDR struct {
	Given  *tenant.ClusterCIDRRequest
	Config *tenant.ClusterCIDRResponse
	Error  error
}

func (m MockClustersClient) Create(_ context.Context, given *tenant.CreateClusterRequest, _ ...grpc.CallOption) (*tenant.Cluster, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) Delete(_ context.Context, given *tenant.DeleteClusterRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) List(_ context.Context, given *tenant.ClusterFilter, _ ...grpc.CallOption) (*tenant.ClusterList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) Update(_ context.Context, given *tenant.Cluster, _ ...grpc.CallOption) (*tenant.Cluster, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) Config(_ context.Context, given *tenant.ClusterConfigRequest, _ ...grpc.CallOption) (*tenant.ClusterConfigResponse, error) {
	for _, o := range m.OnConfig {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Config, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) Profiles(_ context.Context, given *tenant.ClusterProfilesRequest, _ ...grpc.CallOption) (*tenant.ClusterProfilesResponse, error) {
	for _, o := range m.OnProfiles {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Config, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) Discover(_ context.Context, given *tenant.ClusterDiscoveryRequest, _ ...grpc.CallOption) (*tenant.ClusterDiscoveryResponse, error) {
	for _, o := range m.OnDiscover {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Config, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockClustersClient) CIDR(_ context.Context, given *tenant.ClusterCIDRRequest, _ ...grpc.CallOption) (*tenant.ClusterCIDRResponse, error) {
	for _, o := range m.OnCIDR {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Config, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
