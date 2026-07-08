/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/chainguard-dev/clog"
)

var (
	ErrTokenRejected                    = errors.New("token rejected by AWS STS endpoint")
	ErrTokenExpired                     = errors.New("token expired")
	ErrInvalidAudience                  = errors.New("audience header in token does not match expected audience")
	ErrInvalidIdentity                  = errors.New("identity header in token does not match expected identity")
	ErrInvalidEncoding                  = errors.New("invalid token encoding")
	ErrInvalidVerificationConfiguration = errors.New("verification was incorrectly configured")

	// Added with the SignedHeaders binding fix. Callers that branch on
	// VerifyToken errors with errors.Is should handle these explicitly rather
	// than letting them fall into a generic "unknown error" bucket.
	ErrUnsignedBindingHeaders = errors.New("a required header is not covered by the request signature")
	ErrInvalidRequest         = errors.New("token is not a well-formed GetCallerIdentity request")
	ErrTokenFutureDated       = errors.New("token timestamp is in the future")
)

// defaultAllowedClockSkew bounds how far a token's X-Amz-Date may be ahead of
// the verifier's clock before it is rejected as future-dated. It is the default
// for the configurable per-verification skew (see WithAllowedClockSkew) and
// matches AWS SigV4's own tolerance.
const defaultAllowedClockSkew = 5 * time.Minute

// Header keys (canonical form) we read from the decoded request, and the layout
// of the X-Amz-Date timestamp. (audHeader/idHeader are defined in generate.go.)
const (
	authorizationHeader = "Authorization"
	hostHeader          = "Host"
	amzDateHeader       = "X-Amz-Date"
	amzDateFormat       = "20060102T150405Z"
)

// requiredSignedHeaders are the lower-cased header names that must appear in the
// request's SigV4 SignedHeaders before VerifyToken trusts the request: the
// Chainguard binding headers plus the SigV4-intrinsic host and date. Derived
// from generate.go's bindingHeaders (computed once) so the generator and
// verifier cannot drift, and so we don't re-lowercase constants on every call.
var requiredSignedHeaders = func() []string {
	hs := append([]string{hostHeader, amzDateHeader}, bindingHeaders...)
	for i := range hs {
		hs[i] = strings.ToLower(hs[i])
	}
	return hs
}()

// noDuplicateHeaders must each appear at most once in the decoded request: a
// duplicate would let our checks read one value while STS validates or uses
// another, since the request is re-serialized and forwarded verbatim.
var noDuplicateHeaders = append([]string{authorizationHeader, amzDateHeader}, bindingHeaders...)

type VerifiedClaims struct {
	UserID  string `json:"UserId"`
	Arn     string `json:"Arn"`
	Account string `json:"Account"`
}

