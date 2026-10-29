/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"context"
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
	alias           string
	issuer          string
	primary         string
	rotationCapable bool
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

// WithIssuer scopes a refresh token to the issuer that minted it. For
// KindRefresh, Load and RemainingLife then return only a token from that
// issuer, resolved through the family store (see SaveFamily), and Save and
// Delete act on that issuer's families. A caller that passes no issuer sees
// only the per-audience file, exactly as it would without the family store,
// and never a family. Other kinds ignore the option, and so do aliased
// tokens, which never belong to a family.
func WithIssuer(iss string) Option {
	return func(token *token) {
		token.issuer = auth.NormalizeIssuer(iss)
	}
}

// WithPrimaryAudience names the audience a login was requested for. SaveFamily
// always replaces that audience's per-audience copy, even a live token from
// another issuer, which it keeps in every other audience it covers.
func WithPrimaryAudience(aud string) Option {
	return func(token *token) {
		token.primary = aud
	}
}

// WithRotationCapable tells SaveFamily that the token comes from an exchange
// that declared refresh-token rotation (sts.WithRefreshTokenRotation), so the
// issuer rotates its family. SaveFamily then writes no per-audience copies of
// it and deletes the copies it would have overwritten.
//
// An issuer-less reader refreshes a copy outside the store lock. Once the
// family rotates, a copy holds a superseded token, and presenting it revokes
// the family. Without a copy, that reader logs in on its own instead.
func WithRotationCapable() Option {
	return func(token *token) {
		token.rotationCapable = true
	}
}

func newToken(opts ...Option) token {
	t := token{}
	for _, o := range opts {
		o(&t)
	}
	return t
}

// Save saves the given token to cache/audience. A KindRefresh token saved
// WithIssuer is stored as a family covering audience alone, the same as
// SaveFamily with that one audience.
func Save(token []byte, kind Kind, audience string, opts ...Option) error {
	t := newToken(opts...)
	if t.usesFamilies(kind) {
		return t.saveFamily(string(token), []string{audience})
	}
	return t.save(token, kind, audience)
}

func (t token) save(token []byte, kind Kind, audience string) error {
	path, err := t.path(kind, audience)
	if err != nil {
		return err
	}

	if err := writeFileAtomic(path, token); err != nil {
		return fmt.Errorf("writing token file: %w", err)
	}
	return nil
}

