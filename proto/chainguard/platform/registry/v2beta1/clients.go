/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"context"
	"iter"

	"google.golang.org/grpc"

	v2iter "chainguard.dev/sdk/proto/chainguard/platform/iter"
)

// Clients provides access to v2beta1 Registry service clients.
type Clients interface {
	ReposService() ReposServiceClient
	TagsService() TagsServiceClient

	// Iterator methods for pagination - Repos
	ListReposIter(ctx context.Context, req *ListReposRequest) iter.Seq2[*Repo, error]
	ListReposAll(ctx context.Context, req *ListReposRequest) ([]*Repo, error)

	// Iterator methods for pagination - Tags
	ListTagsIter(ctx context.Context, req *ListTagsRequest) iter.Seq2[*Tag, error]
	ListTagsAll(ctx context.Context, req *ListTagsRequest) ([]*Tag, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 Registry clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		reposService: NewReposServiceClient(conn),
		tagsService:  NewTagsServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	reposService ReposServiceClient
	tagsService  TagsServiceClient

	conn *grpc.ClientConn
}

func (c *clients) ReposService() ReposServiceClient {
	return c.reposService
}

func (c *clients) TagsService() TagsServiceClient {
	return c.tagsService
}

func (c *clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ListReposIter returns an iterator over repos matching the request.
func (c *clients) ListReposIter(ctx context.Context, req *ListReposRequest) iter.Seq2[*Repo, error] {
	return v2iter.Paginate(ctx, req, "repos", func(ctx context.Context, r *ListReposRequest) ([]*Repo, string, error) {
		resp, err := c.ReposService().ListRepos(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetRepos(), resp.GetNextPageToken(), nil
	})
}

// ListReposAll fetches all repos matching the request by automatically handling pagination.
// For large result sets, consider using ListReposIter directly to process items incrementally.
func (c *clients) ListReposAll(ctx context.Context, req *ListReposRequest) ([]*Repo, error) {
	return v2iter.All(c.ListReposIter(ctx, req))
}

// ListTagsIter returns an iterator over tags matching the request.
func (c *clients) ListTagsIter(ctx context.Context, req *ListTagsRequest) iter.Seq2[*Tag, error] {
	return v2iter.Paginate(ctx, req, "tags", func(ctx context.Context, r *ListTagsRequest) ([]*Tag, string, error) {
		resp, err := c.TagsService().ListTags(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetTags(), resp.GetNextPageToken(), nil
	})
}

// ListTagsAll fetches all tags matching the request by automatically handling pagination.
// For large result sets, consider using ListTagsIter directly to process items incrementally.
func (c *clients) ListTagsAll(ctx context.Context, req *ListTagsRequest) ([]*Tag, error) {
	return v2iter.All(c.ListTagsIter(ctx, req))
}
