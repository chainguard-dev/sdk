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
	"strconv"
)

type server struct {
	ctx context.Context

	callbackURL string
	rootURL     string

	token chan string

	refreshToken string
	err          error

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

		s.callbackURL = fmt.Sprintf("http://localhost:%d/callback?token=true&error=true", port)
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
		// A lack of a token indicates there was an error.
		if token == "" {
			errMessage := r.URL.Query().Get("error")
			// We discard the error because if the int parse failed we return a generic error.
			statusCode, _ := strconv.Atoi(r.URL.Query().Get("status_code"))

			switch statusCode {
			case http.StatusNotFound:
				fmt.Fprint(w, "Account not found, registration required")
				s.err = errors.New("account not found")
			default:
				fmt.Fprintf(w, "%d %s", statusCode, errMessage)
				// Wrap the error message in Error so callers can distinguish server
				// errors from account not found errors.
				s.err = &Error{Details: remoteServerError, Err: errors.New(errMessage)}
			}

			close(s.token)
			return
		}
		s.refreshToken = r.URL.Query().Get("refresh_token")
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
			return "", s.err
		}

		// We've received a token, but need to block until the success page has been written to
		// the browser
		<-s.token
		return t, nil
	}
}

// RefreshToken is called after Token(), so we don't need any blocking here.
func (s *server) RefreshToken() string {
	return s.refreshToken
}

func (s *server) Close() error {
	return s.l.Close()
}
