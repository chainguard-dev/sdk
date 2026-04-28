/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"fmt"

	events "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1"
	eventstest "chainguard.dev/sdk/proto/chainguard/platform/events/v2beta1/test"
	"chainguard.dev/sdk/proto/chainguard/platform/test"
)

func ExampleMockClients() {
	mock := &eventstest.MockClients{
		SubscriptionsServiceClient: eventstest.MockSubscriptionsServiceClient{
			OnListSubscriptions: []test.On[*events.ListSubscriptionsRequest, *events.ListSubscriptionsResponse]{{
				Given: &events.ListSubscriptionsRequest{},
				Result: &events.ListSubscriptionsResponse{
					Subscriptions: []*events.Subscription{
						{Uid: "sub-123", Sink: "https://example.com/webhook"},
					},
				},
			}},
		},
	}

	resp, err := mock.SubscriptionsService().ListSubscriptions(context.Background(), &events.ListSubscriptionsRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Subscriptions[0].Sink)
	// Output: https://example.com/webhook
}
