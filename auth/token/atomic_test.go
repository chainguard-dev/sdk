/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

// TestWriteFileAtomic_ReaderNeverSeesPartialWrite is the property the Docker
// credential helper depends on: a reader holding no lock observes the previous
// bytes or the new ones, never a prefix. os.WriteFile truncates first, so this
// test fails against it.
func TestWriteFileAtomic_ReaderNeverSeesPartialWrite(t *testing.T) {
	if runtime.GOOS == "windows" {
		// The reader below keeps a handle open in a tight loop. Go opens
		// files with no FILE_SHARE_DELETE, so on Windows the rename can be
		// starved past its retry budget and Save would fail rather than tear
		// -- which is the trade this platform makes, not a torn read. The
		// retry itself is covered by TestRenameWithRetry_* on every platform.
		t.Skip("rename cannot replace a file held open on Windows; see isRenameRetryable")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "oidc-token")

	old := bytes.Repeat([]byte("o"), 64*1024)
	if err := writeFileAtomic(path, old); err != nil {
		t.Fatalf("seeding: %v", err)
	}
	newTok := bytes.Repeat([]byte("n"), 64*1024)

	var wg sync.WaitGroup
	stop := make(chan struct{})
	bad := make(chan string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-stop:
				return
			default:
			}
			got, err := os.ReadFile(path) //nolint:gosec // G304: test's own temp dir
			if err != nil {
				// A rename is never observable as a missing file, so this is
				// a real failure rather than a tolerable race.
				select {
				case bad <- "read failed: " + err.Error():
				default:
				}
				return
			}
			if !bytes.Equal(got, old) && !bytes.Equal(got, newTok) {
				select {
				case bad <- "torn read: saw a file that is neither the old nor the new token":
				default:
				}
				return
			}
		}
	}()

	// Registered before the write loop so a t.Fatalf below still joins the
	// reader: otherwise it keeps calling os.ReadFile on a path in a TempDir
	// the harness is removing, past the end of the test.
	stopReader := sync.OnceFunc(func() {
		close(stop)
		wg.Wait()
	})
	t.Cleanup(stopReader)

	for range 50 {
		if err := writeFileAtomic(path, newTok); err != nil {
			t.Fatalf("writeFileAtomic() = %v", err)
		}
		if err := writeFileAtomic(path, old); err != nil {
			t.Fatalf("writeFileAtomic() = %v", err)
		}
	}
	stopReader()

	select {
	case msg := <-bad:
		t.Error(msg)
	default:
	}
}

// TestWriteFileAtomic_LeavesNoStagingFile pins that the success path cleans up
// after itself: a staging file that survived would hold token bytes.
func TestWriteFileAtomic_LeavesNoStagingFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "refresh-token")
	if err := writeFileAtomic(path, []byte("tok")); err != nil {
		t.Fatalf("writeFileAtomic() = %v", err)
	}
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() = %v", err)
	}
	for _, e := range ents {
		if isTempFile(e.Name()) {
			t.Errorf("staging file %q survived a successful write", e.Name())
		}
	}
	if len(ents) != 1 {
		t.Errorf("entries = %d, want 1 (the token itself)", len(ents))
	}
}

// TestWriteFileAtomic_Perm checks the mode is exact rather than umask-filtered.
// Windows is excluded because Chmod there only toggles the read-only bit and
// the file inherits the cache directory's ACL, which the doc says this design
// claims nothing stronger about.
func TestWriteFileAtomic_Perm(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Chmod only toggles the read-only bit on Windows")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "refresh-token")
	if err := writeFileAtomic(path, []byte("tok")); err != nil {
		t.Fatalf("writeFileAtomic() = %v", err)
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() = %v", err)
	}
	if got := fi.Mode().Perm(); got != 0600 {
		t.Errorf("mode = %o, want 600", got)
	}
}

// TestWriteFileAtomic_OverwritesExisting covers the ordinary rotation case:
// a successor replacing a parent in place.
func TestWriteFileAtomic_OverwritesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "refresh-token")
	if err := writeFileAtomic(path, []byte("parent")); err != nil {
		t.Fatalf("seeding: %v", err)
	}
	if err := writeFileAtomic(path, []byte("successor")); err != nil {
		t.Fatalf("writeFileAtomic() = %v", err)
	}
	got, err := os.ReadFile(path) //nolint:gosec // G304: test's own temp dir
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(got) != "successor" {
		t.Errorf("contents = %q, want %q", got, "successor")
	}
}

