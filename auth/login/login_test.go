/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestOpenBrowserErrorAs(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"nil": {
			err:  nil,
			want: false,
		},
		"success": {
			err: &OpenBrowserError{ //nolint:staticcheck
				errors.New("unit test"),
			},
			want: true,
		},
		"success - login.Error": {
			err: &Error{
				Details: openBrowserError,
				Err:     errors.New("unit test"),
			},
			want: true,
		},
		"failure": {
			err:  errors.New("unit test"),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var want *OpenBrowserError //nolint:staticcheck
			got := errors.As(test.err, &want)
			if got != test.want {
				t.Errorf("As() expected %t, got %t", test.want, got)
			}
			if got {
				t.Log(want.Error())
			}
		})
	}
}

func TestErrorAs(t *testing.T) {
	tests := map[string]struct {
		err  error
		want bool
	}{
		"nil": {
			err:  nil,
			want: false,
		},
		"success - login.Error": {
			err: &Error{
				Details: remoteServerError,
				Err:     errors.New("unit test"),
			},
			want: true,
		},
		"failure": {
			err:  errors.New("unit test"),
			want: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var want *Error
			got := errors.As(test.err, &want)
			if got != test.want {
				t.Errorf("As() expected %t, got %t", test.want, got)
			}
			if got {
				t.Log(want.Error())
			}
		})
	}
}

func TestBuildHeadlessURL(t *testing.T) {
	// Set up a test server for organization verification
	orgServer := httptest.NewServer(expectedBehavior("my-org"))
	defer orgServer.Close()

	for _, tt := range []struct {
		name    string
		opts    []Option
		want    string
		wantErr string
	}{
		{
			name: "no code, error",
			opts: []Option{
				WithIssuer("https://issuer.chaintest.net"),
				WithIdentityProvider("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			},
			wantErr: "headless code is required",
		},
		{
			name: "has code, no idp",
			opts: []Option{
				WithHeadlessCode("code"),
				WithIssuer("https://issuer.chaintest.net"),
			},
			want: "https://issuer.chaintest.net/oauth?headless_code=code",
		},
		{
			name: "has code, has idp",
			opts: []Option{
				WithHeadlessCode("code"),
				WithIssuer("https://issuer.chaintest.net"),
				WithIdentityProvider("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
			},
			want: "https://issuer.chaintest.net/oauth?headless_code=code&idp_id=deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
		},
		{
			name: "has code, has connection",
			opts: []Option{
				WithHeadlessCode("code"),
				WithIssuer("https://issuer.chaintest.net"),
				WithAuth0Connection("google-oauth2"),
			},
			want: "https://issuer.chaintest.net/oauth?connection=google-oauth2&headless_code=code",
		},
		{
			name: "has code, has org name",
			opts: []Option{
				WithHeadlessCode("code"),
				WithIssuer(orgServer.URL),
				WithOrgName("my-org"),
			},
			want: orgServer.URL + "/oauth?headless_code=code&idp_id=my-org",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildHeadlessURL(tt.opts...)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, got none")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error %s, got %s", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("unexpected diff: %s", diff)
			}
		})
	}
}
