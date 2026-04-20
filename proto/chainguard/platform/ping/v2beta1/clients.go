/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"google.golang.org/grpc"
)

// Clients provides access to v2beta1 Ping service clients.
type Clients interface {
	PingService() PingServiceClient

	Close() error
}

// NewClientsFromConnection creates v2beta1 Ping clients from an existing gRPC connection.
// The returned Clients does not own conn and will not close it; callers are
// responsible for closing conn when done.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		ping: NewPingServiceClient(conn),
	}
}

type clients struct {
	ping PingServiceClient
}

var _ Clients = (*clients)(nil)

func (c *clients) PingService() PingServiceClient {
	return c.ping
}

func (c *clients) Close() error {
	return nil
}
