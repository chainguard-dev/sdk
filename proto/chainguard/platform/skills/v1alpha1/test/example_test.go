/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	skillstest "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1/test"
)

// ExampleMockSkillsClients demonstrates constructing a mock skills
// catalog client.
func ExampleMockSkillsClients() {
	mock := skillstest.MockSkillsClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
