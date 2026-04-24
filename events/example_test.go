/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events_test

import (
	"fmt"

	"chainguard.dev/sdk/events"
)

// ExampleDeliveryTypeKey demonstrates the CloudEvents extension key for delivery type.
func ExampleDeliveryTypeKey() {
	fmt.Println(events.DeliveryTypeKey)
	// Output:
	// chainguarddev1delivery
}

// ExampleAudienceCustomer demonstrates the customer audience constant.
func ExampleAudienceCustomer() {
	fmt.Println(events.AudienceCustomer)
	// Output:
	// customer
}

// ExampleAudienceInternal demonstrates the internal audience constant.
func ExampleAudienceInternal() {
	fmt.Println(events.AudienceInternal)
	// Output:
	// internal
}
