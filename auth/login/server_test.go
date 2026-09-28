/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestServerTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
	defer cancel()
	s, err := newServer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	_, err = s.Token()
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("expect timeout error getting token")
	}
}

func TestServerHappyPath(t *testing.T) {
	s, err := newServer(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	callback := strings.ReplaceAll(s.URL(), "token=true", "token=foo")
	http.Get(callback)

	token, err := s.Token()
	if err != nil {
		t.Errorf("Token(): got error = %#v, want nil", err)
	}
	if token != "foo" {
		t.Errorf("Token(): got = %q, want = %q", token, "foo")
	}

	s.Close()
}
