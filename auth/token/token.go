/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sort"
	"strings"
	"time"

	"chainguard.dev/sdk/auth"
)

type Kind string

const (
	KindAccess  Kind = "oidc-token"
	KindRefresh Kind = "refresh-token"
)

var (
	AllKinds = []Kind{KindAccess, KindRefresh}

	parentDir = "chainguard"
)

type token struct {
	alias string
}

type Option func(*token)

// WithAlias allows callers to organize tokens into subdirectories
// under an audience to manage multiple tokens for the same audience
// without overwriting.
func WithAlias(a string) Option {
	return func(token *token) {
		token.alias = a
	}
}

func newToken(opts ...Option) token {
	t := token{}
	for _, o := range opts {
		o(&t)
	}
	return t
}

// Save saves the given token to cache/audience
func Save(token []byte, kind Kind, audience string, opts ...Option) error {
	t := newToken(opts...)
	return t.save(token, kind, audience)
}

func (t token) save(token []byte, kind Kind, audience string) error {
	path, err := t.path(kind, audience)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, token, 0600); err != nil {
		return fmt.Errorf("writing token file: %w", err)
	}
	return nil
}

// Load returns the token for the given audience if it exists,
// or an error if it doesn't.
func Load(kind Kind, audience string, opts ...Option) ([]byte, error) {
	t := newToken(opts...)
	return t.load(kind, audience)
}

func (t token) load(kind Kind, audience string) ([]byte, error) {
	path, err := t.path(kind, audience)
	if err != nil {
		return nil, err
	}

	bs, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading token file: %w", err)
	}
	return bs, nil
}

// Delete removes the token for the given audience, if it exists.
// No error is returned if the token doesn't exist.
func Delete(kind Kind, audience string, opts ...Option) error {
	t := newToken(opts...)
	return t.delete(kind, audience)
}

func (t token) delete(kind Kind, audience string) error {
	path, err := t.path(kind, audience)
	if err != nil {
		return err
	}
	err = os.Remove(path)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("failed to remove token (%s): %w", path, err)
	}
	return nil
}

// DeleteAll removes all Chainguard tokens and empty audience directories.
func DeleteAll() error {
	base, err := cacheFilePath("")
	if err != nil {
		return fmt.Errorf("error locating Chainguard token dir: %w", err)
	}

	// Token directory is expected to be structured as group of audience-specific
	// directories, with token files (oidc-token, refresh-token) or alias directories
	// nested within.
	//
	//  $ tree ~/Library/Caches/chainguard
	// /Users/foo/Library/Caches/chainguard
	// ├── https:--console-api.enforce.dev
	// │   ├── oidc-token
	// │   ├── refresh-token
	// │   └── foo
	// │       ├── oidc-token
	// │       └── refresh-token
	// ├── cgr.dev
	//     └── oidc-token
	var dirs []string
	if err = filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		// Return errors encountered reading the base directory.
		if err != nil {
			return err
		}

		switch {
		case path == base:
			// Skip the base directory, we don't want to remove it
		case d.IsDir():
			// Keep track of directories we'll want to delete, if they end up empty.
			dirs = append(dirs, path)
		case slices.Contains(AllKinds, Kind(d.Name())):
			// Remove recognized token files.
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("removing file %s: %w", path, err)
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Sort directories from longest to shortest to remove nested dirs first.
	sort.Slice(dirs, func(i, j int) bool {
		return len(dirs[i]) > len(dirs[j])
	})

	// Remove empty directories
	for _, d := range dirs {
		ents, err := os.ReadDir(d)
		if err != nil {
			return fmt.Errorf("reading %s: %w", d, err)
		}
		// Remove the directory if it is empty
		if len(ents) == 0 {
			if err := os.Remove(d); err != nil {
				return fmt.Errorf("removing directory %s: %w", d, err)
			}
		}
	}

	return nil
}

// Path is the filepath of the token for the given audience.
func Path(kind Kind, audience string, opts ...Option) (string, error) {
	t := newToken(opts...)
	return t.path(kind, audience)
}

func (t token) path(kind Kind, audience string) (string, error) {
	a := strings.ReplaceAll(audience, "/", "-")
	// Windows does not allow : as a valid character for directory names.
	// For backwards compatibility, keep : in directory names for non-Windows systems.
	// Ref: https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	if runtime.GOOS == "windows" {
		a = strings.ReplaceAll(a, ":", "-")
	}
	// NB: empty elements in Join are ignored, so we don't need to
	// check the existence of t.alias here.
	fp := filepath.Join(a, t.alias, string(kind))
	return cacheFilePath(fp)
}

// RemainingLife returns the amount of time remaining before the token for
// the given audience expires. Returns 0 for expired and non-existent tokens.
func RemainingLife(kind Kind, audience string, less time.Duration, opts ...Option) time.Duration {
	t := newToken(opts...)
	return t.remainingLife(kind, audience, less)
}

func (t token) remainingLife(kind Kind, audience string, less time.Duration) time.Duration {
	tok, err := t.load(kind, audience)
	if err != nil {
		// Not a big deal, life is zero.
		return 0
	}
	var expiry time.Time
	switch kind {
	case KindRefresh:
		expiry, err = auth.ExtractRefreshExpiry(string(tok))
	default:
		expiry, err = auth.ExtractExpiry(string(tok))
	}
	if err != nil {
		fmt.Printf("failed to extract expiry: %v\n", err)
		// Not a big deal, life is zero.
		return 0
	}
	return subtractOrZero(expiry, less)
}

// For testing.
var timeUntil = time.Until

// Safe calculation for duration remaining from a given time, less the given duration.
func subtractOrZero(expiry time.Time, less time.Duration) time.Duration {
	life := timeUntil(expiry.Add(less * -1))
	return max(0, life)
}

func cacheFilePath(file string) (string, error) {
	path, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user cache dir: %w", err)
	}

	// Create the cache directories if needed
	path = filepath.Join(path, parentDir, file)
	_, err = os.Stat(filepath.Dir(path))
	if errors.Is(err, os.ErrNotExist) {
		err = nil // Clear err, we're dealing with it.
		if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return "", fmt.Errorf("failed to create %s: %w", path, err)
		}
	}
	return path, err // err could be non-nil, if err != os.ErrNotExist
}
