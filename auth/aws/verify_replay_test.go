/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws

// Regression tests for the DataDog report (PSEC-1498 / GHSA-wq26-3j2m-4hvv):
// VerifyToken must require the Chainguard-Audience / Chainguard-Identity binding
// headers (and host / x-amz-date) to be covered by the SigV4 signature. The
// tests are deterministic and network-free: they point the verifier at a local
// STS stub that returns 200 (modelling "STS validated the signature and ignores
// unsigned headers").

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"k8s.io/apimachinery/pkg/util/sets"
)

const (
	testAud = "https://issuer.enforce.dev"
	testID  = "deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"

	// globalSTSURL is the canonical GetCallerIdentity request target the
	// generator produces; mintSignedToken signs against it by default.
	globalSTSURL = "https://sts.amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15"

	// testAmzDate is the X-Amz-Date matching time.Unix(1_700_000_000, 0) UTC, so
	// hand-crafted tokens line up with withTimestamp(signAt).
	testAmzDate = "20231114T221320Z"
)

func signAtTime() time.Time { return time.Unix(1_700_000_000, 0).UTC() }

// serializeToken encodes a request into the base64 token form VerifyToken expects.
func serializeToken(t *testing.T, req *http.Request) string {
	t.Helper()
	var b bytes.Buffer
	if err := req.Write(&b); err != nil {
		t.Fatal(err)
	}
	return base64.URLEncoding.EncodeToString(b.Bytes())
}

// setChainguardHeaders sets the headers a legitimate token carries.
func setChainguardHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/json")
	req.Header.Set(audHeader, testAud)
	req.Header.Set(idHeader, testID)
}

// mintSignedToken builds and SigV4-signs a GetCallerIdentity request with the
// given method and target URL at signAt, including the Chainguard binding
// headers BEFORE signing (so they are covered by the signature). Used to
// isolate the secondary hardening checks (method / host / timestamp).
func mintSignedToken(t *testing.T, method, rawURL string, signAt time.Time) string {
	t.Helper()
	req, err := http.NewRequest(method, rawURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	setChainguardHeaders(req)
	creds := aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}
	if err := v4.NewSigner().SignHTTP(t.Context(), creds, req, emptyBodySHA256, "sts", "us-east-1", signAt); err != nil {
		t.Fatal(err)
	}
	return serializeToken(t, req)
}

// tokenWithSignedHeaders crafts a request whose Authorization header advertises
// exactly the given SignedHeaders list. The signature value is irrelevant:
// VerifyToken's SignedHeaders check only parses the list and the STS stub does
// not validate signatures. This models requests where host/x-amz-date were not
// actually signed, which the AWS signer would never emit on its own.
func tokenWithSignedHeaders(t *testing.T, signedHeaderList string) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, globalSTSURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	setChainguardHeaders(req)
	req.Header.Set("X-Amz-Date", testAmzDate)
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKID/20231114/us-east-1/sts/aws4_request, SignedHeaders="+signedHeaderList+", Signature=0000")
	return serializeToken(t, req)
}

// rawToken serializes a request constructed by fn, letting tests craft
// duplicate/multi-valued headers the AWS signer would never emit.
func rawToken(t *testing.T, fn func(*http.Request)) string {
	t.Helper()
	req, err := http.NewRequest(http.MethodPost, globalSTSURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	fn(req)
	return serializeToken(t, req)
}

// okSTS returns a stub STS endpoint that always returns 200 + victim claims, so
// a token reaching the forwarding step (i.e. NOT rejected by a local check) is
// observable as returned claims.
func okSTS(t *testing.T) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, `{"GetCallerIdentityResponse":{"GetCallerIdentityResult":{"UserId":"AIDAVICTIM","Account":"123456789012","Arn":"arn:aws:iam::123456789012:user/victim"}}}`)
	}))
	t.Cleanup(ts.Close)
	return ts
}

