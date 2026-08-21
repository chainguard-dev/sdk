/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/chainguard-dev/clog"

	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"golang.org/x/oauth2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type TokenPair struct {
	AccessToken  string //nolint:gosec // G117: struct field name for STS token response, not a hardcoded secret
	RefreshToken string //nolint:gosec // G117: struct field name for STS token response, not a hardcoded secret
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
// endpoint, and requests tokens with the specified audience(s). The audience
// parameter may contain multiple comma-separated audience values. Its behavior
// can be further customized via optional ExchangerOption parameters.
func New(issuer, audience string, opts ...ExchangerOption) Exchanger {
	i := &impl{
		opts: options{
			issuer:    issuer,
			audiences: strings.Split(audience, ","),
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
		issuer:    issuer,
		audiences: strings.Split(audience, ","),
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
}

type options struct {
	issuer           string
	audiences        []string
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

const (
	// maxRetries is the maximum number of retry attempts for transient errors.
	// With the initial attempt, this means up to 3 total attempts.
	maxRetries = 2

	// retryBackoff is the delay between retry attempts. Sized for sub-second
	// transient network blips (TCP reset, DNS hiccup), not multi-second
	// cold starts.
	retryBackoff = 100 * time.Millisecond
)

// retryable returns true if the gRPC error code indicates a transient failure
// that is safe to retry. Only codes.Unavailable is retried — codes.Internal
// is excluded because it can indicate non-transient server bugs (nil-pointer
// panics, DB inconsistency) that won't resolve on retry.
func retryable(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable:
		return true
	default:
		return false
	}
}

// refreshRetryable returns true for gRPC codes that are transient on the
// refresh-token exchange path. Unlike retryable() (used by Exchange), this also
// retries codes.Internal and codes.DeadlineExceeded: under STS/datastore
// saturation the refresh RPC surfaces a failed identity lookup as
// codes.Internal, and a server-side query/connection deadline surfaces as
// codes.DeadlineExceeded. Both are transient here and safe to retry on a fresh
// connection. codes.Canceled is intentionally excluded — it signals the caller
// (or parent context) gave up, where a retry is pointless.
//
// Note: a client-side deadline (e.g. the caller's own context timeout) is
// caught by the loop's ctx.Done() guard before another attempt is issued, so
// the DeadlineExceeded entry effectively targets server-side query timeouts
// only — where the caller's context is still alive. See CUS-839.
func refreshRetryable(err error) bool {
	switch status.Code(err) {
	case codes.Unavailable, codes.Internal, codes.DeadlineExceeded:
		return true
	default:
		return false
	}
}

// gatewayStatusCode reverses the issuer grpc-gateway's gRPC-code -> HTTP-status
// mapping (grpc-gateway's runtime.HTTPStatusFromCode, which this issuer uses
// with no custom error handler) so the HTTP/1-downgrade path can hand the retry
// predicates the same gRPC code the native gRPC path would have seen.
//
// The mapping is deliberately narrow — only statuses whose origin code is
// retry-relevant are recovered:
//
//   - 503 Service Unavailable  -> Unavailable      (retried by Exchange + Refresh)
//   - 504 Gateway Timeout      -> DeadlineExceeded (retried by Refresh only)
//   - 502 Bad Gateway          -> Unavailable      (an intermediary LB/proxy is
//     unreachable; grpc-gateway never emits 502, so this is unambiguously an
//     infra transient)
//
// 500 is intentionally NOT recovered: grpc-gateway collapses Internal, Unknown,
// and DataLoss all onto 500, so the origin code is ambiguous. Recovering it as
// Unavailable (the prior behavior) made Exchange retry an Internal that
// retryable() deliberately excludes as a possible non-transient server bug.
// Leaving it codes.Unknown means neither predicate retries it — matching
// Exchange's gRPC behavior and erring against amplifying a real server fault.
// Every other status (4xx, 429, ...) is likewise left codes.Unknown.
func gatewayStatusCode(httpStatus int) codes.Code {
	switch httpStatus {
	case http.StatusServiceUnavailable, http.StatusBadGateway:
		return codes.Unavailable
	case http.StatusGatewayTimeout:
		return codes.DeadlineExceeded
	default:
		return codes.Unknown
	}
}

// codeFromGatewayStatusBody recovers the gRPC status code the issuer's
// grpc-gateway serialized into the error body. Its DefaultHTTPErrorHandler
// writes a google.rpc.Status as JSON ({"code":<int>,"message":...}); that
// numeric code is exact, where the HTTP status is not — the gateway collapses
// Internal, Unknown, and DataLoss all onto 500. Keying on the body therefore
// lets the Exchange/Refresh predicates make the same decision the native gRPC
// path would (e.g. Refresh retries a body-Internal that a 500 alone would hide).
//
// Returns codes.OK (the zero value, "no usable code") when the body is not a
// grpc-gateway status — e.g. a raw 5xx page from an intermediary LB/proxy — so
// the caller falls back to the HTTP-status reversal (gatewayStatusCode).
func codeFromGatewayStatusBody(body []byte) codes.Code {
	// Pointer distinguishes an absent "code" field from an explicit 0.
	var s struct {
		Code *int32 `json:"code"`
	}
	if err := json.Unmarshal(body, &s); err != nil || s.Code == nil {
		return codes.OK
	}
	// Range-check on the int32 domain before converting to codes.Code (a
	// uint32): accept only the canonical gRPC codes Canceled..Unauthenticated
	// (1..16), which rejects OK-on-an-error, negatives, and out-of-range garbage
	// from a hostile or non-gateway body.
	n := *s.Code
	if n < int32(codes.Canceled) || n > int32(codes.Unauthenticated) {
		return codes.OK
	}
	return codes.Code(n) //nolint:gosec // G115: n range-checked to [1,16] above
}

// backoffWithJitter returns retryBackoff with uniform jitter in
// [retryBackoff, 2*retryBackoff). The jitter avoids synchronized retry bursts
// (thundering herd) when many clients receive the same transient failure at
// once — likely under the STS/datastore saturation that produces retryable
// codes.Internal in the first place. See CUS-839.
//
// Both Exchange and Refresh use this. Refresh's retry *codes* widened for
// CUS-839; Exchange's retry codes are unchanged, but its backoff intentionally
// moved from a flat retryBackoff to this jittered value here — a strict
// improvement (decorrelates retries on a saturated path), called out so the
// timing change on the hot exchange path is not a surprise.
func backoffWithJitter() time.Duration {
	return retryBackoff + rand.N(retryBackoff) //nolint:gosec // G404: retry-backoff jitter does not require cryptographic randomness
}

// Exchange implements Exchanger
func (i *impl) Exchange(ctx context.Context, token string, opts ...ExchangerOption) (TokenPair, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	var lastErr error
	for attempt := range maxRetries + 1 {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return TokenPair{}, ctx.Err()
			case <-time.After(backoffWithJitter()):
			}
		}

		c, err := oidcNewClients(ctx, o.issuer, fmt.Sprintf("Bearer %s", token), oidc.WithUserAgent(o.userAgent))
		if err != nil {
			if retryable(err) {
				lastErr = err
				continue
			}
			return TokenPair{}, err
		}

		resp, err := c.STS().Exchange(ctx, &oidc.ExchangeRequest{
			Aud:              o.audiences,
			Scope:            o.firstScope, //nolint:staticcheck // Populating for backward compatibility
			Scopes:           o.scope,
			Identity:         o.identity,
			Cap:              o.capabilities,
			IdentityProvider: o.identityProvider,
		})
		c.Close()
		if err != nil {
			if retryable(err) {
				lastErr = err
				continue
			}
			return TokenPair{}, err
		}

		var expiry time.Time
		if resp.GetExpiry() != nil {
			expiry = resp.GetExpiry().AsTime()
		}
		return TokenPair{AccessToken: resp.Token, RefreshToken: resp.RefreshToken, Expiry: expiry}, nil
	}
	return TokenPair{}, lastErr
}

