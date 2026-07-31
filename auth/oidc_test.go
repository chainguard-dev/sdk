/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"encoding/base64"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// makeJWT builds a minimal unsigned JWT with the given raw JSON payload.
func makeJWT(payloadJSON string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(payloadJSON))
	return header + "." + payload + ".sig"
}

func TestExtractAudiences(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		want    []string
		wantErr bool
	}{{
		name:  "single string aud",
		token: makeJWT(`{"aud":"https://console-api.enforce.dev"}`),
		want:  []string{"https://console-api.enforce.dev"},
	}, {
		name:  "array aud",
		token: makeJWT(`{"aud":["aud-a","aud-b"]}`),
		want:  []string{"aud-a", "aud-b"},
	}, {
		name:  "missing aud",
		token: makeJWT(`{"sub":"someone"}`),
		want:  nil,
	}, {
		name:    "malformed token",
		token:   "not-a-jwt",
		wantErr: true,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ExtractAudiences(test.token)
			if (err != nil) != test.wantErr {
				t.Fatalf("ExtractAudiences() error = %v, wantErr = %t", err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("ExtractAudiences() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
