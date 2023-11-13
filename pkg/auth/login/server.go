/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package login implements client login functionality shared between various
// clients
package login

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

type server struct {
	ctx context.Context

	callbackURL string
	rootURL     string

	token chan string

	l net.Listener
}

func newServer(ctx context.Context) (*server, error) {
	s := new(server)

	{
		// NB: port 0 here means select a random port
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, err
		}
		port := listener.Addr().(*net.TCPAddr).Port

		s.token = make(chan string, 1)

		s.callbackURL = fmt.Sprintf("http://localhost:%d/callback?token=true", port)
		s.rootURL = fmt.Sprintf("http://localhost:%d/", port)
		s.l = listener
		s.ctx = ctx
	}

	// nolint:all
	go http.Serve(s.l, s)

	return s, nil
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		fmt.Fprint(w, HTMLAuthSuccessful)
		close(s.token)
		return

	case "/callback":
		token := r.URL.Query().Get("token")
		if token == "" {
			fmt.Fprint(w, "Account not found, registration required.")
			close(s.token)
			return
		}
		s.token <- token
		// We redirect to `/` to print a "successful auth" message and strip
		// the token out of the URI
		http.Redirect(w, r, s.rootURL, http.StatusPermanentRedirect)

	default:
		http.NotFound(w, r)
	}
}

func (s *server) URL() string {
	return s.callbackURL
}

// Token blocks until a token has been received
func (s *server) Token() (string, error) {
	select {
	case <-s.ctx.Done():
		return "", s.ctx.Err()
	case t, ok := <-s.token:
		if !ok {
			// Didn't receive token from redirect to callback
			return "", errors.New("login failed")
		}

		// We've received a token, but need to block until the success page has been written to
		// the browser
		<-s.token
		return t, nil
	}
}

func (s *server) Close() error {
	return s.l.Close()
}
