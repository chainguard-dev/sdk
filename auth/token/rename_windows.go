//go:build windows

/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"errors"
	"syscall"
)

// errSharingViolation is ERROR_SHARING_VIOLATION. Go's syscall package
// defines ERROR_ACCESS_DENIED but not this one, and golang.org/x/sys is only
// an indirect dependency of this module, so the Win32 value is spelled out
// rather than promoting a dependency for a single constant.
const errSharingViolation = syscall.Errno(32)

// isRenameRetryable reports whether a failed rename is the transient Windows
// sharing case rather than a real error.
//
// Go opens files with FILE_SHARE_READ|FILE_SHARE_WRITE and no
// FILE_SHARE_DELETE (syscall.Open's sharemode), so MoveFileEx with
// MOVEFILE_REPLACE_EXISTING onto a file another process has open fails. A
// reader doing os.ReadFile holds that handle only for the length of the read,
// so the condition clears on its own -- which is exactly the case
// writeFileAtomic must not turn into a failed Save: two parallel Docker
// credential-helper Gets, one reading the access token while the other
// re-mints and saves it.
func isRenameRetryable(err error) bool {
	return errors.Is(err, syscall.ERROR_ACCESS_DENIED) ||
		errors.Is(err, errSharingViolation)
}
