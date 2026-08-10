/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package skills_test

import (
	"fmt"
	"strings"

	"chainguard.dev/sdk/skills"
)

func ExampleHardenJobID() {
	// The same inputs always yield the same id (idempotent submit / dedupe).
	a := skills.HardenJobID("org1", "user1", "pdf", "sha256:abc")
	b := skills.HardenJobID("org1", "user1", "pdf", "sha256:abc")
	fmt.Println(a == b)
	// A different content digest yields a different id.
	c := skills.HardenJobID("org1", "user1", "pdf", "sha256:def")
	fmt.Println(a == c)
	// Output:
	// true
	// false
}

func ExampleHardenOperationName() {
	name := skills.HardenOperationName("org1/child", "user1", "pdf", "sha256:abc")
	fmt.Println(strings.HasPrefix(name, "harden/org1/child/"))
	// Output:
	// true
}

func ExampleGroupFromHardenName() {
	name := skills.HardenOperationName("org1/child", "user1", "pdf", "sha256:abc")
	group, err := skills.GroupFromHardenName(name)
	fmt.Println(err)
	fmt.Println(group)
	// Output:
	// <nil>
	// org1/child
}
