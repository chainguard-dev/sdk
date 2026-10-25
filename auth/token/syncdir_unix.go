//go:build !windows

/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import "os"

// syncDir flushes dir's own metadata, so a rename into it survives a power
// loss. Opening a directory read-only and calling Sync is the portable unix
// spelling of fsync(2) on a directory fd.
func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	err = d.Sync()
	if cerr := d.Close(); err == nil {
		err = cerr
	}
	return err
}
