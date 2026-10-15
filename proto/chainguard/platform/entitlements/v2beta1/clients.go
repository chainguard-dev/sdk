/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"google.golang.org/grpc"
)

// Clients provides access to v2beta1 Entitlements service clients.
type Clients interface {
	Entitlements() EntitlementsServiceClient

	Close() error
}

// NewClientsFromConnection creates v2beta1 Entitlements clients from an existing gRPC connection.
// The returned Clients does not own conn and will not close it; callers are
// responsible for closing conn when done.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		entitlements: NewEntitlementsServiceClient(conn),
	}
}

type clients struct {
	entitlements EntitlementsServiceClient
}

var _ Clients = (*clients)(nil)

func (c *clients) Entitlements() EntitlementsServiceClient {
	return c.entitlements
}

func (c *clients) Close() error {
	return nil
}
