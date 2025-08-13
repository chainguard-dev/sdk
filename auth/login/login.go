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
	"strings"

	"github.com/pkg/browser"

	"github.com/chainguard-dev/clog"
)

// OpenBrowserError wraps the error returned from browser.OpenURL,
// since this can take a few different forms depending on the OS.
// Deprecated: use Error.
type OpenBrowserError struct {
	error
}

func (e OpenBrowserError) Error() string {
	if e.error == nil {
		return "login: failed to open browser"
	}
	return "login: failed to open browser: " + e.error.Error()
}

func (e OpenBrowserError) Unwrap() error { return e.error }

const (
	invalidConfigurationError = "invalid configuration"
	openBrowserError          = "failed to open browser"
	localServerError          = "failed to start localhost server"
	remoteServerError         = "error returned from server"
)

// Error is a generic wrapper around client side errors this package may return.
// It helps callers distinguish between local and remote errors.
type Error struct {
	Details string
	Err     error
}

func (e *Error) Error() string {
	if e.Err == nil && e.Details == "" {
		return "login: unknown error"
	}
	if e.Err == nil {
		return fmt.Sprintf("login: %s: unknown error", e.Details)
	}
	return fmt.Sprintf("login: %s: %s", e.Details, e.Err.Error())
}

func (e *Error) Unwrap() error { return e.Err }

func (e *Error) As(target any) bool {
	// For backwards compatibility, ensure Error can be cast as OpenBrowserError
	// under the right circumstances.
	//nolint:staticcheck
	if obe, ok := target.(**OpenBrowserError); ok && e.Details == openBrowserError {
		*obe = &OpenBrowserError{e.Err}
		return true
	}
	return false
}

func BuildHeadlessURL(opts ...Option) (u string, err error) {
	conf, err := newConfigFromOptions(opts...)
	if err != nil {
		return "", &Error{Details: invalidConfigurationError, Err: err}
	}
	if conf.HeadlessCode == "" {
		return "", &Error{Details: invalidConfigurationError, Err: fmt.Errorf("headless code is required")}
	}
	params := make(url.Values)
	params.Set("headless_code", conf.HeadlessCode)
	switch {
	case conf.IDP != "":
		params.Set("idp_id", conf.IDP)
	case conf.OrgName != "":
		// Verified org single sign-on – reuse idp_id for org name.
		params.Set("idp_id", conf.OrgName)
	case conf.Auth0Connection != "":
		// The connection param is only for Auth0 social logins.
		params.Set("connection", conf.Auth0Connection)
	}
	return fmt.Sprintf("%s/oauth?%s", conf.Issuer, params.Encode()), nil
}

func Login(ctx context.Context, opts ...Option) (token string, refreshToken string, err error) {
	conf, err := newConfigFromOptions(opts...)
	if err != nil {
		return "", "", &Error{Details: invalidConfigurationError, Err: err}
	}

	// Start new token server on a random available localhost port
	s, err := newServer(ctx)
	if err != nil {
		return "", "", &Error{Details: localServerError, Err: err}
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
	for _, scope := range conf.Scope {
		params.Add("scope", url.QueryEscape(scope))
	}

	u := fmt.Sprintf("%s/oauth?%s", conf.Issuer, params.Encode())
	clog.DebugContext(ctx, "Authenticating", "url", u)
	if conf.SkipBrowser {
		fmt.Fprintf(conf.MessageWriter, "Please open a browser to %s\n", u)
	} else {
		fmt.Fprintf(conf.MessageWriter, "Opening browser to %s\n", u)
		err = browser.OpenURL(u)
		if err != nil {
			return "", "", &Error{Details: openBrowserError, Err: err}
		}
	}
	token, err = s.Token()
	if err != nil {
		return "", "", err
	}
	refreshToken = s.RefreshToken()
	return token, refreshToken, nil
}
