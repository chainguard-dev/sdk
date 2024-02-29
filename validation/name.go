/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

const namePattern = `^[a-z0-9 ._-]{1,}$`

var (
	namePatternCompiled = regexp.MustCompile(namePattern)

	// ErrInvalidName describes invalid names by sharing the regular expression
	// they must match
	ErrInvalidName = fmt.Errorf("name must match %q", namePattern)
)

func ValidateName(name string) error {
	if namePatternCompiled.MatchString(name) {
		return nil
	}
	return ErrInvalidName
}
