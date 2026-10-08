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

	billing "chainguard.dev/sdk/proto/chainguard/platform/billing/v2beta1"
)

var _ billing.Clients = (*MockBillingClients)(nil)

// MockBillingClients is a test double for billing.Clients.
type MockBillingClients struct {
	OnClose error

	BillingServiceClient MockBillingServiceClient
}

func (m MockBillingClients) BillingService() billing.BillingServiceClient {
	return &m.BillingServiceClient
}

func (m MockBillingClients) Close() error {
	return m.OnClose
}

var _ billing.BillingServiceClient = (*MockBillingServiceClient)(nil)

// MockBillingServiceClient is a test double for billing.BillingServiceClient.
// Each call matches the request against the Given of its On* entries and
// returns the first match's canned result.
type MockBillingServiceClient struct {
	OnCreateCheckoutSession []CheckoutSessionCall
	OnCreatePortalSession   []PortalSessionCall
}

// CheckoutSessionCall is a canned CreateCheckoutSession result.
type CheckoutSessionCall struct {
	Given    *billing.CreateCheckoutSessionRequest
	Response *billing.CheckoutSession
	Error    error
}

// PortalSessionCall is a canned CreatePortalSession result.
type PortalSessionCall struct {
	Given    *billing.CreatePortalSessionRequest
	Response *billing.PortalSession
	Error    error
}

func (m MockBillingServiceClient) CreateCheckoutSession(_ context.Context, in *billing.CreateCheckoutSessionRequest, _ ...grpc.CallOption) (*billing.CheckoutSession, error) {
	for _, c := range m.OnCreateCheckoutSession {
		if cmp.Equal(in, c.Given, protocmp.Transform()) {
			return c.Response, c.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", in)
}

func (m MockBillingServiceClient) CreatePortalSession(_ context.Context, in *billing.CreatePortalSessionRequest, _ ...grpc.CallOption) (*billing.PortalSession, error) {
	for _, c := range m.OnCreatePortalSession {
		if cmp.Equal(in, c.Given, protocmp.Transform()) {
			return c.Response, c.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", in)
}
