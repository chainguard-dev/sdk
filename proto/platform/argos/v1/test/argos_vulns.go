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

	argos "chainguard.dev/sdk/proto/platform/argos/v1"
)

var _ argos.ArgosVulnsClient = (*MockArgosVulnsClient)(nil)

type MockArgosVulnsClient struct {
	argos.ArgosVulnsClient

	OnBatchGet   []ArgosVulnsOnBatchGet
	OnListForOrg []ArgosVulnsOnListForOrg
}

type ArgosVulnsOnBatchGet struct {
	Given    *argos.BatchGetVulnsRequest
	Response *argos.VulnList
	Error    error
}

type ArgosVulnsOnListForOrg struct {
	Given    *argos.ListOrgVulnsRequest
	Response *argos.VulnList
	Error    error
}

func (m MockArgosVulnsClient) BatchGet(_ context.Context, given *argos.BatchGetVulnsRequest, _ ...grpc.CallOption) (*argos.VulnList, error) {
	for _, o := range m.OnBatchGet {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosVulnsClient) ListForOrg(_ context.Context, given *argos.ListOrgVulnsRequest, _ ...grpc.CallOption) (*argos.VulnList, error) {
	for _, o := range m.OnListForOrg {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
