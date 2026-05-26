/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"fmt"
	"strings"
	"testing"

	"chainguard.dev/sdk/helm/images"
	helmv1 "chainguard.dev/sdk/helm/v1"
	"github.com/google/go-cmp/cmp"
)

const (
	digestA = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	digestB = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	digestC = "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	digestD = "sha256:dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	digestE = "sha256:eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

func echoResolve(_ string, ref *helmv1.ChartImage) (string, string, error) {
	fullRef := "test.registry/" + ref.RepoName
	if ref.Tag != "" {
		fullRef += ":" + ref.Tag
	}
	if ref.Digest != "" {
		fullRef += "@" + ref.Digest
	}
	return fullRef, ref.Digest, nil
}

func imageTemplate(imageID string, values map[string]any) *images.Mapping {
	return &images.Mapping{
		Images: map[string]*images.Image{
			imageID: {Values: values},
		},
	}
}

func registryDigestTemplate(imageID string, prefix ...string) *images.Mapping {
	v := map[string]any{"registry": "${registry}", "digest": "${digest}"}
	for i := len(prefix) - 1; i >= 0; i-- {
		v = map[string]any{prefix[i]: v}
	}
	return imageTemplate(imageID, v)
}

func TestPatchChartImages(t *testing.T) {
	tests := []struct {
		name       string
		valuesYAML string
		ci         *helmv1.ChartImages
		resolve    ResolveFunc
		want       string // expected raw YAML output; empty means unchanged from input
		wantErr    string
	}{
		{
			name:       "nil ChartImages passes through unchanged",
			valuesYAML: `key: value` + "\n",
		},
		{
			name: "root-only replaces existing paths",
			valuesYAML: `
image:
  registry: cgr.dev
  repository: chainguard/nginx
  digest: placeholder
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Tag: "latest", Digest: digestA}},
				Template: imageTemplate("nginx", map[string]any{"image": map[string]any{"registry": "${registry}", "repository": "${repo}", "digest": "${digest}"}}),
			},
			want: `image:
  registry: test.registry
  repository: nginx
  digest: ` + digestA + "\n",
		},
		{
			name: "root + single subchart",
			valuesYAML: `
image:
  registry: cgr.dev
  repository: chainguard/nginx
  digest: placeholder
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Tag: "latest", Digest: digestA}},
				Template: imageTemplate("nginx", map[string]any{"image": map[string]any{"registry": "${registry}", "repository": "${repo}", "digest": "${digest}"}}),
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Tag: "7", Digest: digestB}},
						Template: registryDigestTemplate("redis", "image"),
					},
				},
			},
			want: `image:
  registry: test.registry
  repository: nginx
  digest: ` + digestA + `
redis:
  image:
    digest: ` + digestB + `
    registry: test.registry
`,
		},
		{
			name:       "subchart-only adds new keys",
			valuesYAML: `someKey: someValue` + "\n",
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Tag: "7", Digest: digestB}},
						Template: registryDigestTemplate("redis", "image"),
					},
				},
			},
			want: `someKey: someValue
redis:
  image:
    digest: ` + digestB + `
    registry: test.registry
`,
		},
		{
			name:       "nested subcharts two levels deep",
			valuesYAML: `placeholder: true` + "\n",
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Tag: "7", Digest: digestA}},
						Template: registryDigestTemplate("redis", "image"),
						Subcharts: map[string]*helmv1.ChartImages{
							"sentinel": {
								Refs:     map[string]*helmv1.ChartImage{"sentinel": {RepoName: "sentinel", Tag: "7", Digest: digestB}},
								Template: registryDigestTemplate("sentinel", "image"),
							},
						},
					},
				},
			},
			want: `placeholder: true
redis:
  image:
    digest: ` + digestA + `
    registry: test.registry
  sentinel:
    image:
      digest: ` + digestB + `
      registry: test.registry
`,
		},
		{
			name: "duplicate image ID across root and subchart resolves independently",
			valuesYAML: `
init:
  image:
    registry: cgr.dev
    digest: placeholder
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"curl": {RepoName: "curl", Tag: "latest", Digest: digestC}},
				Template: registryDigestTemplate("curl", "init", "image"),
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"curl": {RepoName: "curl", Tag: "latest-dev", Digest: digestD}},
						Template: registryDigestTemplate("curl", "init", "image"),
					},
				},
			},
			want: `init:
  image:
    registry: test.registry
    digest: ` + digestC + `
redis:
  init:
    image:
      digest: ` + digestD + `
      registry: test.registry
`,
		},
		{
			name: "same image ID in multiple sibling subcharts",
			valuesYAML: `
root:
  image:
    registry: cgr.dev
    digest: placeholder
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"helper": {RepoName: "helper", Digest: digestA}},
				Template: registryDigestTemplate("helper", "root", "image"),
				Subcharts: map[string]*helmv1.ChartImages{
					"alpha":   {Refs: map[string]*helmv1.ChartImage{"curl": {RepoName: "curl", Tag: "latest-dev", Digest: digestB}}, Template: registryDigestTemplate("curl", "initImage")},
					"bravo":   {Refs: map[string]*helmv1.ChartImage{"curl": {RepoName: "curl", Tag: "latest-dev", Digest: digestC}}, Template: registryDigestTemplate("curl", "initImage")},
					"charlie": {Refs: map[string]*helmv1.ChartImage{"curl": {RepoName: "curl", Tag: "latest-dev", Digest: digestD}}, Template: registryDigestTemplate("curl", "initImage")},
				},
			},
			want: `root:
  image:
    registry: test.registry
    digest: ` + digestA + `
alpha:
  initImage:
    digest: ` + digestB + `
    registry: test.registry
bravo:
  initImage:
    digest: ` + digestC + `
    registry: test.registry
charlie:
  initImage:
    digest: ` + digestD + `
    registry: test.registry
`,
		},
		{
			name: "multiple root images",
			valuesYAML: `
app:
  image:
    registry: cgr.dev
    digest: placeholder
sidecar:
  image:
    registry: cgr.dev
    digest: placeholder
`[1:],
			ci: &helmv1.ChartImages{
				Refs: map[string]*helmv1.ChartImage{
					"app":     {RepoName: "app", Digest: digestA},
					"sidecar": {RepoName: "sidecar", Digest: digestB},
				},
				Template: &images.Mapping{Images: map[string]*images.Image{
					"app":     {Values: map[string]any{"app": map[string]any{"image": map[string]any{"registry": "${registry}", "digest": "${digest}"}}}},
					"sidecar": {Values: map[string]any{"sidecar": map[string]any{"image": map[string]any{"registry": "${registry}", "digest": "${digest}"}}}},
				}},
			},
			want: `app:
  image:
    registry: test.registry
    digest: ` + digestA + `
sidecar:
  image:
    registry: test.registry
    digest: ` + digestB + "\n",
		},
		{
			name:       "structural subchart with only nested subcharts",
			valuesYAML: `placeholder: true` + "\n",
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"infra": {
						Subcharts: map[string]*helmv1.ChartImages{
							"redis": {
								Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Digest: digestA}},
								Template: registryDigestTemplate("redis", "image"),
							},
						},
					},
				},
			},
			want: `placeholder: true
infra:
  redis:
    image:
      digest: ` + digestA + `
      registry: test.registry
`,
		},
		{
			name: "subchart replaces existing and preserves other values",
			valuesYAML: `
redis:
  replicas: 3
`[1:],
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Digest: digestB}},
						Template: imageTemplate("redis", map[string]any{"image": map[string]any{"digest": "${digest}"}}),
					},
				},
			},
			want: `redis:
  replicas: 3
  image:
    digest: ` + digestB + "\n",
		},
		{
			name: "root template with missing path errors strictly",
			valuesYAML: `
other:
  key: value
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Tag: "latest", Digest: digestA}},
				Template: imageTemplate("nginx", map[string]any{"image": map[string]any{"registry": "${registry}"}}),
			},
			wantErr: "does not exist",
		},
		{
			name: "resolve callback error propagates",
			valuesYAML: `
image:
  registry: cgr.dev
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Digest: digestA}},
				Template: imageTemplate("nginx", map[string]any{"image": map[string]any{"registry": "${registry}"}}),
			},
			resolve: func(string, *helmv1.ChartImage) (string, string, error) {
				return "", "", fmt.Errorf("repo not found")
			},
			wantErr: "repo not found",
		},
		{
			name:       "sibling subcharts at depth 2 resolve independently",
			valuesYAML: `placeholder: true` + "\n",
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Digest: digestA}},
						Template: registryDigestTemplate("redis", "image"),
						Subcharts: map[string]*helmv1.ChartImages{
							"sentinel": {Refs: map[string]*helmv1.ChartImage{"sentinel": {RepoName: "sentinel", Digest: digestB}}, Template: registryDigestTemplate("sentinel", "image")},
							"metrics":  {Refs: map[string]*helmv1.ChartImage{"metrics": {RepoName: "redis-exporter", Digest: digestE}}, Template: registryDigestTemplate("metrics", "image")},
						},
					},
				},
			},
			want: `placeholder: true
redis:
  image:
    digest: ` + digestA + `
    registry: test.registry
  metrics:
    image:
      digest: ` + digestE + `
      registry: test.registry
  sentinel:
    image:
      digest: ` + digestB + `
      registry: test.registry
`,
		},
		{
			name: "preserves comments and formatting",
			valuesYAML: `
# Chart configuration
image:
  # Registry to pull from
  registry: cgr.dev
  repository: chainguard/nginx
  digest: placeholder
`[1:],
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Tag: "latest", Digest: digestA}},
				Template: imageTemplate("nginx", map[string]any{"image": map[string]any{"registry": "${registry}", "repository": "${repo}", "digest": "${digest}"}}),
			},
			want: `# Chart configuration
image:
  # Registry to pull from
  registry: test.registry
  repository: nginx
  digest: ` + digestA + "\n",
		},
		{
			name: "subchart patch only adds image fields",
			valuesYAML: `
# Parent chart values
nameOverride: ""
commonLabels: {}
backend:
  enabled: true
  config:
    replicas: 1
    retention: 120h
    resources: {}
# Subchart defaults
dashboard:
  enabled: true
  replicas: 1
  password: changeme
  extras:
    logging:
      enabled: true
`[1:],
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"dashboard": {
						Refs:     map[string]*helmv1.ChartImage{"dashboard": {RepoName: "dashboard", Digest: digestA}},
						Template: registryDigestTemplate("dashboard", "image"),
					},
				},
			},
			want: `# Parent chart values
nameOverride: ""
commonLabels: {}
backend:
  enabled: true
  config:
    replicas: 1
    retention: 120h
    resources: {}
# Subchart defaults
dashboard:
  enabled: true
  replicas: 1
  password: changeme
  extras:
    logging:
      enabled: true
  image:
    digest: ` + digestA + `
    registry: test.registry
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chartImg := createTestChart(t, "my-chart", tt.valuesYAML)
			resolve := tt.resolve
			if resolve == nil {
				resolve = echoResolve
			}

			patched, _, err := PatchChartImages(chartImg, tt.ci, resolve)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tt.wantErr)
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got: %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("PatchChartImages() error = %v", err)
			}

			got, err := ReadValues(patched)
			if err != nil {
				t.Fatalf("ReadValues: %v", err)
			}

			want := tt.want
			if want == "" {
				want = tt.valuesYAML
			}
			if diff := cmp.Diff(want, string(got)); diff != "" {
				t.Errorf("output (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPatchChartImages_Digests(t *testing.T) {
	tests := []struct {
		name        string
		ci          *helmv1.ChartImages
		wantDigests *helmv1.ChartDigests
	}{
		{
			name: "root only",
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Digest: digestA}},
				Template: registryDigestTemplate("nginx", "image"),
			},
			wantDigests: &helmv1.ChartDigests{Digests: map[string]string{"nginx": digestA}},
		},
		{
			name: "root + subchart + nested",
			ci: &helmv1.ChartImages{
				Refs:     map[string]*helmv1.ChartImage{"nginx": {RepoName: "nginx", Digest: digestA}},
				Template: registryDigestTemplate("nginx", "image"),
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Digest: digestB}},
						Template: registryDigestTemplate("redis", "image"),
						Subcharts: map[string]*helmv1.ChartImages{
							"sentinel": {
								Refs:     map[string]*helmv1.ChartImage{"sentinel": {RepoName: "sentinel", Digest: digestC}},
								Template: registryDigestTemplate("sentinel", "image"),
							},
						},
					},
				},
			},
			wantDigests: &helmv1.ChartDigests{
				Digests: map[string]string{"nginx": digestA},
				Subcharts: map[string]*helmv1.ChartDigests{
					"redis": {
						Digests: map[string]string{"redis": digestB},
						Subcharts: map[string]*helmv1.ChartDigests{
							"sentinel": {Digests: map[string]string{"sentinel": digestC}},
						},
					},
				},
			},
		},
		{
			name: "subchart-only with no root refs",
			ci: &helmv1.ChartImages{
				Subcharts: map[string]*helmv1.ChartImages{
					"redis": {
						Refs:     map[string]*helmv1.ChartImage{"redis": {RepoName: "redis", Digest: digestA}},
						Template: registryDigestTemplate("redis", "image"),
					},
				},
			},
			wantDigests: &helmv1.ChartDigests{
				Subcharts: map[string]*helmv1.ChartDigests{
					"redis": {Digests: map[string]string{"redis": digestA}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chartImg := createTestChart(t, "my-chart", `
image:
  registry: cgr.dev
  digest: placeholder
`[1:])
			_, digests, err := PatchChartImages(chartImg, tt.ci, echoResolve)
			if err != nil {
				t.Fatalf("PatchChartImages: %v", err)
			}
			if diff := cmp.Diff(tt.wantDigests, digests); diff != "" {
				t.Errorf("digests (-want +got):\n%s", diff)
			}
		})
	}
}
