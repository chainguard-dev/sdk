/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
)

var _ terms.ErrorDetail = (*TermsNotAcceptedDetail)(nil)

// MissingTermsDocs satisfies the terms.ErrorDetail interface, allowing
// the version-agnostic terms.IsTermsNotAccepted to extract missing documents
// from v1 error details.
func (d *TermsNotAcceptedDetail) MissingTermsDocs() []terms.Document {
	docs := make([]terms.Document, len(d.GetMissing()))
	for i, m := range d.GetMissing() {
		docs[i] = terms.Document{
			ID:    m.GetId(),
			Label: m.GetLabel(),
			URL:   m.GetUrl(),
		}
	}
	return docs
}
