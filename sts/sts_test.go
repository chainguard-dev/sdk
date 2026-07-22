/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"chainguard.dev/sdk/proto/platform/oidc/v1/test"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestRefresh(t *testing.T) {
	tests := map[string]struct {
		issuer           string
		audience         string
		newOpts          []ExchangerOption
		exchangeOpts     []ExchangerOption
		wantToken        string
		wantRefreshToken string
		wantErr          bool
		clientMock       test.MockOIDCClient
	}{
		"zero options": {
			issuer:   "bar",
			audience: "baz",
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnGetAccessToken: []test.STSOnGetAccessToken{{
						Given: &oidc.ExchangeRefreshTokenRequest{
							Aud: []string{"baz"},
						},
						Exchanged: &oidc.TokenPair{
							Token:        &oidc.RawToken{Token: "token!"},
							RefreshToken: &oidc.RawToken{Token: "refresh token!"},
						},
					}},
				},
			},
			wantToken:        "token!",
			wantRefreshToken: "refresh token!",
		},
		"basic error plumbing": {
			issuer:   "bar",
			audience: "baz",
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnGetAccessToken: []test.STSOnGetAccessToken{{
						Given: &oidc.ExchangeRefreshTokenRequest{
							Aud: []string{"baz"},
						},
						Error: errors.New("unexpected EOF"),
					}},
				},
			},
			wantErr: true,
		},
		"capabilities and scope": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCapabilities("groups.list"),
				WithScope("derp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnGetAccessToken: []test.STSOnGetAccessToken{{
						Given: &oidc.ExchangeRefreshTokenRequest{
							Aud:    []string{"baz"},
							Cap:    []string{"groups.list"},
							Scope:  "derp",
							Scopes: []string{"derp"},
						},
						Exchanged: &oidc.TokenPair{
							Token:        &oidc.RawToken{Token: "token!"},
							RefreshToken: &oidc.RawToken{Token: "refresh token"},
						},
					}},
				},
			},
			wantToken:        "token!",
			wantRefreshToken: "refresh token",
		},
		"multiple scopes": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithScope("derp", "ferp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnGetAccessToken: []test.STSOnGetAccessToken{{
						Given: &oidc.ExchangeRefreshTokenRequest{
							Aud:    []string{"baz"},
							Scope:  "derp",
							Scopes: []string{"derp", "ferp"},
						},
						Exchanged: &oidc.TokenPair{
							Token:        &oidc.RawToken{Token: "token!"},
							RefreshToken: &oidc.RawToken{Token: "refresh token"},
						},
					}},
				},
			},
			wantToken:        "token!",
			wantRefreshToken: "refresh token",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oidcNewClients = func(_ context.Context, _ string, _ string, _ ...oidc.ClientOption) (oidc.Clients, error) {
				return test.clientMock, nil
			}

			exch := New(test.issuer, test.audience, test.newOpts...)
			token, refreshToken, gotErr := exch.Refresh(t.Context(), "foo", test.exchangeOpts...)
			if (gotErr != nil) != test.wantErr {
				t.Fatalf("Refresh() error = %v, wantErr = %v", gotErr, test.wantErr)
			}
			if diff := cmp.Diff(test.wantToken, token); diff != "" {
				t.Errorf("token mismatch (-want, +got):\n%s", diff)
			}
			if diff := cmp.Diff(test.wantRefreshToken, refreshToken); diff != "" {
				t.Errorf("refreshToken mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}
func TestImplExchange(t *testing.T) {
	tests := map[string]struct {
		issuer       string
		audience     string
		newOpts      []ExchangerOption
		exchangeOpts []ExchangerOption
		wantToken    TokenPair
		wantErr      bool
		clientMock   test.MockOIDCClient
	}{
		"zero options": {
			issuer:   "bar",
			audience: "baz",
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud: []string{"baz"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: "refresh token!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!", RefreshToken: "refresh token!"},
		},
		"basic error plumbing": {
			issuer:   "bar",
			audience: "baz",
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud: []string{"baz"},
						},
						Error: errors.New("unexpected EOF"),
					}},
				},
			},
			wantErr: true,
		},
		"capabilities and scope on create": {
			issuer:   "bar",
			audience: "baz",
			newOpts: []ExchangerOption{
				WithCapabilities("groups.list"),
				WithScope("derp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:    []string{"baz"},
							Cap:    []string{"groups.list"},
							Scope:  "derp",
							Scopes: []string{"derp"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: ""},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!"},
		},
		"capabilities and scope on exchange": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCapabilities("groups.list"),
				WithScope("derp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:    []string{"baz"},
							Cap:    []string{"groups.list"},
							Scope:  "derp",
							Scopes: []string{"derp"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: "refreshToken!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!", RefreshToken: "refreshToken!"},
		},
		"multiple scopes": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithScope("derp", "ferp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:    []string{"baz"},
							Scope:  "derp",
							Scopes: []string{"derp", "ferp"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: "refreshToken!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!", RefreshToken: "refreshToken!"},
		},
		"identity": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIdentity("my-identity"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:      []string{"baz"},
							Identity: "my-identity",
						},
						Exchanged: &oidc.RawToken{Token: "token foo", RefreshToken: "refresh token"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token foo", RefreshToken: "refresh token"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oidcNewClients = func(_ context.Context, _ string, _ string, _ ...oidc.ClientOption) (oidc.Clients, error) {
				return test.clientMock, nil
			}

			exch := New(test.issuer, test.audience, test.newOpts...)
			gotTok, gotErr := exch.Exchange(t.Context(), "foo", test.exchangeOpts...)
			if (gotErr != nil) != test.wantErr {
				t.Fatalf("error = %v, wantErr = %v", gotErr, test.wantErr)
			}
			if diff := cmp.Diff(test.wantToken, gotTok); diff != "" {
				t.Errorf("token mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestExchange(t *testing.T) {
	tests := map[string]struct {
		issuer       string
		audience     string
		exchangeOpts []ExchangerOption
		wantToken    TokenPair
		wantErr      bool
		clientMock   test.MockOIDCClient
	}{
		"zero options": {
			issuer:   "bar",
			audience: "baz",
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud: []string{"baz"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: "refresh token!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!", RefreshToken: "refresh token!"},
		},
		"basic error plumbing": {
			issuer:   "bar",
			audience: "baz",
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud: []string{"baz"},
						},
						Error: errors.New("unexpected EOF"),
					}},
				},
			},
			wantErr: true,
		},
		"capabilities and scope on create": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCapabilities("groups.list"),
				WithScope("derp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:    []string{"baz"},
							Cap:    []string{"groups.list"},
							Scope:  "derp",
							Scopes: []string{"derp"},
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!"},
		},
		"capabilities and scope on exchange": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCapabilities("groups.list"),
				WithScope("derp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:    []string{"baz"},
							Cap:    []string{"groups.list"},
							Scope:  "derp",
							Scopes: []string{"derp"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: "refreshToken!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!", RefreshToken: "refreshToken!"},
		},
		"multiple scopes": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithScope("derp", "ferp"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:    []string{"baz"},
							Scope:  "derp",
							Scopes: []string{"derp", "ferp"},
						},
						Exchanged: &oidc.RawToken{Token: "token!", RefreshToken: "refreshToken!"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token!", RefreshToken: "refreshToken!"},
		},
		"identity": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIdentity("my-identity"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:      []string{"baz"},
							Identity: "my-identity",
						},
						Exchanged: &oidc.RawToken{Token: "token foo"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token foo"},
		},
		"identityProvider": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIdentityProvider("my-identity-provider"),
			},
			clientMock: test.MockOIDCClient{
				STSClient: test.MockSTSClient{
					OnExchange: []test.STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:              []string{"baz"},
							IdentityProvider: "my-identity-provider",
						},
						Exchanged: &oidc.RawToken{Token: "token foo"},
					}},
				},
			},
			wantToken: TokenPair{AccessToken: "token foo"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oidcNewClients = func(_ context.Context, _ string, _ string, _ ...oidc.ClientOption) (oidc.Clients, error) {
				return test.clientMock, nil
			}

			gotTok, gotErr := ExchangePair(t.Context(), test.issuer, test.audience, "foo", test.exchangeOpts...)
			if (gotErr != nil) != test.wantErr {
				t.Fatalf("error = %v, wantErr = %v", gotErr, test.wantErr)
			}
			if diff := cmp.Diff(test.wantToken, gotTok); diff != "" {
				t.Errorf("token mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestExchangeRetry(t *testing.T) {
	for _, tt := range []struct {
		name string
		// errors to return before succeeding; len determines how many attempts fail.
		errors       []error
		wantAttempts int
		wantErr      bool
	}{{
		name:         "succeeds on first attempt",
		wantAttempts: 1,
	}, {
		name:         "retries on Unavailable then succeeds",
		errors:       []error{status.Error(codes.Unavailable, "service unavailable")},
		wantAttempts: 2,
	}, {
		name:         "does not retry on Internal",
		errors:       []error{status.Error(codes.Internal, "internal error")},
		wantAttempts: 1,
		wantErr:      true,
	}, {
		name: "retries twice then succeeds",
		errors: []error{
			status.Error(codes.Unavailable, "try 1"),
			status.Error(codes.Unavailable, "try 2"),
		},
		wantAttempts: 3,
	}, {
		name: "exhausts retries and returns last error",
		errors: []error{
			status.Error(codes.Unavailable, "try 1"),
			status.Error(codes.Unavailable, "try 2"),
			status.Error(codes.Unavailable, "try 3"),
		},
		wantAttempts: 3,
		wantErr:      true,
	}, {
		name:         "does not retry on PermissionDenied",
		errors:       []error{status.Error(codes.PermissionDenied, "forbidden")},
		wantAttempts: 1,
		wantErr:      true,
	}, {
		name:         "does not retry on InvalidArgument",
		errors:       []error{status.Error(codes.InvalidArgument, "bad request")},
		wantAttempts: 1,
		wantErr:      true,
	}, {
		name:         "does not retry on non-gRPC error",
		errors:       []error{errors.New("unexpected EOF")},
		wantAttempts: 1,
		wantErr:      true,
	}} {
		t.Run(tt.name, func(t *testing.T) {
			var attempts atomic.Int32
			oidcNewClients = func(_ context.Context, _ string, _ string, _ ...oidc.ClientOption) (oidc.Clients, error) {
				idx := int(attempts.Add(1)) - 1
				if idx < len(tt.errors) {
					return test.MockOIDCClient{
						STSClient: test.MockSTSClient{
							OnExchange: []test.STSOnExchange{{
								Given: &oidc.ExchangeRequest{Aud: []string{"aud"}},
								Error: tt.errors[idx],
							}},
						},
					}, nil
				}
				return test.MockOIDCClient{
					STSClient: test.MockSTSClient{
						OnExchange: []test.STSOnExchange{{
							Given:     &oidc.ExchangeRequest{Aud: []string{"aud"}},
							Exchanged: &oidc.RawToken{Token: "ok"},
						}},
					},
				}, nil
			}

			exch := New("issuer", "aud")
			got, err := exch.Exchange(t.Context(), "tok")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Exchange() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.AccessToken != "ok" {
				t.Errorf("Exchange() token = %q, want %q", got.AccessToken, "ok")
			}
			if gotAttempts := int(attempts.Load()); gotAttempts != tt.wantAttempts {
				t.Errorf("Exchange() attempts = %d, want %d", gotAttempts, tt.wantAttempts)
			}
		})
	}
}

func TestRefreshRetry(t *testing.T) {
	for _, tt := range []struct {
		name         string
		errors       []error
		wantAttempts int
		wantErr      bool
	}{{
		name:         "succeeds on first attempt",
		wantAttempts: 1,
	}, {
		name:         "retries on Unavailable then succeeds",
		errors:       []error{status.Error(codes.Unavailable, "service unavailable")},
		wantAttempts: 2,
	}, {
		name: "exhausts retries",
		errors: []error{
			status.Error(codes.Unavailable, "try 1"),
			status.Error(codes.Unavailable, "try 2"),
			status.Error(codes.Unavailable, "try 3"),
		},
		wantAttempts: 3,
		wantErr:      true,
	}, {
		name:         "does not retry on PermissionDenied",
		errors:       []error{status.Error(codes.PermissionDenied, "forbidden")},
		wantAttempts: 1,
		wantErr:      true,
	}, {
		name:         "does not retry on InvalidArgument",
		errors:       []error{status.Error(codes.InvalidArgument, "bad request")},
		wantAttempts: 1,
		wantErr:      true,
	}, {
		// Unlike Exchange, the refresh path retries Internal: under
		// STS/datastore saturation a failed identity lookup surfaces as
		// codes.Internal, which is transient here. See CUS-839.
		name:         "retries on Internal then succeeds",
		errors:       []error{status.Error(codes.Internal, "internal error")},
		wantAttempts: 2,
	}, {
		name: "exhausts retries on Internal",
		errors: []error{
			status.Error(codes.Internal, "try 1"),
			status.Error(codes.Internal, "try 2"),
			status.Error(codes.Internal, "try 3"),
		},
		wantAttempts: 3,
		wantErr:      true,
	}, {
		// A server-side query/connection deadline surfaces as
		// DeadlineExceeded; transient on the refresh path. See CUS-839.
		name:         "retries on DeadlineExceeded then succeeds",
		errors:       []error{status.Error(codes.DeadlineExceeded, "deadline exceeded")},
		wantAttempts: 2,
	}, {
		name:         "does not retry on Canceled",
		errors:       []error{status.Error(codes.Canceled, "canceled")},
		wantAttempts: 1,
		wantErr:      true,
	}, {
		name:         "does not retry on non-gRPC error",
		errors:       []error{errors.New("unexpected EOF")},
		wantAttempts: 1,
		wantErr:      true,
	}} {
		t.Run(tt.name, func(t *testing.T) {
			var attempts atomic.Int32
			oidcNewClients = func(_ context.Context, _ string, _ string, _ ...oidc.ClientOption) (oidc.Clients, error) {
				idx := int(attempts.Add(1)) - 1
				if idx < len(tt.errors) {
					return test.MockOIDCClient{
						STSClient: test.MockSTSClient{
							OnGetAccessToken: []test.STSOnGetAccessToken{{
								Given: &oidc.ExchangeRefreshTokenRequest{Aud: []string{"aud"}},
								Error: tt.errors[idx],
							}},
						},
					}, nil
				}
				return test.MockOIDCClient{
					STSClient: test.MockSTSClient{
						OnGetAccessToken: []test.STSOnGetAccessToken{{
							Given: &oidc.ExchangeRefreshTokenRequest{Aud: []string{"aud"}},
							Exchanged: &oidc.TokenPair{
								Token:        &oidc.RawToken{Token: "access"},
								RefreshToken: &oidc.RawToken{Token: "refresh"},
							},
						}},
					},
				}, nil
			}

			exch := New("issuer", "aud")
			gotAccess, gotRefresh, err := exch.Refresh(t.Context(), "tok")
			if (err != nil) != tt.wantErr {
				t.Fatalf("Refresh() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if gotAccess != "access" {
					t.Errorf("Refresh() access = %q, want %q", gotAccess, "access")
				}
				if gotRefresh != "refresh" {
					t.Errorf("Refresh() refresh = %q, want %q", gotRefresh, "refresh")
				}
			}
			if gotAttempts := int(attempts.Load()); gotAttempts != tt.wantAttempts {
				t.Errorf("Refresh() attempts = %d, want %d", gotAttempts, tt.wantAttempts)
			}
		})
	}
}

func TestNewHTTP1Transport(t *testing.T) {
	// Do not invoke tr.Proxy here: http.ProxyFromEnvironment latches the
	// proxy environment process-wide on first call, which would break
	// TestHTTP1DowngradeExchangerHonorsProxy's t.Setenv.
	tr := newHTTP1Transport()
	if tr.Proxy == nil {
		t.Error("Proxy is nil; HTTP(S)_PROXY would be ignored (CUS-1012)")
	}
	if tr.TLSNextProto == nil || len(tr.TLSNextProto) != 0 {
		t.Errorf("TLSNextProto = %v, want non-nil empty map to disable HTTP/2", tr.TLSNextProto)
	}
	if tr.ForceAttemptHTTP2 {
		t.Error("ForceAttemptHTTP2 = true, want false")
	}
	// The ALPN offer must not include h2: a cloned DefaultTransport that was
	// previously h2-initialized carries NextProtos = ["h2", "http/1.1"], and
	// a server selecting h2 speaks HTTP/2 frames at this HTTP/1.1-only
	// transport.
	if tr.TLSClientConfig == nil {
		t.Fatal("TLSClientConfig is nil, want NextProtos pinned to http/1.1")
	}
	if got := tr.TLSClientConfig.NextProtos; len(got) != 1 || got[0] != "http/1.1" {
		t.Errorf("TLSClientConfig.NextProtos = %v, want [http/1.1]", got)
	}
}

// TestHTTP1DowngradeExchangerHonorsProxy reproduces CUS-1012: the downgrade
// exchange must route through the proxy from the environment rather than
// dialing the issuer directly. A plain-HTTP issuer is used so the proxied
// request arrives at the proxy in absolute-URI form and can be answered
// there; the .invalid issuer host guarantees a direct dial would fail.
//
// Not parallel, and the only proxy-environment test in this binary:
// http.ProxyFromEnvironment caches the proxy environment process-wide on
// first use, so exactly one t.Setenv configuration can take effect per test
// process. Both exchanger methods are therefore exercised as subtests under
// the single latched environment.
func TestHTTP1DowngradeExchangerHonorsProxy(t *testing.T) {
	const issuerHost = "sts-cus-1012.invalid"
	var proxied atomic.Int32
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxied.Add(1)
		// Absolute-form URL (host present) proves the client treated us as
		// a proxy rather than the origin server.
		if r.URL.Host != issuerHost {
			t.Errorf("proxied request host = %q, want %q", r.URL.Host, issuerHost)
		}
		if r.Proto != "HTTP/1.1" {
			t.Errorf("proxied request proto = %q, want HTTP/1.1", r.Proto)
		}
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/sts/exchange":
			if got := r.Header.Get("Authorization"); got != "Bearer id-token" {
				t.Errorf("Authorization = %q, want %q", got, "Bearer id-token")
			}
			fmt.Fprint(w, `{"token": "cg-token", "refreshToken": "cg-refresh"}`)
		case "/sts/exchange_refresh_token":
			if got := r.Header.Get("Authorization"); got != "Bearer refresh-token" {
				t.Errorf("Authorization = %q, want %q", got, "Bearer refresh-token")
			}
			fmt.Fprint(w, `{"token": {"token": "new-access"}, "refreshToken": {"token": "new-refresh"}}`)
		default:
			t.Errorf("unexpected proxied request path %q", r.URL.Path)
			http.NotFound(w, r)
		}
	}))
	defer proxy.Close()

	for _, k := range []string{"HTTP_PROXY", "HTTPS_PROXY", "http_proxy", "https_proxy"} {
		t.Setenv(k, proxy.URL)
	}
	for _, k := range []string{"NO_PROXY", "no_proxy"} {
		t.Setenv(k, "")
	}

	exch := NewHTTP1DowngradeExchanger("http://"+issuerHost, "aud")

	t.Run("Exchange", func(t *testing.T) {
		pair, err := exch.Exchange(t.Context(), "id-token")
		if err != nil {
			t.Fatalf("Exchange() error = %v", err)
		}
		if pair.AccessToken != "cg-token" {
			t.Errorf("Exchange() access token = %q, want %q", pair.AccessToken, "cg-token")
		}
		if pair.RefreshToken != "cg-refresh" {
			t.Errorf("Exchange() refresh token = %q, want %q", pair.RefreshToken, "cg-refresh")
		}
	})

	t.Run("Refresh", func(t *testing.T) {
		access, refresh, err := exch.Refresh(t.Context(), "refresh-token")
		if err != nil {
			t.Fatalf("Refresh() error = %v", err)
		}
		if access != "new-access" {
			t.Errorf("Refresh() access token = %q, want %q", access, "new-access")
		}
		if refresh != "new-refresh" {
			t.Errorf("Refresh() refresh token = %q, want %q", refresh, "new-refresh")
		}
	})

	if got := int(proxied.Load()); got != 2 {
		t.Errorf("proxied request count = %d, want 2", got)
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// TestNewHTTP1TransportFallback covers a replaced http.DefaultTransport:
// consumers may wrap it with a non-*http.Transport RoundTripper
// (instrumentation, test doubles), and the constructor must fall back to
// stdlib-default construction rather than depend on the concrete type.
func TestNewHTTP1TransportFallback(t *testing.T) {
	orig := http.DefaultTransport
	http.DefaultTransport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("not used")
	})
	defer func() { http.DefaultTransport = orig }()

	tr := newHTTP1Transport()
	if tr.Proxy == nil {
		t.Error("fallback Proxy is nil; HTTP(S)_PROXY would be ignored")
	}
	if tr.TLSNextProto == nil || len(tr.TLSNextProto) != 0 {
		t.Errorf("fallback TLSNextProto = %v, want non-nil empty map", tr.TLSNextProto)
	}
	if tr.TLSClientConfig == nil || len(tr.TLSClientConfig.NextProtos) != 1 || tr.TLSClientConfig.NextProtos[0] != "http/1.1" {
		t.Errorf("fallback NextProtos = %v, want [http/1.1]", tr.TLSClientConfig)
	}
}

