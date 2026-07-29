/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ libraries.ResolutionCacheClient = (*MockResolutionCacheClient)(nil)

type MockResolutionCacheClient struct {
	OnZap       []ResolutionCacheOnZap
	OnSetOptOut []ResolutionCacheOnSetOptOut
	OnStatus    []ResolutionCacheOnStatus
	OnList      []ResolutionCacheOnList
}

type ResolutionCacheOnZap struct {
	Given   *libraries.ZapResolutionCacheRequest
	Summary *libraries.ZapResolutionCacheSummary
	Error   error
}

type ResolutionCacheOnSetOptOut struct {
	Given  *libraries.ResolutionCacheOptOutRequest
	OptOut *libraries.ResolutionCacheOptOut
	Error  error
}

type ResolutionCacheOnStatus struct {
	Given  *libraries.ResolutionCacheStatusFilter
	Status *libraries.ResolutionCacheStatus
	Error  error
}

type ResolutionCacheOnList struct {
	Given *libraries.ResolutionCacheEntryFilter
	List  *libraries.ResolutionCacheEntryList
	Error error
}

func (m MockResolutionCacheClient) Zap(_ context.Context, given *libraries.ZapResolutionCacheRequest, _ ...grpc.CallOption) (*libraries.ZapResolutionCacheSummary, error) {
	for _, o := range m.OnZap {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Summary, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockResolutionCacheClient) SetOptOut(_ context.Context, given *libraries.ResolutionCacheOptOutRequest, _ ...grpc.CallOption) (*libraries.ResolutionCacheOptOut, error) {
	for _, o := range m.OnSetOptOut {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.OptOut, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockResolutionCacheClient) Status(_ context.Context, given *libraries.ResolutionCacheStatusFilter, _ ...grpc.CallOption) (*libraries.ResolutionCacheStatus, error) {
	for _, o := range m.OnStatus {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Status, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockResolutionCacheClient) List(_ context.Context, given *libraries.ResolutionCacheEntryFilter, _ ...grpc.CallOption) (*libraries.ResolutionCacheEntryList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
