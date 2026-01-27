/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

// Any non-whitespace character such as letters, numbers and punctuation, and a
// space. Limited in length to 512 characters, including an empty string.
const descriptionPattern = `^[\S ]{0,512}$`

var (
	descriptionPatternCompiled = regexp.MustCompile(descriptionPattern)

	// ErrInvalidDescription describes invalid names by sharing the regular
	// expression they must match
	ErrInvalidDescription = fmt.Errorf("description must match %q", descriptionPattern)
)

func ValidateDescription(description string) error {
	if descriptionPatternCompiled.MatchString(description) {
		return nil
	}
	return ErrInvalidDescription
}