// TestHTTP1DowngradeALPNOffer verifies at the TLS layer that the downgrade
// client's ClientHello offers exactly ["http/1.1"] via ALPN. A cloned
// DefaultTransport that was previously h2-initialized would offer h2, letting
// a server negotiate HTTP/2 frames at this HTTP/1.1-only client. The exchange
// itself fails on certificate verification, which is irrelevant here — the
// server records the ALPN offer before the client rejects the cert.
func TestHTTP1DowngradeALPNOffer(t *testing.T) {
	var protos atomic.Value
	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	srv.TLS = &tls.Config{
		GetConfigForClient: func(chi *tls.ClientHelloInfo) (*tls.Config, error) {
			protos.Store(append([]string{}, chi.SupportedProtos...))
			return nil, nil
		},
	}
	srv.StartTLS()
	defer srv.Close()

	exch := NewHTTP1DowngradeExchanger(srv.URL, "aud")
	_, _ = exch.Exchange(t.Context(), "tok")

	got, _ := protos.Load().([]string)
	if len(got) != 1 || got[0] != "http/1.1" {
		t.Errorf("ClientHello ALPN offer = %v, want [http/1.1]", got)
	}
}

// TestHTTP1DowngradeExchangerRefusesRedirects guards the CheckRedirect
// behavior: a redirecting STS response must fail the exchange without the
// client following it, since oauth2.Transport would re-send the bearer token
// to the redirect target (golang/oauth2#314).
func TestHTTP1DowngradeExchangerRefusesRedirects(t *testing.T) {
	var followed atomic.Bool
	target := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		followed.Store(true)
	}))
	defer target.Close()

	issuer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer issuer.Close()

	exch := NewHTTP1DowngradeExchanger(issuer.URL, "aud")
	_, err := exch.Exchange(t.Context(), "secret-token")
	if err == nil || !strings.Contains(err.Error(), "302") {
		t.Errorf("Exchange() error = %v, want error carrying the 302 status", err)
	}
	if followed.Load() {
		t.Error("redirect was followed; bearer token would have been re-sent to the target")
	}
}

