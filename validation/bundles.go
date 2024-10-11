/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"fmt"
	"regexp"
)

const bundleAllowList = `^application$|^base$|^byol$|^ai$|^ai-gpu$|^featured$|^fips$`

var (
	bundleAllowListCompiled = regexp.MustCompile(bundleAllowList)

	// ErrInvalidEntry flags keywords that are not in the allow list
	ErrInvalidEntry = fmt.Errorf("only the following keywords are valid %q", bundleAllowList)
)

func ValidateBundles(bundles []string) error {
	for _, bundle := range bundles {
		if !bundleAllowListCompiled.MatchString(bundle) {
			return ErrInvalidEntry
		}
	}
	return nil
}
