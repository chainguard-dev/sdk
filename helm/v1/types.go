/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"fmt"
	"iter"
	"slices"

	"chainguard.dev/sdk/helm/images"
)

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
//
// Subcharts holds nested ChartImages keyed by Helm dependency name.
type ChartImages struct {
	Refs      map[string]*ChartImage  `json:"refs"`
	Template  *images.Mapping         `json:"template"`
	Subcharts map[string]*ChartImages `json:"subcharts,omitempty"`
}

// ChartImage holds a pinned image reference.
type ChartImage struct {
	RepoName string `json:"repoName"`
	Tag      string `json:"tag"`
	Digest   string `json:"digest"`
}

// Reference renders the image reference as prefix/repoName[:tag][@digest],
// where prefix is the registry/org the image lives under. An empty tag,
// digest, or prefix is omitted.
func (img *ChartImage) Reference(prefix string) string {
	r := img.RepoName
	if prefix != "" {
		r = prefix + "/" + r
	}
	if img.Tag != "" {
		r += ":" + img.Tag
	}
	if img.Digest != "" {
		r += "@" + img.Digest
	}
	return r
}

// IsOptionalImage reports whether imageID is marked as optional in the template
// at this level. Returns false if the template is nil or the image is not found.
func (ci *ChartImages) IsOptionalImage(imageID string) bool {
	if ci == nil || ci.Template == nil {
		return false
	}
	img, ok := ci.Template.Images[imageID]
	if !ok || img == nil {
		return false
	}
	return img.Requirement.IsOptional()
}

// WalkFunc is the callback for ChartImages.Walk.
// It receives the image ID, the pinned image ref, and the tokenized template value.
// Return the transformed value, or nil to exclude the field from output.
type WalkFunc func(imageID string, ref *ChartImage, tokens images.TokenList) (any, error)

// Walk recursively walks all images at this level and its subcharts.
// Each image's ref is looked up from Refs and passed to fn along with the image ID.
// Returns an error if an image ID in the template has no corresponding entry in Refs.
// Subchart results are nested under the subchart key in the output.
func (ci *ChartImages) Walk(fn WalkFunc) (map[string]any, error) {
	if ci == nil {
		return nil, nil
	}
	if fn == nil {
		return nil, fmt.Errorf("callback function is nil")
	}

	var result map[string]any
	if ci.Template != nil {
		iv, err := ci.Template.Walk(func(imageID string, tokens images.TokenList) (any, error) {
			ref, ok := ci.Refs[imageID]
			if !ok {
				return nil, fmt.Errorf("image %q: not found in refs", imageID)
			}
			return fn(imageID, ref, tokens)
		})
		if err != nil {
			return nil, err
		}
		result, err = iv.Merge()
		if err != nil {
			return nil, err
		}
	}

	// Recurse into subcharts, nesting under the subchart key.
	depNames := make([]string, 0, len(ci.Subcharts))
	for name := range ci.Subcharts {
		depNames = append(depNames, name)
	}
	slices.Sort(depNames)
	for _, depName := range depNames {
		sub := ci.Subcharts[depName]
		subResult, err := sub.Walk(fn)
		if err != nil {
			return nil, fmt.Errorf("subchart %q: %w", depName, err)
		}
		if len(subResult) > 0 {
			if result == nil {
				result = make(map[string]any)
			}
			result[depName] = subResult
		}
	}

	return result, nil
}

// Images iterates every pinned image in this ChartImages and all nested
// subcharts in unspecified order; a nil receiver yields nothing. Unlike Walk
// it ignores the template, so it also visits images with no template mapping.
func (ci *ChartImages) Images() iter.Seq[*ChartImage] {
	return func(yield func(*ChartImage) bool) {
		var walk func(*ChartImages) bool
		walk = func(c *ChartImages) bool {
			if c == nil {
				return true
			}
			for _, img := range c.Refs {
				if img != nil && !yield(img) {
					return false
				}
			}
			for _, sub := range c.Subcharts {
				if !walk(sub) {
					return false
				}
			}
			return true
		}
		walk(ci)
	}
}

// ChartDigests holds resolved digest strings for a chart level and its subcharts.
type ChartDigests struct {
	Digests   map[string]string        // image ID -> resolved digest
	Subcharts map[string]*ChartDigests // dep name -> subchart digests
}

// WithChartDigests returns a new Lock with the chart ref updated and image digests
// replaced with the provided resolved digests.
func (l *Lock) WithChartDigests(chartRef string, digests *ChartDigests) *Lock {
	return &Lock{
		Chart: &Chart{
			Package: l.Chart.Package,
			Ref:     chartRef,
		},
		Images: l.Images.withChartDigests(digests),
	}
}

func (ci *ChartImages) withChartDigests(digests *ChartDigests) *ChartImages {
	resolvedRefs := make(map[string]*ChartImage, len(ci.Refs))
	for id, img := range ci.Refs {
		resolved := &ChartImage{
			RepoName: img.RepoName,
			Tag:      img.Tag,
			Digest:   img.Digest,
		}
		if digests != nil {
			if d, ok := digests.Digests[id]; ok {
				resolved.Digest = d
			}
		}
		resolvedRefs[id] = resolved
	}

	result := &ChartImages{
		Refs:     resolvedRefs,
		Template: ci.Template,
	}

	if len(ci.Subcharts) > 0 {
		result.Subcharts = make(map[string]*ChartImages, len(ci.Subcharts))
		for name, sub := range ci.Subcharts {
			var subDigests *ChartDigests
			if digests != nil {
				subDigests = digests.Subcharts[name]
			}
			result.Subcharts[name] = sub.withChartDigests(subDigests)
		}
	}

	return result
}