// VerifyToken verifies a Chainguard token backed by an AWS STS
// GetCallerIdentity request and returns the caller's verified AWS identity.
//
// By default only tokens addressed to the global endpoint "sts.amazonaws.com"
// are accepted; tokens signed against a regional, FIPS, or partition (GovCloud,
// China) endpoint must opt in with WithExpectedSTSHost, otherwise VerifyToken
// rejects them with ErrInvalidRequest before contacting STS.
func VerifyToken(ctx context.Context, token string, opts ...VerifyOption) (*VerifiedClaims, error) {
	conf, err := newConfigFromOptions(opts...)
	if err != nil {
		clog.ErrorContextf(ctx, "invalid verification configuration: %v", err)
		return nil, ErrInvalidVerificationConfiguration
	}

	req, err := decodeRequest(token, conf.stsURL)
	if err != nil {
		if errors.Is(err, ErrInvalidVerificationConfiguration) {
			clog.FromContext(ctx).With("sts_url", conf.stsURL).Errorf("invalid verification configuration. invalid sts url: %v", err)
		}
		return nil, err
	}

	// The checks below are defense-in-depth PRE-FILTERS; they do NOT verify the
	// SigV4 signature. req.Host is whatever the serialized request advertised
	// (request-line authority or Host header) and the SignedHeaders list is
	// parsed from the human-readable Authorization header — neither is locally
	// authenticated. The cryptographic guarantee that the audience/identity/host
	// were actually signed comes only from AWS STS recomputing the canonical
	// request and rejecting a mismatch. Requiring these names in SignedHeaders
	// rejects a request signed for a different relying party (which never signed
	// them) before we incur the STS round-trip.

	// A duplicated required header would let our checks read one value while STS
	// validates/uses another (all values are forwarded verbatim).
	for _, h := range noDuplicateHeaders {
		if len(req.Header.Values(h)) > 1 {
			clog.WarnContext(ctx, "verification failed because of a duplicated request header", "header", h)
			return nil, ErrInvalidRequest
		}
	}
	// Constrain to the exact GetCallerIdentity shape our generator produces, so
	// signed requests of other shapes (presigned GET URLs, a different host) are
	// not repurposed here.
	if req.Method != http.MethodPost {
		clog.WarnContext(ctx, "verification failed because of unexpected request method", "method", req.Method)
		return nil, ErrInvalidRequest
	}
	if req.Host != conf.expectedSTSHost {
		clog.WarnContext(ctx, "verification failed because of unexpected request host", "host", req.Host, "expected", conf.expectedSTSHost)
		return nil, ErrInvalidRequest
	}
	// Every header whose value we trust below must be covered by the signature.
	signed := signedHeaderNames(req.Header.Get(authorizationHeader))
	for _, h := range requiredSignedHeaders {
		if !slices.Contains(signed, h) {
			clog.WarnContext(ctx, "verification failed because a required header is not covered by the request signature", "header", h, "signed_headers", signed)
			return nil, ErrUnsignedBindingHeaders
		}
	}

	rawAud := req.Header.Get(audHeader)
	if aud := sigv4CanonicalHeaderValue(rawAud); !conf.allowedAudiences.Has(aud) {
		clog.WarnContext(ctx, "verification failed with audience mismatch", "received", aud, "received_raw", rawAud)
		return nil, ErrInvalidAudience
	}
	rawID := req.Header.Get(idHeader)
	if id := sigv4CanonicalHeaderValue(rawID); id != conf.identity {
		clog.WarnContext(ctx, "verification failed with identity mismatch", "wanted", conf.identity, "received", id, "received_raw", rawID)
		return nil, ErrInvalidIdentity
	}

	timestamp, err := time.Parse(amzDateFormat, req.Header.Get(amzDateHeader))
	if err != nil {
		clog.WarnContextf(ctx, "verification failed because of a poorly formatted x-amz-date header format: %v", err)
		return nil, ErrInvalidEncoding
	}
	expiry, now := timestamp.Add(15*time.Minute), conf.time()
	if timestamp.After(now.Add(conf.allowedClockSkew)) {
		// A future-dated X-Amz-Date would otherwise trivially satisfy the
		// expiry check below; reject it (real AWS STS independently rejects
		// large clock skew).
		clog.WarnContext(ctx, "verification failed because the token timestamp is in the future", "timestamp", timestamp, "now", now)
		return nil, ErrTokenFutureDated
	}
	if expiry.Before(now) {
		// According to AWS docs
		// > The signed portions (using AWS Signatures) of requests are valid within 15 minutes of the timestamp in the request.
		// If the signature timestamp is already older than 15 minutes the token is expired and we reject it.
		// c.f https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html
		clog.WarnContext(ctx, "verification failed because of expired token")
		return nil, ErrTokenExpired
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx)) //nolint:gosec // G704: URL from trusted AWS STS config
	if err != nil {
		clog.ErrorContextf(ctx, "verification failed because of failure to make AWS STS request: %v", err)
		return nil, fmt.Errorf("failed to reach AWS STS endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body from AWS STS endpoint: %w", err)
		}
		clog.ErrorContext(ctx, "verification failed because it was rejected by AWS STS endpoint", "response_code", resp.StatusCode, "response", string(body))
		return nil, ErrTokenRejected
	}

	var response struct {
		GetCallerIdentityResponse struct {
			GetCallerIdentityResult VerifiedClaims
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		clog.ErrorContextf(ctx, "verification failed because json parsing err in response: %v", err)
		return nil, fmt.Errorf("failed to parse json from AWS STS response %w", err)
	}

	return &response.GetCallerIdentityResponse.GetCallerIdentityResult, nil
}

