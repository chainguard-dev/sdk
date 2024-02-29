/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

const bundlePattern = `^[a-z]+(:[a-z]+)?$`

var (
	bundlePatternCompiled = regexp.MustCompile(bundlePattern)

	// ErrInvalidBundles describes invalid bundle(s) by sharing the regular expression
	// they must match
	ErrInvalidBundles = fmt.Errorf("each bundle item must match %q", bundlePattern)
)

func ValidateBundles(bundles []string) error {
	for _, bundle := range bundles {
		if !bundlePatternCompiled.MatchString(bundle) {
			return ErrInvalidBundles
		}
	}
	return nil
}
