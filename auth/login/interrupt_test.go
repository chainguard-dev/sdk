/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"context"
	"errors"
	"net/http"
	"runtime"
	"strings"
	"testing"
	"time"
)

// noRedirect returns an HTTP client that does not follow redirects, so a hit to
// /callback delivers the token to the server without also loading the "/"
// success page (which is what closes s.token).
func noRedirect() *http.Client {
	return &http.Client{
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

// Control: cancelling the context while Token() waits for the browser callback
// must unblock it promptly. This path is already ctx-guarded (server.go outer
// select), so this should pass today — it pins the working half.
func TestToken_CancelBeforeCallback_ReturnsPromptly(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	s, err := newServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	done := make(chan error, 1)
	go func() {
		_, err := s.Token()
		done <- err
	}()

	cancel()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("Token() error: got %#v, want context.Canceled", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Token() did not return after context cancellation while waiting for the callback")
	}
}

// After the browser delivers a token to /callback but before the success page
// ("/") renders and closes s.token, Token() waits on the inner select. A Ctrl+C
// here must not hang — and, because the OAuth round-trip already produced a
// valid token, must not discard it either: Token() returns the token, not an
// error, so a late interrupt doesn't force a full re-login.
func TestToken_CancelAfterTokenBeforeSuccessPage_ReturnsPromptly(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	s, err := newServer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	// Deliver a token to /callback WITHOUT following the redirect to "/", so the
	// token is received but the success page never closes s.token.
	callback := strings.ReplaceAll(s.URL(), "token=true", "token=foo")
	if _, err := noRedirect().Get(callback); err != nil {
		t.Fatalf("callback request failed: %v", err)
	}

	type result struct {
		tok string
		err error
	}
	done := make(chan result, 1)
	go func() {
		tok, err := s.Token()
		done <- result{tok, err}
	}()

	// Wait (bounded) until Token() consumes the delivered token — the cap-1
	// buffer drains from 1 to 0 — so it is now on the inner success-page receive.
	// This is the event we synchronize on, not a fixed sleep.
	deadline := time.Now().Add(2 * time.Second)
	for len(s.token) > 0 {
		if time.Now().After(deadline) {
			t.Fatal("Token() never consumed the delivered token")
		}
		runtime.Gosched()
	}
	cancel()

	select {
	case r := <-done:
		// A late cancel here must not hang, and must not discard the token: the
		// OAuth round-trip already completed, so Token() returns it rather than
		// forcing a full re-login.
		if r.err != nil {
			t.Errorf("Token() error: got %#v, want nil (completed auth preserved)", r.err)
		}
		if r.tok != "foo" {
			t.Errorf("Token() token: got %q, want %q (a received token must not be discarded)", r.tok, "foo")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Token() did not return after context cancellation while waiting for the success page — the interrupt is swallowed")
	}
}
