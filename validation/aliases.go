/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

const aliasPattern = `^[a-zA-Z0-9_.\-/]*[a-zA-Z0-9_.\-]:[a-zA-Z0-9_.\-]+$`

var (
	aliasPatternCompiled = regexp.MustCompile(aliasPattern)

	// ErrInvalidAlias describes invalid alias(es) by sharing the regular expression
	// they must match
	ErrInvalidAlias = fmt.Errorf("each alias must match %q", aliasPattern)
)

func ValidateAliases(aliases []string) error {
	for _, alias := range aliases {
		if !aliasPatternCompiled.MatchString(alias) {
			return ErrInvalidAlias
		}
	}
	return nil
}
