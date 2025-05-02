/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"context"
	"os/exec"

	"golang.org/x/oauth2"
)

type Option func(*options)

type options struct {
	iss string
	aud string
}

func WithIssuer(iss string) Option {
	return func(opts *options) {
		opts.iss = iss
	}
}

func WithAudience(aud string) Option {
	return func(opts *options) {
		opts.aud = aud
	}
}

func NewChainctlTokenSource(ctx context.Context, opts ...Option) oauth2.TokenSource {
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	return &chainctlTokenSource{
		ctx:  ctx,
		opts: o,
	}
}

type chainctlTokenSource struct {
	ctx  context.Context
	opts *options
}

func (ts *chainctlTokenSource) Token() (*oauth2.Token, error) {
	args := []string{"auth", "token"}
	if ts.opts.iss != "" {
		args = append(args, "--issuer", ts.opts.iss)
	}
	if ts.opts.aud != "" {
		args = append(args, "--audience", ts.opts.aud)
	}

	cmd := exec.CommandContext(ts.ctx, "chainctl", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	t := string(out)
	exp, err := ExtractExpiry(t)
	if err != nil {
		return nil, err
	}
	return &oauth2.Token{
		AccessToken: t,
		Expiry:      exp,
	}, nil
}
