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

var _ argos.ArgosPartnersClient = (*MockArgosPartnersClient)(nil)

type MockArgosPartnersClient struct {
	argos.ArgosPartnersClient

	OnGetCapabilities       []ArgosPartnersOnGetCapabilities
	OnUpdateCapabilities    []ArgosPartnersOnUpdateCapabilities
	OnDeleteCapabilities    []ArgosPartnersOnDeleteCapabilities
	OnListMitigationMethods []ArgosPartnersOnListMitigationMethods

	OnCreateMitigation   []ArgosPartnersOnCreateMitigation
	OnGetMitigation      []ArgosPartnersOnGetMitigation
	OnListMitigations    []ArgosPartnersOnListMitigations
	OnUpdateMitigation   []ArgosPartnersOnUpdateMitigation
	OnWithdrawMitigation []ArgosPartnersOnWithdrawMitigation

	OnListPartnerReviews      []ArgosPartnersOnListPartnerReviews
	OnDecidePartnerCapability []ArgosPartnersOnDecidePartnerCapability
}

type ArgosPartnersOnGetCapabilities struct {
	Given    *argos.GetPartnerCapabilitiesRequest
	Response *argos.PartnerCapabilities
	Error    error
}

type ArgosPartnersOnUpdateCapabilities struct {
	Given    *argos.UpdatePartnerCapabilitiesRequest
	Response *argos.UpdatePartnerCapabilitiesResponse
	Error    error
}

type ArgosPartnersOnDeleteCapabilities struct {
	Given    *argos.DeletePartnerCapabilitiesRequest
	Response *argos.PartnerCapabilities
	Error    error
}

type ArgosPartnersOnListMitigationMethods struct {
	Given    *argos.ListMitigationMethodsRequest
	Response *argos.ListMitigationMethodsResponse
	Error    error
}

type ArgosPartnersOnCreateMitigation struct {
	Given    *argos.CreateMitigationRequest
	Response *argos.Mitigation
	Error    error
}

type ArgosPartnersOnGetMitigation struct {
	Given    *argos.GetMitigationRequest
	Response *argos.Mitigation
	Error    error
}

type ArgosPartnersOnListMitigations struct {
	Given    *argos.ListMitigationsRequest
	Response *argos.ListMitigationsResponse
	Error    error
}

type ArgosPartnersOnUpdateMitigation struct {
	Given    *argos.UpdateMitigationRequest
	Response *argos.Mitigation
	Error    error
}

type ArgosPartnersOnWithdrawMitigation struct {
	Given    *argos.WithdrawMitigationRequest
	Response *argos.WithdrawMitigationResponse
	Error    error
}

type ArgosPartnersOnListPartnerReviews struct {
	Given    *argos.ListPartnerReviewsRequest
	Response *argos.ListPartnerReviewsResponse
	Error    error
}

type ArgosPartnersOnDecidePartnerCapability struct {
	Given    *argos.DecidePartnerCapabilityRequest
	Response *argos.DecidePartnerCapabilityResponse
	Error    error
}

func (m MockArgosPartnersClient) GetCapabilities(_ context.Context, given *argos.GetPartnerCapabilitiesRequest, _ ...grpc.CallOption) (*argos.PartnerCapabilities, error) {
	for _, o := range m.OnGetCapabilities {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) UpdateCapabilities(_ context.Context, given *argos.UpdatePartnerCapabilitiesRequest, _ ...grpc.CallOption) (*argos.UpdatePartnerCapabilitiesResponse, error) {
	for _, o := range m.OnUpdateCapabilities {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) DeleteCapabilities(_ context.Context, given *argos.DeletePartnerCapabilitiesRequest, _ ...grpc.CallOption) (*argos.PartnerCapabilities, error) {
	for _, o := range m.OnDeleteCapabilities {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) ListMitigationMethods(_ context.Context, given *argos.ListMitigationMethodsRequest, _ ...grpc.CallOption) (*argos.ListMitigationMethodsResponse, error) {
	for _, o := range m.OnListMitigationMethods {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) ListPartnerReviews(_ context.Context, given *argos.ListPartnerReviewsRequest, _ ...grpc.CallOption) (*argos.ListPartnerReviewsResponse, error) {
	for _, o := range m.OnListPartnerReviews {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) DecidePartnerCapability(_ context.Context, given *argos.DecidePartnerCapabilityRequest, _ ...grpc.CallOption) (*argos.DecidePartnerCapabilityResponse, error) {
	for _, o := range m.OnDecidePartnerCapability {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) CreateMitigation(_ context.Context, given *argos.CreateMitigationRequest, _ ...grpc.CallOption) (*argos.Mitigation, error) {
	for _, o := range m.OnCreateMitigation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) GetMitigation(_ context.Context, given *argos.GetMitigationRequest, _ ...grpc.CallOption) (*argos.Mitigation, error) {
	for _, o := range m.OnGetMitigation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) ListMitigations(_ context.Context, given *argos.ListMitigationsRequest, _ ...grpc.CallOption) (*argos.ListMitigationsResponse, error) {
	for _, o := range m.OnListMitigations {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) UpdateMitigation(_ context.Context, given *argos.UpdateMitigationRequest, _ ...grpc.CallOption) (*argos.Mitigation, error) {
	for _, o := range m.OnUpdateMitigation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockArgosPartnersClient) WithdrawMitigation(_ context.Context, given *argos.WithdrawMitigationRequest, _ ...grpc.CallOption) (*argos.WithdrawMitigationResponse, error) {
	for _, o := range m.OnWithdrawMitigation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
