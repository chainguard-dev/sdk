/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
)

var _ libraries.RemoteRegistriesClient = (*MockRemoteRegistriesClient)(nil)

type MockRemoteRegistriesClient struct {
	OnList    []RemoteRegistriesOnList
	OnUpdate  []RemoteRegistriesOnUpdate
	OnReorder []RemoteRegistriesOnReorder
}

type RemoteRegistriesOnList struct {
	Given *libraries.RemoteRegistryFilter
	List  *libraries.RemoteRegistryList
	Error error
}

type RemoteRegistriesOnUpdate struct {
	Given          *libraries.UpdateRemoteRegistryRequest
	RemoteRegistry *libraries.RemoteRegistry
	Error          error
}

type RemoteRegistriesOnReorder struct {
	Given *libraries.ReorderRemoteRegistriesRequest
	List  *libraries.RemoteRegistryList
	Error error
}

func (m MockRemoteRegistriesClient) List(_ context.Context, given *libraries.RemoteRegistryFilter, _ ...grpc.CallOption) (*libraries.RemoteRegistryList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRemoteRegistriesClient) Update(_ context.Context, given *libraries.UpdateRemoteRegistryRequest, _ ...grpc.CallOption) (*libraries.RemoteRegistry, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.RemoteRegistry, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRemoteRegistriesClient) Reorder(_ context.Context, given *libraries.ReorderRemoteRegistriesRequest, _ ...grpc.CallOption) (*libraries.RemoteRegistryList, error) {
	for _, o := range m.OnReorder {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
