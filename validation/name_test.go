/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := map[string]struct {
		Input  string
		Expect bool
	}{
		"valid":                     {"boo.foo_-", true},
		"empty":                     {"", false},
		"no slashes":                {"foo/bar", false},
		"spaces are alright":        {"chainguard engineering", true},
		"no spaces are alright too": {"chainguarddemo", true},
		"other whitespace is bad":   {"\r\n\tderp", false},
		"no uppercase":              {"DERP", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := ValidateName(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected GroupName(`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
