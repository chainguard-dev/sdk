/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package octosts_test

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync/atomic"
	"testing"
	"time"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	"chainguard.dev/sdk/octosts"
	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"chainguard.dev/sdk/sts"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// stsServer is a minimal in-process implementation of the Octo STS exchange
// endpoint. It records the requests it receives and delegates to mint for the
// response, letting tests drive the real gRPC exchange path instead of stubbing
// the token function.
type stsServer struct {
	oidc.UnimplementedSecurityTokenServiceServer

	// mint returns the response (or error) for the nth exchange (1-based).
	mint func(n int32) (*oidc.RawToken, error)

	calls   atomic.Int32
	lastReq atomic.Pointer[oidc.ExchangeRequest]
}

func (s *stsServer) Exchange(_ context.Context, req *oidc.ExchangeRequest) (*oidc.RawToken, error) {
	n := s.calls.Add(1)
	s.lastReq.Store(req)
	return s.mint(n)
}

// staticServer returns a server that always answers with resp.
func staticServer(resp *oidc.RawToken) *stsServer {
	return &stsServer{mint: func(int32) (*oidc.RawToken, error) { return resp, nil }}
}

// startExchange spins up an in-process Octo STS gRPC server backed by srv and
// returns an sts.Exchanger wired to talk to it over a bufconn. The exchanger is
// configured with the supplied options (e.g. sts.WithScope) so tests can assert
// on what the server received.
func startExchange(t *testing.T, srv *stsServer, opts ...sts.ExchangerOption) sts.Exchanger {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	scheme := delegate.RegisterListenerForTest(lis)
	t.Cleanup(func() { delegate.UnregisterTestListener(scheme) })

	s := grpc.NewServer()
	oidc.RegisterSecurityTokenServiceServer(s, srv)
	go func() {
		// Serve returns ErrServerStopped when GracefulStop wins the race with
		// this goroutine starting up; that is a clean shutdown, not a failure.
		if err := s.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			panic(fmt.Sprintf("stsServer exited: %v", err))
		}
	}()
	t.Cleanup(s.GracefulStop)

	issuer := fmt.Sprintf("%s://%s", scheme, lis.Addr().String())
	return sts.New(issuer, "octo-sts.dev", opts...)
}

func TestTokenSource_Exchange_Success(t *testing.T) {
	ctx := t.Context()
	wantToken := fmt.Sprintf("minted-%d", rand.Int64())

	xchg := startExchange(t, staticServer(&oidc.RawToken{Token: wantToken}))

	base := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "base"})
	tok, err := octosts.NewTokenSourceContext(ctx, base, xchg).Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if tok.AccessToken != wantToken {
		t.Errorf("AccessToken: got = %q, want = %q", tok.AccessToken, wantToken)
	}
	// With no expiry from the exchanger, the source defaults to a ~20 minute
	// refresh margin.
	wantExpiry := time.Now().Add(20 * time.Minute)
	if tok.Expiry.Before(wantExpiry.Add(-time.Minute)) || tok.Expiry.After(wantExpiry.Add(time.Minute)) {
		t.Errorf("Expiry: got = %v, want approximately = %v", tok.Expiry, wantExpiry)
	}
}

func TestTokenSource_Exchange_HonorsExpiry(t *testing.T) {
	ctx := t.Context()
	// A concrete expiry from the exchanger must be used verbatim rather than
	// replaced with the default refresh margin.
	wantExpiry := time.Now().Add(3 * time.Hour).Truncate(time.Second)

	xchg := startExchange(t, staticServer(&oidc.RawToken{
		Token:  "tok",
		Expiry: timestamppb.New(wantExpiry),
	}))

	base := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "base"})
	tok, err := octosts.NewTokenSourceContext(ctx, base, xchg).Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !tok.Expiry.Equal(wantExpiry) {
		t.Errorf("Expiry: got = %v, want = %v", tok.Expiry, wantExpiry)
	}
}

