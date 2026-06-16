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

	libraries "chainguard.dev/sdk/proto/platform/libraries/v1"
)

var _ libraries.AWSMarketplaceSubscriptionsClient = (*MockAWSMarketplaceSubscriptionsClient)(nil)

type MockAWSMarketplaceSubscriptionsClient struct {
	libraries.AWSMarketplaceSubscriptionsClient

	OnCreate []AWSMarketplaceSubscriptionsOnCreate
	OnGet    []AWSMarketplaceSubscriptionsOnGet
	OnList   []AWSMarketplaceSubscriptionsOnList
	OnCancel []AWSMarketplaceSubscriptionsOnCancel
}

type AWSMarketplaceSubscriptionsOnCreate struct {
	Given   *libraries.CreateAWSMarketplaceSubscriptionRequest
	Created *libraries.AWSMarketplaceSubscription
	Error   error
}

type AWSMarketplaceSubscriptionsOnGet struct {
	Given *libraries.GetAWSMarketplaceSubscriptionRequest
	Got   *libraries.AWSMarketplaceSubscription
	Error error
}

type AWSMarketplaceSubscriptionsOnList struct {
	Given *libraries.AWSMarketplaceSubscriptionFilter
	List  *libraries.AWSMarketplaceSubscriptionList
	Error error
}

type AWSMarketplaceSubscriptionsOnCancel struct {
	Given     *libraries.CancelAWSMarketplaceSubscriptionRequest
	Cancelled *libraries.AWSMarketplaceSubscription
	Error     error
}

func (m MockAWSMarketplaceSubscriptionsClient) Create(_ context.Context, given *libraries.CreateAWSMarketplaceSubscriptionRequest, _ ...grpc.CallOption) (*libraries.AWSMarketplaceSubscription, error) {
	for _, o := range m.OnCreate {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAWSMarketplaceSubscriptionsClient) Get(_ context.Context, given *libraries.GetAWSMarketplaceSubscriptionRequest, _ ...grpc.CallOption) (*libraries.AWSMarketplaceSubscription, error) {
	for _, o := range m.OnGet {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Got, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAWSMarketplaceSubscriptionsClient) List(_ context.Context, given *libraries.AWSMarketplaceSubscriptionFilter, _ ...grpc.CallOption) (*libraries.AWSMarketplaceSubscriptionList, error) {
	for _, o := range m.OnList {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockAWSMarketplaceSubscriptionsClient) Cancel(_ context.Context, given *libraries.CancelAWSMarketplaceSubscriptionRequest, _ ...grpc.CallOption) (*libraries.AWSMarketplaceSubscription, error) {
	for _, o := range m.OnCancel {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Cancelled, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
