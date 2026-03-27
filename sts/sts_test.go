/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"chainguard.dev/sdk/proto/platform/oidc/v1/test"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		name:         "does not retry on Internal",
		errors:       []error{status.Error(codes.Internal, "internal error")},
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
