/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"errors"
	"testing"
)

func TestLibraryAllowListEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		entry   *LibraryAllowListEntry
		wantErr error
	}{{
		name:    "no bypass flag is inert",
		entry:   &LibraryAllowListEntry{Purl: "pkg:npm/lodash"},
		wantErr: ErrAllowEntryInert,
	}, {
		name:    "justification without a bypass flag is still inert",
		entry:   &LibraryAllowListEntry{Purl: "pkg:npm/lodash", Justification: "vetted"},
		wantErr: ErrAllowEntryInert,
	}, {
		name:    "malware bypass requires justification",
		entry:   &LibraryAllowListEntry{Purl: "pkg:npm/lodash", BypassMalware: true},
		wantErr: ErrJustificationRequired,
	}, {
		name:  "cooldown bypass alone is valid",
		entry: &LibraryAllowListEntry{Purl: "pkg:npm/lodash", BypassCooldown: true},
	}, {
		name:  "malware bypass with justification is valid",
		entry: &LibraryAllowListEntry{Purl: "pkg:npm/lodash", BypassMalware: true, Justification: "vetted"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Validate() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
