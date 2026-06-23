/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"context"

	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
)

const maxPages = 10

// CheckTermsAcceptance verifies that the group identified by groupID has
// accepted all of the required documents listed in requiredDocIDs.
// If any documents are missing, it returns a gRPC FailedPrecondition error
// carrying a TermsNotAcceptedDetail with the missing document metadata.
func CheckTermsAcceptance(ctx context.Context, groupID string, client TermsServiceClient, requiredDocIDs []string) error {
	accepted := make(map[string]struct{})
	var pageToken string
	for range maxPages {
		resp, err := client.ListTermsAcceptances(ctx, &ListTermsAcceptancesRequest{
			Group:     groupID,
			PageToken: pageToken,
		})
		if err != nil {
			return err
		}
		for _, item := range resp.GetTermsAcceptances() {
			accepted[item.GetDocId()] = struct{}{}
		}

		// Early exit once all required docs are matched.
		allMatched := true
		for _, id := range requiredDocIDs {
			if _, ok := accepted[id]; !ok {
				allMatched = false
				break
			}
		}
		if allMatched {
			return nil
		}

		if resp.GetNextPageToken() == "" {
			break
		}
		pageToken = resp.GetNextPageToken()
	}

	var missing []terms.Document
	for _, id := range requiredDocIDs {
		if _, ok := accepted[id]; !ok {
			missing = append(missing, terms.DocumentMetadata(id))
		}
	}
	if len(missing) > 0 {
		return ErrTermsNotAccepted(missing)
	}
	return nil
}
