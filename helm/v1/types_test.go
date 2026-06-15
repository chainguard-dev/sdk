/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"fmt"
	"slices"
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

func TestChartImage_Reference(t *testing.T) {
	tests := []struct {
		name   string
		img    *ChartImage
		prefix string
		want   string
	}{
		{
			name:   "repo tag and digest",
			img:    &ChartImage{RepoName: "nginx", Tag: "latest", Digest: "sha256:a"},
			prefix: "cgr.dev/test-org",
			want:   "cgr.dev/test-org/nginx:latest@sha256:a",
		},
		{
			name:   "repo only",
			img:    &ChartImage{RepoName: "plain"},
			prefix: "cgr.dev/test-org",
			want:   "cgr.dev/test-org/plain",
		},
		{
			name:   "tag omitted",
			img:    &ChartImage{RepoName: "digestonly", Digest: "sha256:a"},
			prefix: "cgr.dev/test-org",
			want:   "cgr.dev/test-org/digestonly@sha256:a",
		},
		{
			name:   "digest omitted",
			img:    &ChartImage{RepoName: "tagged", Tag: "v1"},
			prefix: "cgr.dev/test-org",
			want:   "cgr.dev/test-org/tagged:v1",
		},
		{
			name: "empty prefix yields bare coordinates",
			img:  &ChartImage{RepoName: "nginx", Tag: "latest", Digest: "sha256:a"},
			want: "nginx:latest@sha256:a",
		},
		{
			name:   "bare registry prefix",
			img:    &ChartImage{RepoName: "nginx", Tag: "latest", Digest: "sha256:a"},
			prefix: "myregistry.internal",
			want:   "myregistry.internal/nginx:latest@sha256:a",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.img.Reference(tc.prefix); got != tc.want {
				t.Errorf("Reference(%q) = %q, want %q", tc.prefix, got, tc.want)
			}
		})
	}
}

func TestChartImages_Images(t *testing.T) {
	tests := []struct {
		name string
		ci   *ChartImages
		// want is the sorted list of "repoName:tag@digest" keys yielded,
		// including duplicates (Images does not deduplicate).
		want []string
	}{
		{
			name: "nil receiver yields nothing",
			ci:   nil,
		},
		{
			name: "root refs only",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"nginx": {RepoName: "nginx", Tag: "latest", Digest: "sha256:a"},
					"redis": {RepoName: "redis", Tag: "7.0", Digest: "sha256:b"},
				},
			},
			want: []string{"nginx:latest@sha256:a", "redis:7.0@sha256:b"},
		},
		{
			name: "nil entries are skipped",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"good":  {RepoName: "nginx", Tag: "latest", Digest: "sha256:a"},
					"ghost": nil,
				},
			},
			want: []string{"nginx:latest@sha256:a"},
		},
		{
			name: "subcharts walked recursively",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"main": {RepoName: "main", Tag: "1.0", Digest: "sha256:a"},
				},
				Subcharts: map[string]*ChartImages{
					"redis": {
						Refs: map[string]*ChartImage{
							"server": {RepoName: "redis", Tag: "7.0", Digest: "sha256:b"},
						},
						Subcharts: map[string]*ChartImages{
							"exporter": {
								Refs: map[string]*ChartImage{
									"exporter": {RepoName: "redis-exporter", Tag: "v1", Digest: "sha256:c"},
								},
							},
						},
					},
				},
			},
			want: []string{"main:1.0@sha256:a", "redis-exporter:v1@sha256:c", "redis:7.0@sha256:b"},
		},
		{
			name: "duplicates across subcharts are preserved",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"shared": {RepoName: "shared", Tag: "1.0", Digest: "sha256:a"},
				},
				Subcharts: map[string]*ChartImages{
					"sub": {
						Refs: map[string]*ChartImage{
							"shared": {RepoName: "shared", Tag: "1.0", Digest: "sha256:a"},
						},
					},
				},
			},
			want: []string{"shared:1.0@sha256:a", "shared:1.0@sha256:a"},
		},
		{
			name: "nil subchart entry does not panic",
			ci: &ChartImages{
				Refs: map[string]*ChartImage{
					"nginx": {RepoName: "nginx", Tag: "latest", Digest: "sha256:a"},
				},
				Subcharts: map[string]*ChartImages{
					"empty": nil,
				},
			},
			want: []string{"nginx:latest@sha256:a"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var got []string
			for img := range tc.ci.Images() {
				got = append(got, fmt.Sprintf("%s:%s@%s", img.RepoName, img.Tag, img.Digest))
			}
			slices.Sort(got)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestChartImages_Images_EarlyTermination(t *testing.T) {
	ci := &ChartImages{
		Refs: map[string]*ChartImage{
			"a": {RepoName: "a"},
			"b": {RepoName: "b"},
			"c": {RepoName: "c"},
		},
	}

	var count int
	for range ci.Images() {
		count++
		break
	}
	if count != 1 {
		t.Errorf("breaking after the first image visited %d images, want 1", count)
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

func TestIsOptionalImage(t *testing.T) {
	tests := []struct {
		name    string
		ci      *ChartImages
		imageID string
		want    bool
	}{
		{
			name:    "nil ChartImages",
			ci:      nil,
			imageID: "nginx",
			want:    false,
		},
		{
			name:    "nil template",
			ci:      &ChartImages{},
			imageID: "nginx",
			want:    false,
		},
		{
			name: "image not in template",
			ci: &ChartImages{
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"nginx": {Values: map[string]any{"image": "${ref}"}},
					},
				},
			},
			imageID: "redis",
			want:    false,
		},
		{
			name: "required image",
			ci: &ChartImages{
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"nginx": {Requirement: images.Required, Values: map[string]any{"image": "${ref}"}},
					},
				},
			},
			imageID: "nginx",
			want:    false,
		},
		{
			name: "unspecified requirement",
			ci: &ChartImages{
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"nginx": {Values: map[string]any{"image": "${ref}"}},
					},
				},
			},
			imageID: "nginx",
			want:    false,
		},
		{
			name: "optional image",
			ci: &ChartImages{
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"sidecar": {Requirement: images.Optional, Values: map[string]any{"image": "${ref}"}},
					},
				},
			},
			imageID: "sidecar",
			want:    true,
		},
		{
			name: "does not check subcharts",
			ci: &ChartImages{
				Template: &images.Mapping{
					Images: map[string]*images.Image{
						"nginx": {Requirement: images.Required, Values: map[string]any{"image": "${ref}"}},
					},
				},
				Subcharts: map[string]*ChartImages{
					"redis": {
						Template: &images.Mapping{
							Images: map[string]*images.Image{
								"exporter": {Requirement: images.Optional, Values: map[string]any{"image": "${ref}"}},
							},
						},
					},
				},
			},
			imageID: "exporter",
			want:    false, // only checks current level, not subcharts
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.ci.IsOptionalImage(tc.imageID)
			if got != tc.want {
				t.Errorf("IsOptionalImage(%q) = %v, want %v", tc.imageID, got, tc.want)
			}
		})
	}
}
