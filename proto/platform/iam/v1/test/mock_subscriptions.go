/*
Copyright 2022 Chainguard, Inc.
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

	events "chainguard.dev/sdk/proto/platform/events/v1"
)

var _ events.SubscriptionsClient = (*MockSubscriptionsClient)(nil)

type MockSubscriptionsClient struct {
	OnCreate []SubscriptionOnCreate
	OnUpdate []SubscriptionOnUpdate
	OnDelete []SubscriptionOnDelete
	OnList   []SubscriptionOnList
}

type SubscriptionOnCreate struct {
	Given   *events.CreateSubscriptionRequest
	Created *events.Subscription
	Error   error
}

type SubscriptionOnUpdate struct {
	Given   *events.Subscription
	Updated *events.Subscription
	Error   error
}

type SubscriptionOnDelete struct {
	Given *events.DeleteSubscriptionRequest
	Error error
}

type SubscriptionOnList struct {
	Given *events.SubscriptionFilter
	List  *events.SubscriptionList
	Error error
}

func (m MockSubscriptionsClient) Create(_ context.Context, given *events.CreateSubscriptionRequest, _ ...grpc.CallOption) (*events.Subscription, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSubscriptionsClient) Update(_ context.Context, given *events.Subscription, _ ...grpc.CallOption) (*events.Subscription, error) {
	for _, o := range m.OnUpdate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Updated, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSubscriptionsClient) Delete(_ context.Context, given *events.DeleteSubscriptionRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	for _, o := range m.OnDelete {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return &emptypb.Empty{}, o.Error
		}
	}
	return &emptypb.Empty{}, fmt.Errorf("mock not found for %v", given)
}

func (m MockSubscriptionsClient) List(_ context.Context, given *events.SubscriptionFilter, _ ...grpc.CallOption) (*events.SubscriptionList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
