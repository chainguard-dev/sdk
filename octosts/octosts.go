/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package octosts

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/compute/metadata"
	"github.com/chainguard-dev/clog"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	"chainguard.dev/sdk/auth"
	"chainguard.dev/sdk/sts"
)

const (
	OctoSTSEndpoint = "https://octo-sts.dev"
)

// Token mints a new octo sts token based on the policy for a given repo.
func Token(ctx context.Context, policyName, org, repo string) (string, error) {
	// To help enable local development, we allow the use of a GitHub token,
	// but *only when not running on GCE*.
	if tok := os.Getenv("GH_TOKEN"); tok != "" && !metadata.OnGCE() {
		clog.Warnf("using GH_TOKEN for token exchange")
		return tok, nil
	}
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" && !metadata.OnGCE() {
		clog.Warnf("using GITHUB_TOKEN for token exchange")
		return tok, nil
	}

	scope := org
	if repo != "" {
		scope = fmt.Sprintf("%s/%s", org, repo)
	}

	xchg := sts.New(
		OctoSTSEndpoint,
		policyName,
		sts.WithScope(scope),
		sts.WithIdentity(policyName),
	)

	ts, err := idtoken.NewTokenSource(ctx, "octo-sts.dev" /* aud */)
	if err != nil {
		return "", err
	}

	token, err := ts.Token()
	if err != nil {
		return "", err
	}

	tok, err := xchg.Exchange(ctx, token.AccessToken)
	if err != nil {
		return "", err
	}

	return tok.AccessToken, nil
}

// Revoke revokes the given security token.
func Revoke(ctx context.Context, tok string) error {
	req, err := http.NewRequest(http.MethodDelete, "https://api.github.com/installation/token", nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Add("Authorization", "Bearer "+tok)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// The token was revoked!
	return nil
}

// NewTokenSource creates an octoSTSTokenSource, similar to sts.NewTokenSource
func NewTokenSource(ts oauth2.TokenSource, xchg sts.Exchanger) oauth2.TokenSource {
	return &octoSTSTokenSource{
		ctx:  context.Background(),
		ts:   ts,
		xchg: xchg,
	}
}

// NewTokenSourceFromValues creates the exchanger and token source required to create an octoSTSTokenSource
// and then calls NewTokenSource with those values.
func NewTokenSourceFromValues(ctx context.Context, policyName, org, repo string) (oauth2.TokenSource, error) {
	scope := org
	if repo != "" {
		scope = fmt.Sprintf("%s/%s", org, repo)
	}

	xchg := sts.New(
		OctoSTSEndpoint,
		policyName,
		sts.WithScope(scope),
		sts.WithIdentity(policyName),
	)

	ts, err := idtoken.NewTokenSource(ctx, "octo-sts.dev" /* aud */)
	if err != nil {
		return nil, err
	}

	return NewTokenSource(ts, xchg), nil
}

type octoSTSTokenSource struct {
	ctx  context.Context
	ts   oauth2.TokenSource
	xchg sts.Exchanger
}

// Token implements oauth2.TokenSource
func (sts *octoSTSTokenSource) Token() (*oauth2.Token, error) {
	tok, err := sts.ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base token: %w", err)
	}

	idt, err := sts.xchg.Exchange(sts.ctx, tok.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange base token: %w", err)
	}

	accessToken := idt.AccessToken
	expiry, err := auth.ExtractExpiry(accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to extract expiry from access token: %w", err)
	}

	return &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      expiry,
	}, nil
}
