/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package deviceflow

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sigstore/sigstore/pkg/oauthflow"
	"golang.org/x/oauth2"
)

// Forked from: https://github.com/sigstore/sigstore/blob/8cd960fb1915c526bd838df6341b027634434985/pkg/oauthflow/device.go
// Changes from source:
// - Remove deprecated functions.
// - Remove PKCE since Auth0 doesn't support it for device flow.
// - Add `client_id` to token endpoint polling requests, since Auth0 requires it.

type deviceResp struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	Interval                int    `json:"interval"`
	ExpiresIn               int    `json:"expires_in"`
}

type tokenResp struct {
	IDToken string `json:"id_token"`
	Error   string `json:"error"`
}

var _ oauthflow.TokenGetter = (*TokenGetter)(nil)

// TokenGetter fetches an OIDC Identity token using the Device Code Grant flow as specified in RFC8628
type TokenGetter struct {
	messagePrinter func(string)
	sleeper        func(time.Duration)
	issuer         string
	codeURL        string
}

type Option func(tg *TokenGetter)

func WithMessagePrinter(fn func(string)) Option {
	return func(tg *TokenGetter) {
		tg.messagePrinter = fn
	}
}

func WithSleeper(fn func(time.Duration)) Option {
	return func(tg *TokenGetter) {
		tg.sleeper = fn
	}
}

// NewTokenGetter creates a new TokenGetter that retrieves an OIDC Identity Token using a Device Code Grant
func NewTokenGetter(issuer string, opts ...Option) *TokenGetter {
	tg := &TokenGetter{
		messagePrinter: func(s string) { fmt.Fprintln(os.Stderr, s) },
		sleeper:        time.Sleep,
		issuer:         issuer,
	}

	for _, opt := range opts {
		opt(tg)
	}
	return tg
}

func (d *TokenGetter) deviceFlow(p *oidc.Provider, clientID, redirectURL string, scopes []string) (string, error) {
	data := url.Values{
		"client_id": []string{clientID},
		"scope":     []string{strings.Join(scopes, " ")},
	}
	if redirectURL != "" {
		// If a redirect uri is provided then use it
		data["redirect_uri"] = []string{redirectURL}
	}

	codeURL, err := d.CodeURL()
	if err != nil {
		return "", err
	}
	/* #nosec */
	resp, err := http.PostForm(codeURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s", resp.Status, b)
	}

	parsed := deviceResp{}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return "", err
	}

	d.messagePrinter(fmt.Sprintf("Enter the verification code %s in your browser at: %s", parsed.UserCode, parsed.VerificationURI))
	d.messagePrinter(fmt.Sprintf("Code will be valid for %d seconds", parsed.ExpiresIn))
	d.sleeper(time.Duration(parsed.Interval) * time.Second)

	for {
		data := url.Values{
			"client_id":   []string{clientID},
			"grant_type":  []string{"urn:ietf:params:oauth:grant-type:device_code"},
			"device_code": []string{parsed.DeviceCode},
		}

		/* #nosec */
		resp, err := http.PostForm(p.Endpoint().TokenURL, data)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		tr := tokenResp{}
		if err := json.Unmarshal(b, &tr); err != nil {
			return "", err
		}

		if tr.IDToken != "" {
			d.messagePrinter("Token received!")
			return tr.IDToken, nil
		}
		switch tr.Error {
		case "access_denied", "expired_token":
			return "", fmt.Errorf("error obtaining token: %s", tr.Error)
		case "authorization_pending":
			d.sleeper(time.Duration(parsed.Interval) * time.Second)
		case "slow_down":
			// Add ten seconds if we got told to slow down
			d.sleeper(time.Duration(parsed.Interval)*time.Second + 10*time.Second)
		default:
			return "", fmt.Errorf("unexpected error in device flow: %s", tr.Error)
		}
	}
}

// GetIDToken gets an OIDC ID Token from the specified provider using the device code grant flow
func (d *TokenGetter) GetIDToken(p *oidc.Provider, cfg oauth2.Config) (*oauthflow.OIDCIDToken, error) {
	idToken, err := d.deviceFlow(p, cfg.ClientID, cfg.RedirectURL, cfg.Scopes)
	if err != nil {
		return nil, err
	}
	verifier := p.Verifier(&oidc.Config{ClientID: cfg.ClientID})
	parsedIDToken, err := verifier.Verify(context.Background(), idToken)
	if err != nil {
		return nil, err
	}

	subj, err := oauthflow.SubjectFromToken(parsedIDToken)
	if err != nil {
		return nil, err
	}

	return &oauthflow.OIDCIDToken{
		RawString: idToken,
		Subject:   subj,
	}, nil
}

// CodeURL fetches the device authorization endpoint URL from the provider's well-known configuration endpoint
func (d *TokenGetter) CodeURL() (string, error) {
	if d.codeURL != "" {
		return d.codeURL, nil
	}

	wellKnown := strings.TrimSuffix(d.issuer, "/") + "/.well-known/openid-configuration"
	/* #nosec */
	httpClient := &http.Client{
		Timeout: 3 * time.Second,
	}
	resp, err := httpClient.Get(wellKnown)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("unable to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("%s: %s", resp.Status, body)
	}

	providerConfig := struct {
		Issuer         string `json:"issuer"`
		DeviceEndpoint string `json:"device_authorization_endpoint"`
	}{}
	if err = json.Unmarshal(body, &providerConfig); err != nil {
		return "", fmt.Errorf("oidc: failed to decode provider discovery object: %w", err)
	}

	if d.issuer != providerConfig.Issuer {
		return "", fmt.Errorf("oidc: issuer did not match the issuer returned by provider, expected %q got %q", d.issuer, providerConfig.Issuer)
	}

	if providerConfig.DeviceEndpoint == "" {
		return "", fmt.Errorf("oidc: device authorization endpoint not returned by provider")
	}

	d.codeURL = providerConfig.DeviceEndpoint
	return d.codeURL, nil
}