// TestGenerateTokenSignsBindingHeaders asserts the generator binds the required
// headers — i.e. they appear in the actual SignedHeaders, checked via the same
// parser the verifier uses (not a substring of the whole request).
func TestGenerateTokenSignsBindingHeaders(t *testing.T) {
	timeNow = signAtTime
	t.Cleanup(func() { timeNow = time.Now })

	tok, err := GenerateToken(t.Context(),
		aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, testAud, "some-identity")
	if err != nil {
		t.Fatal(err)
	}
	raw, err := base64.URLEncoding.DecodeString(tok)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(raw)))
	if err != nil {
		t.Fatal(err)
	}
	signed := signedHeaderNames(req.Header.Get("Authorization"))
	for _, h := range []string{"chainguard-audience", "chainguard-identity", "host", "x-amz-date"} {
		if !slices.Contains(signed, h) {
			t.Fatalf("legitimate token must sign %q; got SignedHeaders=%v", h, signed)
		}
	}
}

// TestVerifyRejectsUnsignedBindingHeaders is the headline regression test for
// the reported bypass: a GetCallerIdentity request signed WITHOUT the Chainguard
// binding headers, with those headers injected after signing, must be rejected
// because the headers are not covered by the signature. Before the fix this
// returned the victim's claims with a nil error.
func TestVerifyRejectsUnsignedBindingHeaders(t *testing.T) {
	ctx := t.Context()
	signAt := signAtTime()

	// A "stolen" signed GetCallerIdentity request minted WITHOUT the Chainguard
	// binding headers -- the shape of a request minted for Vault AWS auth,
	// aws-iam-authenticator, another audience, etc.
	creds := aws.Credentials{AccessKeyID: "AKIAVICTIM", SecretAccessKey: "victimsecret"}
	req, err := http.NewRequest("POST", "https://sts.amazonaws.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.URL.Path = "/"
	req.URL.RawQuery = "Action=GetCallerIdentity&Version=2011-06-15"
	req.Header.Add("Accept", "application/json")
	if err := v4.NewSigner().SignHTTP(ctx, creds, req, emptyBodySHA256, "sts", "us-east-1", signAt); err != nil {
		t.Fatal(err)
	}
	if authz := strings.ToLower(req.Header.Get("Authorization")); strings.Contains(authz, "chainguard") {
		t.Fatalf("precondition failed: chainguard headers unexpectedly in SignedHeaders: %s", authz)
	}

	// The attacker injects the Chainguard headers AFTER signing.
	req.Header.Set(audHeader, testAud)
	req.Header.Set(idHeader, testID)
	forged := serializeToken(t, req)

	claims, err := VerifyToken(ctx, forged,
		WithAudience(sets.New(testAud)), WithIdentity(testID),
		withSTSURL(okSTS(t).URL), withTimestamp(signAt))
	if claims != nil {
		t.Fatalf("SECURITY: VerifyToken returned claims for a token whose binding headers were not signed: %+v", claims)
	}
	if !errors.Is(err, ErrUnsignedBindingHeaders) {
		t.Fatalf("want ErrUnsignedBindingHeaders, got: %v", err)
	}
}

// TestVerifyAcceptsRealGeneratedToken guards against over-rejection: the
// production generator's output must pass every new check.
func TestVerifyAcceptsRealGeneratedToken(t *testing.T) {
	signAt := signAtTime()
	timeNow = func() time.Time { return signAt }
	t.Cleanup(func() { timeNow = time.Now })

	tok, err := GenerateToken(t.Context(),
		aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET"}, testAud, testID)
	if err != nil {
		t.Fatal(err)
	}
	claims, err := VerifyToken(t.Context(), tok,
		WithAudience(sets.New(testAud)), WithIdentity(testID),
		withSTSURL(okSTS(t).URL), withTimestamp(signAt))
	if err != nil || claims == nil {
		t.Fatalf("want a real generated token to verify, got claims=%+v err=%v", claims, err)
	}
}

// TestVerifyHardening exercises the local pre-STS checks: every negative case
// each guard rejects, and the positive cases that must still be accepted.
func TestVerifyHardening(t *testing.T) {
	signAt := signAtTime()
	const skew = defaultAllowedClockSkew

	dupHeader := func(set func(*http.Request)) func(*testing.T) string {
		return func(t *testing.T) string {
			return rawToken(t, func(r *http.Request) {
				setChainguardHeaders(r)
				r.Header.Set("X-Amz-Date", testAmzDate)
				r.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential=AKID/20231114/us-east-1/sts/aws4_request, SignedHeaders=accept;chainguard-audience;chainguard-identity;host;x-amz-date, Signature=0000")
				set(r) // introduce the duplicate/multi-valued header under test
			})
		}
	}

	tests := []struct {
		name    string
		token   func(t *testing.T) string
		extra   []VerifyOption
		now     time.Time
		wantErr error // nil = accept
	}{
		{"non-POST method", func(t *testing.T) string { return mintSignedToken(t, http.MethodGet, globalSTSURL, signAt) }, nil, signAt, ErrInvalidRequest},
		{"unexpected host", func(t *testing.T) string {
			return mintSignedToken(t, http.MethodPost, "https://evil.example.com/?Action=GetCallerIdentity&Version=2011-06-15", signAt)
		}, nil, signAt, ErrInvalidRequest},
		{"configured host rejects global-host token", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) },
			[]VerifyOption{WithExpectedSTSHost("sts.us-gov-west-1.amazonaws.com")}, signAt, ErrInvalidRequest},
		{"configured non-global host accepted", func(t *testing.T) string {
			return mintSignedToken(t, http.MethodPost, "https://sts.us-gov-west-1.amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15", signAt)
		}, []VerifyOption{WithExpectedSTSHost("sts.us-gov-west-1.amazonaws.com")}, signAt, nil},
		{"unsigned x-amz-date", func(t *testing.T) string {
			return tokenWithSignedHeaders(t, "accept;chainguard-audience;chainguard-identity;host")
		}, nil, signAt, ErrUnsignedBindingHeaders},
		{"unsigned host", func(t *testing.T) string {
			return tokenWithSignedHeaders(t, "accept;chainguard-audience;chainguard-identity;x-amz-date")
		}, nil, signAt, ErrUnsignedBindingHeaders},
		{"beyond default clock skew", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) }, nil, signAt.Add(-8 * time.Minute), ErrTokenFutureDated},
		{"within configured clock skew", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) },
			[]VerifyOption{WithAllowedClockSkew(10 * time.Minute)}, signAt.Add(-8 * time.Minute), nil},
		{"exactly at clock skew", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) }, nil, signAt.Add(-skew), nil},
		{"just beyond clock skew", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) }, nil, signAt.Add(-skew - time.Second), ErrTokenFutureDated},
		{"exactly at expiry", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) }, nil, signAt.Add(15 * time.Minute), nil},
		{"just past expiry", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) }, nil, signAt.Add(15*time.Minute + time.Second), ErrTokenExpired},
		{"negative clock skew is rejected as misconfiguration", func(t *testing.T) string { return mintSignedToken(t, http.MethodPost, globalSTSURL, signAt) },
			[]VerifyOption{WithAllowedClockSkew(-time.Minute)}, signAt, ErrInvalidVerificationConfiguration},
		{"duplicate Authorization", func(t *testing.T) string {
			return rawToken(t, func(r *http.Request) {
				setChainguardHeaders(r)
				r.Header.Set("X-Amz-Date", testAmzDate)
				r.Header.Add("Authorization", "AWS4-HMAC-SHA256 Credential=AKID/20231114/us-east-1/sts/aws4_request, SignedHeaders=accept;chainguard-audience;chainguard-identity;host;x-amz-date, Signature=0000")
				r.Header.Add("Authorization", "AWS4-HMAC-SHA256 Credential=AKID/20231114/us-east-1/sts/aws4_request, SignedHeaders=host;x-amz-date, Signature=1111")
			})
		}, nil, signAt, ErrInvalidRequest},
		{"duplicate X-Amz-Date", dupHeader(func(r *http.Request) { r.Header.Add("X-Amz-Date", testAmzDate) }), nil, signAt, ErrInvalidRequest},
		{"multi-valued Chainguard-Audience", dupHeader(func(r *http.Request) { r.Header.Add(audHeader, "https://attacker.example.com") }), nil, signAt, ErrInvalidRequest},
		{"duplicate Chainguard-Identity", dupHeader(func(r *http.Request) { r.Header.Add(idHeader, "another-identity") }), nil, signAt, ErrInvalidRequest},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Per-case options first, then the base overrides last (last-wins) so
			// the STS stub and injected clock still apply even when a case sets
			// WithExpectedSTSHost, which now also derives stsURL.
			opts := append(append([]VerifyOption{}, tc.extra...),
				WithAudience(sets.New(testAud)), WithIdentity(testID),
				withSTSURL(okSTS(t).URL), withTimestamp(tc.now))
			claims, err := VerifyToken(t.Context(), tc.token(t), opts...)
			if tc.wantErr == nil {
				if err != nil || claims == nil {
					t.Fatalf("want acceptance, got claims=%+v err=%v", claims, err)
				}
				return
			}
			if claims != nil || !errors.Is(err, tc.wantErr) {
				t.Fatalf("want %v, got claims=%+v err=%v", tc.wantErr, claims, err)
			}
		})
	}
}

