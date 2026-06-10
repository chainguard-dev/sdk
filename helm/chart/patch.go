/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"slices"

	"chainguard.dev/sdk/helm/images"
	helmv1 "chainguard.dev/sdk/helm/v1"
	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// ErrSkipImage is returned by a [ResolveFunc] to indicate that an image
// should be excluded from patching (e.g., an optional image that is missing).
var ErrSkipImage = errors.New("skip image")

// ResolveFunc resolves a single chart image to its full OCI reference string
// and the resolved digest. Callers provide their own resolution logic
// (e.g., registry lookups, custom assembly checks) via this callback.
//
// The ci parameter is the [helmv1.ChartImages] that owns this image, allowing
// callers to check level-specific metadata (e.g., whether the image is optional
// at this level of the chart tree).
//
// Return [ErrSkipImage] to exclude the image from patching.
type ResolveFunc func(ci *helmv1.ChartImages, imageID string, ref *helmv1.ChartImage) (fullRef, resolvedDigest string, err error)

// PatchChartImages resolves all image references in a [helmv1.ChartImages]
// tree and patches them into the chart's values.yaml.
//
// Root image values are replaced strictly — an error is returned if a
// template path does not exist in the chart's values.yaml. Subchart image
// overrides use [images.WithAddMissing]: existing paths are replaced, missing
// paths are added at the deepest existing ancestor.
//
// Both paths preserve YAML comments and formatting via yamlpatch.
//
// Returns the patched chart image and a ChartDigests tree for attestation.
func PatchChartImages(chart v1.Image, ci *helmv1.ChartImages, resolve ResolveFunc, opts ...images.ResolveOption) (v1.Image, *helmv1.ChartDigests, error) {
	if ci == nil {
		return chart, nil, nil
	}

	p := &patcher{
		resolve: resolve,
		opts:    opts,
		subOpts: append([]images.ResolveOption{images.WithAddMissing(true)}, opts...),
	}

	rootRefs, rootDigests, err := p.resolveRefs(ci)
	if err != nil {
		return nil, nil, err
	}

	chartDigests := &helmv1.ChartDigests{Digests: rootDigests}

	values, err := ReadValues(chart)
	if err != nil {
		return nil, nil, fmt.Errorf("reading values: %w", err)
	}
	if values == nil {
		return nil, nil, fmt.Errorf("chart has no values.yaml")
	}

	patched := values

	// Apply root image template strictly — paths must exist in values.yaml.
	// When images were skipped, filter the template to exclude them.
	if ci.Template != nil && len(rootRefs) > 0 {
		tmpl := ci.Template
		if len(rootRefs) < len(ci.Refs) {
			tmpl = tmpl.ForImages(rootRefs)
		}
		patched, err = tmpl.Resolve(rootRefs, bytes.NewReader(patched), p.opts...)
		if err != nil {
			return nil, nil, fmt.Errorf("resolving root images: %w", err)
		}
	}

	// Resolve and apply subchart images recursively.
	if len(ci.Subcharts) > 0 {
		patched, chartDigests.Subcharts, err = p.patchSubcharts(patched, ci.Subcharts, nil)
		if err != nil {
			return nil, nil, err
		}
	}

	img, err := rewriteChart(chart, patched, isTopLevelValuesYAML)
	if err != nil {
		return nil, nil, err
	}
	return img, chartDigests, nil
}

type patcher struct {
	resolve ResolveFunc
	opts    []images.ResolveOption
	subOpts []images.ResolveOption
}

func (p *patcher) resolveRefs(ci *helmv1.ChartImages) (map[string]string, map[string]string, error) {
	if len(ci.Refs) == 0 {
		return nil, nil, nil
	}
	fullRefs := make(map[string]string, len(ci.Refs))
	digests := make(map[string]string, len(ci.Refs))
	for imageID, ref := range ci.Refs {
		fullRef, resolvedDigest, err := p.resolve(ci, imageID, ref)
		if errors.Is(err, ErrSkipImage) {
			continue
		}
		if err != nil {
			return nil, nil, fmt.Errorf("image %q: %w", imageID, err)
		}
		fullRefs[imageID] = fullRef
		digests[imageID] = resolvedDigest
	}
	return fullRefs, digests, nil
}

// patchSubcharts resolves each subchart's images and patches them into the
// YAML sequentially. path accumulates ancestor dependency keys so that nested
// subcharts target the correct scope (e.g., redis.sentinel.image.registry).
func (p *patcher) patchSubcharts(valuesYAML []byte, subcharts map[string]*helmv1.ChartImages, path []string) ([]byte, map[string]*helmv1.ChartDigests, error) {
	digestsMap := make(map[string]*helmv1.ChartDigests, len(subcharts))
	patched := valuesYAML

	for _, depName := range slices.Sorted(maps.Keys(subcharts)) {
		ci := subcharts[depName]
		subPath := append(slices.Clone(path), depName)

		refs, digests, err := p.resolveRefs(ci)
		if err != nil {
			return nil, nil, fmt.Errorf("subchart %q: %w", depName, err)
		}

		if ci.Template != nil && len(refs) > 0 {
			tmpl := ci.Template
			if len(refs) < len(ci.Refs) {
				tmpl = tmpl.ForImages(refs)
			}
			nested := nestTemplate(tmpl, subPath)
			patched, err = nested.Resolve(refs, bytes.NewReader(patched), p.subOpts...)
			if err != nil {
				return nil, nil, fmt.Errorf("subchart %q: %w", depName, err)
			}
		}

		chartDigests := &helmv1.ChartDigests{Digests: digests}

		if len(ci.Subcharts) > 0 {
			patched, chartDigests.Subcharts, err = p.patchSubcharts(patched, ci.Subcharts, subPath)
			if err != nil {
				return nil, nil, fmt.Errorf("subchart %q: %w", depName, err)
			}
		}

		digestsMap[depName] = chartDigests
	}

	return patched, digestsMap, nil
}

// nestTemplate wraps each image's Values under the given path of dependency
// keys. For path ["redis", "sentinel"], {image: {registry: ...}} becomes
// {redis: {sentinel: {image: {registry: ...}}}}.
func nestTemplate(m *images.Mapping, path []string) *images.Mapping {
	nested := &images.Mapping{Images: make(map[string]*images.Image, len(m.Images))}
	for id, img := range m.Images {
		v := img.Values
		for i := len(path) - 1; i >= 0; i-- {
			v = map[string]any{path[i]: v}
		}
		nested.Images[id] = &images.Image{Values: v}
	}
	return nested
}
