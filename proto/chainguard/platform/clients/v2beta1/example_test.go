/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"fmt"

	clients "chainguard.dev/sdk/proto/chainguard/platform/clients/v2beta1"
)

func ExampleNewClients() {
	// NewClients creates v2beta1 API clients. In production, pass
	// a real API URL, user agent, and credentials.
	// c, err := clients.NewClients(ctx, apiURL, userAgent, cred)
	_ = clients.NewClients
	fmt.Println("func available")
	// Output: func available
}
