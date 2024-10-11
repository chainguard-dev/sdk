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
		"valid":                   {[]string{"application"}, true},
		"not allowed":             {[]string{"kubernetes"}, false},
		"start with colon":        {[]string{":featured"}, false},
		"end with colon":          {[]string{"featured:"}, false},
		"multiple colons":         {[]string{"application:featured:ai"}, false},
		"empty value":             {[]string{""}, false},
		"empty":                   {[]string{}, true},
		"no slashes":              {[]string{"featured/application"}, false},
		"spaces are not alright":  {[]string{"featured application"}, false},
		"other whitespace is bad": {[]string{"\r\n\tfeatured"}, false},
		"no uppercase":            {[]string{"FEATURED"}, false},
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
