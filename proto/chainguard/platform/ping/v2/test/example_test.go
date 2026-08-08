/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	ping "chainguard.dev/sdk/proto/chainguard/platform/ping/v2"
	pingtest "chainguard.dev/sdk/proto/chainguard/platform/ping/v2/test"
)

func ExampleMockPingServiceClients() {
	mock := &pingtest.MockPingServiceClients{
		PingServiceClient: pingtest.MockPingServiceClient{
			OnPing: pingtest.PingCall{
				Given:    &ping.PingRequest{},
				Response: &ping.PingResponse{},
			},
		},
	}

	resp, err := mock.PingService().Ping(context.Background(), &ping.PingRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp != nil)
	// Output: true
}
