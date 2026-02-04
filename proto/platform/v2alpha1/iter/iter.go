/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package iter provides generic pagination iterators for v2alpha1 APIs.
package iter

import (
	"context"
	"fmt"
	"iter"
)

// List is a generic pagination iterator for any v2alpha1 List endpoint.
// It handles cursor-based pagination automatically, yielding items one at a time.
//
// Parameters:
//   - ctx: Context for cancellation
//   - resourceName: Human-readable name for error messages (e.g., "groups", "identities")
//   - fetch: Function that fetches a page given a page token, returning items, next token, and error
//
// Usage:
//
//	return v2iter.List(ctx, "groups", func(pageToken string) ([]*Group, string, error) {
//	    r.PageToken = pageToken
//	    resp, err := c.GroupsService().ListGroups(ctx, r)
//	    if err != nil {
//	        return nil, "", err
//	    }
//	    return resp.GetGroups(), resp.GetNextPageToken(), nil
//	})
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
				yield(zero, fmt.Errorf("listing %s: %w", resourceName, err))
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
