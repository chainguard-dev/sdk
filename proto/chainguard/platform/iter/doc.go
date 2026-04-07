/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package iter provides generic pagination iterators for v2alpha1 APIs.

# Overview

This package simplifies working with paginated API responses by providing
generic iterator functions that handle cursor-based pagination automatically.
Instead of manually managing page tokens and making multiple API calls, you
can use the iterators to process items one at a time in a for-range loop.

# Features

  - Generic pagination support for any v2alpha1 List endpoint
  - Automatic cursor-based pagination handling
  - Context-aware cancellation support
  - Clean error handling with resource-specific error messages
  - Zero-allocation iteration using Go 1.23+ iter.Seq2

# Usage

The primary function is List, which creates an iterator for any paginated
API endpoint. You provide a fetch function that retrieves a single page,
and List handles the pagination logic:

	for group, err := range iter.List(ctx, "groups", func(pageToken string) ([]*v2alpha1.Group, string, error) {
		resp, err := client.ListGroups(ctx, &v2alpha1.ListGroupsRequest{
			PageToken: pageToken,
			PageSize:  iter.DefaultPageSize,
		})
		if err != nil {
			return nil, "", err
		}
		return resp.GetGroups(), resp.GetNextPageToken(), nil
	}) {
		if err != nil {
			return fmt.Errorf("iterating groups: %w", err)
		}
		// Process group
	}

The iterator automatically:
  - Fetches pages as needed
  - Checks context cancellation between pages
  - Wraps errors with the resource name for better diagnostics
  - Stops when there are no more pages

# Integration Patterns

The package is designed to work with v2alpha1 API clients. The typical
integration pattern is to wrap the client's List method in a fetch function
that extracts the items and next page token from the response.

For services with multiple list endpoints, create helper functions:

	func ListGroups(ctx context.Context, client *Client, req *v2alpha1.ListGroupsRequest) iter.Seq2[*v2alpha1.Group, error] {
		return iter.List(ctx, "groups", func(pageToken string) ([]*v2alpha1.Group, string, error) {
			req.PageToken = pageToken
			resp, err := client.ListGroups(ctx, req)
			if err != nil {
				return nil, "", err
			}
			return resp.GetGroups(), resp.GetNextPageToken(), nil
		})
	}

This pattern allows callers to use a simple for-range loop while the
pagination complexity is hidden in the helper function.

# Error Handling

Errors are yielded as the second value in the iterator. When an error occurs,
the iterator stops and yields a zero value for the item along with the error.
Callers should check for errors in the loop:

	for item, err := range iter.List(...) {
		if err != nil {
			return err
		}
		// Process item
	}

Context cancellation is checked between pages, allowing long-running iterations
to be interrupted gracefully.
*/
package iter
