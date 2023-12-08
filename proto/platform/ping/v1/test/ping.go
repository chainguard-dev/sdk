/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	ping "chainguard.dev/sdk/proto/platform/ping/v1"
)

var _ ping.Clients = (*MockPingServiceClients)(nil)

type MockPingServiceClients struct {
	OnClose error

	PingServiceClient MockPingServiceClient
}

func (m MockPingServiceClients) Ping() ping.PingServiceClient {
	return &m.PingServiceClient
}

func (m MockPingServiceClients) Close() error {
	return m.OnClose
}

var _ ping.PingServiceClient = (*MockPingServiceClient)(nil)

type MockPingServiceClient struct {
	ping.PingServiceClient

	OnPing Ping
}

type Ping struct {
	Given    *ping.PingRequest
	Response *ping.Response
	Error    error
}

func (m MockPingServiceClient) Ping(_ context.Context, _ *ping.PingRequest, _ ...grpc.CallOption) (*ping.Response, error) {
	if m.OnPing.Given == nil {
		return nil, fmt.Errorf("OnPing.Given defined to be %v", m.OnPing.Given)
	}
	return &ping.Response{}, nil
}
