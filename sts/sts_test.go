/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"context"
	"errors"
	"testing"

	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"chainguard.dev/sdk/proto/platform/oidc/v1/test"
	"github.com/google/go-cmp/cmp"
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
			token, refreshToken, gotErr := exch.Refresh(context.Background(), "foo", test.exchangeOpts...)
			if gotErr != nil && !test.wantErr {
				t.Error("got err: ", gotErr, "and expected no error")
			}
			if diff := cmp.Diff(token, test.wantToken); diff != "" {
				t.Error("Got unexpected diff in token: ", diff)
			}
			if diff := cmp.Diff(refreshToken, test.wantRefreshToken); diff != "" {
				t.Error("Got unexpected diff in refresh token: ", diff)
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
			gotTok, gotErr := exch.Exchange(context.Background(), "foo", test.exchangeOpts...)
			if gotErr != nil && !test.wantErr {
				t.Error("got err: ", gotErr, "and expected no error")
			}
			if diff := cmp.Diff(test.wantToken, gotTok); diff != "" {
				t.Error("Got unexpected diff in token: ", diff)
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

			gotTok, gotErr := ExchangePair(context.Background(), test.issuer, test.audience, "foo", test.exchangeOpts...)
			if gotErr != nil && !test.wantErr {
				t.Error("got err: ", gotErr, "and expected no error")
			}
			if diff := cmp.Diff(test.wantToken, gotTok); diff != "" {
				t.Error("Got unexpected diff in token: ", diff)
			}
		})
	}
}
