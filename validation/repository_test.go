/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package validation

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateRepositoryName(t *testing.T) {
	for _, tt := range []struct {
		name  string
		valid bool
	}{
		{name: "a", valid: true},
		{name: "team/api", valid: true},
		{name: "team-a/my.api_v2", valid: true},
		{name: "a.b_c-d", valid: true},
		{name: strings.Repeat("a", 191), valid: true},
		{name: strings.Repeat("a", 189) + "/a", valid: true},
		{name: strings.Repeat("a", 192)},
		{name: strings.Repeat("a", 190) + "/a"},
		{name: strings.Repeat("a", 256)},
		{name: "-a"},
		{name: "_a"},
		{name: ".a"},
		{name: "a-"},
		{name: "a_"},
		{name: "a."},
		{name: "a..b"},
		{name: "a._b"},
		{name: "a__b"},
		{name: "a--b"},
		{name: "a---b"},
		{name: "team/a__b"},
		{name: "team/a--b"},
		{name: "a___b"},
		{name: "a\n"},
		{name: "a\x00b"},
		{name: "a\tb"},
		{name: "team/é"},
		{name: "team/\xff"},
		{name: ""},
		{name: "/api"},
		{name: "team/"},
		{name: "team//api"},
		{name: "team/../api"},
		{name: "Team/api"},
		{name: "team/api:latest"},
		{name: "team/api@sha256:abc"},
		{name: "team/a b"},
		{name: "https://example.com/api"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRepositoryName(tt.name)
			if tt.valid {
				if err != nil {
					t.Errorf("ValidateRepositoryName(%q): got = %v, want = nil", tt.name, err)
				}
			} else if !errors.Is(err, ErrInvalidRepositoryName) {
				t.Errorf("ValidateRepositoryName(%q): got = %v, want = ErrInvalidRepositoryName", tt.name, err)
			}
		})
	}
}
