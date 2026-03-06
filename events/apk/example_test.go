/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package apk_test

import (
	"fmt"
	"time"

	"chainguard.dev/sdk/civil"
	"chainguard.dev/sdk/events/apk"
)

// Example demonstrates basic usage of PushEvent and its methods.
func Example() {
	push := apk.PushEvent{
		Repository:    "example-repo",
		RepoID:        "abc123def456",
		Package:       "nginx",
		Origin:        "nginx",
		Version:       "1.25.3-r0",
		Architecture:  "x86_64",
		Checksum:      "Q1abc123def456",
		When:          civil.DateTimeOf(time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)),
		Location:      "SeattleWAUS",
		RemoteAddress: "203.0.113.42",
		UserAgent:     "apk-tools/2.14.0",
	}

	fmt.Println("Package:", push.Package)
	fmt.Println("Version:", push.Version)
	fmt.Println("Base Path:", push.APKBasePath())
	fmt.Println("Full Path:", push.APKPath())

	// Output:
	// Package: nginx
	// Version: 1.25.3-r0
	// Base Path: x86_64/nginx-1.25.3-r0.apk
	// Full Path: abc123def456/x86_64/nginx-1.25.3-r0.apk
}

// ExamplePushEvent demonstrates creating and using a PushEvent.
func ExamplePushEvent() {
	event := apk.PushEvent{
		Repository:    "chainguard",
		RepoID:        "repo-uidp-123",
		Package:       "curl",
		Origin:        "curl",
		Version:       "8.5.0-r0",
		Architecture:  "aarch64",
		Checksum:      "Q1checksum",
		Location:      "LondonGB",
		RemoteAddress: "198.51.100.10",
		UserAgent:     "apk-tools/2.14.0",
	}

	fmt.Println(event.APKPath())

	// Output:
	// repo-uidp-123/aarch64/curl-8.5.0-r0.apk
}

// ExamplePushEvent_aPKBasePath demonstrates the APKBasePath method.
func ExamplePushEvent_APKBasePath() {
	event := apk.PushEvent{
		Package:      "busybox",
		Version:      "1.36.1-r0",
		Architecture: "x86_64",
	}

	fmt.Println(event.APKBasePath())

	// Output:
	// x86_64/busybox-1.36.1-r0.apk
}

// ExamplePushEvent_aPKPath demonstrates the APKPath method.
func ExamplePushEvent_APKPath() {
	event := apk.PushEvent{
		RepoID:       "my-repo-id",
		Package:      "openssl",
		Version:      "3.1.4-r0",
		Architecture: "armv7",
	}

	fmt.Println(event.APKPath())

	// Output:
	// my-repo-id/armv7/openssl-3.1.4-r0.apk
}

// ExamplePushEvent_error demonstrates a PushEvent with an error.
func ExamplePushEvent_error() {
	event := apk.PushEvent{
		Package: "invalid-pkg",
		Error: &apk.Error{
			Status:  403,
			Code:    "FORBIDDEN",
			Message: "insufficient permissions",
		},
	}

	if event.Error != nil {
		fmt.Printf("Push failed: %s (status %d)\n", event.Error.Message, event.Error.Status)
	}

	// Output:
	// Push failed: insufficient permissions (status 403)
}

// ExamplePullEvent demonstrates creating and using a PullEvent.
func ExamplePullEvent() {
	event := apk.PullEvent{
		Repository:    "chainguard",
		RepoID:        "repo-uidp-456",
		Package:       "git",
		Origin:        "git",
		Version:       "2.43.0-r0",
		Architecture:  "x86_64",
		Checksum:      "Q1checksum",
		Location:      "TokyoJP",
		RemoteAddress: "203.0.113.100",
		UserAgent:     "apk-tools/2.14.0",
	}

	fmt.Println(event.APKPath())

	// Output:
	// repo-uidp-456/x86_64/git-2.43.0-r0.apk
}

// ExamplePullEvent_aPKBasePath demonstrates the APKBasePath method for PullEvent.
func ExamplePullEvent_APKBasePath() {
	event := apk.PullEvent{
		Package:      "python3",
		Version:      "3.12.1-r0",
		Architecture: "aarch64",
	}

	fmt.Println(event.APKBasePath())

	// Output:
	// aarch64/python3-3.12.1-r0.apk
}

// ExamplePullEvent_aPKPath demonstrates the APKPath method for PullEvent.
func ExamplePullEvent_APKPath() {
	event := apk.PullEvent{
		RepoID:       "customer-repo",
		Package:      "nodejs",
		Version:      "20.10.0-r0",
		Architecture: "x86_64",
	}

	fmt.Println(event.APKPath())

	// Output:
	// customer-repo/x86_64/nodejs-20.10.0-r0.apk
}

// ExamplePullEvent_proxy demonstrates a PullEvent with proxy information.
func ExamplePullEvent_proxy() {
	event := apk.PullEvent{
		Package:      "alpine-base",
		Version:      "3.19.0-r0",
		Architecture: "x86_64",
		ProxyUIDP:    "customer-uidp-789",
		ProxyHash:    "sha256:abc123",
	}

	if event.ProxyUIDP != "" {
		fmt.Printf("Proxied pull for customer: %s\n", event.ProxyUIDP)
		if event.ProxyHash != "" {
			fmt.Printf("Image hash: %s\n", event.ProxyHash)
		}
	}

	// Output:
	// Proxied pull for customer: customer-uidp-789
	// Image hash: sha256:abc123
}

// ExamplePushedEventType demonstrates using the event type constant.
func ExamplePushedEventType() {
	fmt.Println(apk.PushedEventType)

	// Output:
	// dev.chainguard.apk.push.v1
}

// ExamplePulledEventType demonstrates using the event type constant.
func ExamplePulledEventType() {
	fmt.Println(apk.PulledEventType)

	// Output:
	// dev.chainguard.apk.pull.v1
}
