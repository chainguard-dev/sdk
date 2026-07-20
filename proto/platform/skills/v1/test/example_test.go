/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	skillstest "chainguard.dev/sdk/proto/platform/skills/v1/test"
)

// ExampleMockSkillsClients demonstrates constructing a mock Skills client.
func ExampleMockSkillsClients() {
	mock := skillstest.MockSkillsClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
