/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"google.golang.org/grpc"
)

type Clients interface {
	Guardener() GuardenerClient
	GitHubAssociations() GitHubAssociationsClient

	Close() error
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		guardener:          NewGuardenerClient(conn),
		githubAssociations: NewGitHubAssociationsClient(conn),
		// conn is not set; this client struct does not own closing it.
	}
}

type clients struct {
	guardener          GuardenerClient
	githubAssociations GitHubAssociationsClient

	conn *grpc.ClientConn
}

func (c *clients) Guardener() GuardenerClient {
	return c.guardener
}

func (c *clients) GitHubAssociations() GitHubAssociationsClient {
	return c.githubAssociations
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
