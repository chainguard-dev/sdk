/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login_test

import (
	"fmt"

	"chainguard.dev/sdk/auth/login"
)

// ExampleError demonstrates the Error type returned by login operations.
func ExampleError() {
	err := &login.Error{
		Details: "failed to open browser",
		Err:     fmt.Errorf("exec: not found"),
	}
	fmt.Println(err.Error())
	// Output:
	// login: failed to open browser: exec: not found
}

// ExampleBuildHeadlessURL demonstrates that a missing headless code returns an error.
func ExampleBuildHeadlessURL() {
	_, err := login.BuildHeadlessURL()
	fmt.Println(err != nil)
	// Output:
	// true
}
