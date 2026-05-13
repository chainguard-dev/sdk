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

	argos "chainguard.dev/sdk/proto/platform/argos/v1"
)

var _ argos.ArgosDocumentsClient = (*MockArgosDocumentsClient)(nil)

type MockArgosDocumentsClient struct {
	argos.ArgosDocumentsClient

	OnCreate []ArgosDocumentsOnCreate
	OnList   []ArgosDocumentsOnList
	OnDelete []ArgosDocumentsOnDelete
}

type ArgosDocumentsOnCreate struct {
	Given   *argos.CreateArgosDocumentRequest
	Created *argos.ArgosDocument
	Error   error
}

type ArgosDocumentsOnList struct {
	Given *argos.ArgosDocumentFilter
	List  *argos.ArgosDocumentList
	Error error
}

type ArgosDocumentsOnDelete struct {
	Given *argos.DeleteArgosDocumentRequest
	Error error
}

func (m MockArgosDocumentsClient) Create(_ context.Context, given *argos.CreateArgosDocumentRequest, _ ...grpc.CallOption) (*argos.ArgosDocument, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosDocumentsClient) List(_ context.Context, given *argos.ArgosDocumentFilter, _ ...grpc.CallOption) (*argos.ArgosDocumentList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosDocumentsClient) Delete(_ context.Context, given *argos.DeleteArgosDocumentRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}
