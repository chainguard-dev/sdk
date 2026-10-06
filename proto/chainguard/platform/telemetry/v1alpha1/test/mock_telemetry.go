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

	telemetry "chainguard.dev/sdk/proto/chainguard/platform/telemetry/v1alpha1"
)

var _ telemetry.TelemetryServiceClient = (*MockTelemetryServiceClient)(nil)

type MockTelemetryServiceClient struct {
	OnListContainers  []TelemetryOnListContainers
	OnListImages      []TelemetryOnListImages
	OnSummarizeImages []TelemetryOnSummarizeImages
}

type TelemetryOnListContainers struct {
	Given *telemetry.ListContainersRequest
	List  *telemetry.ListContainersResponse
	Error error
}

type TelemetryOnListImages struct {
	Given *telemetry.ListImagesRequest
	List  *telemetry.ListImagesResponse
	Error error
}

type TelemetryOnSummarizeImages struct {
	Given   *telemetry.SummarizeImagesRequest
	Summary *telemetry.ImageSummary
	Error   error
}

func (m MockTelemetryServiceClient) ListContainers(_ context.Context, given *telemetry.ListContainersRequest, _ ...grpc.CallOption) (*telemetry.ListContainersResponse, error) {
	for _, o := range m.OnListContainers {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockTelemetryServiceClient) ListImages(_ context.Context, given *telemetry.ListImagesRequest, _ ...grpc.CallOption) (*telemetry.ListImagesResponse, error) {
	for _, o := range m.OnListImages {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockTelemetryServiceClient) SummarizeImages(_ context.Context, given *telemetry.SummarizeImagesRequest, _ ...grpc.CallOption) (*telemetry.ImageSummary, error) {
	for _, o := range m.OnSummarizeImages {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Summary, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
