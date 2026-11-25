/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms

import "slices"

// customerUploadsDocuments are the documents an organization must accept
// before it uploads its own language packages or container images to
// Chainguard. They do not cover skills pushed to the registry's uploads
// subdomain, which requires the "skills" feature's documents instead.
// Uploading packages and images is not an entitlement feature, so it has no
// featureDocuments entry. Acceptance is recorded per document, so an
// organization that accepted these for another offering, such as Workspaces,
// does not accept them again.
var customerUploadsDocuments = []string{"secure-platform-terms.v1"}

// CustomerUploadsDocuments returns the IDs of the documents an organization
// must accept before it uploads its own language packages or container
// images. It does not apply to skills pushes.
func CustomerUploadsDocuments() []string {
	return slices.Clone(customerUploadsDocuments)
}
