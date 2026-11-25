/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"chainguard.dev/sdk/helm/images"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/registry"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/static"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

// helmConfig is a config blob as helm push writes it: Chart.yaml as JSON,
// with fields [Meta] does not model.
const helmConfig = `{"name":"my-chart","home":"https://example.com/home","sources":["https://example.com/first","https://example.com/second"],"version":"1.0.0","apiVersion":"v2"}`

var helmAnnotations = map[string]string{
	"org.opencontainers.image.title":  "my-chart",
	"org.opencontainers.image.source": "https://example.com/first",
}

// rawManifest is a manifest to put as written, so the fixture's manifest does
// not pass through the code under test.
type rawManifest []byte

func (m rawManifest) RawManifest() ([]byte, error)        { return m, nil }
func (m rawManifest) MediaType() (types.MediaType, error) { return types.OCIManifestSchema1, nil }

// pushHelmChart pushes a chart shaped as helm push writes one, a Helm config
// blob and the chart layer, and returns it as pulled from the registry.
func pushHelmChart(t *testing.T, repo name.Repository, valuesYAML string) v1.Image {
	t.Helper()

	layers, err := createTestChart(t, "my-chart", valuesYAML).Layers()
	if err != nil {
		t.Fatalf("Layers: %v", err)
	}
	layer := layers[0]
	config := static.NewLayer([]byte(helmConfig), ConfigMediaType)

	descs := make([]v1.Descriptor, 0, 2)
	for _, l := range []v1.Layer{config, layer} {
		if err := remote.WriteLayer(repo, l); err != nil {
			t.Fatalf("WriteLayer: %v", err)
		}
		mt, err := l.MediaType()
		if err != nil {
			t.Fatalf("MediaType: %v", err)
		}
		size, err := l.Size()
		if err != nil {
			t.Fatalf("Size: %v", err)
		}
		digest, err := l.Digest()
		if err != nil {
			t.Fatalf("Digest: %v", err)
		}
		descs = append(descs, v1.Descriptor{MediaType: mt, Size: size, Digest: digest})
	}

	manifest, err := json.Marshal(&v1.Manifest{
		SchemaVersion: 2,
		MediaType:     types.OCIManifestSchema1,
		Config:        descs[0],
		Layers:        descs[1:],
		Annotations:   helmAnnotations,
	})
	if err != nil {
		t.Fatalf("marshaling manifest: %v", err)
	}

	ref := repo.Tag("source")
	if err := remote.Put(ref, rawManifest(manifest)); err != nil {
		t.Fatalf("Put: %v", err)
	}
	img, err := remote.Image(ref)
	if err != nil {
		t.Fatalf("Image: %v", err)
	}
	return img
}

// roundTrip pushes a rewritten chart and returns it as pulled back, so the
// assertions see what a registry client would.
func roundTrip(t *testing.T, repo name.Repository, img v1.Image) v1.Image {
	t.Helper()

	ref := repo.Tag("rewritten")
	if err := remote.Write(ref, img); err != nil {
		t.Fatalf("Write: %v", err)
	}
	got, err := remote.Image(ref)
	if err != nil {
		t.Fatalf("Image: %v", err)
	}
	return got
}

func testRepo(t *testing.T) name.Repository {
	t.Helper()

	s := httptest.NewServer(registry.New())
	t.Cleanup(s.Close)
	repo, err := name.NewRepository(strings.TrimPrefix(s.URL, "http://") + "/charts/my-chart")
	if err != nil {
		t.Fatalf("NewRepository: %v", err)
	}
	return repo
}

func TestRewriteKeepsHelmConfig(t *testing.T) {
	const values = "image:\n  registry: cgr.dev\n  repository: chainguard/nginx\n"
	mapping := &images.Mapping{
		Images: map[string]*images.Image{
			"nginx": {
				Values: map[string]any{
					"image": map[string]any{
						"registry":   "${registry}",
						"repository": "${repo}",
					},
				},
			},
		},
	}
	refs := map[string]string{
		"nginx": "my-registry.io/my-group/nginx:latest@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
	}

	tests := []struct {
		name       string
		rewrite    func(v1.Image) (v1.Image, error)
		wantConfig string
		wantValues string
	}{{
		name: "values rewrite carries the config over",
		rewrite: func(img v1.Image) (v1.Image, error) {
			return ReplaceValues(img, mapping, refs)
		},
		wantConfig: helmConfig,
		wantValues: "image:\n  registry: my-registry.io\n  repository: my-group/nginx\n",
	}, {
		name: "Chart.yaml rewrite writes the new metadata",
		rewrite: func(img v1.Image) (v1.Image, error) {
			return ReplaceChartMeta(img, &Meta{APIVersion: "v2", Name: "my-chart", Version: "2.0.0"})
		},
		wantConfig: `{"apiVersion":"v2","name":"my-chart","version":"2.0.0"}`,
		wantValues: values,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := testRepo(t)
			source := pushHelmChart(t, repo, values)

			rewritten, err := tt.rewrite(source)
			if err != nil {
				t.Fatalf("rewrite: %v", err)
			}
			got := roundTrip(t, repo, rewritten)

			manifest, err := got.Manifest()
			if err != nil {
				t.Fatalf("Manifest: %v", err)
			}
			if manifest.Config.MediaType != ConfigMediaType {
				t.Errorf("config media type: got = %q, want = %q", manifest.Config.MediaType, ConfigMediaType)
			}
			if diff := cmp.Diff(helmAnnotations, manifest.Annotations); diff != "" {
				t.Errorf("annotations (-want +got):\n%s", diff)
			}

			config, err := got.RawConfigFile()
			if err != nil {
				t.Fatalf("RawConfigFile: %v", err)
			}
			if string(config) != tt.wantConfig {
				t.Errorf("config: got = %s, want = %s", config, tt.wantConfig)
			}

			gotValues, err := ReadValues(got)
			if err != nil {
				t.Fatalf("ReadValues: %v", err)
			}
			if string(gotValues) != tt.wantValues {
				t.Errorf("values: got = %q, want = %q", gotValues, tt.wantValues)
			}
		})
	}
}
