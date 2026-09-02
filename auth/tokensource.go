/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/oauth2"
)

// maxStderrBytes caps how much of a failed chainctl invocation's stderr we put
// in the returned error, so a runaway command cannot blow up the message.
const maxStderrBytes = 2048

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
		return nil, commandError(cmd, err)
	}
	t := string(out)
	exp, err := ExtractExpiry(t)
	if err != nil {
		return nil, fmt.Errorf("reading the expiry of the token printed by %q: %w", strings.Join(cmd.Args, " "), err)
	}
	return &oauth2.Token{
		AccessToken: t,
		Expiry:      exp,
	}, nil
}

// commandError describes a failed chainctl invocation.
//
// exec.Cmd.Output captures the command's stderr in exec.ExitError.Stderr, but
// the Error method of that type reports only the exit status. A caller that
// prints the bare error therefore learns that chainctl failed and nothing about
// why, so this pulls the stderr into the message.
//
// The stderr goes through %q because it is output from another program that a
// caller writes to a terminal or a log: quoting escapes any control or ANSI
// bytes in it. The command's stdout never appears in the message — on the
// success path stdout is the bearer token itself.
func commandError(cmd *exec.Cmd, err error) error {
	if exitErr, ok := errors.AsType[*exec.ExitError](err); ok {
		if stderr := strings.TrimSpace(string(exitErr.Stderr)); stderr != "" {
			if len(stderr) > maxStderrBytes {
				stderr = stderr[:maxStderrBytes] + " (truncated)"
			}
			return fmt.Errorf("running %q: %w: %q", strings.Join(cmd.Args, " "), err, stderr)
		}
	}
	return fmt.Errorf("running %q: %w", strings.Join(cmd.Args, " "), err)
}