// TestDeleteAll_RemovesStagingFiles pins the clearing half: a staging file a
// killed process left behind holds token bytes, so a clear that skipped it
// would report success while leaving a credential on disk.
func TestDeleteAll_RemovesStagingFiles(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", base)
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skip("os.UserCacheDir does not honor XDG_CACHE_HOME here")
	}

	audDir := filepath.Join(base, parentDir, "cgr.dev")
	if err := os.MkdirAll(audDir, 0700); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	orphan := filepath.Join(audDir, tempPrefix+"refresh-token-abc123")
	if err := os.WriteFile(orphan, []byte("a token nobody renamed"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := DeleteAll(); err != nil {
		t.Fatalf("DeleteAll() = %v", err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Errorf("staging file survived DeleteAll: Stat() = %v, want not-exist", err)
	}
}

// TestIsTempFile keeps the recognizer and the writer agreeing: DeleteAll finds
// a staging file only if it matches what writeFileAtomic produces.
func TestIsTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "refresh-token")
	if err := writeFileAtomic(path, []byte("tok")); err != nil {
		t.Fatalf("writeFileAtomic() = %v", err)
	}
	// Reconstruct the shape writeFileAtomic uses, and assert the recognizer
	// accepts it and rejects the real token names around it.
	if !isTempFile(tempPrefix + "refresh-token-XYZ") {
		t.Error("isTempFile rejected a name writeFileAtomic can produce")
	}
	for _, k := range AllKinds {
		if isTempFile(string(k)) {
			t.Errorf("isTempFile accepted the real token name %q", k)
		}
	}
	if strings.HasPrefix(string(KindRefresh), tempPrefix) {
		t.Error("tempPrefix collides with a token kind")
	}
}

// TestRenameWithRetry_RetriesTransientThenSucceeds covers the Windows sharing
// case on every platform by injecting the error: a reader's handle clears on
// its own, so a rename that fails twice and then works must be reported as a
// success rather than as a failed Save.
func TestRenameWithRetry_RetriesTransientThenSucceeds(t *testing.T) {
	transient := errors.New("sharing violation")
	calls, slept := 0, time.Duration(0)
	err := renameWithRetry(
		func() error {
			calls++
			if calls <= 2 {
				return transient
			}
			return nil
		},
		func(err error) bool { return errors.Is(err, transient) },
		func(d time.Duration) { slept += d },
	)
	if err != nil {
		t.Errorf("renameWithRetry() = %v, want nil", err)
	}
	if calls != 3 {
		t.Errorf("rename calls = %d, want 3", calls)
	}
	if slept != 6*time.Millisecond {
		t.Errorf("slept = %v, want 6ms (2ms then 4ms)", slept)
	}
}

// TestRenameWithRetry_GivesUpAtBudget pins that the wait is bounded. Past the
// budget the caller must be told the write failed: under rotation that is the
// signal to discard the token, so retrying forever would hide it.
func TestRenameWithRetry_GivesUpAtBudget(t *testing.T) {
	transient := errors.New("sharing violation")
	calls, slept := 0, time.Duration(0)
	err := renameWithRetry(
		func() error { calls++; return transient },
		func(err error) bool { return errors.Is(err, transient) },
		func(d time.Duration) { slept += d },
	)
	if !errors.Is(err, transient) {
		t.Errorf("renameWithRetry() = %v, want the transient error", err)
	}
	if slept > 100*time.Millisecond {
		t.Errorf("slept = %v, want <= the 100ms budget", slept)
	}
	if calls < 2 {
		t.Errorf("rename calls = %d, want more than one attempt", calls)
	}
}

// TestRenameWithRetry_DoesNotRetryRealErrors keeps the backoff from masking a
// genuine failure: only the transient sharing case is worth waiting out.
func TestRenameWithRetry_DoesNotRetryRealErrors(t *testing.T) {
	permanent := errors.New("no such file or directory")
	calls := 0
	err := renameWithRetry(
		func() error { calls++; return permanent },
		func(error) bool { return false },
		func(time.Duration) { t.Error("slept on a non-retryable error") },
	)
	if !errors.Is(err, permanent) {
		t.Errorf("renameWithRetry() = %v, want the original error", err)
	}
	if calls != 1 {
		t.Errorf("rename calls = %d, want 1", calls)
	}
}

// TestIsRenameRetryable_RejectsUnrelated guards the platform predicate against
// treating an ordinary failure as transient.
func TestIsRenameRetryable_RejectsUnrelated(t *testing.T) {
	if isRenameRetryable(errors.New("disk full")) {
		t.Error("isRenameRetryable accepted an unrelated error")
	}
	if isRenameRetryable(os.ErrNotExist) {
		t.Error("isRenameRetryable accepted ErrNotExist")
	}
}

