/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	v1 "chainguard.dev/sdk/proto/platform/common/v1"
)

// ExampleUIDPFilter demonstrates creating a UIDPFilter for filtering by ancestors.
func ExampleUIDPFilter() {
	filter := &v1.UIDPFilter{
		AncestorsOf: "uidp/group/abc123",
	}
	fmt.Printf("Ancestors of: %s\n", filter.GetAncestorsOf())
	// Output: Ancestors of: uidp/group/abc123
}

// ExampleUIDPFilter_descendantsOf demonstrates filtering by descendants.
func ExampleUIDPFilter_descendantsOf() {
	filter := &v1.UIDPFilter{
		DescendantsOf: "uidp/group/xyz789",
	}
	fmt.Printf("Descendants of: %s\n", filter.GetDescendantsOf())
	// Output: Descendants of: uidp/group/xyz789
}

// ExampleUIDPFilter_childrenOf demonstrates filtering by direct children.
func ExampleUIDPFilter_childrenOf() {
	filter := &v1.UIDPFilter{
		ChildrenOf: "uidp/group/parent123",
	}
	fmt.Printf("Children of: %s\n", filter.GetChildrenOf())
	// Output: Children of: uidp/group/parent123
}

// ExampleUIDPFilter_inRoot demonstrates filtering for root-level resources.
func ExampleUIDPFilter_inRoot() {
	filter := &v1.UIDPFilter{
		InRoot: true,
	}
	fmt.Printf("In root: %v\n", filter.GetInRoot())
	// Output: In root: true
}

// ExampleUIDPFilter_ids demonstrates filtering by specific IDs.
func ExampleUIDPFilter_ids() {
	filter := &v1.UIDPFilter{
		Ids: []string{"uidp/group/id1", "uidp/group/id2", "uidp/group/id3"},
	}
	fmt.Printf("Number of IDs: %d\n", len(filter.GetIds()))
	fmt.Printf("First ID: %s\n", filter.GetIds()[0])
	// Output:
	// Number of IDs: 3
	// First ID: uidp/group/id1
}

// ExampleUIDPFilter_getters demonstrates all getter methods on an empty filter.
func ExampleUIDPFilter_getters() {
	filter := &v1.UIDPFilter{}
	fmt.Printf("AncestorsOf: %q\n", filter.GetAncestorsOf())
	fmt.Printf("DescendantsOf: %q\n", filter.GetDescendantsOf())
	fmt.Printf("ChildrenOf: %q\n", filter.GetChildrenOf())
	fmt.Printf("InRoot: %v\n", filter.GetInRoot())
	fmt.Printf("Ids: %v\n", filter.GetIds())
	// Output:
	// AncestorsOf: ""
	// DescendantsOf: ""
	// ChildrenOf: ""
	// InRoot: false
	// Ids: []
}

// ExampleUIDPFilter_combined demonstrates a filter with multiple criteria.
func ExampleUIDPFilter_combined() {
	filter := &v1.UIDPFilter{
		DescendantsOf: "uidp/group/root",
		Ids:           []string{"uidp/group/specific1", "uidp/group/specific2"},
	}
	fmt.Printf("Descendants of: %s\n", filter.GetDescendantsOf())
	fmt.Printf("Specific IDs: %v\n", filter.GetIds())
	// Output:
	// Descendants of: uidp/group/root
	// Specific IDs: [uidp/group/specific1 uidp/group/specific2]
}
