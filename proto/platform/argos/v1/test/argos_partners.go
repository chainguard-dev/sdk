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
