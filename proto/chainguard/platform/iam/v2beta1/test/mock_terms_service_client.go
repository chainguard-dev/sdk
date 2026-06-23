/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	iam "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	"google.golang.org/grpc"
)

var _ iam.TermsServiceClient = (*MockTermsServiceClient)(nil)

type MockTermsServiceClient struct {
	OnAcceptTerms          []TermsOnAcceptTerms
	OnListTermsAcceptances []TermsOnListTermsAcceptances
}

type TermsOnAcceptTerms struct {
	Given *iam.AcceptTermsRequest
	Resp  *iam.AcceptTermsResponse
	Error error
}

type TermsOnListTermsAcceptances struct {
	Given *iam.ListTermsAcceptancesRequest
	Resp  *iam.ListTermsAcceptancesResponse
	Error error
}

func (m *MockTermsServiceClient) AcceptTerms(_ context.Context, given *iam.AcceptTermsRequest, _ ...grpc.CallOption) (*iam.AcceptTermsResponse, error) {
	if len(m.OnAcceptTerms) == 0 {
		return nil, fmt.Errorf("unexpected call to AcceptTerms with %v", given)
	}
	next := m.OnAcceptTerms[0]
	m.OnAcceptTerms = m.OnAcceptTerms[1:]
	return next.Resp, next.Error
}

func (m *MockTermsServiceClient) ListTermsAcceptances(_ context.Context, given *iam.ListTermsAcceptancesRequest, _ ...grpc.CallOption) (*iam.ListTermsAcceptancesResponse, error) {
	if len(m.OnListTermsAcceptances) == 0 {
		return nil, fmt.Errorf("unexpected call to ListTermsAcceptances with %v", given)
	}
	next := m.OnListTermsAcceptances[0]
	m.OnListTermsAcceptances = m.OnListTermsAcceptances[1:]
	return next.Resp, next.Error
}
