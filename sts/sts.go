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
	"runtime"
	"sync"
	"time"

	"chainguard.dev/sdk/auth"
	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"github.com/chainguard-dev/clog"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	Expiry       time.Time
}

// Exchanger is an interface for exchanging a third-party token for a Chainguard
// token.
type Exchanger interface {

	// Exchange performs the actual token exchange, sending "token" to the
	// Chainguard issuer's STS interface, and receiving bytes or an error.
	Exchange(ctx context.Context, token string, opts ...ExchangerOption) (TokenPair, error)

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
//
// Deprecated: use ExchangePair instead. This is kept around only until we migrate all existing callers
// to ExchangePair.
func Exchange(ctx context.Context, issuer, audience, idToken string, opts ...ExchangerOption) (string, error) {
	o, err := ExchangePair(ctx, issuer, audience, idToken, opts...)
	if err != nil {
		return "", err
	}
	return o.AccessToken, nil
}

// ExchangePair performs an OIDC token exchange with the correct Exchanger based on the provided options.
func ExchangePair(ctx context.Context, issuer, audience, idToken string, exchangerOptions ...ExchangerOption) (TokenPair, error) {
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
	tokenSet, err := e.Exchange(ctx, idToken, exchangerOptions...)
	if err != nil {
		return TokenPair{}, fmt.Errorf("exchanging token with %q: %w", issuer, err)
	}
	return tokenSet, nil
}

type impl struct {
	opts options

	once    sync.Once
	clients oidc.Clients
	err     error
}

type options struct {
	issuer           string
	audience         string
	userAgent        string
	firstScope       string
	scope            []string
	capabilities     []string
	identity         string
	http1Downgrade   bool
	identityProvider string
}

var _ Exchanger = (*impl)(nil)

// Stubbed when testing
var oidcNewClients = oidc.NewClients

func (i *impl) newClients() (oidc.Clients, error) {
	i.once.Do(func() {
		i.clients, i.err = oidcNewClients(i.opts.issuer, oidc.WithUserAgent(i.opts.userAgent))

		// This is dumb, but we've already exposed an API everywhere that is really hard to change.
		// Rather than making everyone defer xch.Close() on an Exchanger, we'll just attempt to rely
		// on the runtime to close the client whenever this impl goes out of scope.
		//
		// We wrap this thing in numerous interfaces that don't expose a Close() method like ggcr's
		// keychain or oauth2.TokenSource, so we really don't have a better option.
		runtime.AddCleanup(i, func(clients oidc.Clients) {
			defer clients.Close()
		}, i.clients)
	})

	return i.clients, i.err
}

func (i *impl) callOpts(ctx context.Context, token string) []grpc.CallOption {
	// TODO: we may want to require transport security at some future point.
	if cred := auth.NewFromToken(ctx, fmt.Sprintf("Bearer %s", token), false); cred != nil {
		return []grpc.CallOption{grpc.PerRPCCredentials(cred)}
	}

	clog.FromContext(ctx).Warn("No authentication provided, this may end badly.")
	return []grpc.CallOption{}
}

// Exchange implements Exchanger
func (i *impl) Exchange(ctx context.Context, token string, opts ...ExchangerOption) (TokenPair, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	c, err := i.newClients()
	if err != nil {
		return TokenPair{}, err
	}

	resp, err := c.STS().Exchange(ctx, &oidc.ExchangeRequest{
		Aud:              []string{o.audience},
		Scope:            o.firstScope, //nolint:staticcheck // Populating for backward compatibility
		Scopes:           o.scope,
		Identity:         o.identity,
		Cap:              o.capabilities,
		IdentityProvider: o.identityProvider,
	}, i.callOpts(ctx, token)...)
	if err != nil {
		return TokenPair{}, err
	}

	var expiry time.Time
	if resp.GetExpiry() != nil {
		expiry = resp.GetExpiry().AsTime()
	}
	return TokenPair{AccessToken: resp.Token, RefreshToken: resp.RefreshToken, Expiry: expiry}, nil
}

// Refresh implements Exchanger
func (i *impl) Refresh(ctx context.Context, token string, opts ...ExchangerOption) (string, string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	c, err := i.newClients()
	if err != nil {
		return "", "", err
	}

	resp, err := c.STS().ExchangeRefreshToken(ctx, &oidc.ExchangeRefreshTokenRequest{
		Aud:    []string{o.audience},
		Scope:  o.firstScope, //nolint:staticcheck // Populating for backward compatibility
		Scopes: o.scope,
		Cap:    o.capabilities,
	}, i.callOpts(ctx, token)...)
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

// WithScope sets the scope parameter sent by the Exchanger.
//
// Only one of cluster or scope may be set.
func WithScope(scope ...string) ExchangerOption {
	return func(i *options) {
		i.scope = scope
		// Capture the first scope for backward compatibility.
		// For callers who expect a single WithScope("scope") to
		// be included in the exchange request Scope field.
		if len(scope) > 0 {
			i.firstScope = scope[0]
		}
	}
}

// WithCapabilities sets the capabilities sent by the Exchanger.
func WithCapabilities(capabilities ...string) ExchangerOption {
	return func(i *options) {
		i.capabilities = capabilities
	}
}

// WithIdentity sets the the unique ID of the identity so that STS exchange can
// look up pre-stored verification keys without ambiguity
func WithIdentity(uid string) ExchangerOption {
	return func(i *options) {
		i.identity = uid
	}
}

// WithHTTP1Downgrade signals Exchange to use HTTP1DowngradeExchanger in the STS exchange.
func WithHTTP1Downgrade() ExchangerOption {
	return func(i *options) {
		i.http1Downgrade = true
	}
}

// WithIdentityProvider sets the identity provider to use for the exchange.
func WithIdentityProvider(idp string) ExchangerOption {
	return func(i *options) {
		i.identityProvider = idp
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

func (i *HTTP1DowngradeExchanger) Exchange(ctx context.Context, token string, opts ...ExchangerOption) (TokenPair, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}
	in := &oidc.ExchangeRequest{
		Aud:              []string{o.audience},
		Scope:            o.firstScope, //nolint:staticcheck // Populating for backward compatibility
		Scopes:           o.scope,
		Identity:         o.identity,
		Cap:              o.capabilities,
		IdentityProvider: o.identityProvider,
	}
	out := new(oidc.RawToken)
	if err := i.doHTTP1(ctx, token, "/sts/exchange", in, out, o); err != nil {
		return TokenPair{}, err
	}

	var expiry time.Time
	if out.GetExpiry() != nil {
		expiry = out.GetExpiry().AsTime()
	}
	return TokenPair{AccessToken: out.Token, RefreshToken: out.RefreshToken, Expiry: expiry}, nil
}

func (i *HTTP1DowngradeExchanger) Refresh(ctx context.Context, token string, opts ...ExchangerOption) (string, string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	in := &oidc.ExchangeRefreshTokenRequest{
		Aud:    []string{o.audience},
		Scope:  o.firstScope, //nolint:staticcheck // Populating for backward compatibility
		Scopes: o.scope,
		Cap:    o.capabilities,
	}

	out := new(oidc.TokenPair)
	if err := i.doHTTP1(ctx, token, "sts/access_token", in, out, o); err != nil {
		return "", "", err
	}

	return out.GetToken().Token, out.GetRefreshToken().Token, nil
}
