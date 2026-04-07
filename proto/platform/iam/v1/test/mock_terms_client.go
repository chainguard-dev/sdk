/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ iam.TermsClient = (*MockTermsClient)(nil)

type MockTermsClient struct {
	OnAcceptTerms  []TermsOnAcceptTerms
	OnListAccepted []TermsOnListAccepted
}

type TermsOnAcceptTerms struct {
	Given *iam.AcceptTermsRequest
	Error error
}

type TermsOnListAccepted struct {
	Given *iam.TermsFilter
	List  *iam.TermsList
	Error error
}

func (m *MockTermsClient) AcceptTerms(_ context.Context, given *iam.AcceptTermsRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	if len(m.OnAcceptTerms) == 0 {
		return nil, fmt.Errorf("unexpected call to AcceptTerms with %v", given)
	}
	next := m.OnAcceptTerms[0]
	m.OnAcceptTerms = m.OnAcceptTerms[1:]
	return &emptypb.Empty{}, next.Error
}

func (m *MockTermsClient) ListAccepted(_ context.Context, given *iam.TermsFilter, _ ...grpc.CallOption) (*iam.TermsList, error) {
	if len(m.OnListAccepted) == 0 {
		return nil, fmt.Errorf("unexpected call to ListAccepted with %v", given)
	}
	next := m.OnListAccepted[0]
	m.OnListAccepted = m.OnListAccepted[1:]
	return next.List, next.Error
}
