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
	"gopkg.in/yaml.v3"
)

// ResolveOption allows customizing the resolution process.
type ResolveOption func(*resolveConfig)

type resolveConfig struct {
	omitDigests bool
	addMissing  bool
}

// WithOmitDigests controls whether digests are included in resolved values. Default is false.
func WithOmitDigests(omitDigests bool) ResolveOption {
	return func(cfg *resolveConfig) {
		cfg.omitDigests = omitDigests
	}
}

// WithAddMissing controls whether paths not present in the target YAML are
// added (true) or cause an error (false, the default). Use this when patching
// values that may not have pre-existing entries, such as subchart overrides.
func WithAddMissing(addMissing bool) ResolveOption {
	return func(cfg *resolveConfig) {
		cfg.addMissing = addMissing
	}
}

// Mapping is the top-level schema for chart image value mappings,
// parsed from the cg.json file embedded in Chainguard chart APKs.
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
			if v == nil {
				result[k] = nil
				continue
			}
			transformed, err := walkValue(v, fn)
			if err != nil {
				return nil, err
			}
			if transformed != nil {
				result[k] = transformed
			}
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

// Resolve resolves all image references in the mapping and patches the results
// into the provided values.yaml content, preserving YAML comments and formatting.
//
// By default, all template paths must exist in the YAML — an error is returned
// for missing paths. Use [WithAddMissing] to allow missing paths to be added.
func (m *Mapping) Resolve(refs map[string]string, valuesr io.Reader, opts ...ResolveOption) ([]byte, error) {
	original, err := io.ReadAll(valuesr)
	if err != nil {
		return nil, fmt.Errorf("reading values: %w", err)
	}

	if m == nil || len(refs) == 0 {
		return original, nil
	}

	cfg := &resolveConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	ociRefs := make(map[string]OCIRef, len(refs))
	for id, ref := range refs {
		ociRef, err := NewRef(ref)
		if err != nil {
			return nil, fmt.Errorf("image %q: %w", id, err)
		}
		ociRefs[id] = ociRef
	}

	imageValues, err := m.Walk(Resolve(ociRefs, opts...))
	if err != nil {
		return nil, err
	}

	merged, err := imageValues.Merge()
	if err != nil {
		return nil, err
	}

	paths, err := newYAMLPaths(original)
	if err != nil {
		return nil, fmt.Errorf("parsing values: %w", err)
	}

	patch := paths.buildPatch(yamlpatch.Path{""}, merged)

	if !cfg.addMissing {
		for _, op := range patch {
			if op.Type == yamlpatch.OperationAdd {
				return nil, fmt.Errorf("path %s does not exist in values", op.Path)
			}
		}
	}

	// yamlpatch requires a root mapping node; seed one if the document is empty.
	if paths.children == nil && len(patch) > 0 {
		original = []byte("{}\n")
	}

	return yamlpatch.Apply(original, patch)
}

// yamlPaths mirrors the key hierarchy of a YAML document as a tree.
// buildPatch walks this tree alongside the values map to decide whether
// each operation should Add (missing path) or Replace (existing path).
type yamlPaths struct {
	children map[string]*yamlPaths
}

func newYAMLPaths(data []byte) (*yamlPaths, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	root := &yamlPaths{}
	if doc.Kind == yaml.DocumentNode && len(doc.Content) > 0 {
		root.collect(doc.Content[0])
	}
	return root, nil
}

func (yp *yamlPaths) collect(node *yaml.Node) {
	if node.Kind != yaml.MappingNode {
		return
	}
	yp.children = make(map[string]*yamlPaths, len(node.Content)/2)
	for i := 0; i+1 < len(node.Content); i += 2 {
		child := &yamlPaths{}
		child.collect(node.Content[i+1])
		yp.children[node.Content[i].Value] = child
	}
}

// buildPatch generates yamlpatch operations for the given values.
// Existing paths produce Replace; missing paths produce Add.
func (yp *yamlPaths) buildPatch(path yamlpatch.Path, m map[string]any) yamlpatch.Patch {
	var patch yamlpatch.Patch
	for k, v := range m {
		p := append(slices.Clone(path), k)
		child := yp.children[k]
		nested, isMap := v.(map[string]any)
		switch {
		case isMap && child != nil:
			patch = append(patch, child.buildPatch(p, nested)...)
		case isMap:
			patch = append(patch, yamlpatch.Operation{Type: yamlpatch.OperationAdd, Path: p, Value: v})
		case child != nil:
			patch = append(patch, yamlpatch.Operation{Type: yamlpatch.OperationReplace, Path: p, Value: v})
		default:
			patch = append(patch, yamlpatch.Operation{Type: yamlpatch.OperationAdd, Path: p, Value: v})
		}
	}
	return patch
}
