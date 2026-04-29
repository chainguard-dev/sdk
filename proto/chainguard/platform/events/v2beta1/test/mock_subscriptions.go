/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	events "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

var _ events.SubscriptionsServiceClient = (*MockSubscriptionsServiceClient)(nil)

type MockSubscriptionsServiceClient struct {
	events.SubscriptionsServiceClient
	T *testing.T

	OnGetSubscription    []test.On[*events.GetSubscriptionRequest, *events.Subscription]
	OnCreateSubscription []test.On[*events.CreateSubscriptionRequest, *events.Subscription]
	OnListSubscriptions  []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]
	OnDeleteSubscription []test.On[*events.DeleteSubscriptionRequest, *emptypb.Empty]
}

func (m MockSubscriptionsServiceClient) GetSubscription(_ context.Context, given *events.GetSubscriptionRequest, _ ...grpc.CallOption) (*events.Subscription, error) {
	return test.Match(m.T, m.OnGetSubscription, given, "get-subscription", protocmp.Transform())
}

func (m MockSubscriptionsServiceClient) CreateSubscription(_ context.Context, given *events.CreateSubscriptionRequest, _ ...grpc.CallOption) (*events.Subscription, error) {
	return test.Match(m.T, m.OnCreateSubscription, given, "create-subscription", protocmp.Transform())
}

func (m MockSubscriptionsServiceClient) ListSubscriptions(_ context.Context, given *events.ListSubscriptionsRequest, _ ...grpc.CallOption) (*events.ListSubscriptionsResponse, error) {
	return test.Match(m.T, m.OnListSubscriptions, given, "list-subscriptions", protocmp.Transform())
}

func (m MockSubscriptionsServiceClient) DeleteSubscription(_ context.Context, given *events.DeleteSubscriptionRequest, _ ...grpc.CallOption) (*emptypb.Empty, error) {
	return test.Match(m.T, m.OnDeleteSubscription, given, "delete-subscription", protocmp.Transform())
}
