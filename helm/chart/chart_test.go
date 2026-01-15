/*
Copyright 2025 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"testing"

	"chainguard.dev/sdk/helm/images"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

func TestReadValues(t *testing.T) {
	tests := []struct {
		name      string
		hasValues bool
		values    string
		want      string
	}{
		{
			name:      "with values.yaml",
			hasValues: true,
			values:    "image: nginx\n",
			want:      "image: nginx\n",
		},
		{
			name:      "without values.yaml",
			hasValues: false,
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var chart v1.Image
			if tt.hasValues {
				chart = createTestChart(t, "my-chart", tt.values)
			} else {
				chart = createChartWithoutValues(t, "my-chart")
			}

			got, err := ReadValues(chart)
			if err != nil {
				t.Fatalf("ReadValues: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("ReadValues = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplaceValues(t *testing.T) {
	chart := createTestChart(t, "my-chart", "image:\n  registry: cgr.dev\n  repository: chainguard/nginx\n")

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

	patched, err := ReplaceValues(chart, mapping, refs)
	if err != nil {
		t.Fatalf("ReplaceValues: %v", err)
	}

	got, err := ReadValues(patched)
	if err != nil {
		t.Fatalf("ReadValues: %v", err)
	}

	want := "image:\n  registry: my-registry.io\n  repository: my-group/nginx\n"
	if string(got) != want {
		t.Errorf("after ReplaceValues, ReadValues = %q, want %q", got, want)
	}
}

func TestReplaceValues_ErrorWithoutValues(t *testing.T) {
	chart := createChartWithoutValues(t, "my-chart")

	_, err := ReplaceValues(chart, &images.Mapping{}, map[string]string{})
	if err == nil {
		t.Fatal("expected error for chart without values.yaml")
	}
}

func TestReplaceValues_PreservesOtherFiles(t *testing.T) {
	chart := createTestChart(t, "my-chart", "image:\n  registry: cgr.dev\n  repository: chainguard/nginx\n")

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

	patched, err := ReplaceValues(chart, mapping, refs)
	if err != nil {
		t.Fatalf("ReplaceValues: %v", err)
	}

	layer, err := getChartLayer(patched)
	if err != nil {
		t.Fatalf("getChartLayer: %v", err)
	}

	rc, err := layer.Uncompressed()
	if err != nil {
		t.Fatalf("Uncompressed: %v", err)
	}
	defer rc.Close()

	foundChartYAML := false
	tr := tar.NewReader(rc)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("reading tar: %v", err)
		}
		if header.Name == "my-chart/Chart.yaml" {
			foundChartYAML = true
		}
	}

	if !foundChartYAML {
		t.Error("Chart.yaml not found in patched chart")
	}
}

func TestIsTopLevelValuesYAML(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"chart/values.yaml", true},
		{"my-chart/values.yaml", true},
		{"values.yaml", false},
		{"chart/charts/subchart/values.yaml", false},
		{"chart/templates/values.yaml", false},
		{"chart/values.yml", false},
		{"a/b/c/values.yaml", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTopLevelValuesYAML(tt.path); got != tt.want {
				t.Errorf("isTopLevelValuesYAML(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func createTestChart(t *testing.T, chartName, valuesYAML string) v1.Image {
	t.Helper()

	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)

	chartYAML := "apiVersion: v2\nname: " + chartName + "\nversion: 1.0.0\n"
	if err := tw.WriteHeader(&tar.Header{
		Name: chartName + "/Chart.yaml",
		Size: int64(len(chartYAML)),
		Mode: 0644,
	}); err != nil {
		t.Fatalf("writing Chart.yaml header: %v", err)
	}
	if _, err := tw.Write([]byte(chartYAML)); err != nil {
		t.Fatalf("writing Chart.yaml: %v", err)
	}

	if err := tw.WriteHeader(&tar.Header{
		Name: chartName + "/values.yaml",
		Size: int64(len(valuesYAML)),
		Mode: 0644,
	}); err != nil {
		t.Fatalf("writing values.yaml header: %v", err)
	}
	if _, err := tw.Write([]byte(valuesYAML)); err != nil {
		t.Fatalf("writing values.yaml: %v", err)
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	return createImageFromTar(t, tarBuf.Bytes())
}

func createChartWithoutValues(t *testing.T, chartName string) v1.Image {
	t.Helper()

	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)

	chartYAML := "apiVersion: v2\nname: " + chartName + "\nversion: 1.0.0\n"
	if err := tw.WriteHeader(&tar.Header{
		Name: chartName + "/Chart.yaml",
		Size: int64(len(chartYAML)),
		Mode: 0644,
	}); err != nil {
		t.Fatalf("writing Chart.yaml header: %v", err)
	}
	if _, err := tw.Write([]byte(chartYAML)); err != nil {
		t.Fatalf("writing Chart.yaml: %v", err)
	}

	if err := tw.Close(); err != nil {
		t.Fatalf("closing tar: %v", err)
	}

	return createImageFromTar(t, tarBuf.Bytes())
}

func createImageFromTar(t *testing.T, tarData []byte) v1.Image {
	t.Helper()

	var gzBuf bytes.Buffer
	gzw := gzip.NewWriter(&gzBuf)
	if _, err := gzw.Write(tarData); err != nil {
		t.Fatalf("compressing: %v", err)
	}
	if err := gzw.Close(); err != nil {
		t.Fatalf("closing gzip: %v", err)
	}

	layer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(gzBuf.Bytes())), nil
	}, tarball.WithMediaType(MediaType))
	if err != nil {
		t.Fatalf("creating layer: %v", err)
	}

	img := mutate.MediaType(empty.Image, types.OCIManifestSchema1)
	img, err = mutate.Append(img, mutate.Addendum{
		Layer:     layer,
		MediaType: MediaType,
	})
	if err != nil {
		t.Fatalf("appending layer: %v", err)
	}

	return img
}
