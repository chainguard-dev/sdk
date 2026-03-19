/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package v1 provides types for chart-lock attestation data attached to
Chainguard Helm charts.

A chart-lock is an in-toto attestation attached to a Chainguard Helm chart OCI
artifact. It captures the pinned image references and value mapping templates for
every image in the chart, including subcharts. This enables downstream
consumers to resolve, relocate, or transform a chart's image references without
needing to understand the chart's internal values.yaml structure.

# Reading a Chart-Lock

Extract a chart-lock from a DSSE-wrapped attestation payload with
[ParseChartLockAttestation]. It returns [ErrChartLockNotFound] when the payload
is not a chart-lock:

	for _, att := range attestations {
		payload, _ := att.Payload()
		lock, err := v1.ParseChartLockAttestation(payload)
		if errors.Is(err, v1.ErrChartLockNotFound) {
			continue // not a chart-lock, try next
		}
		if err != nil {
			return err
		}
		// use lock
	}

# Generating Value Overrides

The primary workflow is to walk a chart-lock's image tree and produce Helm
value overrides. [ChartImages.Walk] recursively traverses images at the current
level and all subcharts, nesting subchart results under their dependency key:

	merged, err := lock.Images.Walk(func(id string, ref *v1.ChartImage, tokens images.TokenList) (any, error) {
		// ref has the pinned RepoName, Tag, and Digest.
		// tokens has the lexed template markers from the mapping.
		//
		// Return a transformed value, or nil to exclude the field.
		return myRegistry + "/" + ref.RepoName, nil
	})

The callback receives both the pinned image metadata and the template tokens,
so callers can produce anything from simple strings to structured expressions
depending on the use case.

# Updating Digests After Relocation

When copying a chart to a new registry, image digests may change.
[Lock.WithChartDigests] returns a new lock with updated digests and chart ref,
leaving all other fields intact:

	updated := lock.WithChartDigests(newChartRef, &v1.ChartDigests{
		Digests: map[string]string{"nginx": "sha256:newdigest..."},
	})
*/
package v1