// decodeRequest base64-decodes the token, parses the serialized HTTP request,
// and repoints it at the configured STS endpoint so it can be forwarded.
func decodeRequest(token, stsURL string) (*http.Request, error) {
	decoded, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, ErrInvalidEncoding
	}
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(decoded)))
	if err != nil {
		return nil, ErrInvalidEncoding
	}
	// RequestURI is set by parsing a serialized request but must be cleared
	// before an http.Client will re-send it.
	req.RequestURI = ""
	u, err := url.Parse(stsURL)
	if err != nil {
		return nil, ErrInvalidVerificationConfiguration
	}
	req.URL = u
	return req, nil
}

// sigv4CanonicalHeaderValue canonicalizes a header value to match the form
// that the AWS SigV4 signer itself produces when constructing the canonical
// request, closing the representation-hygiene gap noted in PSEC-1498.
//
// It mirrors the composition at github.com/aws/aws-sdk-go-v2/aws/signer/v4/v4.go:470:
//
//	cleanedValue := strings.TrimSpace(v4Internal.StripExcessSpaces(v))
//
// Step by step:
//
//  1. Collapse runs of internal ASCII space (0x20) characters to a single
//     space, mirroring github.com/aws/aws-sdk-go-v2/aws/signer/internal/v4.StripExcessSpaces.
//     Only 0x20 is collapsed; internal tabs and Unicode space characters are
//     preserved (StripExcessSpaces does not touch them either).
//
//  2. Trim leading/trailing whitespace with strings.TrimSpace, exactly as the
//     signer does.  Because STS must apply the same pipeline to verify the
//     signature, a raw value "admin\t" is signed as "admin", validated by STS
//     as "admin", and attributed here as "admin" — all three agree.
//
// This is applied symmetrically to both the token header values and the
// configured audiences/identity, so both sides are compared in the same form.
func sigv4CanonicalHeaderValue(v string) string {
	// Step 1: collapse runs of internal ASCII 0x20 spaces (strip excess spaces).
	// We only collapse 0x20; tab and other whitespace characters are not touched
	// by StripExcessSpaces and are therefore preserved internally.
	var b strings.Builder
	b.Grow(len(v))
	inSpace := false
	start := true
	for _, c := range []byte(v) {
		if c == ' ' {
			if !start {
				inSpace = true
			}
		} else {
			if inSpace {
				b.WriteByte(' ')
			}
			b.WriteByte(c)
			inSpace = false
			start = false
		}
	}
	// Step 2: trim leading/trailing whitespace with strings.TrimSpace, exactly
	// as v4.go:470 does; the loop above already consumed leading 0x20 via start.
	return strings.TrimSpace(b.String())
}

// signedHeaderNames extracts the lower-cased header names from the SignedHeaders
// component of a SigV4 Authorization header, e.g. given
//
//	AWS4-HMAC-SHA256 Credential=..., SignedHeaders=accept;host;x-amz-date, Signature=...
//
// it returns ["accept", "host", "x-amz-date"]. It returns nil if the header is
// absent, has no SignedHeaders component, or — fail closed — carries more than
// one SignedHeaders component (malformed; the AWS signer never emits that).
//
// Names are lower-cased so callers can compare against lower-case header names
// (AWS canonicalizes SignedHeaders to lower case, but we normalize defensively
// rather than rely on it). Elements are deliberately NOT whitespace-trimmed:
// the canonical form has no spaces, so a spaced element fails the membership
// check (fail closed), and such a request would not validate at STS anyway.
func signedHeaderNames(authorization string) []string {
	var names []string
	for part := range strings.SplitSeq(authorization, ",") {
		rest, ok := strings.CutPrefix(strings.TrimSpace(part), "SignedHeaders=")
		if !ok {
			continue
		}
		if names != nil {
			return nil // more than one SignedHeaders= component is malformed
		}
		names = strings.Split(strings.ToLower(rest), ";")
	}
	return names
}
