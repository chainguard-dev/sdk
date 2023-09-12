/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"net/http/httptest"
	"testing"

	"chainguard.dev/sdk/pkg/uidp"
	"github.com/google/go-cmp/cmp"
)

func TestConfFromOptions(t *testing.T) {
	id := uidp.NewUID().String()

	testIssuer := httptest.NewServer(expectedBehavior("chainguard.dev"))

	tests := map[string]struct {
		Options    []Option
		WantConfig *config
		WantErr    bool
	}{
		"Happy path": {
			Options: []Option{
				WithIssuer("https://example.com"),
				WithIdentityProvider(id),
				WithInviteCode("foo"),
			},
			WantConfig: &config{
				Issuer:     "https://example.com",
				IDP:        id,
				InviteCode: "foo",
			},
		},
		"Bad IDP ID": {
			Options: []Option{
				WithIssuer("https://example.com"),
				WithIdentityProvider("imnotanidp"),
			},
			WantErr: true,
		},
		"Org name": {
			Options: []Option{
				WithIssuer(testIssuer.URL),
				WithOrgName("chainguard.dev"),
			},
			WantConfig: &config{
				Issuer:  testIssuer.URL,
				OrgName: "chainguard.dev",
			},
		},
		"Cannot specify both identity provider and org name": {
			Options: []Option{
				WithOrgName("chainguard.dev"),
				WithIdentityProvider("foo"),
			},
			WantErr: true,
		},
		"Cannot specify both client id and org name": {
			Options: []Option{
				WithOrgName("chainguard.dev"),
				WithClientID("auth0"),
			},
			WantErr: true,
		},
		"No issuer defaults to prod issuer": {
			Options: []Option{
				WithIdentityProvider(id),
			},
			WantConfig: &config{
				Issuer: defaultIssuer,
				IDP:    id,
			},
		},
		"No idp ID or client ID set errors": {
			Options: nil,
			WantErr: true,
		},
		"IDP and ClientID both set errors": {
			Options: []Option{
				WithIdentityProvider(id),
				WithClientID("client_id"),
			},
			WantErr: true,
		},
		"IDP and Auth0Connection both set errors": {
			Options: []Option{
				WithIdentityProvider(id),
				WithAuth0Connection("github"),
			},
			WantErr: true,
		},
	}

	for test, data := range tests {
		t.Run(test, func(t *testing.T) {
			gotConfig, err := newConfigFromOptions(data.Options...)
			if err != nil && !data.WantErr {
				t.Errorf("got unexpected error %#v", err)
				return
			} else if err == nil && data.WantErr {
				t.Error("expected error and got none")
				return
			}

			if diff := cmp.Diff(gotConfig, data.WantConfig); diff != "" {
				t.Errorf("diff in expected config = %s", diff)
			}
		})
	}
}
