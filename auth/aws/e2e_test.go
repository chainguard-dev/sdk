/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/util/sets"
)

func epochTime() time.Time {
	return time.Unix(0, 0)
}

func TestGenerateToken(t *testing.T) {
	timeNow = epochTime

	var creds = aws.Credentials{AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION"}
	validToken, err := GenerateToken(context.Background(), creds, "aud", "identity")
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/good":
			fmt.Fprintf(w, `{
			"GetCallerIdentityResponse": {
				"GetCallerIdentityResult": {
					"UserId": "userid",
					"Account": "123456789012",
					"Arn": "arn"
				}
			}
			}`)
		case "/reject":
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "go away")
		}
	}))

	tests := map[string]struct {
		audience  string
		identity  string
		timestamp time.Time
		stsURL    string
		token     string
		wantErr   error
	}{
		"happy path": {
			audience:  "aud",
			identity:  "identity",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0),
			token:     validToken,
			wantErr:   nil,
		},
		"audience mismatch": {
			audience:  "bad audience",
			identity:  "identity",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0),
			token:     validToken,
			wantErr:   ErrInvalidAudience,
		},
		"identity mismatch": {
			audience:  "aud",
			identity:  "identity is wrong",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0),
			token:     validToken,
			wantErr:   ErrInvalidIdentity,
		},
		"expired token": {
			audience:  "aud",
			identity:  "identity",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0).Add(16 * time.Minute),
			token:     validToken,
			wantErr:   ErrTokenExpired,
		},
		"no identity set": {
			audience:  "aud",
			identity:  "",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0),
			token:     validToken,
			wantErr:   ErrInvalidVerificationConfiguration,
		},
		"bad token encoding": {
			audience:  "aud",
			identity:  "identity",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0),
			token:     "no-base64-encoded",
			wantErr:   ErrInvalidEncoding,
		},
		"base64 encoded but not even an http request": {
			audience:  "aud",
			identity:  "identity",
			stsURL:    ts.URL + "/good",
			timestamp: time.Unix(0, 0),
			token:     "bm90IGFuIGh0dHAgcmVxdWVzdA==", // echo -n "not an http request" | base64
			wantErr:   ErrInvalidEncoding,
		},
		"bad sts URL": {
			audience:  "aud",
			identity:  "identity",
			timestamp: time.Unix(0, 0),
			token:     validToken,
			stsURL:    "\t\n",
			wantErr:   ErrInvalidVerificationConfiguration,
		},
		"rejected token": {
			audience:  "aud",
			identity:  "identity",
			timestamp: time.Unix(0, 0),
			token:     validToken,
			stsURL:    ts.URL + "/reject",
			wantErr:   ErrTokenRejected,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotClaims, err := VerifyToken(
				context.Background(),
				test.token,
				WithAudience(sets.New(test.audience)),
				WithIdentity(test.identity),
				withSTSURL(test.stsURL),
				withTimestamp(test.timestamp))
			if err != test.wantErr { //nolint:errorlint
				t.Errorf("received error %T and wanted %T", err, test.wantErr)
			}
			if err != nil {
				return
			}

			expectedClaims := &VerifiedClaims{
				UserID:  "userid",
				Account: "123456789012",
				Arn:     "arn",
			}
			if diff := cmp.Diff(gotClaims, expectedClaims); diff != "" {
				t.Errorf("got different claims than expected diff = %s", diff)
			}
		})
	}
}
