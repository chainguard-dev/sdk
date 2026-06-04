/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"google.golang.org/grpc"
)

type Clients interface {
	ArgosDocuments() ArgosDocumentsClient
	ArgosOSV() ArgosOSVClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		argosDocuments: NewArgosDocumentsClient(conn),
		argosOSV:       NewArgosOSVClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	argosDocuments ArgosDocumentsClient
	argosOSV       ArgosOSVClient

	conn *grpc.ClientConn
}

func (c *clients) ArgosDocuments() ArgosDocumentsClient {
	return c.argosDocuments
}

func (c *clients) ArgosOSV() ArgosOSVClient {
	return c.argosOSV
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
