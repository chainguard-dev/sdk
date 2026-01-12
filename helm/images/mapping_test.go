/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package images

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func mustNewRef(reference string) OCIRef {
	ref, err := NewRef(reference)
	if err != nil {
		panic(err)
	}
	return ref
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Mapping
		wantErr string
	}{
		// Valid cases
		{
			name: "valid single image",
			input: `{
				"images": {
					"nginx": {
						"values": {
							"image": "${registry}"
						}
					}
				}
			}`,
			want: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": "${registry}",
					}},
				},
			},
		},
		{
			name: "valid multiple images",
			input: `{
				"images": {
					"nginx": {"values": {"image": "${registry}/nginx"}},
					"redis": {"values": {"image": "${registry}/redis"}}
				}
			}`,
			want: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{"image": "${registry}/nginx"}},
					"redis": {Values: map[string]any{"image": "${registry}/redis"}},
				},
			},
		},
		{
			name: "valid nested values",
			input: `{
				"images": {
					"nginx": {
						"values": {
							"image": {
								"registry": "${registry}",
								"tag": "${tag}"
							}
						}
					}
				}
			}`,
			want: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{
							"registry": "${registry}",
							"tag":      "${tag}",
						},
					}},
				},
			},
		},
		{
			name: "valid with non-string values",
			input: `{
				"images": {
					"nginx": {
						"values": {
							"replicas": 3,
							"enabled": true,
							"image": "${registry}"
						}
					}
				}
			}`,
			want: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"replicas": float64(3),
						"enabled":  true,
						"image":    "${registry}",
					}},
				},
			},
		},

		// Invalid JSON
		{
			name:    "invalid json",
			input:   `{not valid json}`,
			wantErr: "parsing image mapping",
		},

		// Missing/empty images
		{
			name:    "missing images field",
			input:   `{}`,
			wantErr: "missing or empty 'images'",
		},
		{
			name:    "empty images",
			input:   `{"images": {}}`,
			wantErr: "missing or empty 'images'",
		},
		{
			name:    "null images",
			input:   `{"images": null}`,
			wantErr: "missing or empty 'images'",
		},

		// Invalid image definitions
		{
			name:    "null image definition",
			input:   `{"images": {"nginx": null}}`,
			wantErr: "image \"nginx\": nil definition",
		},
		{
			name:    "missing values field",
			input:   `{"images": {"nginx": {}}}`,
			wantErr: "image \"nginx\": missing 'values'",
		},
		{
			name:    "null values field",
			input:   `{"images": {"nginx": {"values": null}}}`,
			wantErr: "image \"nginx\": missing 'values'",
		},

		// Invalid markers
		{
			name: "unknown marker",
			input: `{
				"images": {
					"nginx": {"values": {"image": "${unknown}"}}
				}
			}`,
			wantErr: "unknown field",
		},
		{
			name: "unclosed marker",
			input: `{
				"images": {
					"nginx": {"values": {"image": "${registry"}}
				}
			}`,
			wantErr: "unclosed",
		},
		{
			name: "invalid marker in nested value",
			input: `{
				"images": {
					"nginx": {
						"values": {
							"image": {
								"registry": "${invalid}"
							}
						}
					}
				}
			}`,
			wantErr: "unknown field",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tc.input))

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
				t.Fatalf("Parse: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWalkNilCallback(t *testing.T) {
	m := &Mapping{Images: map[string]*Image{
		"test": {Values: map[string]any{"key": "${registry}"}},
	}}
	_, err := m.Walk(nil)
	if err == nil {
		t.Fatal("expected error for nil callback")
	}
}

