/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package receiver_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"chainguard.dev/go-oidctest/pkg/oidctest"
	"chainguard.dev/sdk/events/receiver"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
)

const testGroup = "d2020cd12118237dcd06214184a6d17c1f831a08"

var testBody = []byte(`{"repository":"chainguard-images/static","tag":"latest"}`)

func digestOf(b []byte) string {
	return fmt.Sprintf("sha256:%x", sha256.Sum256(b))
}

func mintToken(t *testing.T, signer jose.Signer, issuer, subject, digest string) string {
	t.Helper()

	tok, err := jwt.Signed(signer).Claims(jwt.Claims{
		Issuer:   issuer,
		IssuedAt: jwt.NewNumericDate(time.Now()),
		Expiry:   jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:  subject,
		Audience: jwt.Audience{"customer"},
	}).Claims(map[string]any{
		"webhook": map[string]string{"digest": digest},
	}).Serialize()
	if err != nil {
		t.Fatalf("Serialize() = %v", err)
	}
	return tok
}

// eventRequest builds a binary-mode CloudEvent delivery. An empty token omits
// the Authorization header entirely, as an unauthenticated caller would.
func eventRequest(body []byte, token string, cloudEvent bool) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if cloudEvent {
		req.Header.Set("ce-specversion", "1.0")
		req.Header.Set("ce-type", "dev.chainguard.test.v1")
		req.Header.Set("ce-source", "/test")
		req.Header.Set("ce-id", "9f0d204d-9579-4e03-89d3-33810a5d4dd6")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func TestNewHTTPHandler(t *testing.T) {
	signer, issuer := oidctest.NewIssuer(t)

	for name, tc := range map[string]struct {
		token      string
		body       []byte
		cloudEvent bool
		wantStatus int
		wantCalled bool
	}{
		"verified delivery": {
			token:      mintToken(t, signer, issuer, "webhook:"+testGroup, digestOf(testBody)),
			body:       testBody,
			cloudEvent: true,
			wantStatus: http.StatusOK,
			wantCalled: true,
		},
		// Regression: this combination used to panic on a nil dereference and
		// surface as a 500 with no explanation.
		"missing authorization header": {
			body:       testBody,
			cloudEvent: true,
			wantStatus: http.StatusUnauthorized,
		},
		"token for another group": {
			token:      mintToken(t, signer, issuer, "webhook:7c9e6679742f4b2ea1c0ff1d3a6b5e04", digestOf(testBody)),
			body:       testBody,
			cloudEvent: true,
			wantStatus: http.StatusForbidden,
		},
		"subject is not a webhook": {
			token:      mintToken(t, signer, issuer, "user:"+testGroup, digestOf(testBody)),
			body:       testBody,
			cloudEvent: true,
			wantStatus: http.StatusForbidden,
		},
		"digest does not match the body": {
			token:      mintToken(t, signer, issuer, "webhook:"+testGroup, digestOf([]byte(`{"different":"payload"}`))),
			body:       testBody,
			cloudEvent: true,
			wantStatus: http.StatusForbidden,
		},
		"not a cloudevent": {
			token:      mintToken(t, signer, issuer, "webhook:"+testGroup, digestOf(testBody)),
			body:       testBody,
			wantStatus: http.StatusBadRequest,
		},
	} {
		t.Run(name, func(t *testing.T) {
			called := false
			h, err := receiver.NewHTTPHandler(t.Context(), issuer, testGroup,
				func(context.Context, cloudevents.Event) error {
					called = true
					return nil
				})
			if err != nil {
				t.Fatalf("NewHTTPHandler() = %v", err)
			}

			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, eventRequest(tc.body, tc.token, tc.cloudEvent))

			if got, want := rec.Code, tc.wantStatus; got != want {
				t.Errorf("status: got = %d, want = %d (body %q)", got, want, rec.Body.String())
			}
			if got, want := called, tc.wantCalled; got != want {
				t.Errorf("handler invoked: got = %v, want = %v", got, want)
			}
		})
	}
}

// TestNewWithoutRequestData pins the behaviour when the Handler from New is
// served by something that does not attach the HTTP request data — notably the
// CloudEvents client's StartReceiver, which the package documented for years.
// It must report an actionable error rather than panicking.
func TestNewWithoutRequestData(t *testing.T) {
	_, issuer := oidctest.NewIssuer(t)

	called := false
	h, err := receiver.New(t.Context(), issuer, testGroup,
		func(context.Context, cloudevents.Event) error {
			called = true
			return nil
		})
	if err != nil {
		t.Fatalf("New() = %v", err)
	}

	event := cloudevents.NewEvent()
	event.SetID("9f0d204d-9579-4e03-89d3-33810a5d4dd6")
	event.SetType("dev.chainguard.test.v1")
	event.SetSource("/test")

	err = h(t.Context(), event)
	if err == nil {
		t.Fatal("want an error for a context with no request data, got nil")
	}

	var result *cehttp.Result
	if !errors.As(err, &result) {
		t.Fatalf("want a *cehttp.Result, got %T: %v", err, err)
	}
	if got, want := result.StatusCode, http.StatusInternalServerError; got != want {
		t.Errorf("status: got = %d, want = %d", got, want)
	}
	if !strings.Contains(err.Error(), "NewHTTPHandler") {
		t.Errorf("error should point at the fix, got = %q", err.Error())
	}
	if called {
		t.Error("handler must not be invoked when the token cannot be read")
	}
}
