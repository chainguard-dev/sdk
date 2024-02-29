/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

func TestValidateTag(t *testing.T) {
	tests := map[string]struct {
		Input  string
		Expect bool
	}{
		"valid":                   {"latest", true},
		"version":                 {"1.2.3-r1", true},
		"uppercase is ok":         {"DERP", true},
		"128 chars is ok":         {"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", true},
		"129 chars is not ok":     {"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", false},
		"empty":                   {"", false},
		"no slashes":              {"foo/bar", false},
		"no colons":               {"foo:bar", false},
		"no spaces in middle":     {"chainguard engineering", false},
		"no leading spaces":       {" chainguardengineering", false},
		"no spaces both sides":    {" chainguardengineering ", false},
		"no trailing spacess":     {"chainguardengineering ", false},
		"other whitespace is bad": {"\r\n\tderp", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := ValidateTag(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected ValidateTag(`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
