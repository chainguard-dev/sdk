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
	. "chainguard.dev/sdk/proto/platform/oidc/v1/test"
	"github.com/google/go-cmp/cmp"
)

func TestImplExchange(t *testing.T) {
	tests := map[string]struct {
		issuer       string
		audience     string
		newOpts      []ExchangerOption
		exchangeOpts []ExchangerOption
		want         string
		wantErr      bool
		clientMock   MockOIDCClient
	}{
		"zero options": {
			issuer:   "bar",
			audience: "baz",
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud: []string{"baz"},
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			want: "token!",
		},
		"basic error plumbing": {
			issuer:   "bar",
			audience: "baz",
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
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
				WithCapabilities("registry.push"),
				WithScope("derp"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:   []string{"baz"},
							Cap:   []string{"registry.push"},
							Scope: "derp",
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			want: "token!",
		},
		"capabilities and scope on exchange": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCapabilities("registry.push"),
				WithScope("derp"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:   []string{"baz"},
							Cap:   []string{"registry.push"},
							Scope: "derp",
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			want: "token!",
		},
		"identity": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIdentity("my-identity"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:      []string{"baz"},
							Identity: "my-identity",
						},
						Exchanged: &oidc.RawToken{Token: "token foo"},
					}},
				},
			},
			want: "token foo",
		},
		"include upstream": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIncludeUpstreamToken(),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:                  []string{"baz"},
							IncludeUpstreamToken: true,
						},
						Exchanged: &oidc.RawToken{Token: "tokenz"},
					}},
				},
			},
			want: "tokenz",
		},
		"cluster": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCluster("kind i presume"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:     []string{"baz"},
							Cluster: "kind i presume",
						},
						Exchanged: &oidc.RawToken{Token: "tokenz"},
					}},
				},
			},
			want: "tokenz",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oidcNewClients = func(_ context.Context, issuer string, token string, opts ...oidc.ClientOption) (oidc.Clients, error) {
				return test.clientMock, nil
			}

			exch := New(test.issuer, test.audience, test.newOpts...)
			gotToken, gotErr := exch.Exchange(context.Background(), "foo", test.exchangeOpts...)
			if gotErr != nil && !test.wantErr {
				t.Error("got err: ", gotErr, "and expected no error")
			}
			if diff := cmp.Diff(gotToken, test.want); diff != "" {
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
		want         string
		wantErr      bool
		clientMock   MockOIDCClient
	}{
		"zero options": {
			issuer:   "bar",
			audience: "baz",
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud: []string{"baz"},
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			want: "token!",
		},
		"basic error plumbing": {
			issuer:   "bar",
			audience: "baz",
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
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
				WithCapabilities("registry.push"),
				WithScope("derp"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:   []string{"baz"},
							Cap:   []string{"registry.push"},
							Scope: "derp",
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			want: "token!",
		},
		"capabilities and scope on exchange": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCapabilities("registry.push"),
				WithScope("derp"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:   []string{"baz"},
							Cap:   []string{"registry.push"},
							Scope: "derp",
						},
						Exchanged: &oidc.RawToken{Token: "token!"},
					}},
				},
			},
			want: "token!",
		},
		"identity": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIdentity("my-identity"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:      []string{"baz"},
							Identity: "my-identity",
						},
						Exchanged: &oidc.RawToken{Token: "token foo"},
					}},
				},
			},
			want: "token foo",
		},
		"include upstream": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithIncludeUpstreamToken(),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:                  []string{"baz"},
							IncludeUpstreamToken: true,
						},
						Exchanged: &oidc.RawToken{Token: "tokenz"},
					}},
				},
			},
			want: "tokenz",
		},
		"cluster": {
			issuer:   "bar",
			audience: "baz",
			exchangeOpts: []ExchangerOption{
				WithCluster("kind i presume"),
			},
			clientMock: MockOIDCClient{
				STSClient: MockSTSClient{
					OnExchange: []STSOnExchange{{
						Given: &oidc.ExchangeRequest{
							Aud:     []string{"baz"},
							Cluster: "kind i presume",
						},
						Exchanged: &oidc.RawToken{Token: "tokenz"},
					}},
				},
			},
			want: "tokenz",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			oidcNewClients = func(_ context.Context, issuer string, token string, opts ...oidc.ClientOption) (oidc.Clients, error) {
				return test.clientMock, nil
			}

			gotToken, gotErr := Exchange(context.Background(), test.issuer, test.audience, "foo", test.exchangeOpts...)
			if gotErr != nil && !test.wantErr {
				t.Error("got err: ", gotErr, "and expected no error")
			}
			if diff := cmp.Diff(gotToken, test.want); diff != "" {
				t.Error("Got unexpected diff in token: ", diff)
			}
		})
	}
}
