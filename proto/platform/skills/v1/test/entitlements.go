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

	skills "chainguard.dev/sdk/proto/platform/skills/v1"
)

var _ skills.EntitlementsClient = (*MockEntitlementsClient)(nil)

type MockEntitlementsClient struct {
	skills.EntitlementsClient

	OnCreate []EntitlementsOnCreate
	OnDelete []EntitlementsOnDelete
	OnList   []EntitlementsOnList
}

type EntitlementsOnCreate struct {
	Given   *skills.CreateSkillEntitlementRequest
	Created *skills.SkillEntitlement
	Error   error
}

type EntitlementsOnDelete struct {
	Given *skills.DeleteSkillEntitlementRequest
	Error error
}

type EntitlementsOnList struct {
	Given *skills.SkillEntitlementFilter
	List  *skills.SkillEntitlementList
	Error error
}

func (m MockEntitlementsClient) Create(_ context.Context, given *skills.CreateSkillEntitlementRequest, _ ...grpc.CallOption) (*skills.SkillEntitlement, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockEntitlementsClient) Delete(_ context.Context, given *skills.DeleteSkillEntitlementRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockEntitlementsClient) List(_ context.Context, given *skills.SkillEntitlementFilter, _ ...grpc.CallOption) (*skills.SkillEntitlementList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
