//go:build !windows

/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

// isRenameRetryable is always false on unix: rename(2) over a path another
// process holds open succeeds, because the open handle keeps the old inode
// alive rather than pinning the directory entry. There is nothing to wait for.
func isRenameRetryable(error) bool { return false }