// Refresh implements Exchanger
func (i *impl) Refresh(ctx context.Context, token string, opts ...ExchangerOption) (string, string, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}

	var lastErr error
	for attempt := range maxRetries + 1 {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", "", ctx.Err()
			case <-time.After(backoffWithJitter()):
			}
		}

		c, err := oidcNewClients(ctx, o.issuer, fmt.Sprintf("Bearer %s", token), oidc.WithUserAgent(o.userAgent))
		if err != nil {
			if refreshRetryable(err) {
				lastErr = err
				clog.FromContext(ctx).Warnf("refresh token exchange attempt %d failed connecting (%v); retrying", attempt+1, status.Code(err))
				continue
			}
			return "", "", err
		}

		resp, err := c.STS().ExchangeRefreshToken(ctx, &oidc.ExchangeRefreshTokenRequest{
			Aud:    o.audiences,
			Scope:  o.firstScope, //nolint:staticcheck // Populating for backward compatibility
			Scopes: o.scope,
			Cap:    o.capabilities,
		})
		c.Close()
		if err != nil {
			if refreshRetryable(err) {
				lastErr = err
				clog.FromContext(ctx).Warnf("refresh token exchange attempt %d failed (%v); retrying", attempt+1, status.Code(err))
				continue
			}
			return "", "", err
		}
		return resp.GetToken().GetToken(), resp.GetRefreshToken().GetToken(), nil
	}
	return "", "", lastErr
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
	// backoffTimer returns a channel that fires after a jittered retry backoff.
	// A per-instance field (not a package global) so a test can drive
	// retryHTTP1's ctx.Done() branch deterministically on its own exchanger
	// without racing other tests under t.Parallel(). Optional: the zero value
	// (nil) is safe — retryHTTP1 reads it through nextBackoff, which supplies
	// the real jittered timer when unset. This matters because ExchangePair
	// builds this struct directly, bypassing the constructor. See backoffWithJitter.
	backoffTimer func() <-chan time.Time
}

