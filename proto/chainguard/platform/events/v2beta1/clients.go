/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"google.golang.org/grpc"
)

// Clients provides access to v2beta1 Events service clients.
type Clients interface {
	SubscriptionsService() SubscriptionsServiceClient

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
