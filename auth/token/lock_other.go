//go:build !windows && (!unix || aix)

/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"errors"
	"os"
)

// tryLock fails on platforms with no flock or LockFileEx, so Refresh there
// returns ErrLockUnavailable rather than exchanging unlocked.
func tryLock(*os.File) (bool, error) { return false, errors.ErrUnsupported }

func unlockFile(*os.File) {}
