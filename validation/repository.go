/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

// Repository names share the indexed MySQL column's 191-byte storage limit.
const maxRepositoryNameBytes = 191

var (
	repositoryNamePattern = regexp.MustCompile(`^[a-z0-9]+([._-][a-z0-9]+)*(/[a-z0-9]+([._-][a-z0-9]+)*)*$`)

	// ErrInvalidRepositoryName indicates an invalid repository path.
	ErrInvalidRepositoryName = fmt.Errorf("repo name must match %q", repositoryNamePattern)
)

// ValidateRepositoryName checks a relative repository path of at most 191 bytes.
// It accepts a stricter subset of OCI names: separators are a single '.', '_' or '-'.
func ValidateRepositoryName(name string) error {
	if len(name) > maxRepositoryNameBytes {
		return fmt.Errorf("%w; name must be at most %d bytes", ErrInvalidRepositoryName, maxRepositoryNameBytes)
	}
	if !repositoryNamePattern.MatchString(name) {
		return ErrInvalidRepositoryName
	}
	return nil
}
