/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms

import (
	"slices"
	"testing"
)

func TestCustomerUploadsDocuments(t *testing.T) {
	want := []string{"secure-platform-terms.v1"}
	got := CustomerUploadsDocuments()
	if !slices.Equal(got, want) {
		t.Errorf("CustomerUploadsDocuments(): got = %v, want = %v", got, want)
	}
	for _, id := range got {
		// chainctl refuses to prompt for a document with no label or URL.
		if doc := DocumentMetadata(id); doc.Label == "" || doc.URL == "" {
			t.Errorf("document %q: got = no label or URL, want = known document", id)
		}
	}
	// The caller gets its own copy.
	got[0] = "changed"
	if again := CustomerUploadsDocuments(); !slices.Equal(again, want) {
		t.Errorf("CustomerUploadsDocuments() after the caller changed its copy: got = %v, want = %v", again, want)
	}
}
