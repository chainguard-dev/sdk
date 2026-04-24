/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	ping "chainguard.dev/sdk/proto/chainguard/platform/ping/v2beta1"
)

var _ ping.Clients = (*MockPingServiceClients)(nil)

// MockPingServiceClients is a test double for ping.Clients.
type MockPingServiceClients struct {
	OnClose error

	PingServiceClient MockPingServiceClient
}

func (m MockPingServiceClients) PingService() ping.PingServiceClient {
	return &m.PingServiceClient
}

func (m MockPingServiceClients) Close() error {
	return m.OnClose
}

var _ ping.PingServiceClient = (*MockPingServiceClient)(nil)

// MockPingServiceClient is a test double for ping.PingServiceClient.
type MockPingServiceClient struct {
	ping.PingServiceClient

	OnPing PingCall
}

// PingCall holds the canned response for a Ping RPC call.
type PingCall struct {
	Given    *ping.PingRequest
	Response *ping.PingResponse
	Error    error
}

func (m MockPingServiceClient) Ping(_ context.Context, in *ping.PingRequest, _ ...grpc.CallOption) (*ping.PingResponse, error) {
	if m.OnPing.Given == nil {
		return nil, fmt.Errorf("OnPing.Given not configured, got request: %v", in)
	}
	return m.OnPing.Response, m.OnPing.Error
}
