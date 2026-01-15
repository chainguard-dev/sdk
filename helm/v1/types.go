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
