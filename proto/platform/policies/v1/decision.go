/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "cmp"

// GetArtifactIdOrDigest returns the artifact this decision applies to,
// preferring artifact_id and falling back to the deprecated digest field.
//
//nolint:staticcheck,revive // reading the deprecated field is the point; the name mirrors the generated accessors.
func (x *Decision) GetArtifactIdOrDigest() string {
	return cmp.Or(x.GetArtifactId(), x.GetDigest())
}

// GetArtifactIdOrDigest returns the artifact to filter decisions by,
// preferring artifact_id and falling back to the deprecated digest field.
//
//nolint:staticcheck,revive // reading the deprecated field is the point; the name mirrors the generated accessors.
func (x *DecisionFilter) GetArtifactIdOrDigest() string {
	return cmp.Or(x.GetArtifactId(), x.GetDigest())
}
