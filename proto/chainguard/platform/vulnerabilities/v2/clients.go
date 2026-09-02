/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"context"
	"iter"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	v2iter "chainguard.dev/sdk/proto/chainguard/platform/iter"
)

// Clients provides access to v2 vulnerabilities service clients.
type Clients interface {
	AdvisoriesService() AdvisoriesServiceClient
	VulnReportsService() VulnReportsServiceClient

	ListAdvisoriesIter(ctx context.Context, req *ListAdvisoriesRequest) iter.Seq2[*Advisory, error]
	ListAdvisoriesAll(ctx context.Context, req *ListAdvisoriesRequest) ([]*Advisory, error)

	ListVulnCountReportsIter(ctx context.Context, req *ListVulnCountReportsRequest) iter.Seq2[*VulnCountReport, error]
	ListVulnCountReportsAll(ctx context.Context, req *ListVulnCountReportsRequest) ([]*VulnCountReport, error)

	Close() error
}

// NewClientsFromConnection creates v2 vulnerabilities clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		advisoriesService:  NewAdvisoriesServiceClient(conn),
		vulnReportsService: NewVulnReportsServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	advisoriesService  AdvisoriesServiceClient
	vulnReportsService VulnReportsServiceClient

	conn *grpc.ClientConn
}

func (c *clients) AdvisoriesService() AdvisoriesServiceClient {
	return c.advisoriesService
}

func (c *clients) VulnReportsService() VulnReportsServiceClient {
	return c.vulnReportsService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListAdvisoriesIter returns an iterator over advisories matching the request.
func (c *clients) ListAdvisoriesIter(ctx context.Context, req *ListAdvisoriesRequest) iter.Seq2[*Advisory, error] {
	return v2iter.Paginate(ctx, req, "advisories", func(ctx context.Context, r *ListAdvisoriesRequest) ([]*Advisory, string, error) {
		resp, err := c.AdvisoriesService().ListAdvisories(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetAdvisories(), resp.GetNextPageToken(), nil
	})
}

// ListAdvisoriesAll fetches all advisories matching the request by automatically handling pagination.
// For large result sets, consider using ListAdvisoriesIter directly to process items incrementally.
func (c *clients) ListAdvisoriesAll(ctx context.Context, req *ListAdvisoriesRequest) ([]*Advisory, error) {
	return v2iter.All(c.ListAdvisoriesIter(ctx, req))
}

// vulnCountReportsIterPageSize is the page size the iterators request when
// the caller leaves it unset. It matches the server's default, the value
// tuned against the datastore page byte budget, rather than the shared
// iterator's generic 50, which would exhaust the iterator's page cap on
// large listings and silently truncate an All result.
const vulnCountReportsIterPageSize = 1500

// ListVulnCountReportsIter returns an iterator over vuln count reports matching the request.
func (c *clients) ListVulnCountReportsIter(ctx context.Context, req *ListVulnCountReportsRequest) iter.Seq2[*VulnCountReport, error] {
	if req == nil {
		req = &ListVulnCountReportsRequest{}
	}
	if req.GetPageSize() == 0 {
		req = proto.Clone(req).(*ListVulnCountReportsRequest)
		req.PageSize = vulnCountReportsIterPageSize
	}
	return v2iter.Paginate(ctx, req, "vuln_count_reports", func(ctx context.Context, r *ListVulnCountReportsRequest) ([]*VulnCountReport, string, error) {
		resp, err := c.VulnReportsService().ListVulnCountReports(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetVulnCountReports(), resp.GetNextPageToken(), nil
	})
}

// ListVulnCountReportsAll fetches all vuln count reports matching the request by automatically
// handling pagination. For large result sets, consider using ListVulnCountReportsIter directly
// to process items incrementally.
func (c *clients) ListVulnCountReportsAll(ctx context.Context, req *ListVulnCountReportsRequest) ([]*VulnCountReport, error) {
	return v2iter.All(c.ListVulnCountReportsIter(ctx, req))
}