func TestHTTP1DowngradeExchangerSurfacesErrorBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, `{"code":16,"message":"JWT validation failed"}`, http.StatusBadRequest)
	}))
	defer srv.Close()

	exch := NewHTTP1DowngradeExchanger(srv.URL, "aud")
	_, err := exch.Exchange(t.Context(), "bad-token")
	if err == nil || !strings.Contains(err.Error(), "JWT validation failed") {
		t.Errorf("Exchange() error = %v, want error containing the response body", err)
	}
}

// gatewayFakeSTS implements the generated SecurityTokenServiceServer so the
// downgrade exchanger can be exercised against the real grpc-gateway HTTP
// bindings — the same code path the issuer serves. Recording the decoded
// requests pins the client/server contract: every field the exchanger sends
// must survive the gateway's binding into the RPC request message.
type gatewayFakeSTS struct {
	oidc.UnimplementedSecurityTokenServiceServer
	// mu orders handler-goroutine writes before test-goroutine reads; the
	// in-process HTTP round-trip alone is not a memory-model sync point.
	mu          sync.Mutex
	exchangeReq *oidc.ExchangeRequest
	refreshReq  *oidc.ExchangeRefreshTokenRequest
}

func (f *gatewayFakeSTS) Exchange(_ context.Context, req *oidc.ExchangeRequest) (*oidc.RawToken, error) {
	f.mu.Lock()
	f.exchangeReq = req
	f.mu.Unlock()
	if len(req.GetAud()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Audience must be specified")
	}
	return &oidc.RawToken{Token: "access", RefreshToken: "refresh"}, nil
}

