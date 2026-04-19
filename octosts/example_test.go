/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package octosts_test

import (
	"fmt"

	"chainguard.dev/sdk/octosts"
)

// ExampleOctoSTSEndpoint demonstrates the default OctoSTS service endpoint.
func ExampleOctoSTSEndpoint() {
	fmt.Println(octosts.OctoSTSEndpoint)
	// Output:
	// https://octo-sts.dev
}
