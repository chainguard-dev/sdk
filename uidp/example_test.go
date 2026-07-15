/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uidp_test

import (
	"fmt"

	"chainguard.dev/sdk/uidp"
)

// ExampleNewUID demonstrates creating a new globally-unique UID.
func ExampleNewUID() {
	u := uidp.NewUID()
	fmt.Println(len(string(u)) == 40)
	// Output:
	// true
}

// ExampleNewSUID demonstrates creating a new scoped UID.
func ExampleNewSUID() {
	s := uidp.NewSUID()
	fmt.Println(len(string(s)) == 16)
	// Output:
	// true
}

// ExampleNewUIDP demonstrates creating a root UIDP (no parent path).
func ExampleNewUIDP() {
	p := uidp.NewUIDP("")
	fmt.Println(len(string(p)) == 40)
	// Output:
	// true
}

// ExampleUIDP_Reparent demonstrates reparenting a UIDP under a new parent.
func ExampleUIDP_Reparent() {
	child := uidp.UIDP("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/bbbbbbbbbbbbbbbb")
	parent := uidp.UIDP("cccccccccccccccccccccccccccccccccccccccc")
	result, err := child.Reparent(parent)
	fmt.Println(err)
	fmt.Println(result)
	// Output:
	// <nil>
	// cccccccccccccccccccccccccccccccccccccccc/bbbbbbbbbbbbbbbb
}

// ExampleIsAncestor demonstrates checking whether a UIDP is a strict ancestor.
func ExampleIsAncestor() {
	fmt.Println(uidp.IsAncestor("a", "a/b"))
	fmt.Println(uidp.IsAncestor("a", "a"))
	// Output:
	// true
	// false
}

// ExampleIsAncestorOrSelf demonstrates the inclusive ancestor check.
func ExampleIsAncestorOrSelf() {
	fmt.Println(uidp.IsAncestorOrSelf("a", "a"))
	fmt.Println(uidp.IsAncestorOrSelf("a", "a/b"))
	fmt.Println(uidp.IsAncestorOrSelf("a", "b"))
	// Output:
	// true
	// true
	// false
}

// ExampleInRoot demonstrates checking whether a UIDP is at the root level.
func ExampleInRoot() {
	fmt.Println(uidp.InRoot("a"))
	fmt.Println(uidp.InRoot("a/b"))
	// Output:
	// true
	// false
}

// ExampleParent demonstrates retrieving the parent of a UIDP.
func ExampleParent() {
	fmt.Println(uidp.Parent("a/b/c"))
	fmt.Println(uidp.Parent("a"))
	// Output:
	// a/b
	// /
}

// ExampleParents demonstrates retrieving all ancestors of a UIDP.
func ExampleParents() {
	fmt.Println(uidp.Parents("a/b/c/d"))
	// Output:
	// [a/b/c a/b a]
}

// ExampleAncestry demonstrates retrieving a UIDP and all its ancestors.
func ExampleAncestry() {
	fmt.Println(uidp.Ancestry("a/b/c/d"))
	// Output:
	// [a/b/c/d a/b/c a/b a]
}

// ExampleRoot demonstrates retrieving the root segment of a UIDP.
func ExampleRoot() {
	fmt.Println(uidp.Root("a/b/c"))
	fmt.Println(uidp.Root("a"))
	// Output:
	// a
	// a
}

// ExampleValid demonstrates validating UIDP strings.
func ExampleValid() {
	fmt.Println(uidp.Valid("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))
	fmt.Println(uidp.Valid("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa/bbbbbbbbbbbbbbbb"))
	fmt.Println(uidp.Valid("not-a-uidp"))
	// Output:
	// true
	// true
	// false
}