// Load returns the token for the given audience if it exists,
// or an error if it doesn't.
func Load(kind Kind, audience string, opts ...Option) ([]byte, error) {
	t := newToken(opts...)
	if t.usesFamilies(kind) {
		return t.loadSelected(audience)
	}
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
//
// For KindRefresh it also removes audience from the coverage of every family
// (of the given issuer only, when WithIssuer is passed), deleting a family
// left covering nothing.
func Delete(kind Kind, audience string, opts ...Option) error {
	t := newToken(opts...)
	if err := t.delete(kind, audience); err != nil {
		return err
	}
	if kind == KindRefresh && t.alias == "" {
		return t.uncover(audience)
	}
	return nil
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
	// Remove any staging file an atomic write left behind for this kind. It
	// holds token bytes, and a process killed between the write and the
	// rename does not clean up after itself, so a logout that skipped it
	// would report success while the credential stayed on disk. DeleteAll
	// sweeps these too; this covers the per-audience paths that reach here
	// instead -- auth logout, auth delete-account, ClearPlatformTokens.
	return removeStagingFiles(path)
}

// DeleteAll removes all Chainguard tokens and empty audience directories,
// including every refresh-token family. The store's lock file holds no token
// and is left in place.
func DeleteAll() error {
	base, err := cacheFilePath("")
	if err != nil {
		return fmt.Errorf("error locating Chainguard token dir: %w", err)
	}
	families := filepath.Join(base, familiesDir)

	// Taken so a refresh in flight cannot write its successor back after the
	// family file is removed. Removing never presents a token, so a lock that
	// cannot be had does not stop the clear.
	if _, err := os.Stat(families); err == nil {
		if unlock, err := lockStore(context.Background()); err == nil {
			defer unlock()
		}
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
		case slices.Contains(AllKinds, Kind(d.Name())), isTempFile(d.Name()), filepath.Dir(path) == families:
			// Remove recognized token files, and any staging file an atomic
			// write left behind: it holds token bytes, and a process killed
			// between the write and the rename does not get to clean up after
			// itself. A clear that left one behind would leave a credential on
			// disk that the walk below then reports as an empty directory.
			if err := os.Remove(path); err != nil { //nolint:gosec // G122: walking user's own token directory
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

// ListAudiences returns the set of audiences for which an access token
// (KindAccess) currently exists in the local token cache.
//
// Audiences are recovered from each token's "aud" claim rather than from the
// on-disk directory name, since the directory name is a lossy encoding of the
// audience (e.g. "/" is replaced with "-"). The returned slice is sorted and
// de-duplicated. Files that cannot be read or parsed as tokens are skipped.
func ListAudiences() ([]string, error) {
	base, err := cacheFilePath("")
	if err != nil {
		return nil, fmt.Errorf("error locating Chainguard token dir: %w", err)
	}

	seen := map[string]struct{}{}
	if err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// The cache dir may not exist yet; treat as no audiences.
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}
			return err
		}
		if d.IsDir() || d.Name() != string(KindAccess) {
			return nil
		}
		raw, err := os.ReadFile(path) //nolint:gosec // G304: reading user's own token cache
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		auds, err := auth.ExtractAudiences(string(raw))
		if err != nil {
			// Skip files we can't parse as a token.
			return nil //nolint:nilerr
		}
		for _, aud := range auds {
			if aud != "" {
				seen[aud] = struct{}{}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	auds := make([]string, 0, len(seen))
	for a := range seen {
		auds = append(auds, a)
	}
	sort.Strings(auds)
	return auds, nil
}

// Path is the filepath of the token for the given audience.
func Path(kind Kind, audience string, opts ...Option) (string, error) {
	t := newToken(opts...)
	return t.path(kind, audience)
}

func (t token) path(kind Kind, audience string) (string, error) {
	// NB: empty elements in Join are ignored, so we don't need to
	// check the existence of t.alias here.
	fp := filepath.Join(audienceDir(audience), t.alias, string(kind))
	return cacheFilePath(fp)
}

// audienceDir is the name of audience's directory under the cache root.
func audienceDir(audience string) string {
	a := strings.ReplaceAll(audience, "/", "-")
	// Windows does not allow : as a valid character for directory names.
	// For backwards compatibility, keep : in directory names for non-Windows systems.
	// Ref: https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-file
	if runtime.GOOS == "windows" {
		a = strings.ReplaceAll(a, ":", "-")
	}
	return a
}

// RemainingLife returns the amount of time remaining before the token for
// the given audience expires. Returns 0 for expired and non-existent tokens.
func RemainingLife(kind Kind, audience string, less time.Duration, opts ...Option) time.Duration {
	t := newToken(opts...)
	return t.remainingLife(kind, audience, less)
}

func (t token) remainingLife(kind Kind, audience string, less time.Duration) time.Duration {
	var tok []byte
	var err error
	if t.usesFamilies(kind) {
		tok, err = t.loadSelected(audience)
	} else {
		tok, err = t.load(kind, audience)
	}
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
		// Never stdout: the Docker credential helper's Get writes its JSON
		// response there, and reaches this through useRefreshToken, so a
		// token that fails to parse would corrupt the helper's output rather
		// than simply reporting zero life.
		fmt.Fprintf(os.Stderr, "failed to extract expiry: %v\n", err)
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
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			return "", fmt.Errorf("failed to create %s: %w", path, err)
		}
	}
	return path, err // err could be non-nil, if err != os.ErrNotExist
}
