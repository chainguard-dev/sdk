/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	vulnerabilities "chainguard.dev/sdk/proto/platform/vulnerabilities/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ vulnerabilities.AdvisoriesClient = (*MockAdvisoriesClient)(nil)

type MockAdvisoriesClient struct {
	vulnerabilities.AdvisoriesClient

	OnList   []AdvisoriesOnList
	OnCreate []AdvisoriesOnCreate
	OnUpdate []AdvisoriesOnUpdate
	OnDelete []AdvisoriesOnDelete
}

type AdvisoriesOnList struct {
	Given *vulnerabilities.AdvisoryFilter
	List  *vulnerabilities.AdvisoriesList
	Error error
}

type AdvisoriesOnCreate struct {
	Given   *vulnerabilities.Advisory
	Created *vulnerabilities.Advisory
	Error   error
}

type AdvisoriesOnUpdate struct {
	Given   *vulnerabilities.Advisory
	Updated *vulnerabilities.Advisory
	Error   error
}

type AdvisoriesOnDelete struct {
	Given *vulnerabilities.DeleteAdvisoryRequest
	Error error
}

func (m MockAdvisoriesClient) List(_ context.Context, given *vulnerabilities.AdvisoryFilter, _ ...grpc.CallOption) (*vulnerabilities.AdvisoriesList, error) { //nolint: revive
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAdvisoriesClient) Create(_ context.Context, given *vulnerabilities.CreateAdvisoryRequest, _ ...grpc.CallOption) (*vulnerabilities.Advisory, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAdvisoriesClient) Delete(_ context.Context, given *vulnerabilities.DeleteAdvisoryRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockAdvisoriesClient) Update(_ context.Context, given *vulnerabilities.Advisory, _ ...grpc.CallOption) (*vulnerabilities.Advisory, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
