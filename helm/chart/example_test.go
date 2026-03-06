/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart_test

import (
	"fmt"
	"log"

	"chainguard.dev/sdk/helm/chart"
	"chainguard.dev/sdk/helm/images"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

// ExampleReadValues demonstrates reading the values.yaml file from a Helm chart OCI image.
func ExampleReadValues() {
	ref, err := name.ParseReference("cgr.dev/chainguard/helm-charts/nginx:latest")
	if err != nil {
		log.Fatal(err)
	}

	chartImage, err := remote.Image(ref)
	if err != nil {
		log.Fatal(err)
	}

	values, err := chart.ReadValues(chartImage)
	if err != nil {
		log.Fatal(err)
	}

	if values == nil {
		fmt.Println("Chart has no values.yaml")
		return
	}

	fmt.Printf("values.yaml content:\n%s\n", values)
}

// ExampleReplaceValues demonstrates replacing values.yaml with image reference resolution.
func ExampleReplaceValues() {
	ref, err := name.ParseReference("cgr.dev/chainguard/helm-charts/nginx:latest")
	if err != nil {
		log.Fatal(err)
	}

	chartImage, err := remote.Image(ref)
	if err != nil {
		log.Fatal(err)
	}

	// Define the mapping template for image references
	mapping := &images.Mapping{
		Images: map[string]*images.Image{
			"nginx": {
				Values: map[string]any{
					"image": map[string]any{
						"registry":   "${registry}",
						"repository": "${repo}",
						"tag":        "${tag}",
					},
				},
			},
		},
	}

	// Provide the actual image references to resolve
	refs := map[string]string{
		"nginx": "cgr.dev/chainguard/nginx:latest@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}

	// Replace values.yaml with resolved image references
	patched, err := chart.ReplaceValues(chartImage, mapping, refs)
	if err != nil {
		log.Fatal(err)
	}

	// Write the patched chart to a new location
	dstRef, err := name.ParseReference("cgr.dev/my-org/nginx-chart:patched")
	if err != nil {
		log.Fatal(err)
	}

	if err := remote.Write(dstRef, patched); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chart values.yaml updated and pushed successfully")
}

// ExampleMediaType demonstrates using the Helm chart media type constant.
func Example_mediaType() {
	// The MediaType constant identifies Helm chart content layers in OCI images
	fmt.Println(chart.MediaType)
	// Output: application/vnd.cncf.helm.chart.content.v1.tar+gzip
}
