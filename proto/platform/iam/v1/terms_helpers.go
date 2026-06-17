/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"context"

	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
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

// DocumentMetadata returns the display metadata for a known document ID,
// falling back to a TermsDocument with only the ID set for unknown IDs.
// Delegates to the shared terms package to avoid duplicating the canonical map.
func DocumentMetadata(id string) TermsDocument {
	doc := terms.DocumentMetadata(id)
	return TermsDocument{
		ID:    doc.ID,
		Label: doc.Label,
		URL:   doc.URL,
	}
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
