/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"fmt"
	"iter"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func missingMock[T any](t *testing.T, verb string, given any) (T, error) { //nolint: unparam
	var zero T
	if t != nil {
		t.Helper()
		t.Errorf("%s mock not found for %v", verb, given)
	}
	return zero, fmt.Errorf("%s mock not found for %v", verb, given)
}

type On[G any, R any] struct {
	Given  G
	Result R
	Error  error
}

func Match[G any, R any](t *testing.T, patterns []On[G, R], given G, verb string, opts ...cmp.Option) (R, error) {
	for _, o := range patterns {
		if cmp.Equal(o.Given, given, opts...) {
			return o.Result, o.Error
		}
	}
	if t != nil {
		for i, o := range patterns {
			t.Logf("skipping %s pattern %d mismatched: %s", verb, i, cmp.Diff(o.Given, given, opts...))
		}
	}
	return missingMock[R](t, verb, given)
}

// MockIter converts a slice and error into an [iter.Seq2], for use in
// mock client List*Iter implementations that delegate to List*All.
func MockIter[T any](items []T, err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		if err != nil {
			var zero T
			yield(zero, err)
			return
		}
		for _, item := range items {
			if !yield(item, nil) {
				return
			}
		}
	}
}
