/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package receiver

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"chainguard.dev/sdk/events"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/coreos/go-oidc/v3/oidc"
)

// Handler is a function that handles a CloudEvent.
type Handler func(ctx context.Context, event cloudevents.Event) error

// NewHTTPHandler returns an http.Handler that verifies each delivery was sent
// by Chainguard and intended for the specified Group, then invokes fn.
//
// This is the entry point most callers want: serve it directly and the
// Authorization header, event decoding, and failure status codes are all
// handled for you.
//
//	h, err := receiver.NewHTTPHandler(ctx, "https://issuer.enforce.dev", group, fn)
//	if err != nil {
//		return err
//	}
//	http.ListenAndServe(":8080", h)
//
// Prefer this over wiring the Handler from New into the CloudEvents client:
// that client's StartReceiver does not carry the HTTP request context through
// to the handler, so the Authorization header is unreachable there. See New.
func NewHTTPHandler(ctx context.Context, issuer, group string, fn Handler) (http.Handler, error) {
	verify, err := New(ctx, issuer, group, fn)
	if err != nil {
		return nil, err
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		event, err := binding.ToEvent(r.Context(), cehttp.NewMessageFromHttpRequest(r))
		if err != nil {
			http.Error(w, "not a well-formed CloudEvent", http.StatusBadRequest)
			return
		}

		// The verifier reads the Authorization header out of the request data
		// on the context, so it has to be attached before the call.
		if err := verify(cehttp.WithRequestDataAtContext(r.Context(), r), *event); err != nil {
			// Verification failures carry the status they want returned;
			// anything else is a fault in fn rather than a bad delivery.
			status := http.StatusInternalServerError
			var result *cehttp.Result
			if errors.As(err, &result) {
				status = result.StatusCode
			}
			http.Error(w, err.Error(), status)
			return
		}
		w.WriteHeader(http.StatusOK)
	}), nil
}

// New returns a new Handler that verifies the Event was sent by Chainguard,
// intended for the specified Group, then invokes the provided Handler.
//
// The returned Handler reads the delivery's Authorization header from the
// HTTP request data on its context, so whatever serves it MUST attach that
// with cehttp.WithRequestDataAtContext. Most callers should use
// [NewHTTPHandler], which does this for them.
//
// In particular the CloudEvents client cannot serve this Handler: its
// StartReceiver hands the handler the receive loop's context rather than the
// request's, so the request data is never present and every delivery fails.
// Passing cehttp.WithRequestDataAtContextMiddleware to the client does not
// help — the middleware attaches the data to the request context, which the
// client drops before invoking the handler.
//
// TODO(jason): Accept options for configuring the issuer and accepted event types.
func New(ctx context.Context, issuer, group string, fn Handler) (Handler, error) {
	// Construct a verifier that ensures tokens are issued by the Chainguard
	// issuer we expect and are intended for a customer webhook.
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: "customer"})

	return func(ctx context.Context, event cloudevents.Event) error {
		// We expect Chainguard webhooks to pass an Authorization header, which
		// is only reachable through the request data on the context.
		rd := cehttp.RequestDataFromContext(ctx)
		if rd == nil {
			return cloudevents.NewHTTPResult(http.StatusInternalServerError,
				"no HTTP request data on the context, so the Authorization header cannot be read: serve this handler with receiver.NewHTTPHandler, or attach the data yourself with cehttp.WithRequestDataAtContext. The CloudEvents client's StartReceiver does not propagate the request context and cannot serve this handler")
		}
		auth := strings.TrimPrefix(rd.Header.Get("Authorization"), "Bearer ")
		if auth == "" {
			return cloudevents.NewHTTPResult(http.StatusUnauthorized, "Unauthorized")
		}

		claims := events.WebhookCustomClaims{}

		// Verify that the token is well-formed, and in fact intended for us!
		if tok, err := verifier.Verify(ctx, auth); err != nil {
			return cloudevents.NewHTTPResult(http.StatusForbidden, "unable to verify token: %w", err)
		} else if !strings.HasPrefix(tok.Subject, "webhook:") {
			return cloudevents.NewHTTPResult(http.StatusForbidden, "subject should be from the Chainguard webhook component, got: %q", tok.Subject)
		} else if got := strings.TrimPrefix(tok.Subject, "webhook:"); got != group {
			return cloudevents.NewHTTPResult(http.StatusForbidden, "this token is intended for %q, wanted one for %q", got, group)
		} else if err := tok.Claims(&claims); err != nil {
			return cloudevents.NewHTTPResult(http.StatusForbidden, "this token does not contain the Chainguard custom webhook claims: %v", err)
		}

		h := sha256.New()
		h.Write(event.Data())
		bs := h.Sum(nil)
		if got, want := fmt.Sprintf("sha256:%x", bs), claims.Webhook.Digest; got != want {
			return cloudevents.NewHTTPResult(http.StatusForbidden, "this token is intended for a message with digest %q, got message with digest %q", want, got)
		}

		return fn(ctx, event)
	}, nil
}
