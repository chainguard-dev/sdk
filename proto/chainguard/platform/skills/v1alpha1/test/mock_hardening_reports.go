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

	skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"
)

var _ skills.SkillsHardeningReportsClient = (*MockHardeningReportsClient)(nil)

type MockHardeningReportsClient struct {
	OnGetHardeningReport    []HardeningReportsOnGet
	OnListHardeningReports  []HardeningReportsOnList
	OnGetTaxonomy           []HardeningReportsOnGetTaxonomy
	OnListHardeningFindings []HardeningReportsOnListFindings
	OnUpdateTaxonomy        []HardeningReportsOnUpdateTaxonomy
	OnUpdateHardeningReport []HardeningReportsOnUpdate
	OnDeleteHardeningReport []HardeningReportsOnDelete
	OnListRationales        []HardeningReportsOnListRationales
	OnUpdateRationales      []HardeningReportsOnUpdateRationales
}

type HardeningReportsOnListRationales struct {
	Given *skills.ListRationalesRequest
	List  *skills.ListRationalesResponse
	Error error
}

type HardeningReportsOnUpdateRationales struct {
	Given  *skills.UpdateRationalesRequest
	Result *skills.UpdateRationalesResponse
	Error  error
}

type HardeningReportsOnGet struct {
	Given  *skills.GetHardeningReportRequest
	Report *skills.HardeningReport
	Error  error
}

type HardeningReportsOnList struct {
	Given *skills.ListHardeningReportsRequest
	List  *skills.ListHardeningReportsResponse
	Error error
}

type HardeningReportsOnGetTaxonomy struct {
	Given    *skills.GetTaxonomyRequest
	Taxonomy *skills.Taxonomy
	Error    error
}

type HardeningReportsOnListFindings struct {
	Given *skills.ListHardeningFindingsRequest
	List  *skills.ListHardeningFindingsResponse
	Error error
}

type HardeningReportsOnUpdateTaxonomy struct {
	Given    *skills.UpdateTaxonomyRequest
	Taxonomy *skills.Taxonomy
	Error    error
}

type HardeningReportsOnUpdate struct {
	Given  *skills.UpdateHardeningReportRequest
	Report *skills.HardeningReport
	Error  error
}

type HardeningReportsOnDelete struct {
	Given *skills.DeleteHardeningReportRequest
	Error error
}

func (m MockHardeningReportsClient) GetHardeningReport(_ context.Context, given *skills.GetHardeningReportRequest, _ ...grpc.CallOption) (*skills.HardeningReport, error) {
	for _, o := range m.OnGetHardeningReport {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Report, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) ListHardeningReports(_ context.Context, given *skills.ListHardeningReportsRequest, _ ...grpc.CallOption) (*skills.ListHardeningReportsResponse, error) {
	for _, o := range m.OnListHardeningReports {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) GetTaxonomy(_ context.Context, given *skills.GetTaxonomyRequest, _ ...grpc.CallOption) (*skills.Taxonomy, error) {
	for _, o := range m.OnGetTaxonomy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Taxonomy, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) ListHardeningFindings(_ context.Context, given *skills.ListHardeningFindingsRequest, _ ...grpc.CallOption) (*skills.ListHardeningFindingsResponse, error) {
	for _, o := range m.OnListHardeningFindings {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) UpdateTaxonomy(_ context.Context, given *skills.UpdateTaxonomyRequest, _ ...grpc.CallOption) (*skills.Taxonomy, error) {
	for _, o := range m.OnUpdateTaxonomy {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Taxonomy, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) UpdateHardeningReport(_ context.Context, given *skills.UpdateHardeningReportRequest, _ ...grpc.CallOption) (*skills.HardeningReport, error) {
	for _, o := range m.OnUpdateHardeningReport {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Report, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) DeleteHardeningReport(_ context.Context, given *skills.DeleteHardeningReportRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteHardeningReport {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			if o.Error != nil {
				return nil, o.Error
			}
			return &emptypb.Empty{}, nil
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) ListRationales(_ context.Context, given *skills.ListRationalesRequest, _ ...grpc.CallOption) (*skills.ListRationalesResponse, error) {
	for _, o := range m.OnListRationales {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockHardeningReportsClient) UpdateRationales(_ context.Context, given *skills.UpdateRationalesRequest, _ ...grpc.CallOption) (*skills.UpdateRationalesResponse, error) {
	for _, o := range m.OnUpdateRationales {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Result, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
