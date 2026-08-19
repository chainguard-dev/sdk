/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// fakeChainctl writes a shell script named chainctl into a fresh directory and
// makes that directory the whole PATH, so chainctlTokenSource runs the script.
// An empty body leaves the directory without a chainctl at all, which is how an
// operator with no chainctl installed reaches the same code.
func fakeChainctl(t *testing.T, body string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake chainctl is a shell script")
	}
	dir := t.TempDir()
	if body != "" {
		if err := os.WriteFile(filepath.Join(dir, "chainctl"), []byte("#!/bin/sh\n"+body+"\n"), 0o700); err != nil {
			t.Fatalf("writing the fake chainctl: %v", err)
		}
	}
	t.Setenv("PATH", dir)
}

func TestChainctlTokenSourceErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
		opts []Option
		// want lists substrings the error must contain.
		want []string
		// avoid lists substrings the error must not contain.
		avoid []string
	}{{
		name: "stderr from a failed chainctl reaches the caller",
		body: "echo 'not logged in: run chainctl auth login' 1>&2; exit 1",
		want: []string{"chainctl auth token", "exit status 1", "not logged in: run chainctl auth login"},
	}, {
		name: "the failing command reports the issuer and audience it was given",
		body: "exit 1",
		opts: []Option{WithIssuer("https://issuer.chainops.dev"), WithAudience("rebuilder-api")},
		want: []string{"chainctl auth token --issuer https://issuer.chainops.dev --audience rebuilder-api", "exit status 1"},
	}, {
		name: "a silent failure still names the command",
		body: "exit 3",
		want: []string{"chainctl auth token", "exit status 3"},
	}, {
		name: "a missing chainctl names the command",
		body: "",
		want: []string{"chainctl auth token"},
	}, {
		name:  "control bytes in stderr are escaped",
		body:  "printf '\\033[31mboom\\033[0m' 1>&2; exit 1",
		want:  []string{"boom", `\x1b`},
		avoid: []string{"\x1b"},
	}, {
		name:  "unparsable stdout is reported without echoing what chainctl printed",
		body:  "echo shhh-this-could-be-a-token",
		want:  []string{"chainctl auth token", "expiry"},
		avoid: []string{"shhh-this-could-be-a-token"},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			fakeChainctl(t, tc.body)

			tok, err := NewChainctlTokenSource(t.Context(), tc.opts...).Token()
			if err == nil {
				t.Fatalf("Token(): got = %v, want an error", tok)
			}
			for _, want := range tc.want {
				if !strings.Contains(err.Error(), want) {
					t.Errorf("Token() error: got = %q, want it to contain %q", err, want)
				}
			}
			for _, avoid := range tc.avoid {
				if strings.Contains(err.Error(), avoid) {
					t.Errorf("Token() error: got = %q, want it to not contain %q", err, avoid)
				}
			}
		})
	}
}

func TestChainctlTokenSourceTruncatesStderr(t *testing.T) {
	// echo is a shell builtin, so this needs nothing on the emptied PATH.
	fakeChainctl(t, "echo "+strings.Repeat("x", 9000)+" 1>&2; exit 1")

	_, err := NewChainctlTokenSource(t.Context()).Token()
	if err == nil {
		t.Fatal("Token(): got = nil, want an error")
	}
	// The script writes 9000 bytes of stderr. Allow room for the command and the
	// exit status around the capped stderr, but nothing near the whole 9000.
	if got, limit := len(err.Error()), maxStderrBytes+200; got > limit {
		t.Errorf("error length: got = %d, want <= %d", got, limit)
	}
	if !strings.Contains(err.Error(), "(truncated)") {
		t.Errorf("Token() error: got = %q, want it to say the stderr was truncated", err)
	}
}

func TestChainctlTokenSourceSucceeds(t *testing.T) {
	jwt := makeJWT(`{"exp":2000000000}`)
	fakeChainctl(t, "printf %s "+jwt)

	tok, err := NewChainctlTokenSource(t.Context()).Token()
	if err != nil {
		t.Fatalf("Token(): got error = %v, want nil", err)
	}
	if tok.AccessToken != jwt {
		t.Errorf("AccessToken: got = %q, want = %q", tok.AccessToken, jwt)
	}
	if got := tok.Expiry.Unix(); got != 2000000000 {
		t.Errorf("Expiry: got = %d, want = %d", got, 2000000000)
	}
}