// TestVerifyRejectedTokenNeverReachesSTS proves a locally-rejected token is not
// forwarded to STS (the expensive round-trip is gated behind the cheap checks).
func TestVerifyRejectedTokenNeverReachesSTS(t *testing.T) {
	signAt := signAtTime()
	ts := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Errorf("a locally-rejected token must not be forwarded to AWS STS")
	}))
	t.Cleanup(ts.Close)
	token := mintSignedToken(t, http.MethodGet, globalSTSURL, signAt) // non-POST → rejected before forwarding
	claims, err := VerifyToken(t.Context(), token,
		WithAudience(sets.New(testAud)), WithIdentity(testID),
		withSTSURL(ts.URL), withTimestamp(signAt))
	if claims != nil || !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("want ErrInvalidRequest for a non-POST token, got claims=%+v err=%v", claims, err)
	}
}

// TestWithExpectedSTSHostDerivesURL asserts the option derives the connect URL
// from the host (not just the expected-host check), so the two can't drift and
// a non-global endpoint is actually forwarded to — the config a stub-based
// hardening test can't observe.
func TestWithExpectedSTSHostDerivesURL(t *testing.T) {
	conf, err := newConfigFromOptions(
		WithAudience(sets.New(testAud)), WithIdentity(testID),
		WithExpectedSTSHost("sts.us-gov-west-1.amazonaws.com"))
	if err != nil {
		t.Fatal(err)
	}
	if conf.expectedSTSHost != "sts.us-gov-west-1.amazonaws.com" {
		t.Fatalf("expectedSTSHost = %q", conf.expectedSTSHost)
	}
	if want := "https://sts.us-gov-west-1.amazonaws.com" + stsQuery; conf.stsURL != want {
		t.Fatalf("stsURL = %q, want %q", conf.stsURL, want)
	}
}

