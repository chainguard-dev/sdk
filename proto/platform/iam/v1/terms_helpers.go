/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"
)

// TermsDocument describes a legal document, including its stable ID and optional
// display metadata (label and URL) for presentation in a UI.
type TermsDocument struct {
	// ID is the stable identifier for this document (e.g. "guardener-tos.v1").
	ID string
	// Label is a human-readable name for display purposes.
	Label string
	// URL is a link to the full document text.
	URL string
}

// knownDocuments holds display metadata for well-known legal document IDs.
var knownDocuments = map[string]TermsDocument{
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
}

// DocumentMetadata returns the display metadata for a known document ID,
// falling back to a TermsDocument with only the ID set for unknown IDs.
func DocumentMetadata(id string) TermsDocument {
	if doc, ok := knownDocuments[id]; ok {
		return doc
	}
	return TermsDocument{ID: id}
}

// CheckTermsAcceptance verifies that the group identified by groupID has
// accepted all of the required documents listed in requiredDocIDs.
// If any documents are missing, it returns a gRPC FailedPrecondition error
// carrying a TermsNotAcceptedDetail with the missing document IDs.
func CheckTermsAcceptance(ctx context.Context, groupID string, client TermsClient, requiredDocIDs []string) error {
	resp, err := client.ListAccepted(ctx, &TermsFilter{Group: groupID})
	if err != nil {
		return err
	}

	accepted := make(map[string]struct{}, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		accepted[item.GetDocId()] = struct{}{}
	}

	var missing []TermsDocument
	for _, id := range requiredDocIDs {
		if _, ok := accepted[id]; !ok {
			missing = append(missing, DocumentMetadata(id))
		}
	}
	if len(missing) > 0 {
		return ErrTermsNotAccepted(missing)
	}
	return nil
}
