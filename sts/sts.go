/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"

	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const GulfstreamAudience = "gulfstream"

// Exchanger is an interface for exchanging a third-party token for a Chainguard
// token.
type Exchanger interface {
	// Exchange performs the actual token exchange, sending "token" to the
	// Chainguard issuer's STS interface, and receiving bytes or an error.
	Exchange(ctx context.Context, token string, opts ...ExchangerOption) (string, error)

	// Refresh exchanges a refresh token for a new access token and refresh token.
	Refresh(ctx context.Context, token string, opts ...ExchangerOption) (accessToken string, refreshToken string, err error)
}

// New creates a new Exchanger that works against the provided issuer's STS
// endpoint, and requests token with the specified audience.  It's behavior
// can be further customized via optional ExchangerOption parameters.
func New(issuer, audience string, opts ...ExchangerOption) Exchanger {
	i := &impl{
		opts: options{
			issuer:   issuer,
			audience: audience,
		},
	}
	for _, opt := range opts {
		opt(&i.opts)
	}

	return i
}

// Exchange performs an OIDC token exchange with the correct Exchanger based on the provided options.
func Exchange(ctx context.Context, issuer, audience, idToken string, exchangerOptions ...ExchangerOption) (string, error) {
	opts := options{
		issuer:   issuer,
		audience: audience,
	}
	for _, eo := range exchangerOptions {
		eo(&opts)
	}

	var e Exchanger
	if opts.http1Downgrade {
		e = &HTTP1DowngradeExchanger{opts: opts}
	} else {
		e = &impl{opts: opts}
	}
	return e.Exchange(ctx, idToken, exchangerOptions...)
}

type impl struct {
	opts options
}

type options struct {
	issuer               string
	audience             string
	cluster              string
	userAgent            string
	scope                string
	capabilities         []string
	identity             string
	includeUpstreamToken bool
	http1Downgrade       bool
}

var _ Exchanger = (*impl)(nil)

// Stubbed when testing
var oidcNewClients = oidc.NewClients

// Exchange implements Exchanger
func (i *impl) Exchange(ctx context.Context, token string, opts ...ExchangerOption) (string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	c, err := oidcNewClients(ctx, o.issuer, fmt.Sprintf("Bearer %s", token), oidc.WithUserAgent(o.userAgent))
	if err != nil {
		return "", err
	}
	defer c.Close()

	resp, err := c.STS().Exchange(ctx, &oidc.ExchangeRequest{
		Aud:                  []string{o.audience},
		Scope:                o.scope,
		Cluster:              o.cluster,
		Identity:             o.identity,
		IncludeUpstreamToken: o.includeUpstreamToken,
		Cap:                  o.capabilities,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// Refresh implements Exchanger
func (i *impl) Refresh(ctx context.Context, token string, opts ...ExchangerOption) (string, string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	c, err := oidcNewClients(ctx, o.issuer, fmt.Sprintf("Bearer %s", token), oidc.WithUserAgent(o.userAgent))
	if err != nil {
		return "", "", err
	}
	defer c.Close()

	resp, err := c.STS().ExchangeAccessToken(ctx, &oidc.ExchangeAccessTokenRequest{
		Aud:   []string{o.audience},
		Scope: o.scope,
		Cap:   o.capabilities,
	})
	if err != nil {
		return "", "", err
	}
	return resp.GetToken().GetToken(), resp.GetRefreshToken().GetToken(), err
}

// ExchangerOption is a way of customizing the behavior of the Exchanger
// constructed via New()
type ExchangerOption func(*options)

// WithUserAgent sets the user agent sent by the Exchanger.
func WithUserAgent(agent string) ExchangerOption {
	return func(i *options) {
		i.userAgent = agent
	}
}

// WithCluster sets the cluster parameter sent by the Exchanger.
//
// Only one of cluster or scope may be set.
func WithCluster(cluster string) ExchangerOption {
	return func(i *options) {
		i.cluster = cluster
	}
}

// WithScope sets the scope parameter sent by the Exchanger.
//
// Only one of cluster or scope may be set.
func WithScope(scope string) ExchangerOption {
	return func(i *options) {
		i.scope = scope
	}
}

// WithCapabilities sets the capabilities sent by the Exchanger.
func WithCapabilities(cap ...string) ExchangerOption {
	return func(i *options) {
		i.capabilities = cap
	}
}

// WithIdentity sets the the unique ID of the identity so that STS exchange can
// look up pre-stored verification keys without ambiguity
func WithIdentity(uid string) ExchangerOption {
	return func(i *options) {
		i.identity = uid
	}
}

// WithIncludeUpstreamToken requests that the upstream token be included in the returned
// STS token.
func WithIncludeUpstreamToken() ExchangerOption {
	return func(i *options) {
		i.includeUpstreamToken = true
	}
}

// WithHTTP1Downgrade signals Exchange to use HTTP1DowngradeExchanger in the STS exchange.
func WithHTTP1Downgrade() ExchangerOption {
	return func(i *options) {
		i.http1Downgrade = true
	}
}

type HTTP1DowngradeExchanger struct {
	opts options
}

func NewHTTP1DowngradeExchanger(issuer, audience string, opts ...ExchangerOption) *HTTP1DowngradeExchanger {
	i := &HTTP1DowngradeExchanger{
		opts: options{
			issuer:   issuer,
			audience: audience,
		},
	}
	for _, opt := range opts {
		opt(&i.opts)
	}
	return i
}

func (i *HTTP1DowngradeExchanger) doHTTP1(ctx context.Context,
	auth string,
	path string, in proto.Message, out proto.Message, opts options) error {
	body, err := protojson.Marshal(in)
	if err != nil {
		return err
	}
	u, err := url.JoinPath(i.opts.issuer, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if opts.userAgent != "" {
		req.Header.Set("User-Agent", opts.userAgent)
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		for k, v := range md {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
	}

	// Explicitly disable HTTP/2 support by setting the
	// client Transport's TLSNextProto to an empty map.
	// ref: https://pkg.go.dev/net/http#hdr-HTTP_2
	client := &http.Client{
		Transport: &oauth2.Transport{
			Base: &http.Transport{
				TLSNextProto: map[string]func(string, *tls.Conn) http.RoundTripper{},
			},
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth}),
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return protojson.Unmarshal(b, out)
}

func (i *HTTP1DowngradeExchanger) Exchange(ctx context.Context, token string, opts ...ExchangerOption) (string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}
	in := &oidc.ExchangeRequest{
		Aud:                  []string{o.audience},
		Scope:                o.scope,
		Cluster:              o.cluster,
		Identity:             o.identity,
		IncludeUpstreamToken: o.includeUpstreamToken,
		Cap:                  o.capabilities,
	}
	out := new(oidc.RawToken)
	if err := i.doHTTP1(ctx, token, "/sts/exchange", in, out, o); err != nil {
		return "", err
	}
	return out.Token, nil
}

func (i *HTTP1DowngradeExchanger) Refresh(ctx context.Context, token string, opts ...ExchangerOption) (string, string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	in := &oidc.ExchangeAccessTokenRequest{
		Aud:   []string{o.audience},
		Scope: o.scope,
		Cap:   o.capabilities,
	}

	out := new(oidc.TokenPair)
	if err := i.doHTTP1(ctx, token, "sts/access_token", in, out, o); err != nil {
		return "", "", err
	}

	return out.GetToken().Token, out.GetRefreshToken().Token, nil
}