// nextBackoff returns the retry-backoff channel, defaulting to the real
// jittered timer when backoffTimer is unset. Every construction path
// (NewHTTP1DowngradeExchanger and the direct literal in ExchangePair) is
// therefore safe against a nil field; only tests override it.
func (i *HTTP1DowngradeExchanger) nextBackoff() <-chan time.Time {
	if i.backoffTimer != nil {
		return i.backoffTimer()
	}
	return time.After(backoffWithJitter())
}

func NewHTTP1DowngradeExchanger(issuer, audience string, opts ...ExchangerOption) *HTTP1DowngradeExchanger {
	i := &HTTP1DowngradeExchanger{
		opts: options{
			issuer:    issuer,
			audiences: strings.Split(audience, ","),
		},
	}
	for _, opt := range opts {
		opt(&i.opts)
	}
	return i
}

// newHTTP1Transport returns http.DefaultTransport with HTTP/2 disabled and
// nothing else changed. The downgrade transport must match the default
// client's behavior in every other respect — in particular proxy support
// (HTTP_PROXY/HTTPS_PROXY via ProxyFromEnvironment), which the customers
// this exchanger exists for typically depend on. Cloning rather than
// constructing field-by-field keeps that property across Go versions.
func newHTTP1Transport() *http.Transport {
	var t *http.Transport
	if dt, ok := http.DefaultTransport.(*http.Transport); ok {
		t = dt.Clone()
	} else {
		// Consumers may legitimately replace http.DefaultTransport with a
		// wrapping RoundTripper (instrumentation, test doubles); this public
		// SDK must not depend on its concrete type. Mirror the stdlib
		// defaults instead — proxy support is the load-bearing property.
		t = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}
	t.ForceAttemptHTTP2 = false
	// Disable HTTP/2 by setting TLSNextProto to a non-nil empty map.
	// ref: https://pkg.go.dev/net/http#hdr-HTTP_2
	t.TLSNextProto = map[string]func(string, *tls.Conn) http.RoundTripper{}
	// If the process has already used http.DefaultTransport for an HTTPS
	// request, its TLS config advertises "h2" via ALPN and the clone inherits
	// that; a server selecting h2 would then speak HTTP/2 frames at this
	// HTTP/1.1-only transport. Advertise exactly what we speak.
	if t.TLSClientConfig == nil {
		t.TLSClientConfig = &tls.Config{}
	}
	t.TLSClientConfig.NextProtos = []string{"http/1.1"}
	return t
}

