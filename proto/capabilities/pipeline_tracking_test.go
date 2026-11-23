/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPipelineTrackingRoles(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name    string
		caps    []Capability
		allowed Capability
		denied  Capability
	}{
		{"submit only", LibrariesPipelineTrackingSubmitterCaps, Capability_CAP_LIBRARIES_PIPELINE_TRACKING_SUBMIT, Capability_CAP_LIBRARIES_PIPELINE_TRACKING_READ},
		{"read only", LibrariesPipelineTrackingReaderCaps, Capability_CAP_LIBRARIES_PIPELINE_TRACKING_READ, Capability_CAP_LIBRARIES_PIPELINE_TRACKING_SUBMIT},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			want := SortCaps([]Capability{Capability_CAP_INTERNAL, tt.allowed})
			if diff := cmp.Diff(want, tt.caps); diff != "" {
				t.Errorf("role capabilities (-want,+got): %s", diff)
			}
			if !InternalOnly(tt.allowed) {
				t.Errorf("capability %v: got = external, want = internal-only", tt.allowed)
			}
			if slices.Contains(tt.caps, tt.denied) {
				t.Errorf("role grants unrelated tracking capability %v", tt.denied)
			}
			for _, customerCaps := range [][]Capability{OwnerCaps, EditorCaps, ViewerCaps} {
				if slices.Contains(customerCaps, tt.allowed) {
					t.Errorf("customer role grants internal tracking capability %v", tt.allowed)
				}
			}
		})
	}
}
