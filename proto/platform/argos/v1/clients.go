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
	ArgosVEX() ArgosVEXClient
	ArgosVulns() ArgosVulnsClient
	ArgosPartners() ArgosPartnersClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		argosDocuments: NewArgosDocumentsClient(conn),
		argosOSV:       NewArgosOSVClient(conn),
		argosVEX:       NewArgosVEXClient(conn),
		argosVulns:     NewArgosVulnsClient(conn),
		argosPartners:  NewArgosPartnersClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	argosDocuments ArgosDocumentsClient
	argosOSV       ArgosOSVClient
	argosVEX       ArgosVEXClient
	argosVulns     ArgosVulnsClient
	argosPartners  ArgosPartnersClient

	conn *grpc.ClientConn
}

func (c *clients) ArgosDocuments() ArgosDocumentsClient {
	return c.argosDocuments
}

func (c *clients) ArgosOSV() ArgosOSVClient {
	return c.argosOSV
}

func (c *clients) ArgosVEX() ArgosVEXClient {
	return c.argosVEX
}

func (c *clients) ArgosVulns() ArgosVulnsClient {
	return c.argosVulns
}

func (c *clients) ArgosPartners() ArgosPartnersClient {
	return c.argosPartners
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
