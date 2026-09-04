/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"google.golang.org/grpc"
)

type Clients interface {
	Skills() SkillsClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		skills: NewSkillsClient(conn),
		// conn is not set; this client struct does not own closing it.
	}
}

type clients struct {
	skills SkillsClient

	conn *grpc.ClientConn
}

func (c *clients) Skills() SkillsClient {
	return c.skills
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
