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
//   - A ready-to-serve http.Handler, with HTTP status codes for every failure mode
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
// Call NewHTTPHandler with your OIDC issuer URL, group ID, and handler function, then
// serve the result. Verification, event decoding, and failure status codes are handled
// for you:
//
//	h, err := receiver.NewHTTPHandler(ctx, "https://issuer.enforce.dev", "my-group-id",
//		func(ctx context.Context, event cloudevents.Event) error {
//			// Process the verified event
//			return nil
//		})
//	if err != nil {
//		log.Fatalf("failed to create receiver: %v", err)
//	}
//
//	if err := http.ListenAndServe(":8080", h); err != nil {
//		log.Fatalf("server exited: %v", err)
//	}
//
// # Serving the Handler from New directly
//
// New returns the verifying [Handler] on its own, for callers composing it into an
// existing pipeline. It reads the delivery's Authorization header from the HTTP request
// data on its context, so whatever serves it must attach that data with
// cehttp.WithRequestDataAtContext.
//
// The CloudEvents client cannot serve it. Its StartReceiver hands the handler the
// receive loop's context rather than the request's, so the request data is never
// present and every delivery fails. Passing cehttp.WithRequestDataAtContextMiddleware
// to the client does not help: the middleware attaches the data to the request context,
// which the client drops before invoking the handler. Use NewHTTPHandler instead.
//
// # Integration Patterns
//
// The receiver expects:
//   - An Authorization header with a Bearer token
//   - A CloudEvent payload in the request body
//   - The token's digest claim to match the SHA-256 hash of the event data
//
// # Error Handling
//
// Verification failures are reported as CloudEvents HTTP results, which NewHTTPHandler
// turns into responses:
//   - 400 Bad Request: the body is not a well-formed CloudEvent
//   - 401 Unauthorized: missing Authorization header
//   - 403 Forbidden: invalid token, wrong group, wrong subject, or digest mismatch
//   - 500 Internal Server Error: the handler was served without HTTP request data on
//     the context, so the Authorization header could not be read
//
// Your handler function should return an error for processing failures, which becomes a
// 500 response.
package receiver