func TestWalk(t *testing.T) {
	const validDigest = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
	ref := OCIRef{
		Registry:     "cgr.dev",
		Repo:         "chainguard/nginx",
		RegistryRepo: "cgr.dev/chainguard/nginx",
		Tag:          "latest",
		Digest:       validDigest,
		PseudoTag:    fmt.Sprintf("latest@%s", validDigest),
		FullRef:      fmt.Sprintf("cgr.dev/chainguard/nginx:latest@%s", validDigest),
	}

	tests := []struct {
		name     string
		mapping  *Mapping
		refs     map[string]OCIRef // used when callback is nil
		callback WalkFunc          // optional custom callback
		want     ImageValues
		wantErr  string
	}{
		// Lexer edge cases (tested through Walk)
		{
			name: "simple variable",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "${registry}"}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {"key": "cgr.dev"}},
		},
		{
			name: "literal only",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "hello"}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {"key": "hello"}},
		},
		{
			name: "interpolated",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "${registry}/${repo}:${tag}"}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {"key": "cgr.dev/chainguard/nginx:latest"}},
		},
		{
			name: "all variables",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "${registry} ${repo} ${registry_repo} ${tag} ${digest} ${pseudo_tag} ${ref}"}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {"key": fmt.Sprintf("cgr.dev chainguard/nginx cgr.dev/chainguard/nginx latest %s latest@%s cgr.dev/chainguard/nginx:latest@%s", validDigest, validDigest, validDigest)}},
		},
		{
			name: "unknown variable",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "${unknown}"}},
			}},
			refs:    map[string]OCIRef{"test": ref},
			wantErr: "unknown field",
		},
		{
			name: "unclosed marker",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "${registry"}},
			}},
			refs:    map[string]OCIRef{"test": ref},
			wantErr: "unclosed",
		},
		{
			name: "dollar without brace is literal",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "$100 price"}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {"key": "$100 price"}},
		},

		// Walk-specific: nested structures
		{
			name: "nested map",
			mapping: &Mapping{Images: map[string]*Image{
				"nginx": {Values: map[string]any{
					"image": map[string]any{
						"registry": "${registry}",
						"repo":     "${repo}",
					},
				}},
			}},
			refs: map[string]OCIRef{"nginx": ref},
			want: ImageValues{"nginx": {
				"image": map[string]any{
					"registry": "cgr.dev",
					"repo":     "chainguard/nginx",
				},
			}},
		},
		{
			name: "deeply nested",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{
					"l1": map[string]any{
						"l2": map[string]any{
							"l3": "${registry}",
						},
					},
				}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {
				"l1": map[string]any{
					"l2": map[string]any{
						"l3": "cgr.dev",
					},
				},
			}},
		},

		// Walk-specific: arrays
		{
			name: "array of strings",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{
					"images": []any{"${registry}", "${repo}"},
				}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {
				"images": []any{"cgr.dev", "chainguard/nginx"},
			}},
		},
		{
			name: "array of maps",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{
					"containers": []any{
						map[string]any{"image": "${registry}"},
						map[string]any{"image": "${repo}"},
					},
				}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {
				"containers": []any{
					map[string]any{"image": "cgr.dev"},
					map[string]any{"image": "chainguard/nginx"},
				},
			}},
		},

		// Walk-specific: non-string values preserved
		{
			name: "preserves non-strings",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{
					"replicas": float64(3),
					"enabled":  true,
					"config":   nil,
					"image":    "${registry}",
				}},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{"test": {
				"replicas": float64(3),
				"enabled":  true,
				"config":   nil,
				"image":    "cgr.dev",
			}},
		},

		// Walk-specific: multiple images
		{
			name: "multiple images",
			mapping: &Mapping{Images: map[string]*Image{
				"nginx": {Values: map[string]any{"image": "${registry}/nginx"}},
				"redis": {Values: map[string]any{"image": "${registry}/redis"}},
			}},
			refs: map[string]OCIRef{
				"nginx": {Registry: "cgr.dev"},
				"redis": {Registry: "docker.io"},
			},
			want: ImageValues{
				"nginx": {"image": "cgr.dev/nginx"},
				"redis": {"image": "docker.io/redis"},
			},
		},

		// Walk-specific: edge cases
		{
			name:    "nil mapping",
			mapping: nil,
			refs:    map[string]OCIRef{},
			want:    nil,
		},
		{
			name: "nil image",
			mapping: &Mapping{Images: map[string]*Image{
				"test": nil,
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{},
		},
		{
			name: "nil values",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: nil},
			}},
			refs: map[string]OCIRef{"test": ref},
			want: ImageValues{},
		},
		{
			name:    "empty images",
			mapping: &Mapping{Images: map[string]*Image{}},
			refs:    map[string]OCIRef{},
			want:    ImageValues{},
		},

		// Resolve with derived fields (catalog-syncer use case)
		{
			name: "NewRef computes derived fields",
			mapping: &Mapping{Images: map[string]*Image{
				"nginx": {Values: map[string]any{
					"image":        "${ref}",
					"registryRepo": "${registry_repo}",
					"pseudoTag":    "${pseudo_tag}",
				}},
			}},
			refs: map[string]OCIRef{
				"nginx": mustNewRef(fmt.Sprintf("customer.io/nginx:v1@%s", validDigest)),
			},
			want: ImageValues{"nginx": {
				"image":        fmt.Sprintf("customer.io/nginx:v1@%s", validDigest),
				"registryRepo": "customer.io/nginx",
				"pseudoTag":    fmt.Sprintf("v1@%s", validDigest),
			}},
		},

		// Custom callback tests
		{
			name: "custom callback with string transform",
			mapping: &Mapping{Images: map[string]*Image{
				"nginx": {Values: map[string]any{"key": "${registry}/${repo}"}},
			}},
			callback: func(_ string, tokens TokenList) (any, error) {
				parts := tokens.Map(func(f RefField) any {
					return "<<" + string(f) + ">>"
				})
				var sb strings.Builder
				for _, p := range parts {
					sb.WriteString(p.(string))
				}
				return sb.String(), nil
			},
			want: ImageValues{"nginx": {"key": "<<registry>>/<<repo>>"}},
		},
		{
			name: "custom callback with non-string return",
			mapping: &Mapping{Images: map[string]*Image{
				"test": {Values: map[string]any{"key": "${registry}/${repo}"}},
			}},
			callback: func(_ string, tokens TokenList) (any, error) {
				return tokens.Map(func(f RefField) any {
					return struct{ Name string }{Name: string(f)}
				}), nil
			},
			want: ImageValues{"test": {"key": []any{
				struct{ Name string }{Name: "registry"},
				"/",
				struct{ Name string }{Name: "repo"},
			}}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cb := tc.callback
			if cb == nil {
				cb = Resolve(tc.refs)
			}
			got, err := tc.mapping.Walk(cb)

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
				t.Fatalf("Walk: %v", err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMapping_Resolve(t *testing.T) {
	tests := []struct {
		name       string
		mapping    *Mapping
		refs       map[string]string
		valuesYAML string
		want       string
		wantErr    string
	}{
		{
			name: "preserves comments",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{
							"registry": "${registry}",
						},
					}},
				},
			},
			refs: map[string]string{
				"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `# My config
image:
  # Registry to pull from
  registry: docker.io
  tag: latest
`,
			want: `# My config
image:
  # Registry to pull from
  registry: cgr.dev
  tag: latest
`,
		},
		{
			name: "replaces multiple fields including pseudo_tag and registry_repo",
			mapping: &Mapping{
				Images: map[string]*Image{
					"app": {Values: map[string]any{
						"image": map[string]any{
							"registry":     "${registry}",
							"repository":   "${repo}",
							"registryRepo": "${registry_repo}",
							"tag":          "${tag}",
							"digest":       "${digest}",
							"pseudoTag":    "${pseudo_tag}",
							"ref":          "${ref}",
						},
					}},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1.0@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `image:
  registry: docker.io
  repository: library/app
  registryRepo: docker.io/library/app
  tag: old
  digest: ""
  pseudoTag: ""
  ref: ""
`,
			want: `image:
  registry: cgr.dev
  repository: chainguard/app
  registryRepo: cgr.dev/chainguard/app
  tag: v1.0
  digest: sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234
  pseudoTag: v1.0@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234
  ref: cgr.dev/chainguard/app:v1.0@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234
`,
		},
		{
			name: "merges multiple images",
			mapping: &Mapping{
				Images: map[string]*Image{
					"main": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
					"sidecar": {Values: map[string]any{
						"sidecar": map[string]any{"image": map[string]any{"registry": "${registry}"}},
					}},
				},
			},
			refs: map[string]string{
				"main":    "cgr.dev/chainguard/main:latest@sha256:aaaa1234aaaa1234aaaa1234aaaa1234aaaa1234aaaa1234aaaa1234aaaa1234",
				"sidecar": "cgr.dev/chainguard/sidecar:latest@sha256:bbbb1234bbbb1234bbbb1234bbbb1234bbbb1234bbbb1234bbbb1234bbbb1234",
			},
			valuesYAML: `image:
  registry: docker.io
sidecar:
  image:
    registry: docker.io
`,
			want: `image:
  registry: cgr.dev
sidecar:
  image:
    registry: cgr.dev
`,
		},
		{
			name: "nested structure with comments",
			mapping: &Mapping{
				Images: map[string]*Image{
					"proxy": {Values: map[string]any{
						"proxy": map[string]any{
							"image": map[string]any{
								"registry": "${registry}",
								"repo":     "${repo}",
							},
						},
					}},
				},
			},
			refs: map[string]string{
				"proxy": "cgr.dev/chainguard/proxy:v2@sha256:cccc1234cccc1234cccc1234cccc1234cccc1234cccc1234cccc1234cccc1234",
			},
			valuesYAML: `# Proxy configuration
proxy:
  # Enable the proxy
  enabled: true
  image:
    # Image registry
    registry: quay.io
    # Image repository
    repo: original/proxy
`,
			want: `# Proxy configuration
proxy:
  # Enable the proxy
  enabled: true
  image:
    # Image registry
    registry: cgr.dev
    # Image repository
    repo: chainguard/proxy
`,
		},
		{
			name:    "nil mapping returns original",
			mapping: nil,
			refs: map[string]string{
				"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `# Original
image:
  registry: docker.io
`,
			want: `# Original
image:
  registry: docker.io
`,
		},
		{
			name: "empty refs returns original",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{},
			valuesYAML: `image:
  registry: docker.io
`,
			want: `image:
  registry: docker.io
`,
		},
		{
			name: "invalid ref returns error",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{
				"nginx": "not-a-valid-ref",
			},
			valuesYAML: `image:
  registry: docker.io
`,
			wantErr: "must include a digest",
		},
		{
			name: "missing image ref returns error",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{
				"other": "cgr.dev/chainguard/other:latest@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `image:
  registry: docker.io
`,
			wantErr: `no ref found for image "nginx"`,
		},
		{
			name: "path not in values returns error",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{
				"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `other:
  key: value
`,
			wantErr: "does not exist",
		},
		{
			name: "mixed literal and marker string",
			mapping: &Mapping{
				Images: map[string]*Image{
					"app": {Values: map[string]any{
						"image": "${registry}/${repo}:${tag}",
					}},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1.0@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `image: docker.io/library/app:old
`,
			want: `image: cgr.dev/chainguard/app:v1.0
`,
		},
		{
			name: "empty values.yaml with no patches returns empty",
			mapping: &Mapping{
				Images: map[string]*Image{},
			},
			refs:       map[string]string{},
			valuesYAML: ``,
			want:       ``,
		},
		{
			name: "empty values.yaml with patches returns error",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{
				"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: ``,
			wantErr:    "yaml",
		},
		{
			name: "malformed yaml returns error",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{
				"nginx": "cgr.dev/chainguard/nginx:latest@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `image: [invalid yaml`,
			wantErr:    "yaml",
		},
		{
			name: "empty tag marker errors",
			mapping: &Mapping{
				Images: map[string]*Image{
					"nginx": {Values: map[string]any{
						"image": map[string]any{"tag": "${tag}"},
					}},
				},
			},
			refs: map[string]string{
				"nginx": "cgr.dev/chainguard/nginx@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `image:
  tag: old
`,
			wantErr: "empty value",
		},
		{
			name: "arrays in values are preserved",
			mapping: &Mapping{
				Images: map[string]*Image{
					"app": {Values: map[string]any{
						"image": map[string]any{"registry": "${registry}"},
					}},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `image:
  registry: docker.io
ports:
  - 80
  - 443
`,
			want: `image:
  registry: cgr.dev
ports:
  - 80
  - 443
`,
		},
		{
			name: "keys with special characters",
			mapping: &Mapping{
				Images: map[string]*Image{
					"app": {Values: map[string]any{
						"foo/bar": map[string]any{"registry": "${registry}"},
						"baz~qux": map[string]any{"tag": "${tag}"},
					}},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1@sha256:abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234abcd1234",
			},
			valuesYAML: `foo/bar:
  registry: docker.io
baz~qux:
  tag: old
`,
			want: `foo/bar:
  registry: cgr.dev
baz~qux:
  tag: v1
`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.mapping.Resolve(tc.refs, strings.NewReader(tc.valuesYAML))
			if tc.wantErr != "" {
				if err == nil {
					t.Fatalf("expected error containing %q, got nil", tc.wantErr)
				}
				if !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("expected error containing %q, got %q", tc.wantErr, err.Error())
				}
				return
			}
			if err != nil {
				t.Fatalf("Resolve: %v", err)
			}
			if diff := cmp.Diff(tc.want, string(got)); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestImageValuesMerge(t *testing.T) {
	tests := []struct {
		name    string
		input   ImageValues
		want    map[string]any
		wantErr bool
	}{
		{
			name: "disjoint images",
			input: ImageValues{
				"nginx": {"image": map[string]any{"registry": "cgr.dev"}},
				"redis": {"redis": map[string]any{"tag": "latest"}},
			},
			want: map[string]any{
				"image": map[string]any{"registry": "cgr.dev"},
				"redis": map[string]any{"tag": "latest"},
			},
		},
		{
			name: "nested merge",
			input: ImageValues{
				"nginx": {"image": map[string]any{"registry": "cgr.dev"}},
				"redis": {"image": map[string]any{"tag": "latest"}},
			},
			want: map[string]any{
				"image": map[string]any{"registry": "cgr.dev", "tag": "latest"},
			},
		},
		{
			name: "overwrite scalar",
			input: ImageValues{
				"first":  {"key": "old"},
				"second": {"key": "new"},
			},
			want: map[string]any{"key": "new"},
		},
		{
			name: "conflict map vs scalar",
			input: ImageValues{
				"first":  {"key": "scalar"},
				"second": {"key": map[string]any{"nested": "value"}},
			},
			wantErr: true,
		},
		{
			name:  "empty",
			input: ImageValues{},
			want:  map[string]any{},
		},
		{
			name:  "nil",
			input: nil,
			want:  map[string]any{},
		},
		{
			name:  "single image",
			input: ImageValues{"only": {"key": "value"}},
			want:  map[string]any{"key": "value"},
		},
		{
			name:  "image with nil values",
			input: ImageValues{"a": {"key": "value"}, "b": nil},
			want:  map[string]any{"key": "value"},
		},
		{
			name:  "image with empty values",
			input: ImageValues{"a": {"key": "value"}, "b": {}},
			want:  map[string]any{"key": "value"},
		},
		{
			name: "deep nesting",
			input: ImageValues{
				"a": {"l1": map[string]any{"l2": map[string]any{"l3": "from_a"}}},
				"b": {"l1": map[string]any{"l2": map[string]any{"other": "from_b"}}},
			},
			want: map[string]any{
				"l1": map[string]any{"l2": map[string]any{"l3": "from_a", "other": "from_b"}},
			},
		},
		{
			name: "arrays overwrite",
			input: ImageValues{
				"a_first":  {"items": []any{"a", "b"}},
				"b_second": {"items": []any{"c"}},
			},
			want: map[string]any{"items": []any{"c"}}, // b_second wins (sorted after a_first)
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.input.Merge()
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("Merge: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
