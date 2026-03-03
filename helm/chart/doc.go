/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package chart provides utilities for working with Helm charts stored as OCI images.

# Overview

This package enables reading and modifying Helm chart OCI images, specifically
focusing on the values.yaml file. Helm charts can be stored in OCI registries
using the standard Helm chart media type, and this package provides tools to
extract and replace the values.yaml configuration without unpacking the entire
chart.

# Features

  - Read values.yaml from Helm chart OCI images
  - Replace values.yaml with transformed content
  - Preserve all other chart files and metadata
  - Support for image reference templating via chainguard.dev/sdk/helm/images

# Usage

Reading values from a chart:

	chart, err := remote.Image(ref)
	if err != nil {
		return err
	}

	values, err := chart.ReadValues(chart)
	if err != nil {
		return err
	}

	// values contains the raw YAML content
	fmt.Printf("values.yaml:\n%s\n", values)

Replacing values with image reference resolution:

	mapping := &images.Mapping{
		Images: map[string]*images.Image{
			"nginx": {
				Values: map[string]any{
					"image": map[string]any{
						"registry":   "${registry}",
						"repository": "${repo}",
					},
				},
			},
		},
	}

	refs := map[string]string{
		"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abc123...",
	}

	patched, err := chart.ReplaceValues(chart, mapping, refs)
	if err != nil {
		return err
	}

	// patched is a new OCI image with updated values.yaml
	if err := remote.Write(dstRef, patched); err != nil {
		return err
	}

# Integration Patterns

This package is designed to work with:

  - github.com/google/go-containerregistry for OCI image operations
  - chainguard.dev/sdk/helm/images for image reference templating
  - Standard Helm chart OCI images (media type: application/vnd.cncf.helm.chart.content.v1.tar+gzip)

The package preserves all chart metadata, annotations, and files other than
values.yaml when performing replacements, ensuring chart integrity.
*/
package chart
