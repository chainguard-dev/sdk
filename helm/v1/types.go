/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package v1 contains types for chart-lock data.
package v1

import "chainguard.dev/sdk/helm/images"

// PredicateType is the predicate type for chart-lock attestations.
const PredicateType = "https://chainguard.dev/attestation/chart-lock/v1"

// Lock contains the locked chart metadata and image references.
type Lock struct {
	Chart  *Chart       `json:"chart"`
	Images *ChartImages `json:"images"`
}

// Chart identifies the chart this lock is for.
type Chart struct {
	Package string `json:"package"`
	Ref     string `json:"ref"`
}

// ChartImages contains the locked image references and mapping template.
type ChartImages struct {
	Refs     map[string]*ChartImage `json:"refs"`
	Template *images.Mapping        `json:"template"`
}

// ChartImage holds a pinned image reference.
type ChartImage struct {
	RepoName string `json:"repoName"`
	Tag      string `json:"tag"`
	Digest   string `json:"digest"`
}

// WithResolvedDigests returns a new Lock with the chart ref updated and image digests
// replaced with the provided resolved digests. The images map keys should match the
// image IDs in the lock's Refs, and values should be the resolved digest strings.
func (l *Lock) WithResolvedDigests(chartRef string, digests map[string]string) *Lock {
	resolvedRefs := make(map[string]*ChartImage, len(l.Images.Refs))
	for id, img := range l.Images.Refs {
		resolved := &ChartImage{
			RepoName: img.RepoName,
			Tag:      img.Tag,
			Digest:   img.Digest,
		}
		if d, ok := digests[id]; ok {
			resolved.Digest = d
		}
		resolvedRefs[id] = resolved
	}

	return &Lock{
		Chart: &Chart{
			Package: l.Chart.Package,
			Ref:     chartRef,
		},
		Images: &ChartImages{
			Refs:     resolvedRefs,
			Template: l.Images.Template,
		},
	}
}
