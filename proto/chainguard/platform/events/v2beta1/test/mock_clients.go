/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	events "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1"
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

var _ events.Clients = (*MockClients)(nil)
