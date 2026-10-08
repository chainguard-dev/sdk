/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms

import "testing"

func TestFeatureDocumentsAreKnown(t *testing.T) {
	for feature, ids := range featureDocuments {
		for _, id := range ids {
			if doc := DocumentMetadata(id); doc.Label == "" || doc.URL == "" {
				t.Errorf("feature %q requires %q, which has no label or URL", feature, id)
			}
		}
	}
}

func TestFeatureDocuments(t *testing.T) {
	if got := FeatureDocuments("sandbox_checks"); len(got) != 2 || got[0] != "code-analysis-terms.v1" || got[1] != "msla.v1" {
		t.Errorf("sandbox_checks: got %v", got)
	}
	if got := FeatureDocuments("sandbox_workspaces"); len(got) != 2 || got[0] != "secure-platform-terms.v1" || got[1] != "msla.v1" {
		t.Errorf("sandbox_workspaces: got %v", got)
	}
	if got := FeatureDocuments("unknown"); got != nil {
		t.Errorf("unknown feature: got %v, want nil", got)
	}
}
