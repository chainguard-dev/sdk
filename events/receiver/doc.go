/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package receiver provides secure CloudEvent webhook receivers for Chainguard events.
//
// # Overview
//
// This package implements a CloudEvent handler that verifies webhook events sent by
// Chainguard are authentic and intended for your organization. It validates OIDC tokens
// embedded in webhook requests to ensure events come from Chainguard's webhook component
// and are addressed to the correct group.
//
// # Features
//
//   - OIDC token verification using the Chainguard issuer
//   - Group-based authorization to ensure events are intended for your organization
//   - Message digest validation to prevent replay attacks
//   - Integration with the CloudEvents SDK for standard event handling
//   - Automatic HTTP status code responses for authentication and authorization failures
//
// # Security Model
//
// Each webhook request includes an OIDC token in the Authorization header. The token
// contains:
//   - A subject claim identifying the webhook component (must start with "webhook:")
//   - A group identifier indicating the intended recipient organization
//   - A message digest (SHA-256) of the event payload to prevent tampering
//
// The receiver validates all three components before invoking your handler function.
//
// # Usage
//
// Create a receiver by calling New with your OIDC issuer URL, group ID, and handler
// function. The returned handler can be used with the CloudEvents HTTP protocol:
//
//	handler, err := receiver.New(ctx, "https://issuer.enforce.dev", "my-group-id",
//		func(ctx context.Context, event cloudevents.Event) error {
//			// Process the verified event
//			return nil
//		})
//	if err != nil {
//		log.Fatalf("failed to create receiver: %v", err)
//	}
//
//	// Use with CloudEvents HTTP receiver
//	c, err := cloudevents.NewClientHTTP()
//	if err != nil {
//		log.Fatalf("failed to create client: %v", err)
//	}
//	if err := c.StartReceiver(ctx, handler); err != nil {
//		log.Fatalf("failed to start receiver: %v", err)
//	}
//
// # Integration Patterns
//
// The receiver integrates with the CloudEvents SDK's HTTP protocol. It expects:
//   - An Authorization header with a Bearer token
//   - A CloudEvent payload in the request body
//   - The token's digest claim to match the SHA-256 hash of the event data
//
// On authentication or authorization failure, the receiver returns an appropriate
// HTTP status code (401 Unauthorized or 403 Forbidden) without invoking your handler.
//
// # Error Handling
//
// The receiver returns CloudEvents HTTP results for all authentication and authorization
// failures. These are automatically converted to HTTP responses by the CloudEvents SDK:
//   - 401 Unauthorized: Missing Authorization header
//   - 403 Forbidden: Invalid token, wrong group, wrong subject, or digest mismatch
//
// Your handler function should return an error for processing failures. The CloudEvents
// SDK will convert these to appropriate HTTP responses.
package receiver
