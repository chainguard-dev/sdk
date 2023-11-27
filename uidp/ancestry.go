/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uidp

import (
	"path"
	"regexp"
	"strings"
)

// IsAncestor checks whether the "parent" UIDP is an ancestor (non-inclusive)
// of the given "child" UIDP.
func IsAncestor(parent, child string) bool {
	return strings.HasPrefix(child, parent+"/")
}

// InRoot checks whether the UIDP is in the root, as opposed to within a group
func InRoot(child string) bool {
	return !strings.Contains(child, "/")
}

// IsAncestorOrSelf checks whether the "parent" UIDP is an ancestor (inclusive)
// of the given "child" UIDP.
func IsAncestorOrSelf(parent, child string) bool {
	return child == parent || IsAncestor(parent, child)
}

// Parent returns the "parent" UIDP for a child UIDP. Returns / if parent is root.
//
// Example:
//
//	Parent("a/b/c") returns "a/b"
//	Parent("a") returns "/"
func Parent(child string) string {
	// path.Dir returns "." if the child is already at the root
	if p := path.Dir(child); p != "." {
		return p
	}
	return "/"
}

// Parents returns all "parent" UIDP for a child UIDP. Returns empty slice if parent is root.
//
// Example:
//
//	Parents("a/b/c/d") returns ["a/b/c", "a/b", "a"]
func Parents(child string) []string {
	parents := make([]string, 0, strings.Count(child, "/"))
	for p := Parent(child); p != "/"; p = Parent(p) {
		parents = append(parents, p)
	}
	return parents
}

// Ancestry returns all parent UIDPs and the child. Returns only the child if it is root.
//
// Example:
//
//	Ancestry("a/b/c/d") returns ["a/b/c/d", "a/b/c", "a/b", "a"]
func Ancestry(child string) []string {
	ancestry := make([]string, 0, strings.Count(child, "/")+1)
	ancestry = append(ancestry, child)
	for p := Parent(child); p != "/"; p = Parent(p) {
		ancestry = append(ancestry, p)
	}
	return ancestry
}

// Valid returns true for valid UIDP values.
// The base segment of a UIDP is 20 hex-encoded bytes (40 characters).  This may
// be followed by zero or more parts with 8 hex-encoded bytes (16 characters).
var Valid = regexp.MustCompile(`^[0-9a-f]{40}(?:/[0-9a-f]{16})*$`).MatchString
