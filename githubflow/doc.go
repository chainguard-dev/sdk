/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package githubflow runs the GitHub App user (user-to-server) OAuth
// "loopback" flow shared by clients such as chainctl and the guardener
// OAuth test binary.
//
// It starts a temporary HTTP server on 127.0.0.1, opens the browser to
// GitHub's authorize endpoint with PKCE (S256), and captures the
// authorization code. The code is exchanged for a token server-side,
// so this package never holds the GitHub App client secret or the
// resulting user token — it only returns the code, the loopback
// redirect URI it used, and the PKCE verifier, all of which the server
// replays at exchange time.
package githubflow