func TestTokenSource_Exchange_ForwardsRequest(t *testing.T) {
	ctx := t.Context()

	tests := []struct {
		name     string
		identity string
		scope    string
	}{{
		name:     "org-scoped",
		identity: fmt.Sprintf("identity-%d", rand.Int64()),
		scope:    fmt.Sprintf("org-%d", rand.Int64()),
	}, {
		name:     "repo-scoped",
		identity: fmt.Sprintf("identity-%d", rand.Int64()),
		scope:    fmt.Sprintf("org-%d/repo-%d", rand.Int64(), rand.Int64()),
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := staticServer(&oidc.RawToken{Token: "tok"})
			xchg := startExchange(t, srv,
				sts.WithScope(tt.scope),
				sts.WithIdentity(tt.identity),
			)

			base := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "base"})
			if _, err := octosts.NewTokenSourceContext(ctx, base, xchg).Token(); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			req := srv.lastReq.Load()
			if got := req.GetIdentity(); got != tt.identity {
				t.Errorf("identity: got = %q, want = %q", got, tt.identity)
			}
			if diff := cmp.Diff([]string{tt.scope}, req.GetScopes()); diff != "" {
				t.Errorf("scopes mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestTokenSource_Exchange_PropagatesError(t *testing.T) {
	ctx := t.Context()

	xchg := startExchange(t, &stsServer{mint: func(int32) (*oidc.RawToken, error) {
		return nil, errors.New("octosts unavailable")
	}})

	base := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "base"})
	tok, err := octosts.NewTokenSourceContext(ctx, base, xchg).Token()
	if err == nil {
		t.Fatal("expected error but got none")
	}
	if tok != nil {
		t.Errorf("token: got = %v, want = nil", tok)
	}
}

func TestTokenSource_BaseTokenError(t *testing.T) {
	ctx := t.Context()
	wantErr := fmt.Errorf("base-token-error-%d", rand.Int64())

	srv := staticServer(&oidc.RawToken{Token: "unused"})
	xchg := startExchange(t, srv)

	tok, err := octosts.NewTokenSourceContext(ctx, &errTokenSource{err: wantErr}, xchg).Token()
	if err == nil {
		t.Fatal("expected error but got none")
	}
	if tok != nil {
		t.Errorf("token: got = %v, want = nil", tok)
	}
	if !errors.Is(err, wantErr) {
		t.Errorf("error: got = %v, want = %v", err, wantErr)
	}
	// The base token failed, so the exchange must not have been attempted.
	if got := srv.calls.Load(); got != 0 {
		t.Errorf("exchange calls: got = %d, want = 0", got)
	}
}

func TestTokenSource_ReuseCachesToken(t *testing.T) {
	// NewTokenSource itself performs no caching; callers wrap it in
	// oauth2.ReuseTokenSource (as NewTokenSourceFromValues does). This proves
	// the returned token's expiry keeps the cache warm across calls, so a
	// second Token() does not hit the exchanger again. The server mints a
	// distinct token per exchange, so a cache hit is observable as a stable
	// AccessToken.
	ctx := t.Context()

	srv := &stsServer{mint: func(n int32) (*oidc.RawToken, error) {
		return &oidc.RawToken{Token: fmt.Sprintf("token-%d", n)}, nil
	}}
	xchg := startExchange(t, srv)

	base := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: "base"})
	ts := oauth2.ReuseTokenSource(nil, octosts.NewTokenSourceContext(ctx, base, xchg))

	first, err := ts.Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := ts.Token()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if first.AccessToken != second.AccessToken {
		t.Errorf("AccessToken: got = %q, want cached = %q", second.AccessToken, first.AccessToken)
	}
	if got := srv.calls.Load(); got != 1 {
		t.Errorf("exchange calls: got = %d, want = 1", got)
	}
}

// errTokenSource is an oauth2.TokenSource that always fails, used to exercise
// the base-token error path.
type errTokenSource struct{ err error }

func (e *errTokenSource) Token() (*oauth2.Token, error) { return nil, e.err }
