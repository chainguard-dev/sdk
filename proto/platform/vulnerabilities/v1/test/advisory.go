/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	v1 "chainguard.dev/sdk/proto/platform/vulnerabilities/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ v1.AdvisoriesClient = (*MockAdvisoriesClient)(nil)

type MockAdvisoriesClient struct {
	v1.AdvisoriesClient

	OnList   []AdvisoriesOnList
	OnCreate []AdvisoriesOnCreate
	OnDelete []AdvisoriesOnDelete
}

var _ v1.AdvisoriesClient = (*MockAdvisoriesClient)(nil)

type AdvisoriesOnList struct {
	Given *v1.AdvisoryFilter
	List  *v1.AdvisoriesList
	Error error
}

type AdvisoriesOnCreate struct {
	Given   *v1.Advisory
	Created *v1.Advisory
	Error   error
}

type AdvisoriesOnDelete struct {
	Given *v1.DeleteAdvisoryRequest
	Error error
}

func (m MockAdvisoriesClient) List(_ context.Context, given *v1.AdvisoryFilter, _ ...grpc.CallOption) (*v1.AdvisoriesList, error) { //nolint: revive
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAdvisoriesClient) Create(_ context.Context, given *v1.Advisory, _ ...grpc.CallOption) (*v1.Advisory, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAdvisoriesClient) Delete(_ context.Context, given *v1.DeleteAdvisoryRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}
