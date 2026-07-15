/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package githubflow_test

import (
	"context"
	"fmt"
	"io"

	"chainguard.dev/sdk/githubflow"
)

// ExampleObtain demonstrates calling the loopback PKCE flow.
// In real usage a server-minted OAuth state is passed and the function
// opens the browser; here we pass an already-cancelled context so the
// call returns immediately without starting a browser session.
func ExampleObtain() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately so no browser window opens

	_, err := githubflow.Obtain(ctx, "my-client-id", "server-state", 9876, io.Discard)
	fmt.Println(err != nil)
	// Output:
	// true
}

// ExampleResult demonstrates accessing fields of a Result value.
func ExampleResult() {
	r := githubflow.Result{
		Code:        "gh_code",
		RedirectURI: "http://127.0.0.1:9876/callback",
		Verifier:    "pkce-verifier",
	}
	fmt.Println(r.Code)
	fmt.Println(r.RedirectURI)
	fmt.Println(r.Verifier)
	// Output:
	// gh_code
	// http://127.0.0.1:9876/callback
	// pkce-verifier
}
