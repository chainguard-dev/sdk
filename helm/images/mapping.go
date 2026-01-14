/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package images

import (
	"encoding/json"
	"fmt"
	"io"
	"slices"

	yamlpatch "github.com/palantir/pkg/yamlpatch"
)

// Mapping is the top-level schema for chart image values mappings.
// This is the structure of values.images.json embedded in chart APKs.
type Mapping struct {
	Images map[string]*Image `json:"images"`
}

// Image defines how a single container image maps to values.yaml paths.
type Image struct {
	Values map[string]any `json:"values"`
}

// WalkFunc is the callback signature for Mapping.Walk.
// It receives the image ID and tokenized string value, returning the transformed value.
type WalkFunc func(imageID string, tokens TokenList) (any, error)

// ImageValues maps image IDs to their resolved values.
type ImageValues map[string]map[string]any

// Parse parses and validates a Mapping from JSON.
// Returns an error if the JSON is invalid, required fields are missing,
// or any ${...} markers reference unknown fields.
func Parse(r io.Reader) (*Mapping, error) {
	var m Mapping
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, fmt.Errorf("parsing image mapping: %w", err)
	}

	for id, img := range m.Images {
		if img == nil {
			return nil, fmt.Errorf("image %q: nil definition", id)
		}
		if img.Values == nil {
			return nil, fmt.Errorf("image %q: missing 'values' field", id)
		}
		if err := validateMarkers(img.Values); err != nil {
			return nil, fmt.Errorf("image %q: %w", id, err)
		}
	}

	return &m, nil
}

func validateMarkers(v any) error {
	switch val := v.(type) {
	case string:
		_, err := lex(val)
		return err
	case map[string]any:
		for k, v := range val {
			if err := validateMarkers(v); err != nil {
				return fmt.Errorf("key %q: %w", k, err)
			}
		}
	case []any:
		for i, v := range val {
			if err := validateMarkers(v); err != nil {
				return fmt.Errorf("index %d: %w", i, err)
			}
		}
	}
	return nil
}

// Merge combines all image values into a single map.
// Images are merged in sorted order by ID for deterministic results.
// Returns an error if there's a conflict (e.g., merging a map with a non-map at the same path).
func (v ImageValues) Merge() (map[string]any, error) {
	ids := make([]string, 0, len(v))
	for id := range v {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	result := make(map[string]any)
	for _, id := range ids {
		if err := merge(result, v[id]); err != nil {
			return nil, fmt.Errorf("image %q: %w", id, err)
		}
	}
	return result, nil
}

// Walk traverses all string values in the mapping, lexing them into tokens.
// The callback receives the image ID and tokens for each string value.
// Returns transformed values for each image.
func (m *Mapping) Walk(fn WalkFunc) (ImageValues, error) {
	if m == nil {
		return nil, nil
	}
	if fn == nil {
		return nil, fmt.Errorf("callback function is nil")
	}
	result := make(ImageValues, len(m.Images))
	for id, img := range m.Images {
		if img == nil || img.Values == nil {
			continue
		}
		transformed, err := walkValue(img.Values, func(s string) (any, error) {
			tokens, err := lex(s)
			if err != nil {
				return nil, err
			}
			return fn(id, tokens)
		})
		if err != nil {
			return nil, fmt.Errorf("image %q: %w", id, err)
		}
		values, ok := transformed.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("image %q: values must be a map, got %T", id, transformed)
		}
		result[id] = values
	}
	return result, nil
}

func walkValue(v any, fn func(string) (any, error)) (any, error) {
	switch val := v.(type) {
	case string:
		return fn(val)
	case map[string]any:
		result := make(map[string]any, len(val))
		for k, v := range val {
			transformed, err := walkValue(v, fn)
			if err != nil {
				return nil, err
			}
			result[k] = transformed
		}
		return result, nil
	case []any:
		result := make([]any, len(val))
		for i, v := range val {
			transformed, err := walkValue(v, fn)
			if err != nil {
				return nil, err
			}
			result[i] = transformed
		}
		return result, nil
	default:
		return val, nil
	}
}

func merge(dst, src map[string]any) error {
	for key, srcVal := range src {
		dstVal, exists := dst[key]
		if !exists {
			dst[key] = srcVal
			continue
		}
		srcMap, srcIsMap := srcVal.(map[string]any)
		dstMap, dstIsMap := dstVal.(map[string]any)
		if srcIsMap && dstIsMap {
			if err := merge(dstMap, srcMap); err != nil {
				return err
			}
			continue
		}
		if srcIsMap != dstIsMap {
			return fmt.Errorf("conflict at key %q: cannot merge map with non-map", key)
		}
		dst[key] = srcVal
	}
	return nil
}

// Resolve resolves the mapping with refs and merges into valuesr, preserving
// comments
func (m *Mapping) Resolve(refs map[string]string, valuesr io.Reader) ([]byte, error) {
	original, err := io.ReadAll(valuesr)
	if err != nil {
		return nil, fmt.Errorf("reading values: %w", err)
	}

	if m == nil || len(refs) == 0 {
		return original, nil
	}

	ociRefs := make(map[string]OCIRef, len(refs))
	for id, ref := range refs {
		ociRef, err := NewRef(ref)
		if err != nil {
			return nil, fmt.Errorf("image %q: %w", id, err)
		}
		ociRefs[id] = ociRef
	}

	imageValues, err := m.Walk(Resolve(ociRefs))
	if err != nil {
		return nil, err
	}

	merged, err := imageValues.Merge()
	if err != nil {
		return nil, err
	}

	var patch yamlpatch.Patch
	var walk func(path yamlpatch.Path, m map[string]any)
	walk = func(path yamlpatch.Path, m map[string]any) {
		for k, v := range m {
			p := append(slices.Clone(path), k)
			if nested, ok := v.(map[string]any); ok {
				walk(p, nested)
			} else {
				patch = append(patch, yamlpatch.Operation{
					Type:  yamlpatch.OperationReplace,
					Path:  p,
					Value: v,
				})
			}
		}
	}
	walk(yamlpatch.Path{""}, merged)

	return yamlpatch.Apply(original, patch)
}
