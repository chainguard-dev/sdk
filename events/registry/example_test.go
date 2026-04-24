/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package registry_test

import (
	"fmt"

	"chainguard.dev/sdk/events/registry"
)

// ExamplePulledEventType demonstrates the CloudEvents type for registry pulls.
func ExamplePulledEventType() {
	fmt.Println(registry.PulledEventType)
	// Output:
	// dev.chainguard.registry.pull.v1
}

// ExamplePushedEventType demonstrates the CloudEvents type for registry pushes.
func ExamplePushedEventType() {
	fmt.Println(registry.PushedEventType)
	// Output:
	// dev.chainguard.registry.push.v1
}

// ExamplePullEvent demonstrates constructing a PullEvent.
func ExamplePullEvent() {
	evt := registry.PullEvent{
		Repository: "cgr.dev/chainguard/nginx",
		Digest:     "sha256:abc123",
		Method:     "GET",
	}
	fmt.Println(evt.Repository)
	fmt.Println(evt.Method)
	// Output:
	// cgr.dev/chainguard/nginx
	// GET
}
