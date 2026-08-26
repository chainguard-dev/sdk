/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "cmp"

// GetArtifactIdOrDigest returns the artifact this override waives, preferring
// artifact_id and falling back to the deprecated digest field.
//
// A client is routinely newer or older than the server it talks to, so both
// fields have to be tolerated on read. Consolidating that here keeps every
// consumer from re-deriving it — and from forgetting to.
//
//nolint:staticcheck,revive // reading the deprecated field is the point; the name mirrors the generated accessors.
func (x *Override) GetArtifactIdOrDigest() string {
	return cmp.Or(x.GetArtifactId(), x.GetDigest())
}
