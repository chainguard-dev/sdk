/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "testing"

// TestOverrideGetArtifactIdOrDigest guards the case that matters most in
// practice: a client and the server it talks to are routinely on different
// versions, so an override arrives carrying either field. Reading artifact_id
// alone renders every pre-migration override with a blank artifact.
func TestOverrideGetArtifactIdOrDigest(t *testing.T) {
	const (
		digest = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
		purl   = "pkg:npm/left-pad@1.3.0"
	)

	for _, tt := range []struct {
		name     string
		override *Override
		want     string
	}{{
		name:     "artifact_id alone",
		override: &Override{ArtifactId: purl},
		want:     purl,
	}, {
		// What a server that predates artifact_id sends.
		name: "deprecated digest alone",
		//nolint:staticcheck // exercising the back-compat path.
		override: &Override{Digest: digest},
		want:     digest,
	}, {
		// artifact_id wins, so a stale mirror cannot redirect the waiver.
		name: "both set",
		//nolint:staticcheck // exercising the precedence rule.
		override: &Override{ArtifactId: purl, Digest: digest},
		want:     purl,
	}, {
		name:     "neither set",
		override: &Override{},
		want:     "",
	}, {
		name:     "nil receiver",
		override: nil,
		want:     "",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.override.GetArtifactIdOrDigest(); got != tt.want {
				t.Errorf("GetArtifactIdOrDigest() = %q, want %q", got, tt.want)
			}
		})
	}
}
