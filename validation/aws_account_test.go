/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

func TestValidateAWSAccount(t *testing.T) {
	tests := map[string]struct {
		Input  string
		Expect bool
	}{
		"valid":     {"123456789012", true},
		"empty":     {"", false},
		"letters":   {"12345678901a", false},
		"too short": {"123", false},
		"too long":  {"1234567890123", false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := ValidateAWSAccount(tt.Input)
			if (got == nil) != tt.Expect {
				t.Errorf("Expected GroupName(`%s`) to return (err == nil) == %v, but got %v", tt.Input, tt.Expect, got)
			}
		})
	}
}
