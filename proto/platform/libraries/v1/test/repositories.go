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
	"google.golang.org/protobuf/types/known/emptypb"

	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
)

var _ libraries.RepositoriesClient = (*MockRepositoriesClient)(nil)

type MockRepositoriesClient struct {
	OnCreate []RepositoriesOnCreate
	OnGet    []RepositoriesOnGet
	OnList   []RepositoriesOnList
	OnDelete []RepositoriesOnDelete
}

type RepositoriesOnCreate struct {
	Given      *libraries.CreateRepositoryRequest
	Repository *libraries.Repository
	Error      error
}

type RepositoriesOnGet struct {
	Given      *libraries.GetRepositoryRequest
	Repository *libraries.Repository
	Error      error
}

type RepositoriesOnList struct {
	Given *libraries.RepositoryFilter
	List  *libraries.RepositoryList
	Error error
}

type RepositoriesOnDelete struct {
	Given *libraries.DeleteRepositoryRequest
	Error error
}

func (m MockRepositoriesClient) Create(_ context.Context, given *libraries.CreateRepositoryRequest, _ ...grpc.CallOption) (*libraries.Repository, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Repository, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRepositoriesClient) Get(_ context.Context, given *libraries.GetRepositoryRequest, _ ...grpc.CallOption) (*libraries.Repository, error) {
	for _, o := range m.OnGet {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Repository, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRepositoriesClient) List(_ context.Context, given *libraries.RepositoryFilter, _ ...grpc.CallOption) (*libraries.RepositoryList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockRepositoriesClient) Delete(_ context.Context, given *libraries.DeleteRepositoryRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
