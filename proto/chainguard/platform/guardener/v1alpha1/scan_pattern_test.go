/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// TestScanHTTPPatterns pins the ListScans and GetScan HTTP bindings: the group
// UIDP is a mid-path multi-segment wildcard (a UIDP may contain "/") anchored
// by the trailing "scans" literal (and, for GetScan, the id segment).
// grpc-gateway supports a mid-path "**" when a fixed number of segments
// follows (the registry surface relies on the same construct), but only ONE
// deep wildcard per pattern — so ListScans' repository filter, which is also
// hierarchical (GitLab namespaces nest arbitrarily deep), rides in the ?repo=
// query parameter instead. This test guards the URL contract documented on
// the RPCs.
func TestScanHTTPPatterns(t *testing.T) {
	for _, tc := range []struct {
		name    string
		pattern runtime.Pattern
		path    string
		want    map[string]string // nil = no match
	}{
		{
			name:    "list root group",
			pattern: pattern_Guardener_ListScans_0,
			path:    "guardener/v1alpha1/deadbeef/scans",
			want:    map[string]string{"group": "deadbeef"},
		},
		{
			name:    "list nested group",
			pattern: pattern_Guardener_ListScans_0,
			path:    "guardener/v1alpha1/deadbeef/cafe1234/scans",
			want:    map[string]string{"group": "deadbeef/cafe1234"},
		},
		{
			// The wildcard group matches zero segments; the handler then
			// rejects the empty group as invalid.
			name:    "list no group",
			pattern: pattern_Guardener_ListScans_0,
			path:    "guardener/v1alpha1/scans",
			want:    map[string]string{"group": ""},
		},
		{
			// A trailing id segment belongs to GetScan, not ListScans: the
			// list wildcard consumes segments until only "scans" remains,
			// which then can't match "someid".
			name:    "list rejects id segment",
			pattern: pattern_Guardener_ListScans_0,
			path:    "guardener/v1alpha1/deadbeef/scans/someid",
		},
		{
			name:    "get root group",
			pattern: pattern_Guardener_GetScan_0,
			path:    "guardener/v1alpha1/deadbeef/scans/someid",
			want:    map[string]string{"group": "deadbeef", "id": "someid"},
		},
		{
			name:    "get nested group",
			pattern: pattern_Guardener_GetScan_0,
			path:    "guardener/v1alpha1/deadbeef/cafe1234/scans/someid",
			want:    map[string]string{"group": "deadbeef/cafe1234", "id": "someid"},
		},
		{
			// Without an id the fixed tail ("scans" + id) can't be satisfied
			// with a "scans" literal in place: the wildcard would have to
			// leave two segments, making the literal position "deadbeef".
			name:    "get rejects bare scans",
			pattern: pattern_Guardener_GetScan_0,
			path:    "guardener/v1alpha1/deadbeef/scans",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.pattern.MatchAndEscape(
				strings.Split(tc.path, "/"), "", runtime.UnescapingModeDefault)
			if tc.want == nil {
				if err == nil {
					t.Fatalf("MatchAndEscape(%q) = %v, want no match", tc.path, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("MatchAndEscape(%q): %v", tc.path, err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("MatchAndEscape(%q) mismatch (-want +got):\n%s", tc.path, diff)
			}
		})
	}
}
