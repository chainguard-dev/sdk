/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms

// Document describes a legal document, including its stable ID and optional
// display metadata (label and URL) for presentation in a UI.
type Document struct {
	ID    string
	Label string
	URL   string
}

var knownDocuments = map[string]Document{
	"guardener-tos.v1": {
		ID:    "guardener-tos.v1",
		Label: "Terms of Service",
		URL:   "https://www.chainguard.dev/legal/guardener",
	},
	"sfdpa.v1": {
		ID:    "sfdpa.v1",
		Label: "Data Privacy Agreement",
		URL:   "https://www.chainguard.dev/legal/supplemental-dpa",
	},
	"skills-tos.v1": {
		ID:    "skills-tos.v1",
		Label: "Skills Registry Terms of Service",
		URL:   "https://www.chainguard.dev/legal/agent-skills-disclosure",
	},
}

// DocumentMetadata returns the display metadata for a known document ID,
// falling back to a Document with only the ID set for unknown IDs.
func DocumentMetadata(id string) Document {
	if doc, ok := knownDocuments[id]; ok {
		return doc
	}
	return Document{ID: id}
}
