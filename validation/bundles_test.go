/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

func TestValidateBundles(t *testing.T) {
	tests := map[string]struct {
		Input  []string
		Expect bool
	}{
		"valid":                   {[]string{"boo"}, true},
		"valid colon":             {[]string{"bundle:boo", "tag:something"}, true},
		"start with colon":        {[]string{":boo"}, false},
		"end with colon":          {[]string{"boo:"}, false},
		"multiple colons":         {[]string{"boo:haha:ra"}, false},
		"empty":                   {[]string{""}, false},
		"no slashes":              {[]string{"foo/bar"}, false},
		"spaces are not alright":  {[]string{"chainguard engineering"}, false},
		"other whitespace is bad": {[]string{"\r\n\tderp"}, false},
		"no uppercase":            {[]string{"DERP"}, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := ValidateBundles(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected Bundles(`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
