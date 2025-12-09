/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

func TestValidateAliases(t *testing.T) {
	tests := map[string]struct {
		Input  []string
		Expect bool
	}{
		"valid":                       {[]string{"image:latest"}, true},
		"valid with underscore":       {[]string{"my_image:1.0.0"}, true},
		"valid with dash":             {[]string{"my-image:1.0.0"}, true},
		"valid with slash":            {[]string{"myrepo/myimage:1.0.0"}, true},
		"valid with dot":              {[]string{"my.image:1.0.0"}, true},
		"valid with multiple slashes": {[]string{"gcr.io/myrepo/myimage:1.0.0"}, true},
		"invalid missing tag":         {[]string{"myimage"}, false},
		"invalid empty string":        {[]string{""}, false},
		"invalid special chars":       {[]string{"myimage@sha256:12345"}, false},
		"invalid space":               {[]string{"my image:latest"}, false},
		"invalid multiple colons":     {[]string{"myrepo:myimage:latest"}, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := ValidateAliases(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected Bundles(`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
