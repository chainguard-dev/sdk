/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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

// deliverToken simulates the browser completing the OIDC flow: it parses the
// callback (redirect) URL out of the oauth URL that Login built and hits it with
// a token, which unblocks Login's local token server.
func deliverToken(t *testing.T, oauthURL, token string) {
	t.Helper()
	u, err := url.Parse(oauthURL)
	if err != nil {
		t.Errorf("parsing oauth URL %q: %v", oauthURL, err)
		return
	}
	cb, err := url.Parse(u.Query().Get("redirect"))
	if err != nil {
		t.Errorf("parsing callback URL: %v", err)
		return
	}
	q := cb.Query()
	q.Set("token", token)
	cb.RawQuery = q.Encode()

	resp, err := http.Get(cb.String()) //nolint:gosec // local test callback server
	if err != nil {
		t.Errorf("delivering token to callback: %v", err)
		return
	}
	_ = resp.Body.Close()
}

func TestLoginURLOpener(t *testing.T) {
	t.Run("url opener is used instead of the default browser", func(t *testing.T) {
		var opened string
		opener := func(u string) error {
			opened = u
			// Drive the flow to completion so Login returns.
			deliverToken(t, u, "opener-token")
			return nil
		}

		token, _, err := Login(t.Context(),
			WithIssuer("https://issuer.chaintest.net"),
			WithClientID("test-client"),
			WithMessageWriter(io.Discard),
			WithURLOpener(opener),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "opener-token" {
			t.Errorf("token: got = %q, want = %q", token, "opener-token")
		}
		if !strings.HasPrefix(opened, "https://issuer.chaintest.net/oauth?") {
			t.Errorf("opener received unexpected URL: %q", opened)
		}
	})

	t.Run("skip browser takes precedence over url opener", func(t *testing.T) {
		openerCalled := false
		opener := func(string) error {
			openerCalled = true
			return nil
		}

		// SkipBrowser prints the URL to the message writer instead of opening a
		// browser; deliver the token from there so Login can return.
		token, _, err := Login(t.Context(),
			WithIssuer("https://issuer.chaintest.net"),
			WithClientID("test-client"),
			WithSkipBrowser(),
			WithURLOpener(opener),
			WithMessageWriter(writerFunc(func(p []byte) (int, error) {
				if i := strings.Index(string(p), "http"); i >= 0 {
					deliverToken(t, strings.TrimSpace(string(p)[i:]), "skip-token")
				}
				return len(p), nil
			})),
		)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "skip-token" {
			t.Errorf("token: got = %q, want = %q", token, "skip-token")
		}
		if openerCalled {
			t.Error("url opener must not be called when SkipBrowser is set")
		}
	})
}

// writerFunc adapts a function to io.Writer.
type writerFunc func(p []byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

func TestBuildHeadlessURL(t *testing.T) {
	// Set up a test server for organization verification
	orgServer := httptest.NewServer(expectedBehavior("my-org"))
	defer orgServer.Close()

	for _, tt := range []struct {
		name    string
		opts    []Option
		want    string
		wantErr string
	}{{
		name: "no code, error",
		opts: []Option{
			WithIssuer("https://issuer.chaintest.net"),
			WithIdentityProvider("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		},
		wantErr: "headless code is required",
	}, {
		name: "has code, no idp",
		opts: []Option{
			WithHeadlessCode("code"),
			WithIssuer("https://issuer.chaintest.net"),
		},
		want: "https://issuer.chaintest.net/oauth?headless_code=code",
	}, {
		name: "has code, has idp",
		opts: []Option{
			WithHeadlessCode("code"),
			WithIssuer("https://issuer.chaintest.net"),
			WithIdentityProvider("deadbeefdeadbeefdeadbeefdeadbeefdeadbeef"),
		},
		want: "https://issuer.chaintest.net/oauth?headless_code=code&idp_id=deadbeefdeadbeefdeadbeefdeadbeefdeadbeef",
	}, {
		name: "has code, has connection",
		opts: []Option{
			WithHeadlessCode("code"),
			WithIssuer("https://issuer.chaintest.net"),
			WithAuth0Connection("google-oauth2"),
		},
		want: "https://issuer.chaintest.net/oauth?connection=google-oauth2&headless_code=code",
	}, {
		name: "has code, has connection and client id",
		opts: []Option{
			WithHeadlessCode("code"),
			WithIssuer("https://issuer.chaintest.net"),
			WithAuth0Connection("email"),
			WithClientID("clientid"),
		},
		want: "https://issuer.chaintest.net/oauth?client_id=clientid&connection=email&headless_code=code",
	}, {
		name: "has code, has org name",
		opts: []Option{
			WithHeadlessCode("code"),
			WithIssuer(orgServer.URL),
			WithOrgName("my-org"),
		},
		want: orgServer.URL + "/oauth?headless_code=code&idp_id=my-org",
	}, {
		// ScreenHint is intentionally NOT included in the headless URL —
		// the headless flow doesn't render Auth0's Universal Login UI,
		// so the hint has no effect. Setting WithScreenHint must be a
		// no-op here; the parallel interactive flow (Login) is where it
		// matters. Locks in the BuildHeadlessURL contract.
		name: "screen_hint is not propagated to headless URL",
		opts: []Option{
			WithHeadlessCode("code"),
			WithIssuer("https://issuer.chaintest.net"),
			WithAuth0Connection("email"),
			WithClientID("clientid"),
			WithScreenHint("signup"),
		},
		want: "https://issuer.chaintest.net/oauth?client_id=clientid&connection=email&headless_code=code",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BuildHeadlessURL(tt.opts...)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("got nil error, want error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got error %s, want error containing %s", err, tt.wantErr)
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
