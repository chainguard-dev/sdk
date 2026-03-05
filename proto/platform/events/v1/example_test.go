/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	v1 "chainguard.dev/sdk/proto/platform/events/v1"
)

// ExampleIdentityMetadata_CloudEventsExtension demonstrates how to retrieve
// CloudEvents extensions from IdentityMetadata.
func ExampleIdentityMetadata_CloudEventsExtension() {
	metadata := &v1.IdentityMetadata{}

	// IdentityMetadata does not provide any extensions
	value, ok := metadata.CloudEventsExtension("group")
	fmt.Printf("Extension found: %v, value: %q\n", ok, value)
	// Output:
	// Extension found: false, value: ""
}

// ExampleIdentityMetadata_CloudEventsSubject demonstrates how to retrieve
// the CloudEvents subject from IdentityMetadata.
func ExampleIdentityMetadata_CloudEventsSubject() {
	metadata := &v1.IdentityMetadata{}

	subject := metadata.CloudEventsSubject()
	fmt.Printf("Subject: %q\n", subject)
	// Output:
	// Subject: ""
}

// ExampleSubscription_CloudEventsExtension demonstrates how to retrieve
// CloudEvents extensions from a Subscription.
func ExampleSubscription_CloudEventsExtension() {
	sub := &v1.Subscription{
		Id: "abc123/def456/subscription789",
	}

	// The "group" extension returns the parent UIDP
	value, ok := sub.CloudEventsExtension("group")
	fmt.Printf("Group extension found: %v, value: %q\n", ok, value)

	// Unknown extensions return false
	value, ok = sub.CloudEventsExtension("unknown")
	fmt.Printf("Unknown extension found: %v, value: %q\n", ok, value)
	// Output:
	// Group extension found: true, value: "abc123/def456"
	// Unknown extension found: false, value: ""
}

// ExampleSubscription_CloudEventsSubject demonstrates how to retrieve
// the CloudEvents subject from a Subscription.
func ExampleSubscription_CloudEventsSubject() {
	sub := &v1.Subscription{
		Id: "abc123/def456/subscription789",
	}

	subject := sub.CloudEventsSubject()
	fmt.Printf("Subject: %q\n", subject)
	// Output:
	// Subject: "abc123/def456/subscription789"
}

// ExampleDeleteSubscriptionRequest_CloudEventsExtension demonstrates how to
// retrieve CloudEvents extensions from a DeleteSubscriptionRequest.
func ExampleDeleteSubscriptionRequest_CloudEventsExtension() {
	req := &v1.DeleteSubscriptionRequest{
		Id: "abc123/def456/subscription789",
	}

	// The "group" extension returns the parent UIDP
	value, ok := req.CloudEventsExtension("group")
	fmt.Printf("Group extension found: %v, value: %q\n", ok, value)

	// Unknown extensions return false
	value, ok = req.CloudEventsExtension("unknown")
	fmt.Printf("Unknown extension found: %v, value: %q\n", ok, value)
	// Output:
	// Group extension found: true, value: "abc123/def456"
	// Unknown extension found: false, value: ""
}

// ExampleDeleteSubscriptionRequest_CloudEventsSubject demonstrates how to
// retrieve the CloudEvents subject from a DeleteSubscriptionRequest.
func ExampleDeleteSubscriptionRequest_CloudEventsSubject() {
	req := &v1.DeleteSubscriptionRequest{
		Id: "abc123/def456/subscription789",
	}

	subject := req.CloudEventsSubject()
	fmt.Printf("Subject: %q\n", subject)
	// Output:
	// Subject: "abc123/def456/subscription789"
}

// ExampleDeleteSubscriptionRequest_CloudEventsRedact demonstrates the
// CloudEventsRedact method which returns nil for DeleteSubscriptionRequest.
func ExampleDeleteSubscriptionRequest_CloudEventsRedact() {
	req := &v1.DeleteSubscriptionRequest{
		Id: "abc123/def456/subscription789",
	}

	redacted := req.CloudEventsRedact()
	fmt.Printf("Redacted value: %v\n", redacted)
	// Output:
	// Redacted value: <nil>
}
