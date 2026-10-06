/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"google.golang.org/grpc"
)

type Clients interface {
	Telemetry() TelemetryServiceClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		telemetry: NewTelemetryServiceClient(conn),
		// conn is not set; this client struct does not own closing it.
	}
}

type clients struct {
	telemetry TelemetryServiceClient

	conn *grpc.ClientConn
}

func (c *clients) Telemetry() TelemetryServiceClient {
	return c.telemetry
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
