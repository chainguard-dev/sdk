/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package chart provides utilities for reading and modifying Chainguard Helm
charts stored as OCI artifacts or on the filesystem.

Chainguard Helm charts embed metadata that describes how container image
references map to values.yaml paths. This package uses that metadata to read,
inspect, and patch charts without requiring knowledge of the underlying
template format.

# Inspecting a Chart Package

To read chart metadata and image mappings from a chart filesystem (typically
an expanded Chainguard chart APK), use [Read]:

	chartData, err := chart.Read(expandedAPK.TarFS)
	if err != nil {
		return err
	}

	// chartData.Meta has Chart.yaml fields (name, version, dependencies, etc.)
	// chartData.Mapping has the parsed image mapping template
	// chartData.Subcharts has recursive subchart data, keyed by dependency alias

# Extracting values.yaml from an OCI Artifact

	img, err := remote.Image(ref)
	if err != nil {
		return err
	}

	values, err := chart.ReadValues(img)
	if err != nil {
		return err
	}
	// values is nil if the chart has no values.yaml

# Patching Image References into a Chart

To replace image references in a chart's values.yaml using a mapping template
and a set of concrete refs, use [ReplaceValues]. It returns a new OCI artifact
with the updated values.yaml — all other files, annotations, and metadata are
preserved:

	refs := map[string]string{
		"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abc123...",
	}
	patched, err := chart.ReplaceValues(img, mapping, refs)
	if err != nil {
		return err
	}
	if err := remote.Write(dstRef, patched); err != nil {
		return err
	}

# Detecting Helm Charts

Use the [MediaType] constant to identify Helm chart content layers when
inspecting OCI manifests:

	mt, _ := layers[0].MediaType()
	if mt == chart.MediaType {
		// this is a Helm chart
	}
*/
package chart
