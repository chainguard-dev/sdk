/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package pathtree_test

import (
	"fmt"

	"chainguard.dev/sdk/pathtree"
)

// ExampleNew demonstrates creating a new Tree.
func ExampleNew() {
	t := pathtree.New()
	fmt.Println(t != nil)
	// Output:
	// true
}

// ExampleTree_Add demonstrates adding a value to the tree.
func ExampleTree_Add() {
	t := pathtree.New()
	err := t.Add("/foo/bar", "my-value", "my-label")
	fmt.Println(err)
	// Output:
	// <nil>
}

// ExampleTree_Get demonstrates looking up a value in the tree.
func ExampleTree_Get() {
	t := pathtree.New()
	_ = t.Add("/foo/bar", "my-value", "my-label")
	val, err := t.Get("/foo/bar")
	fmt.Println(err)
	fmt.Println(val)
	// Output:
	// <nil>
	// my-value
}
