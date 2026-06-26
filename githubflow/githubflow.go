/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package githubflow runs the GitHub App user (user-to-server) OAuth "loopback"
// flow shared by clients such as chainctl and the guardener OAuth test binary.
//
// It starts a temporary HTTP server on 127.0.0.1, opens the browser to GitHub's
// authorize endpoint with PKCE (S256), and captures the authorization code. The
// code is exchanged for a token server-side, so this package never holds the
// GitHub App client secret or the resulting user token — it only returns the
// code, the loopback redirect URI it used, and the PKCE verifier, all of which
// the server replays at exchange time.
package githubflow

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// Result is the outcome of a successful loopback authorization.
type Result struct {
	// Code is the GitHub authorization code.
	Code string
	// RedirectURI is the exact loopback redirect URI used; the server-side code
	// exchange must replay it.
	RedirectURI string
	// Verifier is the PKCE code verifier; the server-side exchange must replay it.
	Verifier string
}

// Obtain runs the loopback PKCE flow for the given GitHub App OAuth client ID on
// the given loopback port, writing user-facing progress to out. The authorize
// request uses PKCE (S256) so an intercepted loopback code can't be redeemed
// without the verifier, forces the account picker (`prompt=select_account`), and
// disables signup (`allow_signup=false`).
//
// state is the server-minted OAuth state from BeginGitHubOAuth: it is placed in
// the authorize request, GitHub echoes it back to the loopback, and Obtain
// verifies the echo (CSRF). The caller forwards the same state on the mutating
// RPC, where the server verifies its signature before exchanging the code.
func Obtain(ctx context.Context, clientID, state string, port int, out io.Writer) (Result, error) {
	if clientID == "" {
		return Result{}, errors.New("githubflow: empty GitHub OAuth client id")
	}
	if state == "" {
		return Result{}, errors.New("githubflow: empty OAuth state")
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return Result{}, fmt.Errorf("githubflow: listen on loopback port %d (is it free / registered on the App?): %w", port, err)
	}
	defer listener.Close()
	redirectURI := fmt.Sprintf("http://127.0.0.1:%d/callback", port)

	verifier := oauth2.GenerateVerifier()

	cfg := &oauth2.Config{ClientID: clientID, RedirectURL: redirectURI, Endpoint: githuboauth.Endpoint}
	authURL := cfg.AuthCodeURL(state,
		oauth2.S256ChallengeOption(verifier),
		oauth2.SetAuthURLParam("prompt", "select_account"),
		oauth2.SetAuthURLParam("allow_signup", "false"),
	)

	type result struct {
		code string
		err  error
	}
	resCh := make(chan result, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		var res result
		switch {
		case q.Get("error") != "":
			http.Error(w, "Authorization was denied. You can close this window.", http.StatusBadRequest)
			res = result{err: fmt.Errorf("github authorization error: %s: %s", q.Get("error"), q.Get("error_description"))}
		case q.Get("state") != state:
			http.Error(w, "Authorization state mismatch. You can close this window.", http.StatusBadRequest)
			res = result{err: errors.New("oauth state mismatch")}
		case q.Get("code") == "":
			http.Error(w, "Missing authorization code. You can close this window.", http.StatusBadRequest)
			res = result{err: errors.New("no authorization code in callback")}
		default:
			fmt.Fprintln(w, "Authorization complete. You can close this window and return to the terminal.")
			res = result{code: q.Get("code")}
		}
		// Flush the response onto the wire before signalling the result: Obtain
		// tears the server down as soon as it receives on resCh, so without this
		// the browser can get a reset (and a refused retry) instead of the page.
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		resCh <- res
	})

	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = srv.Serve(listener) }()

	fmt.Fprintf(out, "Opening your browser to authorize with GitHub.\nIf it does not open, visit:\n  %s\n", authURL)
	_ = browser.OpenURL(authURL)

	select {
	case <-ctx.Done():
		_ = srv.Close()
		return Result{}, fmt.Errorf("githubflow: timed out waiting for GitHub authorization: %w", ctx.Err())
	case res := <-resCh:
		// Graceful shutdown lets the in-flight callback handler finish delivering
		// its response page before the listener closes, so the user sees the
		// success/error page rather than a connection reset.
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		if res.err != nil {
			return Result{}, res.err
		}
		return Result{Code: res.code, RedirectURI: redirectURI, Verifier: verifier}, nil
	}
}
