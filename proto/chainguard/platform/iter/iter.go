/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package iter provides generic pagination iterators for v2alpha1 APIs.
package iter //nolint:revive // redefines-builtin-id: collides with stdlib iter, but renaming would break API

import (
	"context"
	"fmt"
	"iter"
	"reflect"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const (
	// DefaultPageSize is the default number of items per page when not specified.
	// TODO(colin): find a more appropriate place for this
	DefaultPageSize = 50
)

// PagedRequest is a constraint for protobuf list request types that support
// cursor-based pagination via page_size and page_token fields.
type PagedRequest interface {
	proto.Message
	GetPageSize() int32
	GetPageToken() string
}

// Paginate creates a paginated iterator from a list RPC. It handles nil request
// guards, cloning the request to avoid caller mutation, defaulting page size,
// and setting page_token/page_size via protobuf reflection on each page fetch.
//
// The fetch callback receives the prepared request and returns the page of results,
// the next page token, and any error.
//
// Usage:
//
//	return v2iter.Paginate(ctx, req, "groups", func(ctx context.Context, r *ListGroupsRequest) ([]*Group, string, error) {
//	    resp, err := c.GroupsService().ListGroups(ctx, r)
//	    if err != nil {
//	        return nil, "", err
//	    }
//	    return resp.GetGroups(), resp.GetNextPageToken(), nil
//	})
func Paginate[Req PagedRequest, Res any](ctx context.Context, req Req, resourceName string, fetch func(context.Context, Req) ([]Res, string, error)) iter.Seq2[Res, error] {
	if reflect.ValueOf(req).IsNil() {
		req = reflect.New(reflect.TypeOf(req).Elem()).Interface().(Req)
	}
	r := proto.Clone(req).(Req)
	pageSize := r.GetPageSize()
	if pageSize == 0 {
		pageSize = DefaultPageSize
	}
	setPageFields(r, pageSize, "")
	return List(ctx, resourceName, func(pageToken string) ([]Res, string, error) {
		setPageFields(r, pageSize, pageToken)
		return fetch(ctx, r)
	})
}

// setPageFields sets the page_size and page_token fields on a protobuf message
// using protobuf reflection.
func setPageFields(msg proto.Message, pageSize int32, pageToken string) {
	m := msg.ProtoReflect()
	fields := m.Descriptor().Fields()
	if f := fields.ByName("page_size"); f != nil {
		m.Set(f, protoreflect.ValueOfInt32(pageSize))
	}
	if f := fields.ByName("page_token"); f != nil {
		m.Set(f, protoreflect.ValueOfString(pageToken))
	}
}

// All collects all items from a paginated iterator into a slice.
//
// Usage:
//
//	items, err := v2iter.All(c.ListGroupsIter(ctx, req))
func All[T any](seq iter.Seq2[T, error]) ([]T, error) {
	var result []T
	for item, err := range seq {
		if err != nil {
			return nil, err
		}
		result = append(result, item)
	}
	return result, nil
}

// List is a low-level generic pagination iterator for any v2alpha1 List endpoint.
// It handles cursor-based pagination automatically, yielding items one at a time.
// Most callers should use Paginate instead, which also handles request setup.
func List[T any](ctx context.Context, resourceName string, fetch func(pageToken string) ([]T, string, error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		var pageToken string
		for {
			if err := ctx.Err(); err != nil {
				yield(zero, err)
				return
			}

			items, nextToken, err := fetch(pageToken)
			if err != nil {
				yield(zero, fmt.Errorf("listing %q: %w", resourceName, err))
				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			if nextToken == "" {
				return
			}
			pageToken = nextToken
		}
	}
}
