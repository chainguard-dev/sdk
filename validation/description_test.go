/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"testing"
)

// TestValidateDescription makes sure the validation regexp  allows the
// expected ranges of characters.
func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect bool
	}{
		{
			name:   "empty",
			input:  "",
			expect: true,
		},
		{
			name:   "a regular description",
			input:  "this is the project that builds the future of computing",
			expect: true,
		},
		{
			name:   "newline isn't allowed",
			input:  "this is\ndescription",
			expect: false,
		},
		{
			name:   "tab isn't allowed",
			input:  "this is\tdescription",
			expect: false,
		},
		{
			name:   "carriage-return isn't allowed",
			input:  "this is\rdescription",
			expect: false,
		},
		{
			name:   "exactly 512 characters is allowed",
			input:  "This is a test sentence that needs to be exactly five hundred and twelve characters long. We are testing the validation function to ensure it correctly handles the maximum allowed length for descriptions. The description field has a limit of 512 characters as defined by the validation pattern. This sentence continues with more text to reach the exact character count needed. Adding some final words to make sure we hit precisely five hundred twelve characters in total for this test case here.",
			expect: true,
		},
		{
			name:   "513 characters is not allowed",
			input:  "This is a test sentence that needs to be exactly five hundred and thirteen characters long. We are testing the validation function to ensure it correctly rejects descriptions that exceed the maximum allowed length. The description field has a limit of 512 characters as defined by the validation pattern. This sentence continues with more text to reach the exact character count needed. Adding some final words to make sure we hit precisely five hundred thirteen characters in total for this test case here so th.",
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateDescription(tt.input)
			if (got == nil) != tt.expect {
				t.Errorf("Expected (`%s`) to return (err == nil) == %v, but got %v", tt.input, tt.expect, got)
			}
		})
	}
}
