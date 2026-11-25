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
	HardeningReports() SkillsHardeningReportsClient
	Harden() SkillsHardenClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		skills:           NewSkillsClient(conn),
		hardeningReports: NewSkillsHardeningReportsClient(conn),
		harden:           NewSkillsHardenClient(conn),
		// conn is not set; this client struct does not own closing it.
	}
}

type clients struct {
	skills           SkillsClient
	hardeningReports SkillsHardeningReportsClient
	harden           SkillsHardenClient

	conn *grpc.ClientConn
}

func (c *clients) Skills() SkillsClient {
	return c.skills
}

func (c *clients) HardeningReports() SkillsHardeningReportsClient {
	return c.hardeningReports
}

func (c *clients) Harden() SkillsHardenClient {
	return c.harden
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
