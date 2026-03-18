/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"fmt"
	"strings"
	"testing"

	"chainguard.dev/sdk/helm/images"
	"github.com/google/go-cmp/cmp"
)

func TestLock_WithChartDigests(t *testing.T) {
	originalDigest := "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	resolvedDigest := "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"

	tests := []struct {
		name     string
		lock     *Lock
		chartRef string
		digests  *ChartDigests
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
			digests:  &ChartDigests{Digests: map[string]string{"nginx": resolvedDigest}},
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
			got := tt.lock.WithChartDigests(tt.chartRef, tt.digests)

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

func TestChartImages_Walk(t *testing.T) {
	registryOnly := func(_ string, ref *ChartImage, tokens images.TokenList) (any, error) {
		for _, tok := range tokens {
			f, ok := tok.(images.RefField)
			if !ok {
				continue
			}
			switch f {
			case images.Registry:
				return "override.dev", nil
			case images.Repo:
				return ref.RepoName, nil
			}
		}
		return nil, nil
	}

	tests := []struct {
		name    string
		ci      *ChartImages
		fn      WalkFunc
		want    map[string]any
		wantErr string
	}{
		{
			name: "nil ChartImages",
			ci:   nil,
			fn:   registryOnly,
		},
		{
			name: "single image",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"nginx": {RepoName: "nginx", Tag: "latest"},
				},
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"nginx": {Values: map[string]any{
							"image": map[string]any{"registry": "${registry}"},
						}},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"image": map[string]any{"registry": "override.dev"},
			},
		},
		{
			name: "subchart values nested under dependency key",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"main": {RepoName: "main-app"},
				},
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"main": {Values: map[string]any{
							"image": map[string]any{"registry": "${registry}"},
						}},
					},
				},
				Subcharts: map[string]*ChartImages{
					"redis": {
						Refs: map[string]*ChartImage{
							"server": {RepoName: "redis-server"},
						},
						Template: &images.Mapping{
							Images: map[string]*images.Image{
								"server": {Values: map[string]any{
									"image": map[string]any{"registry": "${registry}"},
								}},
							},
						},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"image": map[string]any{"registry": "override.dev"},
				"redis": map[string]any{
					"image": map[string]any{"registry": "override.dev"},
				},
			},
		},
		{
			name: "nested subcharts",
			ci: &ChartImages{
				Subcharts: map[string]*ChartImages{
					"fulcio": {
						Refs: map[string]*ChartImage{
							"server": {RepoName: "fulcio-server"},
						},
						Template: &images.Mapping{
							Images: map[string]*images.Image{
								"server": {Values: map[string]any{
									"registry": "${registry}",
								}},
							},
						},
						Subcharts: map[string]*ChartImages{
							"ctlog": {
								Refs: map[string]*ChartImage{
									"ct": {RepoName: "ct-server"},
								},
								Template: &images.Mapping{
									Images: map[string]*images.Image{
										"ct": {Values: map[string]any{
											"registry": "${registry}",
										}},
									},
								},
							},
						},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"fulcio": map[string]any{
					"registry": "override.dev",
					"ctlog": map[string]any{
						"registry": "override.dev",
					},
				},
			},
		},
		{
			name: "ref passed to callback",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"app": {RepoName: "my-app"},
				},
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"app": {Values: map[string]any{
							"repo": "${repo}",
						}},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"repo": "my-app",
			},
		},
		{
			name: "unknown image errors",
			ci: &ChartImages{
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"orphan": {Values: map[string]any{
							"registry": "${registry}",
						}},
					},
				},
			},
			fn:      registryOnly,
			wantErr: `image "orphan": not found in refs`,
		},
		{
			name:    "nil callback",
			ci:      &ChartImages{},
			fn:      nil,
			wantErr: "callback function is nil",
		},
		{
			name: "nil template with subcharts",
			ci: &ChartImages{
				Subcharts: map[string]*ChartImages{
					"redis": {
						Refs:     map[string]*ChartImage{"server": {RepoName: "redis"}},
						Template: &images.Mapping{Images: map[string]*images.Image{"server": {Values: map[string]any{"registry": "${registry}"}}}},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"redis": map[string]any{"registry": "override.dev"},
			},
		},
		{
			name: "callback returns nil excludes field",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"app": {RepoName: "app"},
				},
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"app": {Values: map[string]any{
							"registry": "${registry}",
							"tag":      "${tag}",
						}},
					},
				},
			},
			fn: func(_ string, _ *ChartImage, tokens images.TokenList) (any, error) {
				for _, tok := range tokens {
					if f, ok := tok.(images.RefField); ok && f == images.Registry {
						return "override.dev", nil
					}
				}
				return nil, nil
			},
			want: map[string]any{
				"registry": "override.dev",
			},
		},
		{
			name: "callback error propagated",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{"app": {RepoName: "app"}},
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"app": {Values: map[string]any{"k": "${registry}"}},
					},
				},
			},
			fn: func(_ string, _ *ChartImage, _ images.TokenList) (any, error) {
				return nil, fmt.Errorf("boom")
			},
			wantErr: "boom",
		},
		{
			name: "subchart callback error includes subchart name",
			ci: &ChartImages{
				Subcharts: map[string]*ChartImages{
					"redis": {
						Refs: map[string]*ChartImage{"server": {RepoName: "redis"}},
						Template: &images.Mapping{
							Images: map[string]*images.Image{
								"server": {Values: map[string]any{"k": "${registry}"}},
							},
						},
					},
				},
			},
			fn: func(_ string, _ *ChartImage, _ images.TokenList) (any, error) {
				return nil, fmt.Errorf("boom")
			},
			wantErr: `subchart "redis"`,
		},
		{
			name: "empty subchart produces no output",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{"app": {RepoName: "app"}},
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"app": {Values: map[string]any{"registry": "${registry}"}},
					},
				},
				Subcharts: map[string]*ChartImages{
					"empty": {
						Template: &images.Mapping{},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"registry": "override.dev",
			},
		},
		{
			name: "multiple subcharts sorted deterministically",
			ci: &ChartImages{
				Subcharts: map[string]*ChartImages{
					"bravo": {
						Refs:     map[string]*ChartImage{"b": {RepoName: "bravo"}},
						Template: &images.Mapping{Images: map[string]*images.Image{"b": {Values: map[string]any{"registry": "${registry}"}}}},
					},
					"alpha": {
						Refs:     map[string]*ChartImage{"a": {RepoName: "alpha"}},
						Template: &images.Mapping{Images: map[string]*images.Image{"a": {Values: map[string]any{"registry": "${registry}"}}}},
					},
				},
			},
			fn: registryOnly,
			want: map[string]any{
				"alpha": map[string]any{"registry": "override.dev"},
				"bravo": map[string]any{"registry": "override.dev"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.ci.Walk(tc.fn)
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got: %v", tc.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Visit: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
