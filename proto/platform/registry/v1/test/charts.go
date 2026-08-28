/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ registry.ChartsClient = (*MockChartsClient)(nil)

type MockChartsClient struct {
	registry.ChartsClient

	OnAddChart  []ChartOnAdd
	OnFindChart []ChartOnFind
}

type ChartOnAdd struct {
	Given    *registry.AddChartRequest
	Response *registry.AddChartResponse
	Error    error
}

type ChartOnFind struct {
	Given    *registry.FindChartRequest
	Response *registry.FindChartResponse
	Error    error
}

func (m *MockChartsClient) AddChart(_ context.Context, given *registry.AddChartRequest, _ ...grpc.CallOption) (*registry.AddChartResponse, error) {
	for _, o := range m.OnAddChart {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m *MockChartsClient) FindChart(_ context.Context, given *registry.FindChartRequest, _ ...grpc.CallOption) (*registry.FindChartResponse, error) {
	for _, o := range m.OnFindChart {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