// doHTTP1 sends the request fields as a form-encoded POST body. The gateway
// bindings for the STS endpoints carry no `body` mapping, so the server
// populates the request message via ParseForm + PopulateQueryParameters and
// never reads a JSON body — form encoding is the contract the deployed
// issuer honors (see TestGateway_ExchangeFormBody in the oidc module).
func (i *HTTP1DowngradeExchanger) doHTTP1(ctx context.Context,
	auth string,
	path string, form url.Values, out proto.Message, userAgent string) error {
	u, err := url.JoinPath(i.opts.issuer, path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	if md, ok := metadata.FromOutgoingContext(ctx); ok {
		for k, v := range md {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
	}

	client := &http.Client{
		Transport: &oauth2.Transport{
			Base:   newHTTP1Transport(),
			Source: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth}),
		},
		// STS endpoints never redirect; refusing to follow prevents
		// oauth2.Transport from re-sending the bearer token to a redirect
		// target (golang/oauth2#314).
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := client.Do(req) //nolint:gosec // G704: URL from internal STS config
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// Surface the error body (bounded); STS returns the actual failure
		// reason there, and the status line alone is not actionable.
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		detail := resp.Status
		if msg := bytes.TrimSpace(b); len(msg) > 0 {
			detail = fmt.Sprintf("%s: %s", resp.Status, msg)
		}
		// Tag the error with the gRPC status code so the retry loops in
		// Exchange/Refresh — which key on it via retryable/refreshRetryable —
		// make the *same* decision they would on the native gRPC path. Prefer
		// the exact code grpc-gateway serialized into the body; fall back to
		// reversing the HTTP status for a 5xx from an intermediary (LB/proxy)
		// that isn't a gateway status body. codes.Unknown from either leaves the
		// error a plain errors.New, which neither predicate retries. See CUS-1189.
		code := codeFromGatewayStatusBody(b)
		if code == codes.OK {
			code = gatewayStatusCode(resp.StatusCode)
		}
		if code != codes.Unknown {
			return status.Error(code, detail)
		}
		return errors.New(detail)
	}

	// Bounded like the error path; token-pair responses are a few KiB.
	b, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return err
	}
	return protojson.Unmarshal(b, out)
}

func (i *HTTP1DowngradeExchanger) Exchange(ctx context.Context, token string, opts ...ExchangerOption) (TokenPair, error) {
	o := i.opts
	for _, opt := range opts {
		opt(&o)
	}
	form := url.Values{}
	for _, aud := range o.audiences {
		form.Add("aud", aud)
	}
	if o.firstScope != "" {
		// Populating for backward compatibility.
		form.Set("scope", o.firstScope)
	}
	for _, scope := range o.scope {
		form.Add("scopes", scope)
	}
	if o.identity != "" {
		form.Set("identity", o.identity)
	}
	for _, cap := range o.capabilities {
		form.Add("cap", cap)
	}
	if o.identityProvider != "" {
		form.Set("identity_provider", o.identityProvider)
	}
	out := new(oidc.RawToken)
	if err := i.retryHTTP1(ctx, token, "/sts/exchange", form, out, o.userAgent, retryable); err != nil {
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

	form := url.Values{}
	for _, aud := range o.audiences {
		form.Add("aud", aud)
	}
	if o.firstScope != "" {
		// Populating for backward compatibility.
		form.Set("scope", o.firstScope)
	}
	for _, scope := range o.scope {
		form.Add("scopes", scope)
	}
	for _, cap := range o.capabilities {
		form.Add("cap", cap)
	}

	out := new(oidc.TokenPair)
	if err := i.retryHTTP1(ctx, token, "/sts/exchange_refresh_token", form, out, o.userAgent, refreshRetryable); err != nil {
		return "", "", err
	}

	return out.GetToken().GetToken(), out.GetRefreshToken().GetToken(), nil
}

// retryHTTP1 wraps doHTTP1 in the same retry loop the gRPC exchanger (impl)
// applies: up to maxRetries+1 attempts, jittered backoff between them, and the
// caller's isRetryable predicate deciding which failures are transient. This
// gives the HTTP/1-downgrade path parity with the gRPC path — a transient
// failure (identified by the gRPC code doHTTP1 recovers from the response,
// exactly as the gRPC path would see it) is retried rather than surfaced on the
// first attempt. The predicate differs by call site: Exchange passes retryable
// (codes.Unavailable only), Refresh passes refreshRetryable (also Internal and
// DeadlineExceeded), matching impl.Exchange and impl.Refresh respectively.
func (i *HTTP1DowngradeExchanger) retryHTTP1(ctx context.Context, token, path string, form url.Values, out proto.Message, userAgent string, isRetryable func(error) bool) error {
	var lastErr error
	for attempt := range maxRetries + 1 {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				// lastErr is always set here: the loop only reaches attempt>0
				// after a retryable failure recorded it. Preserve it as string
				// context so a backoff-timeout error still shows the "why" (e.g.
				// the last 503). Only ctx.Err() is %w-wrapped, so errors.Is /
				// status.Code still see the cancellation, not the transient error.
				return fmt.Errorf("%w (last STS attempt: %s)", ctx.Err(), lastErr.Error())
			case <-i.nextBackoff():
			}
		}
		err := i.doHTTP1(ctx, token, path, form, out, userAgent)
		if err == nil {
			return nil
		}
		if !isRetryable(err) {
			return err
		}
		lastErr = err
		clog.WarnContextf(ctx, "STS HTTP/1 exchange attempt %d failed (%v); retrying", attempt+1, status.Code(err))
	}
	return lastErr
}
