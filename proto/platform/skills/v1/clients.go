/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"google.golang.org/grpc"
)

type Clients interface {
	Entitlements() EntitlementsClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		entitlements: NewEntitlementsClient(conn),
		// conn is not set; this client struct does not own closing it.
	}
}

type clients struct {
	entitlements EntitlementsClient

	conn *grpc.ClientConn
}

func (c *clients) Entitlements() EntitlementsClient {
	return c.entitlements
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
