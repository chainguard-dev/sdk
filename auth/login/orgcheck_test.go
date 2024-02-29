/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func expectedBehavior(verifiedOrg string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/orgcheck" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "name query param required")
			return
		}
		fmt.Fprint(w, name == verifiedOrg)
	}
}

func brokenIssuer(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func TestOrgCheck(t *testing.T) {
	tests := map[string]struct {
		org     string
		handler http.HandlerFunc
		want    bool
		wantErr bool
	}{
		"verified": {
			org:     "example.com",
			handler: expectedBehavior("example.com"),
			want:    true,
		},
		"not verified": {
			org:     "not.example.com",
			handler: expectedBehavior("example.com"),
			want:    false,
		},
		"broken issuer should error": {
			org:     "example.com",
			handler: brokenIssuer,
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			s := httptest.NewServer(test.handler)
			defer s.Close()

			got, err := orgCheck(test.org, s.URL)
			if (err != nil) != test.wantErr {
				t.Errorf("err == %s and wantErr = %v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("got = %v and want = %v", got, test.want)
			}
		})
	}
}
