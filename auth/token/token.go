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
)

// DefaultKind is the default token kind to use when no kind is specified.

var (
	parentDir = "chainguard"
)

// Save saves the given token to cache/audience
func Save(token []byte, kind Kind, audience string) error {
	path, err := Path(kind, audience)
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
func Load(kind Kind, audience string) ([]byte, error) {
	path, err := Path(kind, audience)
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
func Delete(kind Kind, audience string) error {
	path, err := Path(kind, audience)
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
	files, err := os.ReadDir(base)
	if err != nil {
		return fmt.Errorf("error reading Chainguard token dir: %w", err)
	}
	// Token directory is expected to be structured as group of audience-specific
	// directories, with a single file containing the token
	//
	//  $ tree ~/Library/Caches/chainguard
	// /Users/foo/Library/Caches/chainguard
	// ├── https:--console-api.enforce.dev
	// │   └── oidc-token
	// ├── https:--cgr.dev
	//    └── oidc-token
	for _, file := range files {
		if !file.IsDir() {
			// Encountered a file in the directory. Skip.
			continue
		}
		for _, kind := range AllKinds {
			// Try to remove a token, ignore file not exist errors
			tokenFile := filepath.Join(base, file.Name(), string(kind))
			if err := os.Remove(tokenFile); err != nil && !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("failed to remove %s: %w", tokenFile, err)
			}
		}
		// Remove the (hopefully empty) audience directory.
		// Ignore failures since other tools may have stored files in this cache.
		dir := filepath.Join(base, file.Name())
		_ = os.Remove(dir)
	}
	return nil
}

// Path is the filepath of the token for the given audience.
func Path(kind Kind, audience string) (string, error) {
	a := strings.ReplaceAll(audience, "/", "-")
	// Windows does not allow : as a valid character for directory names.
	// For backwards compatibility, keep : in directory names for non-Windows systems.
	// Ref: https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	if runtime.GOOS == "windows" {
		a = strings.ReplaceAll(a, ":", "-")
	}
	fp := filepath.Join(a, string(kind))
	return cacheFilePath(fp)
}

// RemainingLife returns the amount of time remaining before the token for
// the given audience expires. Returns 0 for expired and non-existent tokens.
func RemainingLife(kind Kind, audience string, less time.Duration) time.Duration {
	tok, err := Load(kind, audience)
	if err != nil {
		// Not a big deal, life is zero.
		return 0
	}
	expiry, err := auth.ExtractExpiry(string(tok))
	if err != nil {
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
	if life < 0 {
		return 0
	}
	return life
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
