/*
Copyright 2021 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
	"fmt"
	"net/url"
	"time"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"google.golang.org/grpc"
	"knative.dev/pkg/logging"

	"chainguard.dev/sdk/pkg/auth"
)

type Clients interface {
	Clusters() ClustersClient
	Records() RecordsClient
	RecordContexts() RecordContextsClient
	Sboms() SbomsClient
	Risks() RisksClient
	Signatures() SignaturesClient
	PolicyResults() PolicyResultsClient
	VulnReports() VulnReportsClient

	Nodes() NodesClient
	Namespaces() NamespacesClient
	Workloads() WorkloadsClient

	Close() error
}

func NewClients(ctx context.Context, addr string, token string) (Clients, error) {
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tenant service address, must be a url: %w", err)
	}

	target, opts := delegate.GRPCOptions(*uri)

	// TODO: we may want to require transport security at some future point.
	if cred := auth.NewFromToken(ctx, token, false); cred != nil {
		opts = append(opts, grpc.WithPerRPCCredentials(cred))
	} else {
		logging.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	}

	var cancel context.CancelFunc
	if _, timeoutSet := ctx.Deadline(); !timeoutSet {
		ctx, cancel = context.WithTimeout(ctx, 300*time.Second)
		defer cancel()
	}
	conn, err := grpc.DialContext(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the iam server: %w", err)
	}

	return &clients{
		clusters:       NewClustersClient(conn),
		records:        NewRecordsClient(conn),
		recordContexts: NewRecordContextsClient(conn),
		sboms:          NewSbomsClient(conn),
		vulnReports:    NewVulnReportsClient(conn),
		risks:          NewRisksClient(conn),
		signatures:     NewSignaturesClient(conn),
		nodes:          NewNodesClient(conn),
		namespaces:     NewNamespacesClient(conn),
		workloads:      NewWorkloadsClient(conn),
		policyResults:  NewPolicyResultsClient(conn),

		conn: conn,
	}, nil
}

func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		clusters:       NewClustersClient(conn),
		records:        NewRecordsClient(conn),
		recordContexts: NewRecordContextsClient(conn),
		sboms:          NewSbomsClient(conn),
		vulnReports:    NewVulnReportsClient(conn),
		risks:          NewRisksClient(conn),
		signatures:     NewSignaturesClient(conn),
		policyResults:  NewPolicyResultsClient(conn),
		nodes:          NewNodesClient(conn),
		namespaces:     NewNamespacesClient(conn),
		workloads:      NewWorkloadsClient(conn),
		// conn is not set, this client struct does not own closing it.
	}
}

type clients struct {
	clusters ClustersClient
	records  RecordsClient

	recordContexts RecordContextsClient
	sboms          SbomsClient
	risks          RisksClient
	signatures     SignaturesClient
	policyResults  PolicyResultsClient
	vulnReports    VulnReportsClient

	nodes      NodesClient
	namespaces NamespacesClient
	workloads  WorkloadsClient

	conn *grpc.ClientConn
}

func (c *clients) Clusters() ClustersClient {
	return c.clusters
}

func (c *clients) Records() RecordsClient {
	return c.records
}

func (c *clients) RecordContexts() RecordContextsClient {
	return c.recordContexts
}

func (c *clients) Sboms() SbomsClient {
	return c.sboms
}

func (c *clients) Risks() RisksClient {
	return c.risks
}

func (c *clients) Signatures() SignaturesClient {
	return c.signatures
}

func (c *clients) PolicyResults() PolicyResultsClient {
	return c.policyResults
}

func (c *clients) Nodes() NodesClient {
	return c.nodes
}

func (c *clients) Namespaces() NamespacesClient {
	return c.namespaces
}

func (c *clients) Workloads() WorkloadsClient {
	return c.workloads
}

func (c *clients) VulnReports() VulnReportsClient {
	return c.vulnReports
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
