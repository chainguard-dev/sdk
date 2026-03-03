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

var _ libraries.NpmPackagesClient = (*MockNpmPackagesClient)(nil)

type MockNpmPackagesClient struct {
	OnList         []NpmPackagesOnList
	OnListVersions []NpmPackagesOnListVersions
}

type NpmPackagesOnList struct {
	Given *libraries.NpmPackageFilter
	List  *libraries.NpmPackageList
	Error error
}

type NpmPackagesOnListVersions struct {
	Given *libraries.NpmPackageVersionFilter
	List  *libraries.NpmPackageVersionList
	Error error
}

func (m MockNpmPackagesClient) List(_ context.Context, given *libraries.NpmPackageFilter, _ ...grpc.CallOption) (*libraries.NpmPackageList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockNpmPackagesClient) ListVersions(_ context.Context, given *libraries.NpmPackageVersionFilter, _ ...grpc.CallOption) (*libraries.NpmPackageVersionList, error) {
	for _, o := range m.OnListVersions {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
