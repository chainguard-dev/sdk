//go:build windows

/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

// syncDir is a deliberate no-op on Windows, which exposes no way to flush a
// directory's metadata: opening a directory with CreateFile requires
// FILE_FLAG_BACKUP_SEMANTICS and FlushFileBuffers does not accept the handle.
// Unlike the file lock, where a no-op would leave Windows with the bug the
// lock exists to fix, there is nothing here to call and nothing to emulate.
// Do not "fix" this by syncing the parent or by ignoring an error from Sync.
func syncDir(string) error { return nil }