func (f *gatewayFakeSTS) ExchangeRefreshToken(_ context.Context, req *oidc.ExchangeRefreshTokenRequest) (*oidc.TokenPair, error) {
	f.mu.Lock()
	f.refreshReq = req
	f.mu.Unlock()
	if len(req.GetAud()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Audience must be specified")
	}
	return &oidc.TokenPair{
		Token:        &oidc.RawToken{Token: "access"},
		RefreshToken: &oidc.RawToken{Token: "refresh"},
	}, nil
}

func newGatewaySTSServer(t *testing.T) (*gatewayFakeSTS, *httptest.Server) {
	t.Helper()
	fake := &gatewayFakeSTS{}
	mux := runtime.NewServeMux()
	if err := oidc.RegisterSecurityTokenServiceHandlerServer(t.Context(), mux, fake); err != nil {
		t.Fatalf("RegisterSecurityTokenServiceHandlerServer() = %v", err)
	}
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return fake, srv
}

// TestHTTP1DowngradeExchangeGatewayContract verifies the exchange request
// fields reach the server through the gateway's POST /sts/exchange binding.
// The binding has no body mapping — the gateway populates the request via
// ParseForm, which reads a body only when it is form-encoded. A JSON body
// is silently dropped and the issuer rejects the exchange with "Audience
// must be specified" even though the client sent an audience.
func TestHTTP1DowngradeExchangeGatewayContract(t *testing.T) {
	fake, srv := newGatewaySTSServer(t)

	exch := NewHTTP1DowngradeExchanger(srv.URL, "aud1,aud2",
		WithIdentity("identity-uid"),
		WithScope("scope-a", "scope-b"),
		WithCapabilities("cap1"),
		WithIdentityProvider("idp-uidp"),
	)
	pair, err := exch.Exchange(t.Context(), "external-token")
	if err != nil {
		t.Fatalf("Exchange() = %v", err)
	}
	if pair.AccessToken != "access" || pair.RefreshToken != "refresh" {
		t.Errorf("Exchange() = %+v, want access/refresh tokens", pair)
	}

	want := &oidc.ExchangeRequest{
		Aud:              []string{"aud1", "aud2"},
		Scope:            "scope-a",
		Scopes:           []string{"scope-a", "scope-b"},
		Identity:         "identity-uid",
		Cap:              []string{"cap1"},
		IdentityProvider: "idp-uidp",
	}
	fake.mu.Lock()
	got := fake.exchangeReq
	fake.mu.Unlock()
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("server-side ExchangeRequest (-want +got):\n%s", diff)
	}
}

// TestHTTP1DowngradeRefreshGatewayContract does the same for the refresh
// path, which must POST to the bound /sts/exchange_refresh_token route.
func TestHTTP1DowngradeRefreshGatewayContract(t *testing.T) {
	fake, srv := newGatewaySTSServer(t)

	exch := NewHTTP1DowngradeExchanger(srv.URL, "aud1", WithScope("scope-a"))
	access, refresh, err := exch.Refresh(t.Context(), "refresh-token")
	if err != nil {
		t.Fatalf("Refresh() = %v", err)
	}
	if access != "access" || refresh != "refresh" {
		t.Errorf("Refresh() = (%q, %q), want (access, refresh)", access, refresh)
	}

	want := &oidc.ExchangeRefreshTokenRequest{
		Aud:    []string{"aud1"},
		Scope:  "scope-a",
		Scopes: []string{"scope-a"},
	}
	fake.mu.Lock()
	got := fake.refreshReq
	fake.mu.Unlock()
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("server-side ExchangeRefreshTokenRequest (-want +got):\n%s", diff)
	}
}
