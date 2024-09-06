/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package login implements client login functionality shared between various
// clients
package login

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/browser"
)

// OpenBrowserError wraps the error returned from browser.OpenURL,
// since this can take a few different forms depending on the OS.
type OpenBrowserError struct {
	err error
}

func (e *OpenBrowserError) Error() string { return "failed to open browser: " + e.err.Error() }
func (e *OpenBrowserError) Unwrap() error { return e.err }

func BuildHeadlessURL(opts ...Option) (u string, err error) {
	conf, err := newConfigFromOptions(opts...)
	if err != nil {
		return "", err
	}
	if conf.HeadlessCode == "" {
		return "", fmt.Errorf("headless code is required")
	}
	params := make(url.Values)
	if conf.IDP != "" {
		params.Set("idp_id", conf.IDP)
	}
	params.Set("headless_code", conf.HeadlessCode)
	return fmt.Sprintf("%s/oauth?%s", conf.Issuer, params.Encode()), nil
}

func Login(ctx context.Context, opts ...Option) (token string, refreshToken string, err error) {
	conf, err := newConfigFromOptions(opts...)
	if err != nil {
		return "", "", err
	}

	// Start new token server on a random available localhost port
	s, err := newServer(ctx)
	if err != nil {
		return "", "", err
	}
	defer s.Close()

	params := make(url.Values)
	params.Set("exit", "redirect")
	params.Set("redirect", s.URL())
	if conf.IDP != "" {
		params.Set("idp_id", conf.IDP)
	}
	if conf.OrgName != "" {
		// NB: we reuse the idp_id query param for verified organization SSO
		params.Set("idp_id", conf.OrgName)
	}
	if conf.InviteCode != "" {
		params.Set("invite", conf.InviteCode)
	}
	if conf.ClientID != "" {
		params.Set("client_id", conf.ClientID)
	}
	if conf.Auth0Connection != "" {
		params.Set("connection", conf.Auth0Connection)
	}
	if conf.SkipRegistration {
		params.Set("skip_registration", "true")
	}
	if conf.Identity != "" {
		params.Set("identity", conf.Identity)
	}
	if len(conf.Audience) > 0 {
		params.Set("audience", strings.Join(conf.Audience, ","))
	}
	if conf.CreateRefreshToken {
		params.Set("create_refresh_token", "true")
	}
	u := fmt.Sprintf("%s/oauth?%s", conf.Issuer, params.Encode())
	if conf.SkipBrowser {
		fmt.Fprintf(os.Stderr, "Please open a browser to %s\n", u)
	} else {
		fmt.Fprintf(os.Stderr, "Opening browser to %s\n", u)
		err = browser.OpenURL(u)
		if err != nil {
			return "", "", &OpenBrowserError{err}
		}
	}
	token, err = s.Token()
	if err != nil {
		return "", "", err
	}
	refreshToken = s.RefreshToken()
	return token, refreshToken, nil
}