func TestSignedHeaderNames(t *testing.T) {
	tests := []struct {
		name  string
		authz string
		want  []string
	}{
		{"canonical", "AWS4-HMAC-SHA256 Credential=AKID/20231114/us-east-1/sts/aws4_request, SignedHeaders=accept;host;x-amz-date, Signature=abc", []string{"accept", "host", "x-amz-date"}},
		{"uppercase names lowered", "AWS4-HMAC-SHA256 Credential=x, SignedHeaders=Host;X-Amz-Date, Signature=abc", []string{"host", "x-amz-date"}},
		{"no SignedHeaders component", "AWS4-HMAC-SHA256 Credential=x, Signature=abc", nil},
		{"empty input", "", nil},
		{"empty list", "AWS4-HMAC-SHA256 Credential=x, SignedHeaders=, Signature=abc", []string{""}},
		{"whitespace not trimmed (fails closed)", "x, SignedHeaders=host; x-amz-date, Signature=y", []string{"host", " x-amz-date"}},
		{"multiple SignedHeaders fails closed", "x, SignedHeaders=accept;chainguard-audience;chainguard-identity;host;x-amz-date, SignedHeaders=host;x-amz-date, Signature=y", nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := signedHeaderNames(tc.authz)
			if !slices.Equal(got, tc.want) {
				t.Fatalf("signedHeaderNames(%q) = %#v, want %#v", tc.authz, got, tc.want)
			}
		})
	}
}
