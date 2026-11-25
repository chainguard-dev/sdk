/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"
)

var _ skills.SkillsHardenClient = (*MockHardenClient)(nil)

type MockHardenClient struct {
	OnHardenSkill           []HardenOnHardenSkill
	OnGetHardenOperation    []HardenOnGetOperation
	OnCancelHardenOperation []HardenOnCancelOperation
	OnListHardenOperations  []HardenOnListOperations
}

type HardenOnHardenSkill struct {
	Given     *skills.HardenSkillRequest
	Operation *longrunningpb.Operation
	Error     error
}

type HardenOnGetOperation struct {
	Given     *skills.GetHardenOperationRequest
	Operation *longrunningpb.Operation
	Error     error
}

type HardenOnCancelOperation struct {
	Given     *skills.CancelHardenOperationRequest
	Operation *longrunningpb.Operation
	Error     error
}

type HardenOnListOperations struct {
	Given *skills.ListHardenOperationsRequest
	List  *skills.ListHardenOperationsResponse
	Error error
}

func (m MockHardenClient) HardenSkill(_ context.Context, given *skills.HardenSkillRequest, _ ...grpc.CallOption) (*longrunningpb.Operation, error) {
	for _, o := range m.OnHardenSkill {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Operation, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardenClient) GetHardenOperation(_ context.Context, given *skills.GetHardenOperationRequest, _ ...grpc.CallOption) (*longrunningpb.Operation, error) {
	for _, o := range m.OnGetHardenOperation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Operation, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardenClient) CancelHardenOperation(_ context.Context, given *skills.CancelHardenOperationRequest, _ ...grpc.CallOption) (*longrunningpb.Operation, error) {
	for _, o := range m.OnCancelHardenOperation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Operation, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardenClient) ListHardenOperations(_ context.Context, given *skills.ListHardenOperationsRequest, _ ...grpc.CallOption) (*skills.ListHardenOperationsResponse, error) {
	for _, o := range m.OnListHardenOperations {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
