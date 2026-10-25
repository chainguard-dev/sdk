/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// renameBudget bounds how long writeFileAtomic waits out a Windows sharing
// violation before giving up, and renameFirstBackoff is the first pause; each
// retry doubles it. A reader's os.ReadFile holds its handle for microseconds,
// so this is sized to ride out a handful of them rather than to wait on a
// process that is parked on the file: past the budget the caller is told the
// write failed, which under rotation is the signal to discard rather than to
// keep retrying silently.
const (
	renameBudget       = 100 * time.Millisecond
	renameFirstBackoff = 2 * time.Millisecond
)

// renameWithRetry calls rename, retrying while retryable says the failure is a
// transient sharing violation and renameBudget has time left. The rename,
// the predicate and the sleep are parameters so the backoff can be tested on
// every platform rather than only on the one that can produce the error.
func renameWithRetry(rename func() error, retryable func(error) bool, sleep func(time.Duration)) error {
	err := rename()
	for wait, spent := renameFirstBackoff, time.Duration(0); err != nil && retryable(err) && spent+wait <= renameBudget; wait *= 2 {
		sleep(wait)
		spent += wait
		err = rename()
	}
	return err
}

// tokenFilePerm is the mode every file the cache writes ends up with: readable
// and writable by the owner alone, since each one is a bearer credential.
const tokenFilePerm fs.FileMode = 0o600

// tempPrefix marks the staging files writeFileAtomic renames from. The leading
// dot keeps them out of directory listings that skip dot files, and the fixed
// prefix is what lets DeleteAll recognize one that outlived the process that
// wrote it.
const tempPrefix = ".tmp-"

// syncDirFn is syncDir, indirected so a test can make it fail: the property
// that a sync failure does not fail the write is not otherwise reachable on a
// filesystem that supports the call.
var syncDirFn = syncDir

// isTempFile reports whether name is one of writeFileAtomic's staging files.
func isTempFile(name string) bool {
	return strings.HasPrefix(name, tempPrefix)
}

// stagingPrefix is the fixed leading part of the staging files belonging to
// the token file called name. writeFileAtomic appends random text to it and
// removeStagingFiles matches on it, both through this function: if the writer
// and the cleaner disagreed about the shape of the name, a clear would report
// success while leaving a file of token bytes on disk.
func stagingPrefix(name string) string {
	return tempPrefix + name + "-"
}

// writeFileAtomic writes data to path with mode tokenFilePerm so that a reader
// holding no lock sees either the previous file or the whole new one, never a
// partial write. Every token the cache stores goes through it, not refresh tokens
// alone: the Docker credential helper hands the bytes of a just-loaded access
// token straight to Docker, so a torn access token is a bad credential rather
// than a failed read.
//
// The data goes to a staging file in the same directory, which is written,
// chmod'd, synced and renamed over path, after which the directory itself is
// synced. The staging file is removed on any failure. The directory is opened
// as an os.Root, so neither the staging file nor the rename can land outside
// it through a symlink.
//
// The mode is exact on unix. On Windows os.Chmod only toggles the read-only
// bit, so the file inherits the ACL of the user's cache directory, as the
// cache's files already do.
func writeFileAtomic(path string, data []byte) error {
	dir, name := filepath.Split(path)
	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer root.Close()

	tmp := stagingPrefix(name) + rand.Text()
	f, err := root.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	// Remove the staging file unless the rename below claims it. The flag is
	// deliberate: keying the cleanup off a named return would make it depend on
	// every later error assignment staying an assignment, and one refactor to
	// an `if err := f.Sync(); err != nil` would shadow it, leave the defer
	// looking at a nil error, and silently strand a file of token bytes in the
	// audience directory with no test failing.
	renamed := false
	defer func() {
		if !renamed {
			_ = root.Remove(tmp)
		}
	}()
	_, err = f.Write(data)
	if err == nil {
		// OpenFile's mode is filtered by the umask; Chmod sets it exactly.
		err = f.Chmod(tokenFilePerm)
	}
	if err == nil {
		err = f.Sync()
	}
	if cerr := f.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return err
	}
	// Retry a Windows sharing violation: on that platform a rename over a
	// file another process has open fails, and a parallel credential-helper
	// read must not be read as a persistent write failure. No-op elsewhere,
	// where rename over an open file succeeds.
	if err = renameWithRetry(
		func() error { return root.Rename(tmp, name) },
		isRenameRetryable, time.Sleep,
	); err != nil {
		return err
	}
	renamed = true

	// Sync the directory, not just the file. Without this a power loss after
	// the rename can restore the old directory entry, which under refresh-token
	// rotation is a token the server has already superseded: the next refresh
	// presents it and the family is revoked as a replay. The file's own Sync
	// above does not cover the rename.
	//
	// Best-effort, and deliberately so. By this point the rename has already
	// succeeded and the bytes are already synced, so every reader sees the new
	// token either way; the only thing a failed sync costs is ordering across
	// a power loss. Some filesystems cannot fsync a directory at all --
	// overlayfs has returned EINVAL for it -- and on those, surfacing the
	// error would turn "your cache is on an overlay" into "you cannot log in".
	// A weaker durability guarantee is the better failure.
	_ = syncDirFn(dir)
	return nil
}

// removeStagingFiles deletes any staging file left for the token at path by a
// writer that died between creating it and renaming it over the target.
// A missing directory is not an error: there is nothing to clean up.
//
// The directory is listed and the names compared literally rather than matched
// with filepath.Glob. dir is not a pattern the package controls -- it is built
// from os.UserCacheDir, the audience with only "/" replaced, and the
// caller-chosen WithAlias -- so a glob metacharacter anywhere in it is read as
// syntax. An audience of https://[::1]:8080 alone is enough: [::1] parses as a
// character class, matches no real directory, and the orphaned token bytes
// survive a logout that reports success.
// The directory is opened as an os.Root and the entries removed by name
// through it, for the same reason writeFileAtomic stages through one: a name
// in the listing is not trusted input, and resolving it with filepath.Join
// would follow a symlink out of the cache.
func removeStagingFiles(path string) error {
	dir, name := filepath.Split(path)
	root, err := os.OpenRoot(dir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("opening token directory %s: %w", dir, err)
	}
	defer root.Close()

	d, err := root.Open(".")
	if err != nil {
		return fmt.Errorf("opening token directory %s: %w", dir, err)
	}
	entries, err := d.ReadDir(-1)
	d.Close()
	if err != nil {
		return fmt.Errorf("reading token directory %s: %w", dir, err)
	}

	prefix := stagingPrefix(name)
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), prefix) {
			continue
		}
		if err := root.Remove(e.Name()); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("failed to remove staging file (%s): %w", e.Name(), err)
		}
	}
	return nil
}
