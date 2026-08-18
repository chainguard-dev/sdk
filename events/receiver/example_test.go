/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package receiver_test

import (
	"context"
	"net/http"

	"chainguard.dev/sdk/events/receiver"
	"github.com/chainguard-dev/clog"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Example demonstrates serving a webhook receiver that verifies Chainguard
// events are authentic and intended for your organization.
func Example() {
	ctx := context.Background()

	// Verify events came from Chainguard's issuer and are intended for your
	// group, then hand the verified event to the callback.
	h, err := receiver.NewHTTPHandler(ctx, "https://issuer.enforce.dev", "my-group-id",
		func(ctx context.Context, event cloudevents.Event) error {
			clog.InfoContextf(ctx, "Received event: %s", event.Type())
			return nil
		})
	if err != nil {
		clog.FatalContextf(ctx, "failed to create receiver: %v", err)
	}

	// Deliveries that fail verification never reach the callback; they are
	// answered with 401 or 403 directly.
	server := &http.Server{Addr: ":8080", Handler: h}
	if err := server.ListenAndServe(); err != nil {
		clog.FatalContextf(ctx, "server exited: %v", err)
	}
}

// Example_customHandler demonstrates dispatching on the event type, and
// mounting the receiver on a path alongside other routes.
func Example_customHandler() {
	ctx := context.Background()

	h, err := receiver.NewHTTPHandler(ctx, "https://issuer.enforce.dev", "my-group-id",
		func(ctx context.Context, event cloudevents.Event) error {
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

	mux := http.NewServeMux()
	mux.Handle("POST /webhook", h)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := &http.Server{Addr: ":8080", Handler: mux}
	if err := server.ListenAndServe(); err != nil {
		clog.FatalContextf(ctx, "server exited: %v", err)
	}
}
