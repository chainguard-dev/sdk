/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chainguard-dev/clog"
)

// lockTimeout bounds the wait for the store lock. It must outlast a healthy
// holder's exchange, which chainctl bounds at 30 seconds: a waiter that gave
// up sooner would fail while the holder was about to finish. The lock is
// released when its holder exits, so a timeout means a live holder.
var lockTimeout = 45 * time.Second

const (
	lockFirstBackoff = 25 * time.Millisecond
	lockMaxBackoff   = 250 * time.Millisecond
)

// lockStore takes the store-wide lock, polling until lockTimeout passes or
// ctx ends, and returns the function that releases it. Every failure wraps
// ErrLockUnavailable, including a filesystem that cannot lock at all: going
// ahead unlocked would let two processes present the same token, and the
// second presentation revokes the family.
func lockStore(ctx context.Context) (func(), error) {
	path, err := cacheFilePath(lockFileName)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLockUnavailable, err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, tokenFilePerm) //nolint:gosec // G304: the user's own token cache
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLockUnavailable, err)
	}
	deadline := time.Now().Add(lockTimeout)
	for wait := lockFirstBackoff; ; wait = min(2*wait, lockMaxBackoff) {
		held, err := tryLock(f)
		if err != nil {
			f.Close()
			clog.WarnContextf(ctx, "cannot lock the refresh-token store %s: %v", path, err)
			return nil, fmt.Errorf("%w: %w", ErrLockUnavailable, err)
		}
		if held {
			return func() {
				unlockFile(f)
				f.Close()
			}, nil
		}
		left := time.Until(deadline)
		if left <= 0 {
			f.Close()
			clog.WarnContextf(ctx, "gave up waiting %s for the refresh-token store lock %s", lockTimeout, path)
			return nil, fmt.Errorf("%w: held by another process for %s", ErrLockUnavailable, lockTimeout)
		}
		select {
		case <-ctx.Done():
			f.Close()
			return nil, fmt.Errorf("%w: %w", ErrLockUnavailable, ctx.Err())
		case <-time.After(min(wait, left)):
		}
	}
}