// TestDelete_RemovesStagingFiles pins the per-audience clearing path, which
// is what chainctl auth logout, auth delete-account and ClearPlatformTokens
// reach. DeleteAll sweeping them is not enough: a staging file holds token
// bytes, so a logout that left one behind would report success with the
// credential still on disk.
func TestDelete_RemovesStagingFiles(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", base)
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skip("os.UserCacheDir does not honor XDG_CACHE_HOME here")
	}

	path, err := Path(KindRefresh, "cgr.dev")
	if err != nil {
		t.Fatalf("Path() = %v", err)
	}
	if err := Save([]byte("tok"), KindRefresh, "cgr.dev"); err != nil {
		t.Fatalf("Save() = %v", err)
	}
	dir, name := filepath.Split(path)
	orphan := filepath.Join(dir, tempPrefix+name+"-deadbeef")
	if err := os.WriteFile(orphan, []byte("a token nobody renamed"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := Delete(KindRefresh, "cgr.dev"); err != nil {
		t.Fatalf("Delete() = %v", err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Errorf("staging file survived Delete: Stat() = %v, want not-exist", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("token survived Delete: Stat() = %v, want not-exist", err)
	}
}

// TestDelete_LeavesOtherKindsStagingFiles keeps the prefix match from
// over-reaching: the staging names of every kind share the ".tmp-" prefix and
// sit in one directory, so deleting the refresh token must not take an access
// token's staging file with it.
func TestDelete_LeavesOtherKindsStagingFiles(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", base)
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skip("os.UserCacheDir does not honor XDG_CACHE_HOME here")
	}

	refreshPath, err := Path(KindRefresh, "cgr.dev")
	if err != nil {
		t.Fatalf("Path() = %v", err)
	}
	dir, _ := filepath.Split(refreshPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		t.Fatalf("MkdirAll() = %v", err)
	}
	otherKind := filepath.Join(dir, tempPrefix+string(KindAccess)+"-deadbeef")
	if err := os.WriteFile(otherKind, []byte("access staging"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := Delete(KindRefresh, "cgr.dev"); err != nil {
		t.Fatalf("Delete() = %v", err)
	}
	if _, err := os.Stat(otherKind); err != nil {
		t.Errorf("Delete(KindRefresh) removed an access-token staging file: %v", err)
	}
}

// TestWriteFileAtomic_DirSyncFailureDoesNotFailTheWrite pins the behaviour
// that matters on a filesystem which cannot fsync a directory -- overlayfs
// has returned EINVAL for it. By the time
// the sync runs the rename has already succeeded and the bytes are already
// synced, so surfacing the error would turn "your cache is on an overlay"
// into "you cannot log in" while costing only power-loss ordering.
func TestWriteFileAtomic_DirSyncFailureDoesNotFailTheWrite(t *testing.T) {
	orig := syncDirFn
	t.Cleanup(func() { syncDirFn = orig })
	called := false
	syncDirFn = func(string) error {
		called = true
		return errors.New("invalid argument")
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "refresh-token")
	if err := writeFileAtomic(path, []byte("successor")); err != nil {
		t.Errorf("writeFileAtomic() = %v, want nil despite the directory sync failing", err)
	}
	if !called {
		t.Error("syncDir was not attempted")
	}
	got, err := os.ReadFile(path) //nolint:gosec // G304: test's own temp dir
	if err != nil {
		t.Fatalf("ReadFile() = %v", err)
	}
	if string(got) != "successor" {
		t.Errorf("contents = %q, want %q", got, "successor")
	}
}

// TestDelete_RemovesStagingFilesWithGlobMetacharacters pins that the staging
// sweep matches names literally. The directory holding a token is built from
// the user's cache dir, the audience with only "/" replaced, and the alias, so
// it can contain glob syntax that the package never chose: an IPv6 audience is
// enough, since "[::1]" is a character class. Matched as a pattern it selects
// nothing, and the orphaned staging file -- which holds token bytes -- outlives
// a logout that reports success.
func TestDelete_RemovesStagingFilesWithGlobMetacharacters(t *testing.T) {
	base := t.TempDir()
	t.Setenv("XDG_CACHE_HOME", base)
	if runtime.GOOS == "darwin" || runtime.GOOS == "windows" {
		t.Skip("os.UserCacheDir does not honor XDG_CACHE_HOME here")
	}

	const audience = "https://[::1]:8080"
	path, err := Path(KindRefresh, audience)
	if err != nil {
		t.Fatalf("Path() = %v", err)
	}
	if err := Save([]byte("tok"), KindRefresh, audience); err != nil {
		t.Fatalf("Save() = %v", err)
	}
	dir, name := filepath.Split(path)
	orphan := filepath.Join(dir, stagingPrefix(name)+"DEADBEEF")
	if err := os.WriteFile(orphan, []byte("a token nobody renamed"), 0600); err != nil {
		t.Fatalf("WriteFile() = %v", err)
	}

	if err := Delete(KindRefresh, audience); err != nil {
		t.Fatalf("Delete() = %v", err)
	}
	if _, err := os.Stat(orphan); !os.IsNotExist(err) {
		t.Errorf("staging file survived Delete under a glob-metacharacter audience: Stat() = %v, want not-exist", err)
	}
}

// TestWriteFileAtomic_RemovesStagingFileWhenRenameFails covers the deferred
// cleanup, which is the only thing standing between a failed write and a file
// of token bytes left beside the target. A directory at the target path makes
// the rename fail without making the write fail.
func TestWriteFileAtomic_RemovesStagingFileWhenRenameFails(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "refresh-token")
	if err := os.Mkdir(path, 0700); err != nil {
		t.Fatalf("Mkdir() = %v", err)
	}

	if err := writeFileAtomic(path, []byte("successor")); err == nil {
		t.Fatal("writeFileAtomic() = nil, want an error renaming over a directory")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() = %v", err)
	}
	for _, e := range entries {
		if isTempFile(e.Name()) {
			t.Errorf("a failed write left the staging file %q behind", e.Name())
		}
	}
}
