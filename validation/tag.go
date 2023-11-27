/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

// From https://github.com/opencontainers/distribution-spec/blob/main/spec.md#pulling-manifests
const tagPattern = `^[a-zA-Z0-9_][a-zA-Z0-9._-]{0,127}$`

var (
	tagPatternCompiled = regexp.MustCompile(tagPattern)

	// ErrInvalidTag describes invalid tags by sharing the regular expression
	// they must match
	ErrInvalidTag = fmt.Errorf("tag must match %q", tagPattern)
)

func ValidateTag(tag string) error {
	if tagPatternCompiled.MatchString(tag) {
		return nil
	}
	return ErrInvalidTag
}
