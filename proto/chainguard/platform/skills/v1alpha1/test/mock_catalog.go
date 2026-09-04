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

var _ skills.SkillsClient = (*MockCatalogClient)(nil)

type MockCatalogClient struct {
	OnListSkills   []SkillsOnList
	OnSearchSkills []SkillsOnSearch
	OnUpdateSkill  []SkillsOnUpdate
	OnDeleteSkill  []SkillsOnDelete
}

type SkillsOnList struct {
	Given *skills.ListSkillsRequest
	List  *skills.ListSkillsResponse
	Error error
}

type SkillsOnSearch struct {
	Given  *skills.SearchSkillsRequest
	Result *skills.SearchSkillsResponse
	Error  error
}

type SkillsOnUpdate struct {
	Given   *skills.UpdateSkillRequest
	Updated *skills.Skill
	Error   error
}

type SkillsOnDelete struct {
	Given *skills.DeleteSkillRequest
	Error error
}

func (m MockCatalogClient) ListSkills(_ context.Context, given *skills.ListSkillsRequest, _ ...grpc.CallOption) (*skills.ListSkillsResponse, error) {
	for _, o := range m.OnListSkills {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockCatalogClient) SearchSkills(_ context.Context, given *skills.SearchSkillsRequest, _ ...grpc.CallOption) (*skills.SearchSkillsResponse, error) {
	for _, o := range m.OnSearchSkills {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Result, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockCatalogClient) UpdateSkill(_ context.Context, given *skills.UpdateSkillRequest, _ ...grpc.CallOption) (*skills.Skill, error) {
	for _, o := range m.OnUpdateSkill {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockCatalogClient) DeleteSkill(_ context.Context, given *skills.DeleteSkillRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDeleteSkill {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			if o.Error != nil {
				return nil, o.Error
			}
			return &emptypb.Empty{}, nil
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
