/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package images provides utilities for parsing and transforming Helm chart
// image value mappings.
//
// The primary use case is processing cg.json files embedded in chart APKs,
// which define how container image references map to Helm values.yaml paths.
//
// # Template Format
//
// Image values use ${...} markers for variable substitution:
//
//	{
//	  "images": {
//	    "nginx": {
//	      "values": {
//	        "image": {
//	          "registry": "${registry}",
//	          "repository": "${repo}",
//	          "tag": "${tag}"
//	        }
//	      }
//	    }
//	  }
//	}
//
// Available variables: registry, repo, registry_repo, tag, digest, pseudo_tag, ref.
//
// # Basic Usage
//
// Parse a mapping and resolve values:
//
//	m, err := images.Parse(reader)
//	ref, err := images.NewRef("cgr.dev/chainguard/nginx:latest@sha256:...")
//	refs := map[string]images.OCIRef{"nginx": ref}
//	values, err := m.Walk(images.Resolve(refs))
//	merged, err := values.Merge()
//
// # Custom Callbacks
//
// For advanced use cases, provide a custom WalkFunc:
//
//	values, err := m.Walk(func(imageID string, tokens images.TokenList) (any, error) {
//	    return tokens.Map(func(f images.RefField) any {
//	        // Return custom structs, etc.
//	    }), nil
//	})
package images

// ChainguardChartMetadataFilename is the name of the Chainguard metadata file embedded in chart APKs.
const ChainguardChartMetadataFilename = "cg.json"
