/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package receiver_test

import (
	"context"
	"log"

	"chainguard.dev/sdk/events/receiver"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Example demonstrates creating a secure webhook receiver that verifies
// Chainguard events are authentic and intended for your organization.
func Example() {
	ctx := context.Background()

	// Create a receiver that verifies events from Chainguard's issuer
	// and ensures they are intended for your group.
	handler, err := receiver.New(ctx, "https://issuer.enforce.dev", "my-group-id",
		func(_ context.Context, event cloudevents.Event) error {
			// Process the verified event
			log.Printf("Received event: %s", event.Type())
			return nil
		})
	if err != nil {
		log.Fatalf("failed to create receiver: %v", err)
	}

	// Use the handler with CloudEvents HTTP receiver
	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	// Start receiving events
	if err := c.StartReceiver(ctx, handler); err != nil {
		log.Fatalf("failed to start receiver: %v", err)
	}
}

// Example_customHandler demonstrates using a custom handler function
// to process different event types.
func Example_customHandler() {
	ctx := context.Background()

	handler, err := receiver.New(ctx, "https://issuer.enforce.dev", "my-group-id",
		func(_ context.Context, event cloudevents.Event) error {
			// Handle different event types
			switch event.Type() {
			case "dev.chainguard.image.created":
				log.Printf("New image created: %s", event.Subject())
			case "dev.chainguard.policy.violated":
				log.Printf("Policy violation detected: %s", event.Subject())
			default:
				log.Printf("Unknown event type: %s", event.Type())
			}
			return nil
		})
	if err != nil {
		log.Fatalf("failed to create receiver: %v", err)
	}

	c, err := cloudevents.NewClientHTTP()
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	if err := c.StartReceiver(ctx, handler); err != nil {
		log.Fatalf("failed to start receiver: %v", err)
	}
}
