/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package harness_test

import (
	"fmt"
	"strings"

	"chainguard.dev/sdk/harness"
)

func ExampleHuntName() {
	// The same inputs always yield the same name (idempotent CreateHunt).
	a := harness.HuntName("org1", "user1", "sha256:abc")
	b := harness.HuntName("org1", "user1", "sha256:abc")
	fmt.Println(a == b)
	// A different bundle digest is a different hunt.
	c := harness.HuntName("org1", "user1", "sha256:def")
	fmt.Println(a == c)
	// So is the same bundle submitted by a different member.
	d := harness.HuntName("org1", "user2", "sha256:abc")
	fmt.Println(a == d)
	fmt.Println(strings.HasPrefix(harness.HuntName("org1/child", "user1", "sha256:abc"), "hunt/org1/child/"))
	// Output:
	// true
	// false
	// false
	// true
}

func ExampleGroupFromHuntName() {
	group, err := harness.GroupFromHuntName(harness.HuntName("org1/child", "user1", "sha256:abc"))
	fmt.Println(err)
	fmt.Println(group)
	// Output:
	// <nil>
	// org1/child
}

func ExampleHuntRetryName() {
	failed := harness.HuntName("org/child", "user", "sha256:abc")
	retry, err := harness.HuntRetryName(failed)
	if err != nil {
		fmt.Println(err)
		return
	}
	repeated, err := harness.HuntRetryName(failed)
	if err != nil {
		fmt.Println(err)
		return
	}
	group, err := harness.GroupFromHuntName(retry)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(group, retry != failed, retry == repeated)
	// Output: org/child true true
}
