/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package receiver_test

import (
	"context"

	"chainguard.dev/sdk/events/receiver"
	"github.com/chainguard-dev/clog"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Example demonstrates creating a secure webhook receiver that verifies
// Chainguard events are authentic and intended for your organization.
func Example() {
	ctx := context.Background()

	// Create a receiver that verifies events from Chainguard's issuer
	// and ensures they are intended for your group.
	handler, err := receiver.New(ctx, "https://issuer.enforce.dev", "my-group-id",
		func(ctx context.Context, event cloudevents.Event) error {
			// Process the verified event
			clog.InfoContextf(ctx, "Received event: %s", event.Type())
			return nil
		})
	if err != nil {
		clog.FatalContextf(ctx, "failed to create receiver: %v", err)
	}

	// Use the handler with CloudEvents HTTP receiver
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		clog.FatalContextf(ctx, "failed to create client: %v", err)
	}

	// Start receiving events
	if err := c.StartReceiver(ctx, handler); err != nil {
		clog.FatalContextf(ctx, "failed to start receiver: %v", err)
	}
}

// Example_customHandler demonstrates using a custom handler function
// to process different event types.
func Example_customHandler() {
	ctx := context.Background()

	handler, err := receiver.New(ctx, "https://issuer.enforce.dev", "my-group-id",
		func(ctx context.Context, event cloudevents.Event) error {
			// Handle different event types
			switch event.Type() {
			case "dev.chainguard.image.created":
				clog.InfoContextf(ctx, "New image created: %s", event.Subject())
			case "dev.chainguard.policy.violated":
				clog.InfoContextf(ctx, "Policy violation detected: %s", event.Subject())
			default:
				clog.InfoContextf(ctx, "Unknown event type: %s", event.Type())
			}
			return nil
		})
	if err != nil {
		clog.FatalContextf(ctx, "failed to create receiver: %v", err)
	}

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		clog.FatalContextf(ctx, "failed to create client: %v", err)
	}

	if err := c.StartReceiver(ctx, handler); err != nil {
		clog.FatalContextf(ctx, "failed to start receiver: %v", err)
	}
}
