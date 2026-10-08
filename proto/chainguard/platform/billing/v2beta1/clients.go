/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"google.golang.org/grpc"
)

// Clients provides access to v2beta1 Billing service clients.
type Clients interface {
	BillingService() BillingServiceClient

	Close() error
}

// NewClientsFromConnection creates v2beta1 Billing clients from an existing
// gRPC connection. The returned Clients does not own conn and will not close
// it; callers are responsible for closing conn when done.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		billing: NewBillingServiceClient(conn),
	}
}

type clients struct {
	billing BillingServiceClient
}

var _ Clients = (*clients)(nil)

func (c *clients) BillingService() BillingServiceClient {
	return c.billing
}

func (c *clients) Close() error {
	return nil
}
