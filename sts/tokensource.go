/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
)

// NewTokenSource creates an oauth2.TokenSource by wrapping another TokenSource
// in a Chainguard STS exchange brokered by the provided Exchanger.
func NewTokenSource(ts oauth2.TokenSource, xchg Exchanger) oauth2.TokenSource {
	return &stsTokenSource{
		ctx:  context.Background(),
		ts:   ts,
		xchg: xchg,
	}
}

type stsTokenSource struct {
	ctx  context.Context
	ts   oauth2.TokenSource
	xchg Exchanger
}

// Token implements oauth2.TokenSource
func (sts *stsTokenSource) Token() (*oauth2.Token, error) {
	tok, err := sts.ts.Token()
	if err != nil {
		return nil, fmt.Errorf("fetching base token: %w", err)
	}
	idt, err := sts.xchg.Exchange(sts.ctx, tok.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("exchanging base token: %w", err)
	}
	return &oauth2.Token{
		AccessToken: idt.AccessToken,
		Expiry:      time.Now().Add(time.Hour),
	}, nil
}

func NewTokenSourceFromValues(ctx context.Context, aud string, identity string, ts oauth2.TokenSource) oauth2.TokenSource {
	return &stsTokenSource{
		ctx:  ctx,
		ts:   ts,
		xchg: New(aud, identity),
	}
}
