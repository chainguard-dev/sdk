/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package deviceflow_test

import (
	"fmt"

	"chainguard.dev/sdk/auth/deviceflow"
)

// ExampleNewTokenGetter demonstrates creating a TokenGetter for a given issuer.
func ExampleNewTokenGetter() {
	tg := deviceflow.NewTokenGetter("https://issuer.enforce.dev")
	fmt.Println(tg != nil)
	// Output:
	// true
}
