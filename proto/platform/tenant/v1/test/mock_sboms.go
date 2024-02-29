/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	tenant "chainguard.dev/sdk/proto/platform/tenant/v1"
)

var _ tenant.SbomsClient = (*MockSbomsClient)(nil)

type MockSbomsClient struct {
	OnList []SbomsOnList
}

type SbomsOnList struct {
	Given *tenant.Sbom2Filter
	List  *tenant.Sbom2List
	Error error
}

func (m MockSbomsClient) List(_ context.Context, given *tenant.Sbom2Filter, _ ...grpc.CallOption) (*tenant.Sbom2List, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
