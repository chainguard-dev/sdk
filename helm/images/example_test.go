/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package images_test

import (
	"fmt"
	"strings"

	"chainguard.dev/sdk/helm/images"
)

// ExampleParse demonstrates parsing a chart image mapping from JSON.
func ExampleParse() {
	const mapping = `{
		"images": {
			"nginx": {
				"values": {
					"image.repository": "${registry}/${repo}",
					"image.tag": "${tag}"
				}
			}
		}
	}`
	m, err := images.Parse(strings.NewReader(mapping))
	fmt.Println(err)
	fmt.Println(len(m.Images))
	// Output:
	// <nil>
	// 1
}

// ExampleChainguardChartMetadataFilename demonstrates the metadata filename constant.
func ExampleChainguardChartMetadataFilename() {
	fmt.Println(images.ChainguardChartMetadataFilename)
	// Output:
	// cg.json
}
