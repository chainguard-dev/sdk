/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uidp

import (
	"reflect"
	"testing"
)

func TestAncestorExclusive(t *testing.T) {
	tests := []struct {
		name   string
		parent string
		child  string
		want   bool
	}{{
		name:   "direct parent",
		parent: "a/b",
		child:  "a/b/c",
		want:   true,
	}, {
		name:   "grand parent",
		parent: "a",
		child:  "a/b/c",
		want:   true,
	}, {
		name:   "self",
		parent: "a/b/c",
		child:  "a/b/c",
		want:   false, // This should be the main difference from below
	}, {
		name:   "child",
		parent: "a/b/c/d",
		child:  "a/b/c",
		want:   false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsAncestor(test.parent, test.child)
			if got != test.want {
				t.Errorf("IsAncestor(%q, %q) = %v, wanted %v", test.parent, test.child, got, test.want)
			}
		})
	}
}

func TestAncestorInclusive(t *testing.T) {
	tests := []struct {
		name   string
		parent string
		child  string
		want   bool
	}{{
		name:   "direct parent",
		parent: "a/b",
		child:  "a/b/c",
		want:   true,
	}, {
		name:   "grand parent",
		parent: "a",
		child:  "a/b/c",
		want:   true,
	}, {
		name:   "self",
		parent: "a/b/c",
		child:  "a/b/c",
		want:   true, // This should be the main difference from above
	}, {
		name:   "child",
		parent: "a/b/c/d",
		child:  "a/b/c",
		want:   false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := IsAncestorOrSelf(test.parent, test.child)
			if got != test.want {
				t.Errorf("IsAncestor(%q, %q) = %v, wanted %v", test.parent, test.child, got, test.want)
			}
		})
	}
}

func TestInRoot(t *testing.T) {
	tests := []struct {
		id   string
		want bool
	}{
		{"x", true},
		{"xyz", true},
		{"x/y", false},
		{"x/y/z", false},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			got := InRoot(test.id)
			if got != test.want {
				t.Errorf("InRoot(%q) = %v, wanted %v", test.id, got, test.want)
			}
		})
	}
}

func TestParent(t *testing.T) {
	tests := []struct {
		id   string
		want string
	}{
		{"x", "/"},
		{"xyz", "/"},
		{"x/y", "x"},
		{"x/y/z", "x/y"},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			got := Parent(test.id)
			if got != test.want {
				t.Errorf("Parent(%q) = %q, wanted %q", test.id, got, test.want)
			}
		})
	}
}

func TestParents(t *testing.T) {
	tests := []struct {
		id   string
		want []string
	}{
		{"x", []string{}},
		{"xyz", []string{}},
		{"x/y", []string{"x"}},
		{"x/y/z", []string{"x/y", "x"}},
		{"", []string{}},
		{"/", []string{}},
		{"//", []string{}},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			got := Parents(test.id)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Parent(%q) = %q, wanted %q", test.id, got, test.want)
			}
		})
	}
}

func TestAncestry(t *testing.T) {
	tests := []struct {
		id   string
		want []string
	}{
		{"x", []string{"x"}},
		{"xyz", []string{"xyz"}},
		{"x/y", []string{"x/y", "x"}},
		{"x/y/z", []string{"x/y/z", "x/y", "x"}},
		{"", []string{""}},
		{"/", []string{"/"}},
		{"//", []string{"//"}},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			got := Ancestry(test.id)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Parent(%q) = %q, wanted %q", test.id, got, test.want)
			}
		})
	}
}
