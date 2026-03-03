/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"testing"

	"chainguard.dev/sdk/helm/images"
)

func TestLock_WithResolvedDigests(t *testing.T) {
	originalDigest := "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	resolvedDigest := "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	tests := []struct {
		name     string
		lock     *Lock
		chartRef string
		digests  map[string]string
		want     *Lock
	}{
		{
			name: "resolves image digests",
			lock: &Lock{
				Chart: &Chart{Package: "nginx", Ref: "cgr.dev/chainguard/charts/nginx@sha256:original"},
				Images: &ChartImages{
					Refs: map[string]*ChartImage{
						"nginx": {RepoName: "nginx", Tag: "latest", Digest: originalDigest},
					},
					Template: &images.Mapping{},
				},
			},
			chartRef: "cgr.dev/customer/charts/nginx@sha256:newchart",
			digests:  map[string]string{"nginx": resolvedDigest},
			want: &Lock{
				Chart: &Chart{Package: "nginx", Ref: "cgr.dev/customer/charts/nginx@sha256:newchart"},
				Images: &ChartImages{
					Refs: map[string]*ChartImage{
						"nginx": {RepoName: "nginx", Tag: "latest", Digest: resolvedDigest},
					},
					Template: &images.Mapping{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.lock.WithResolvedDigests(tt.chartRef, tt.digests)

			// Check Chart
			if got.Chart.Package != tt.want.Chart.Package {
				t.Errorf("Chart.Package = %q, want %q", got.Chart.Package, tt.want.Chart.Package)
			}
			if got.Chart.Ref != tt.want.Chart.Ref {
				t.Errorf("Chart.Ref = %q, want %q", got.Chart.Ref, tt.want.Chart.Ref)
			}

			// Check Images.Refs
			if len(got.Images.Refs) != len(tt.want.Images.Refs) {
				t.Errorf("len(Images.Refs) = %d, want %d", len(got.Images.Refs), len(tt.want.Images.Refs))
			}
			for id, wantImg := range tt.want.Images.Refs {
				gotImg, ok := got.Images.Refs[id]
				if !ok {
					t.Errorf("missing image %q", id)
					continue
				}
				if gotImg.RepoName != wantImg.RepoName {
					t.Errorf("Images.Refs[%q].RepoName = %q, want %q", id, gotImg.RepoName, wantImg.RepoName)
				}
				if gotImg.Tag != wantImg.Tag {
					t.Errorf("Images.Refs[%q].Tag = %q, want %q", id, gotImg.Tag, wantImg.Tag)
				}
				if gotImg.Digest != wantImg.Digest {
					t.Errorf("Images.Refs[%q].Digest = %q, want %q", id, gotImg.Digest, wantImg.Digest)
				}
			}
		})
	}
}
