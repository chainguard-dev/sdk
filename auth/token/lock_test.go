/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// Environment for TestHelperProcess, the child side of the two-process tests.
const (
	envHelper      = "TOKEN_LOCK_HELPER"
	envHelperAud   = "TOKEN_LOCK_HELPER_AUD"
	envHelperLog   = "TOKEN_LOCK_HELPER_LOG"
	envHelperWait  = "TOKEN_LOCK_HELPER_TIMEOUT"
	exitLockFailed = 3
)

// TestHelperProcess is not a test: the tests below run the test binary again
// as a child process, which runs helperMain.
func TestHelperProcess(t *testing.T) {
	if os.Getenv(envHelper) != "1" {
		t.Skip("helper process for the two-process lock tests")
	}
	os.Exit(helperMain())
}

// helperMain performs one Refresh. Its exchange appends the token it was
// handed to a shared log, then stands in for the issuer's round trip by
// holding the lock 200ms, long enough that the other child is certain to
// contend for it, and returns the token with a prime appended as the
// successor.
func helperMain() int {
	if d, err := time.ParseDuration(os.Getenv(envHelperWait)); err == nil {
		lockTimeout = d
	}
	err := Refresh(context.Background(), os.Getenv(envHelperAud), func(_ context.Context, tok string) (string, error) {
		f, err := os.OpenFile(os.Getenv(envHelperLog), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
		if err != nil {
			return "", err
		}
		_, err = fmt.Fprintln(f, tok)
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			return "", err
		}
		time.Sleep(200 * time.Millisecond)
		return tok + "'", nil
	}, WithIssuer(issuerA))
	switch {
	case errors.Is(err, ErrLockUnavailable):
		return exitLockFailed
	case err != nil:
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// startHelper starts a child refreshing aud and returns it unwaited.
func startHelper(t *testing.T, aud, logPath string, timeout time.Duration) *exec.Cmd {
	t.Helper()
	cmd := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestHelperProcess$") //nolint:gosec // G204: re-executing this test binary
	cmd.Env = append(os.Environ(),
		envHelper+"=1",
		envHelperAud+"="+aud,
		envHelperLog+"="+logPath,
		envHelperWait+"="+timeout.String(),
	)
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	return cmd
}

func exchangeLog(t *testing.T, path string) []string {
	t.Helper()
	b, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		t.Fatal(err)
	}
	return strings.Fields(string(b))
}

// TestRefresh_TwoProcesses races two processes refreshing one family. With
// the lock, the second loads the first's successor after the first releases
// it, so no token is ever presented twice: presenting a superseded token is
// replay, and the issuer revokes the family for it.
func TestRefresh_TwoProcesses(t *testing.T) {
	for _, tc := range []struct {
		name string
		auds [2]string
	}{
		{name: "same audience", auds: [2]string{audReg, audReg}},
		{name: "cross-audience", auds: [2]string{audAPI, audReg}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			isolateCache(t)
			presented := familyToken(t, login(issuerA, "alice", audAPI, audReg))
			mustSaveFamily(t, presented, issuerA, audAPI, audReg)
			logPath := filepath.Join(t.TempDir(), "exchanges")

			children := []*exec.Cmd{
				startHelper(t, tc.auds[0], logPath, time.Minute),
				startHelper(t, tc.auds[1], logPath, time.Minute),
			}
			for _, c := range children {
				if err := c.Wait(); err != nil {
					t.Fatalf("child: %v", err)
				}
			}

			want := []string{presented, presented + "'"}
			if diff := cmp.Diff(want, exchangeLog(t, logPath)); diff != "" {
				t.Errorf("tokens presented (-want +got):\n%s", diff)
			}
			for _, aud := range []string{audAPI, audReg} {
				if got, err := loadString(t, aud, WithIssuer(issuerA)); err != nil || got != presented+"''" {
					t.Errorf("stored for %s: got %q, %v; want the second successor", aud, got, err)
				}
			}
		})
	}
}

// TestRefresh_LockTimeout pins that a process that cannot get the lock
// exchanges nothing.
func TestRefresh_LockTimeout(t *testing.T) {
	isolateCache(t)
	presented := familyToken(t, login(issuerA, "alice", audAPI))
	mustSaveFamily(t, presented, issuerA, audAPI)
	logPath := filepath.Join(t.TempDir(), "exchanges")

	unlock, err := lockStore(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	err = startHelper(t, audAPI, logPath, 300*time.Millisecond).Wait()
	if exit, ok := errors.AsType[*exec.ExitError](err); !ok || exit.ExitCode() != exitLockFailed {
		t.Fatalf("child: %v, want exit code %d (ErrLockUnavailable)", err, exitLockFailed)
	}
	if got := exchangeLog(t, logPath); len(got) != 0 {
		t.Errorf("tokens presented: %q, want none", got)
	}
	if got, err := loadString(t, audAPI, WithIssuer(issuerA)); err != nil || got != presented {
		t.Errorf("stored: got %q, %v; want it unchanged", got, err)
	}
}
