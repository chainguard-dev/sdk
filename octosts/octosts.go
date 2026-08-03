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
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/chainguard-dev/clog"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"

	"chainguard.dev/sdk/sts"
)

const (
	OctoSTSEndpoint        = "https://octo-sts.dev"
	defaultExchangeTimeout = 67 * time.Second
)

func exchangeTimeout() time.Duration {
	if v := os.Getenv("EXCHANGE_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultExchangeTimeout
}

// Token mints a new octo sts token based on the policy for a given repo.
func Token(ctx context.Context, policyName, org, repo string) (string, error) {
	ts, err := NewTokenSourceFromValues(ctx, policyName, org, repo)
	if err != nil {
		return "", fmt.Errorf("failed to create token source: %w", err)
	}

	tok, err := ts.Token()
	if err != nil {
		return "", fmt.Errorf("failed to fetch token: %w", err)
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

	resp, err := http.DefaultClient.Do(req) //nolint:gosec // G704: URL from GitHub API
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
	return NewTokenSourceContext(context.Background(), ts, xchg)
}

// NewTokenSource creates an octoSTSTokenSource, similar to sts.NewTokenSource.
// The context is used for all token exchanges for the lifetime of the token source,
// so it should be long-lived.
func NewTokenSourceContext(ctx context.Context, ts oauth2.TokenSource, xchg sts.Exchanger) oauth2.TokenSource {
	return &octoSTSTokenSource{
		ctx:  ctx,
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

	// To help enable local development, we allow the use of a GitHub token,
	// but *only when not running on GCE*.
	if tok := os.Getenv("GH_TOKEN"); tok != "" && !metadata.OnGCE() {
		clog.WarnContextf(ctx, "using GH_TOKEN for token exchange")
		return oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: tok,
		}), nil
	}
	if tok := os.Getenv("GITHUB_TOKEN"); tok != "" && !metadata.OnGCE() {
		clog.WarnContextf(ctx, "using GITHUB_TOKEN for token exchange")
		return oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: tok,
		}), nil
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

	return oauth2.ReuseTokenSource(nil, NewTokenSourceContext(ctx, ts, xchg)), nil
}

type octoSTSTokenSource struct {
	ctx  context.Context
	ts   oauth2.TokenSource
	xchg sts.Exchanger
}

// Token implements oauth2.TokenSource
func (sts *octoSTSTokenSource) Token() (*oauth2.Token, error) {
	ctx := sts.ctx
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, exchangeTimeout())
		defer cancel()
	}

	tok, err := sts.ts.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch base token: %w", err)
	}

	idt, err := sts.xchg.Exchange(ctx, tok.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange base token: %w", err)
	}

	accessToken := idt.AccessToken

	exp := idt.Expiry
	if exp.IsZero() {
		// If the Exchanger doesn't provide an expiry time, we set a default expiry time of 20 minutes from now.
		// This is an approximation, as we don't have the actual expiry time from the Exchanger.
		// Tokens are usually valid for 1 hour, but we refresh at the 20-minute mark so we
		// pick up fresh tokens more often, rather than staying stuck with a noisy quota
		// neighbor for the full hour.
		exp = time.Now().Add(20 * time.Minute)
	}

	return &oauth2.Token{
		AccessToken: accessToken,
		Expiry:      exp,
	}, nil
}
