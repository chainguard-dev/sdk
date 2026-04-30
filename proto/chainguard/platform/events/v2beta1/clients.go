/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"context"
	"iter"

	"google.golang.org/grpc"

	v2iter "chainguard.dev/sdk/proto/chainguard/platform/iter"
)

// Clients provides access to v2beta1 Events service clients.
type Clients interface {
	SubscriptionsService() SubscriptionsServiceClient

	// Iterator methods for pagination - Subscriptions
	ListSubscriptionsIter(ctx context.Context, req *ListSubscriptionsRequest) iter.Seq2[*Subscription, error]
	ListSubscriptionsAll(ctx context.Context, req *ListSubscriptionsRequest) ([]*Subscription, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 Events clients from an existing gRPC connection.
// The returned Clients does not own conn and will not close it; callers are
// responsible for closing conn when done.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		subscriptions: NewSubscriptionsServiceClient(conn),
	}
}

type clients struct {
	subscriptions SubscriptionsServiceClient
}

var _ Clients = (*clients)(nil)

func (c *clients) SubscriptionsService() SubscriptionsServiceClient {
	return c.subscriptions
}

func (c *clients) Close() error {
	return nil
}

// ListSubscriptionsIter returns an iterator over subscriptions matching the request.
func (c *clients) ListSubscriptionsIter(ctx context.Context, req *ListSubscriptionsRequest) iter.Seq2[*Subscription, error] {
	return v2iter.Paginate(ctx, req, "subscriptions", func(ctx context.Context, r *ListSubscriptionsRequest) ([]*Subscription, string, error) {
		resp, err := c.SubscriptionsService().ListSubscriptions(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetSubscriptions(), resp.GetNextPageToken(), nil
	})
}

// ListSubscriptionsAll fetches all subscriptions matching the request by automatically handling pagination.
// For large result sets, consider using ListSubscriptionsIter directly to process items incrementally.
func (c *clients) ListSubscriptionsAll(ctx context.Context, req *ListSubscriptionsRequest) ([]*Subscription, error) {
	return v2iter.All(c.ListSubscriptionsIter(ctx, req))
}
