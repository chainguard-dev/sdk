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

var _ argos.ArgosOSVClient = (*MockArgosOSVClient)(nil)

type MockArgosOSVClient struct {
	argos.ArgosOSVClient

	OnQuery      []ArgosOSVOnQuery
	OnQueryBatch []ArgosOSVOnQueryBatch
	OnGetVuln    []ArgosOSVOnGetVuln
}

type ArgosOSVOnQuery struct {
	Given    *argos.OSVQueryRequest
	Response *argos.OSVQueryResponse
	Error    error
}

type ArgosOSVOnQueryBatch struct {
	Given    *argos.OSVQueryBatchRequest
	Response *argos.OSVQueryBatchResponse
	Error    error
}

type ArgosOSVOnGetVuln struct {
	Given  *argos.GetOSVRequest
	Record *argos.OSVRecord
	Error  error
}

func (m MockArgosOSVClient) Query(_ context.Context, given *argos.OSVQueryRequest, _ ...grpc.CallOption) (*argos.OSVQueryResponse, error) {
	for _, o := range m.OnQuery {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosOSVClient) QueryBatch(_ context.Context, given *argos.OSVQueryBatchRequest, _ ...grpc.CallOption) (*argos.OSVQueryBatchResponse, error) {
	for _, o := range m.OnQueryBatch {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosOSVClient) GetVuln(_ context.Context, given *argos.GetOSVRequest, _ ...grpc.CallOption) (*argos.OSVRecord, error) {
	for _, o := range m.OnGetVuln {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Record, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
