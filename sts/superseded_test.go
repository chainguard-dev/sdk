/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"chainguard.dev/sdk/auth/token"
)

// The issuer's wire contract, restated rather than read from the package under
// test: a test sharing the constant would pass even if both drifted from what
// the issuer actually sends.
const (
	supersededTestDomain = "oidc.chainguard.dev"
	supersededTestReason = "REFRESH_TOKEN_SUPERSEDED"
)

func supersededStatus() *errdetails.ErrorInfo {
	return &errdetails.ErrorInfo{Domain: supersededTestDomain, Reason: supersededTestReason}
}

// TestAsSupersededRefreshToken pins the classification itself, including the
// property the re-authentication path depends on: marking the error must not
// hide the gRPC status, or refreshTokenRejected stops recognizing it and an
// interactive client fails instead of logging in again.
func TestAsSupersededRefreshToken(t *testing.T) {
	supersededErr := func(c codes.Code, info *errdetails.ErrorInfo) error {
		st := status.New(c, "refused")
		if info != nil {
			withDetail, err := st.WithDetails(info)
			if err != nil {
				t.Fatalf("attaching ErrorInfo: %v", err)
			}
			st = withDetail
		}
		return st.Err()
	}

	for _, tc := range []struct {
		name       string
		err        error
		superseded bool
	}{
		{"nil", nil, false},
		{"not a status error", errors.New("dial tcp: connection refused"), false},
		{"the refusal", supersededErr(codes.Unauthenticated, supersededStatus()), true},
		{"no details", supersededErr(codes.Unauthenticated, nil), false},
		{"wrong domain", supersededErr(codes.Unauthenticated, &errdetails.ErrorInfo{Domain: "other.example", Reason: supersededTestReason}), false},
		{"wrong reason", supersededErr(codes.Unauthenticated, &errdetails.ErrorInfo{Domain: supersededTestDomain, Reason: "SOMETHING_ELSE"}), false},
		{"right reason, wrong code", supersededErr(codes.NotFound, supersededStatus()), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := AsSupersededRefreshToken(tc.err)
			if errors.Is(got, token.ErrPresentedTokenSuperseded) != tc.superseded {
				t.Errorf("AsSupersededRefreshToken(%v) superseded = %v, want %v", tc.err, !tc.superseded, tc.superseded)
			}
			if tc.err == nil {
				if got != nil {
					t.Errorf("AsSupersededRefreshToken(nil) = %v, want nil", got)
				}
				return
			}
			if !errors.Is(got, tc.err) {
				t.Errorf("AsSupersededRefreshToken() dropped the original error: %v", got)
			}
			// The status has to survive, or re-authentication stops working.
			if want := status.Code(tc.err); status.Code(got) != want {
				t.Errorf("status.Code() = %v after marking, want %v", status.Code(got), want)
			}
			// Clients re-authenticate on the code, so marking the error must
			// leave UNAUTHENTICATED visible or an interactive client fails
			// instead of logging in again.
			if tc.superseded && status.Code(got) != codes.Unauthenticated {
				t.Errorf("status.Code() = %v after marking, want Unauthenticated", status.Code(got))
			}
		})
	}
}

// TestAsSupersededRefreshToken_HTTP1Downgrade covers the transport that
// reconstructs the error instead of receiving it. doHTTP1 rebuilds a failed
// exchange from the gateway's status body, and if it kept only the code and
// the message this classifier would take its "no such detail" branch on this
// path alone -- keeping a refresh token the issuer refused, whose next
// presentation is replay and revokes the family. chainctl reaches this path
// under CHAINCTL_AUTH_STS_HTTP1_DOWNGRADE and microvm's delegated token source
// uses it outright, so it is not a theoretical transport.
func TestAsSupersededRefreshToken_HTTP1Downgrade(t *testing.T) {
	// The body grpc-gateway serializes for the issuer's refusal.
	body := `{"code":16,"message":"refresh token superseded","details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","domain":"` +
		supersededTestDomain + `","reason":"` + supersededTestReason + `"}]}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(srv.Close)

	ex := NewHTTP1DowngradeExchanger(srv.URL, "cgr.test")
	_, _, err := ex.Refresh(t.Context(), "a-refresh-token")
	if err == nil {
		t.Fatal("Refresh() = nil, want the refusal")
	}
	if got := status.Code(err); got != codes.Unauthenticated {
		t.Errorf("status.Code() = %v, want Unauthenticated", got)
	}
	if marked := AsSupersededRefreshToken(err); !errors.Is(marked, token.ErrPresentedTokenSuperseded) {
		t.Errorf("AsSupersededRefreshToken() did not recognize the refusal over HTTP/1: %v", marked)
	}
}

// TestDetailsFromGatewayStatusBody keeps the recovery from turning an
// unparsable body into a failure: the caller still gets its coded error.
func TestDetailsFromGatewayStatusBody(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		want int
	}{
		{"the refusal", `{"code":16,"details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","domain":"d","reason":"r"}]}`, 1},
		{"no details", `{"code":16,"message":"nope"}`, 0},
		{"not json", `502 Bad Gateway`, 0},
		{"empty", ``, 0},
		{"unknown field alongside details", `{"code":16,"somethingNew":1,"details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","domain":"d","reason":"r"}]}`, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := len(detailsFromGatewayStatusBody([]byte(tc.body))); got != tc.want {
				t.Errorf("detailsFromGatewayStatusBody() returned %d details, want %d", got, tc.want)
			}
		})
	}
}
