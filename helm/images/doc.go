/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package images provides utilities for parsing and resolving image value
mappings embedded in Chainguard Helm chart packages.

Chainguard Helm charts include metadata that describes how container image
references map to values.yaml paths. This package parses that metadata and
provides an API for resolving image references into concrete Helm values
without requiring callers to understand the underlying template format.

# Patching a values.yaml File

The simplest way to use a mapping is [Mapping.Resolve], which substitutes all
image references and patches the results into a YAML values file while
preserving comments and formatting:

	m, err := images.Parse(reader)
	if err != nil {
		return err
	}

	refs := map[string]string{
		"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abc123...",
	}
	patched, err := m.Resolve(refs, bytes.NewReader(originalValues))

# Custom Transformations

When the output is not a simple string — for example, generating Terraform
expressions or building relocation overrides — use [Mapping.Walk] with a
custom callback. Walk lexes each template value into a [TokenList] and invokes
the callback per image:

	values, err := m.Walk(func(imageID string, tokens images.TokenList) (any, error) {
		// Inspect tokens to decide what to emit.
		// Return nil to exclude a field from the output.
		parts := tokens.Map(func(f images.RefField) any {
			return myCustomTransform(f)
		})
		return parts, nil
	})

[TokenList.Map] applies a function to each [RefField] in the token list,
preserving literal segments. The [RefField] constants (Registry, Repo, Tag,
Digest, etc.) identify which component of an image reference a given template
position expects.

[Resolve] is a convenience that returns a [WalkFunc] performing standard string
substitution from a map of [OCIRef] values. It can be passed directly to Walk.

Results from Walk are returned as [ImageValues], which can be combined into a
single map with [ImageValues.Merge].
*/
package images

// ChainguardChartMetadataFilename is the name of the Chainguard metadata file embedded in chart APKs.
const ChainguardChartMetadataFilename = "cg.json"
