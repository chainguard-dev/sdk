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
	ImagesService() ImagesServiceClient
	OverlaysService() OverlaysServiceClient
	OverlayBindingsService() OverlayBindingsServiceClient

	// Iterator methods for pagination - Repos
	ListReposIter(ctx context.Context, req *ListReposRequest) iter.Seq2[*Repo, error]
	ListReposAll(ctx context.Context, req *ListReposRequest) ([]*Repo, error)

	// Iterator methods for pagination - Tags
	ListTagsIter(ctx context.Context, req *ListTagsRequest) iter.Seq2[*Tag, error]
	ListTagsAll(ctx context.Context, req *ListTagsRequest) ([]*Tag, error)

	// Iterator methods for pagination - Overlays
	ListOverlaysIter(ctx context.Context, req *ListOverlaysRequest) iter.Seq2[*Overlay, error]
	ListOverlaysAll(ctx context.Context, req *ListOverlaysRequest) ([]*Overlay, error)

	// Iterator methods for pagination - OverlayBindings
	ListOverlayBindingsIter(ctx context.Context, req *ListOverlayBindingsRequest) iter.Seq2[*OverlayBinding, error]
	ListOverlayBindingsAll(ctx context.Context, req *ListOverlayBindingsRequest) ([]*OverlayBinding, error)

	Close() error
}

// NewClientsFromConnection creates v2beta1 Registry clients from an existing gRPC connection.
func NewClientsFromConnection(conn *grpc.ClientConn) Clients {
	return &clients{
		reposService:           NewReposServiceClient(conn),
		tagsService:            NewTagsServiceClient(conn),
		imagesService:          NewImagesServiceClient(conn),
		overlaysService:        NewOverlaysServiceClient(conn),
		overlayBindingsService: NewOverlayBindingsServiceClient(conn),
		// conn is not set, this client struct does not own closing it
	}
}

type clients struct {
	reposService           ReposServiceClient
	tagsService            TagsServiceClient
	imagesService          ImagesServiceClient
	overlaysService        OverlaysServiceClient
	overlayBindingsService OverlayBindingsServiceClient

	conn *grpc.ClientConn
}

func (c *clients) ReposService() ReposServiceClient {
	return c.reposService
}

func (c *clients) TagsService() TagsServiceClient {
	return c.tagsService
}

func (c *clients) ImagesService() ImagesServiceClient {
	return c.imagesService
}

func (c *clients) OverlaysService() OverlaysServiceClient {
	return c.overlaysService
}

func (c *clients) OverlayBindingsService() OverlayBindingsServiceClient {
	return c.overlayBindingsService
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

// ListOverlaysIter returns an iterator over overlays matching the request.
func (c *clients) ListOverlaysIter(ctx context.Context, req *ListOverlaysRequest) iter.Seq2[*Overlay, error] {
	return v2iter.Paginate(ctx, req, "overlays", func(ctx context.Context, r *ListOverlaysRequest) ([]*Overlay, string, error) {
		resp, err := c.OverlaysService().ListOverlays(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetOverlays(), resp.GetNextPageToken(), nil
	})
}

// ListOverlaysAll fetches all overlays matching the request by automatically handling pagination.
// For large result sets, consider using ListOverlaysIter directly to process items incrementally.
func (c *clients) ListOverlaysAll(ctx context.Context, req *ListOverlaysRequest) ([]*Overlay, error) {
	return v2iter.All(c.ListOverlaysIter(ctx, req))
}

// ListOverlayBindingsIter returns an iterator over overlay bindings matching the request.
func (c *clients) ListOverlayBindingsIter(ctx context.Context, req *ListOverlayBindingsRequest) iter.Seq2[*OverlayBinding, error] {
	return v2iter.Paginate(ctx, req, "overlay bindings", func(ctx context.Context, r *ListOverlayBindingsRequest) ([]*OverlayBinding, string, error) {
		resp, err := c.OverlayBindingsService().ListOverlayBindings(ctx, r)
		if err != nil {
			return nil, "", err
		}
		return resp.GetOverlayBindings(), resp.GetNextPageToken(), nil
	})
}

// ListOverlayBindingsAll fetches all overlay bindings matching the request by automatically handling pagination.
// For large result sets, consider using ListOverlayBindingsIter directly to process items incrementally.
func (c *clients) ListOverlayBindingsAll(ctx context.Context, req *ListOverlayBindingsRequest) ([]*OverlayBinding, error) {
	return v2iter.All(c.ListOverlayBindingsIter(ctx, req))
}
