/*
Copyright 2022 Chainguard, Inc.
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

var _ tenant.RecordsClient = (*MockRecordsClient)(nil)

type MockRecordsClient struct {
	OnList []RecordsOnList
}

type RecordsOnList struct {
	Given *tenant.RecordFilter
	List  *tenant.RecordList
	Error error
}

func (m MockRecordsClient) List(_ context.Context, given *tenant.RecordFilter, _ ...grpc.CallOption) (*tenant.RecordList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
