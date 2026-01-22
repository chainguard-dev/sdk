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

func TestNewRef(t *testing.T) {
	const validDigest = "sha256:e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	tests := []struct {
		name      string
		reference string
		want      OCIRef
		wantErr   string
	}{
		// Success cases - with digest
		{
			name:      "tag and digest",
			reference: fmt.Sprintf("cgr.dev/chainguard/nginx:latest@%s", validDigest),
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "chainguard/nginx",
				Tag:          "latest",
				Digest:       validDigest,
				RegistryRepo: "cgr.dev/chainguard/nginx",
				FullRef:      fmt.Sprintf("cgr.dev/chainguard/nginx:latest@%s", validDigest),
			},
		},
		{
			name:      "digest only",
			reference: fmt.Sprintf("cgr.dev/chainguard/nginx@%s", validDigest),
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "chainguard/nginx",
				Digest:       validDigest,
				RegistryRepo: "cgr.dev/chainguard/nginx",
				FullRef:      fmt.Sprintf("cgr.dev/chainguard/nginx@%s", validDigest),
			},
		},
		{
			name:      "explicit latest with digest",
			reference: fmt.Sprintf("cgr.dev/nginx:latest@%s", validDigest),
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "nginx",
				Tag:          "latest",
				Digest:       validDigest,
				RegistryRepo: "cgr.dev/nginx",
				FullRef:      fmt.Sprintf("cgr.dev/nginx:latest@%s", validDigest),
			},
		},
		{
			name:      "nested repo with digest",
			reference: fmt.Sprintf("cgr.dev/chainguard/images/static@%s", validDigest),
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "chainguard/images/static",
				Digest:       validDigest,
				RegistryRepo: "cgr.dev/chainguard/images/static",
				FullRef:      fmt.Sprintf("cgr.dev/chainguard/images/static@%s", validDigest),
			},
		},
		{
			name:      "registry with port and digest",
			reference: fmt.Sprintf("localhost:5000/myimage@%s", validDigest),
			want: OCIRef{
				Registry:     "localhost:5000",
				Repo:         "myimage",
				Digest:       validDigest,
				RegistryRepo: "localhost:5000/myimage",
				FullRef:      fmt.Sprintf("localhost:5000/myimage@%s", validDigest),
			},
		},
		{
			name:      "docker hub with digest",
			reference: fmt.Sprintf("nginx@%s", validDigest),
			want: OCIRef{
				Registry:     "index.docker.io",
				Repo:         "library/nginx",
				Digest:       validDigest,
				RegistryRepo: "index.docker.io/library/nginx",
				FullRef:      fmt.Sprintf("index.docker.io/library/nginx@%s", validDigest),
			},
		},
		{
			name:      "semver tag with digest",
			reference: fmt.Sprintf("cgr.dev/nginx:v1.25.0@%s", validDigest),
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "nginx",
				Tag:          "v1.25.0",
				Digest:       validDigest,
				RegistryRepo: "cgr.dev/nginx",
				FullRef:      fmt.Sprintf("cgr.dev/nginx:v1.25.0@%s", validDigest),
			},
		},
		{
			name:      "tag with special chars and digest",
			reference: fmt.Sprintf("cgr.dev/nginx:v1.25.0-alpine_3.18@%s", validDigest),
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "nginx",
				Tag:          "v1.25.0-alpine_3.18",
				Digest:       validDigest,
				RegistryRepo: "cgr.dev/nginx",
				FullRef:      fmt.Sprintf("cgr.dev/nginx:v1.25.0-alpine_3.18@%s", validDigest),
			},
		},

		// Success cases - without digest
		{
			name:      "tag only",
			reference: "cgr.dev/chainguard/nginx:latest",
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "chainguard/nginx",
				Tag:          "latest",
				RegistryRepo: "cgr.dev/chainguard/nginx",
				FullRef:      "cgr.dev/chainguard/nginx:latest",
			},
		},
		{
			name:      "implicit latest tag",
			reference: "cgr.dev/chainguard/nginx",
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "chainguard/nginx",
				Tag:          "latest",
				RegistryRepo: "cgr.dev/chainguard/nginx",
				FullRef:      "cgr.dev/chainguard/nginx:latest",
			},
		},
		{
			name:      "nested repo path with tag",
			reference: "cgr.dev/chainguard/images/static:v1",
			want: OCIRef{
				Registry:     "cgr.dev",
				Repo:         "chainguard/images/static",
				Tag:          "v1",
				RegistryRepo: "cgr.dev/chainguard/images/static",
				FullRef:      "cgr.dev/chainguard/images/static:v1",
			},
		},
		{
			name:      "registry with port and tag",
			reference: "localhost:5000/myimage:v1",
			want: OCIRef{
				Registry:     "localhost:5000",
				Repo:         "myimage",
				Tag:          "v1",
				RegistryRepo: "localhost:5000/myimage",
				FullRef:      "localhost:5000/myimage:v1",
			},
		},
		{
			name:      "docker hub with tag",
			reference: "nginx:latest",
			want: OCIRef{
				Registry:     "index.docker.io",
				Repo:         "library/nginx",
				Tag:          "latest",
				RegistryRepo: "index.docker.io/library/nginx",
				FullRef:      "index.docker.io/library/nginx:latest",
			},
		},

		// Error cases - invalid reference
		{
			name:      "invalid reference",
			reference: "no can do boss",
			wantErr:   "could not parse reference",
		},
		{
			name:      "empty reference",
			reference: "",
			wantErr:   "could not parse reference",
		},

		// Error cases - invalid digest
		{
			name:      "invalid digest format",
			reference: "cgr.dev/nginx@notadigest",
			wantErr:   "could not parse reference",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewRef(tc.reference)
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
				t.Fatalf("NewRef: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
