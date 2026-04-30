/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"iter"

	events "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

type MockClients struct {
	OnClose error

	SubscriptionsServiceClient MockSubscriptionsServiceClient
}

// Close implements [v2beta1.Clients].
func (m *MockClients) Close() error {
	return m.OnClose
}

// SubscriptionsService implements [v2beta1.Clients].
func (m *MockClients) SubscriptionsService() events.SubscriptionsServiceClient {
	return &m.SubscriptionsServiceClient
}

// ListSubscriptionsAll implements [v2beta1.Clients].
func (m *MockClients) ListSubscriptionsAll(ctx context.Context, req *events.ListSubscriptionsRequest) ([]*events.Subscription, error) {
	resp, err := m.SubscriptionsServiceClient.ListSubscriptions(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.GetSubscriptions(), nil
}

// ListSubscriptionsIter implements [v2beta1.Clients].
func (m *MockClients) ListSubscriptionsIter(ctx context.Context, req *events.ListSubscriptionsRequest) iter.Seq2[*events.Subscription, error] {
	return test.MockIter(m.ListSubscriptionsAll(ctx, req))
}

var _ events.Clients = (*MockClients)(nil)
