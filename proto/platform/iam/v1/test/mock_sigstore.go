/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ iam.SigstoreServiceClient = (*MockSigstoreClient)(nil)

type MockSigstoreClient struct {
	OnCreate []SigstoreOnCreate
	OnUpdate []SigstoreOnUpdate
	OnDelete []SigstoreOnDelete
	OnList   []SigstoreOnList
}

type SigstoreOnCreate struct {
	Given   *iam.CreateSigstoreRequest
	Created *iam.Sigstore
	Error   error
}

type SigstoreOnUpdate struct {
	Given   *iam.Sigstore
	Updated *iam.Sigstore
	Error   error
}

type SigstoreOnDelete struct {
	Given *iam.DeleteSigstoreRequest
	Error error
}

type SigstoreOnList struct {
	Given *iam.SigstoreFilter
	List  *iam.SigstoreList
	Error error
}

func (m MockSigstoreClient) Create(_ context.Context, given *iam.CreateSigstoreRequest, _ ...grpc.CallOption) (*iam.Sigstore, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSigstoreClient) Update(_ context.Context, given *iam.Sigstore, _ ...grpc.CallOption) (*iam.Sigstore, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSigstoreClient) Delete(_ context.Context, given *iam.DeleteSigstoreRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockSigstoreClient) List(_ context.Context, given *iam.SigstoreFilter, _ ...grpc.CallOption) (*iam.SigstoreList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
