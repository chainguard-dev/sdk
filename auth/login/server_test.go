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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
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
	s, err := newServer(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	callback := strings.ReplaceAll(s.URL(), "token=true", "token=foo")
	http.Get(callback)

	token, err := s.Token()
	if err != nil {
		t.Errorf("expected no error, got %#v", err)
	}
	if token != "foo" {
		t.Errorf("expected token == foo, but got token == %q", token)
	}

	s.Close()
}
