/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"strings"
	"testing"

	"chainguard.dev/sdk/helm/images"
	"github.com/google/go-cmp/cmp"
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
	}{{
		name:      "with values.yaml",
		hasValues: true,
		values:    "image: nginx\n",
		want:      "image: nginx\n",
	}, {
		name:      "without values.yaml",
		hasValues: false,
		want:      "",
	}}

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
				t.Errorf("ReadValues: got = %q, wanted = %q", got, tt.want)
			}
		})
	}
}

func TestReadChangelog(t *testing.T) {
	tests := []struct {
		name         string
		hasChangelog bool
		changelog    string
		want         string
	}{{
		name:         "with CHANGELOG.md",
		hasChangelog: true,
		changelog:    "# Changelog\n\n## 1.0.0\n- initial release\n",
		want:         "# Changelog\n\n## 1.0.0\n- initial release\n",
	}, {
		name:         "without CHANGELOG.md",
		hasChangelog: false,
		want:         "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var chart v1.Image
			if tt.hasChangelog {
				chart = createTestChartWithChangelog(t, "my-chart", "image: nginx\n", tt.changelog)
			} else {
				chart = createTestChart(t, "my-chart", "image: nginx\n")
			}

			got, err := ReadChangelog(chart)
			if err != nil {
				t.Fatalf("ReadChangelog: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("ReadChangelog: got = %q, wanted = %q", got, tt.want)
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
		t.Errorf("after ReplaceValues, ReadValues: got = %q, wanted = %q", got, want)
	}
}

func TestReplaceValues_ErrorWithoutValues(t *testing.T) {
	chart := createChartWithoutValues(t, "my-chart")

	_, err := ReplaceValues(chart, &images.Mapping{}, map[string]string{})
	if err == nil {
		t.Fatal("expected error for chart without values.yaml")
	}
}

func TestReplaceValues_WithOmitDigests(t *testing.T) {
	tests := []struct {
		name       string
		valuesYAML string
		mapping    *images.Mapping
		refs       map[string]string
		opts       []images.ResolveOption
		want       string
		wantErr    bool
	}{
		{
			name:       "default includes digest in pseudo_tag",
			valuesYAML: "image:\n  pseudoTag: \"\"\n",
			mapping: &images.Mapping{
				Images: map[string]*images.Image{
					"app": {
						Values: map[string]any{
							"image": map[string]any{
								"pseudoTag": "${pseudo_tag}",
							},
						},
					},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1.0.0@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			opts: nil,
			want: "image:\n  pseudoTag: v1.0.0@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n",
		},
		{
			name:       "WithOmitDigests(false) includes digest in pseudo_tag",
			valuesYAML: "image:\n  pseudoTag: \"\"\n",
			mapping: &images.Mapping{
				Images: map[string]*images.Image{
					"app": {
						Values: map[string]any{
							"image": map[string]any{
								"pseudoTag": "${pseudo_tag}",
							},
						},
					},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1.0.0@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			opts: []images.ResolveOption{images.WithOmitDigests(false)},
			want: "image:\n  pseudoTag: v1.0.0@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef\n",
		},
		{
			name:       "WithOmitDigests(true) omits digest from pseudo_tag",
			valuesYAML: "image:\n  pseudoTag: \"\"\n",
			mapping: &images.Mapping{
				Images: map[string]*images.Image{
					"app": {
						Values: map[string]any{
							"image": map[string]any{
								"pseudoTag": "${pseudo_tag}",
							},
						},
					},
				},
			},
			refs: map[string]string{
				"app": "cgr.dev/chainguard/app:v1.0.0@sha256:0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			opts: []images.ResolveOption{images.WithOmitDigests(true)},
			want: "image:\n  pseudoTag: v1.0.0\n",
		},
		{
			name:       "WithOmitDigests(true) with multiple fields",
			valuesYAML: "image:\n  registry: \"\"\n  repo: \"\"\n  pseudoTag: \"\"\n",
			mapping: &images.Mapping{
				Images: map[string]*images.Image{
					"nginx": {
						Values: map[string]any{
							"image": map[string]any{
								"registry":  "${registry}",
								"repo":      "${repo}",
								"pseudoTag": "${pseudo_tag}",
							},
						},
					},
				},
			},
			refs: map[string]string{
				"nginx": "my-registry.io/my-group/nginx:latest@sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
			},
			opts: []images.ResolveOption{images.WithOmitDigests(true)},
			want: "image:\n  registry: my-registry.io\n  repo: my-group/nginx\n  pseudoTag: latest\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart := createTestChart(t, "my-chart", tt.valuesYAML)

			patched, err := ReplaceValues(chart, tt.mapping, tt.refs, tt.opts...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReplaceValues() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			got, err := ReadValues(patched)
			if err != nil {
				t.Fatalf("ReadValues: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("after ReplaceValues, ReadValues:\ngot  = %q\nwant = %q", got, tt.want)
			}
		})
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
	}{{
		path: "chart/values.yaml",
		want: true,
	}, {
		path: "my-chart/values.yaml",
		want: true,
	}, {
		path: "values.yaml",
		want: false,
	}, {
		path: "chart/charts/subchart/values.yaml",
		want: false,
	}, {
		path: "chart/templates/values.yaml",
		want: false,
	}, {
		path: "chart/values.yml",
		want: false,
	}, {
		path: "a/b/c/values.yaml",
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTopLevelValuesYAML(tt.path); got != tt.want {
				t.Errorf("isTopLevelValuesYAML(%q): got = %v, wanted = %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsTopLevelChartYAML(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{{
		path: "chart/Chart.yaml",
		want: true,
	}, {
		path: "my-chart/Chart.yaml",
		want: true,
	}, {
		path: "Chart.yaml",
		want: false,
	}, {
		path: "chart/charts/subchart/Chart.yaml",
		want: false,
	}, {
		path: "chart/templates/Chart.yaml",
		want: false,
	}, {
		path: "chart/Chart.yml",
		want: false,
	}, {
		path: "a/b/c/Chart.yaml",
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTopLevelChartYAML(tt.path); got != tt.want {
				t.Errorf("isTopLevelChartYAML(%q): got = %v, wanted = %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsTopLevelChangelog(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{{
		path: "chart/CHANGELOG.md",
		want: true,
	}, {
		path: "my-chart/CHANGELOG.md",
		want: true,
	}, {
		path: "CHANGELOG.md",
		want: false,
	}, {
		path: "chart/charts/subchart/CHANGELOG.md",
		want: false,
	}, {
		path: "chart/docs/CHANGELOG.md",
		want: false,
	}, {
		path: "chart/CHANGELOG.txt",
		want: false,
	}, {
		path: "a/b/c/CHANGELOG.md",
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			if got := isTopLevelChangelog(tt.path); got != tt.want {
				t.Errorf("isTopLevelChangelog(%q): got = %v, wanted = %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestReadChartFile(t *testing.T) {
	tests := []struct {
		name        string
		pathMatcher PathMatcher
		want        string
	}{{
		name:        "read Chart.yaml",
		pathMatcher: isTopLevelChartYAML,
		want:        "apiVersion: v2\nname: my-chart\nversion: 1.0.0\n",
	}, {
		name:        "read values.yaml",
		pathMatcher: isTopLevelValuesYAML,
		want:        "image: nginx\n",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			chart := createTestChart(t, "my-chart", "image: nginx\n")

			got, err := readChartFile(chart, tt.pathMatcher)
			if err != nil {
				t.Fatalf("readChartFile: %v", err)
			}

			if string(got) != tt.want {
				t.Errorf("readChartFile: got = %q, wanted = %q", got, tt.want)
			}
		})
	}
}

func TestReadChartFile_SizeLimit(t *testing.T) {
	oversized := strings.Repeat("a", maxFileSize+1)

	// The matched file is rejected when it exceeds maxFileSize.
	t.Run("oversized matched file is rejected", func(t *testing.T) {
		chart := createTestChart(t, "my-chart", oversized)

		_, err := readChartFile(chart, isTopLevelValuesYAML)
		if err == nil {
			t.Fatal("readChartFile: expected error for oversized matched file, got nil")
		}
		if !strings.Contains(err.Error(), "exceeds maximum size") {
			t.Errorf("readChartFile: got error %q, want it to mention exceeding maximum size", err)
		}
	})

	// An oversized file the caller isn't reading must be skipped, not rejected,
	// so charts shipping large unrelated files don't break ReadValues/
	// ReadChartMeta. Here the oversized values.yaml is unmatched (no CHANGELOG.md
	// to match), so the scan skips it and returns nil without error.
	t.Run("oversized unmatched file is skipped", func(t *testing.T) {
		chart := createTestChart(t, "my-chart", oversized)

		got, err := readChartFile(chart, isTopLevelChangelog)
		if err != nil {
			t.Fatalf("readChartFile: unexpected error for oversized unmatched file: %v", err)
		}
		if got != nil {
			t.Errorf("readChartFile: expected nil (no match), got %d bytes", len(got))
		}
	})
}

func TestReadChartMeta(t *testing.T) {
	tests := []struct {
		name  string
		chart v1.Image
		want  *Meta
	}{{
		name:  "valid Chart.yaml",
		chart: createTestChart(t, "my-chart", "image: nginx\n"),
		want: &Meta{
			APIVersion: "v2",
			Name:       "my-chart",
			Version:    "1.0.0",
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta, err := ReadChartMeta(tt.chart)
			if err != nil {
				t.Fatalf("ReadChartMeta: %v", err)
			}

			if cmp.Diff(tt.want, meta) != "" {
				t.Errorf("ReadChartMeta: got = %+v, wanted = %+v", meta, tt.want)
			}
		})
	}
}

func TestReplaceChartMeta(t *testing.T) {
	meta1 := &Meta{
		APIVersion:  "v2",
		Name:        "my-chart",
		Version:     "2.0.0",
		AppVersion:  "1.5.0",
		Description: "Updated chart",
	}
	meta2 := &Meta{
		APIVersion: "v2",
		AppVersion: "1.5.0",
	}

	tests := []struct {
		name         string
		chart        v1.Image
		newMeta      *Meta
		expectedMeta *Meta
	}{
		{
			name:         "replace Chart.yaml metadata full",
			chart:        createTestChart(t, "test-chart", "image: nginx\n"),
			newMeta:      meta1,
			expectedMeta: meta1,
		},
		{
			name:         "replace Chart.yaml metadata partial",
			chart:        createTestChart(t, "test-chart", "image: nginx\n"),
			newMeta:      meta2,
			expectedMeta: meta2,
		},
		{
			name:         "replace Chart.yaml metadata empty",
			chart:        createTestChart(t, "test-chart", "image: nginx\n"),
			newMeta:      &Meta{},
			expectedMeta: &Meta{},
		},
		{
			name:         "replace Chart.yaml metadata nil",
			chart:        createTestChart(t, "test-chart", "image: nginx\n"),
			newMeta:      nil,
			expectedMeta: &Meta{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patched, err := ReplaceChartMeta(tt.chart, tt.newMeta)
			if err != nil {
				t.Fatalf("ReplaceChartMeta: %v", err)
			}

			got, err := ReadChartMeta(patched)
			if err != nil {
				t.Fatalf("ReadChartMeta after replace: %v", err)
			}

			if cmp.Diff(tt.expectedMeta, got) != "" {
				t.Errorf("ReadChartMeta: got = %+v, wanted = %+v", got, tt.newMeta)
			}
		})
	}
}

func TestReplaceChartMeta_PreservesOtherFiles(t *testing.T) {
	chart := createTestChart(t, "my-chart", "image: nginx\n")

	newMeta := &Meta{
		APIVersion: "v2",
		Name:       "my-chart",
		Version:    "2.0.0",
	}

	patched, err := ReplaceChartMeta(chart, newMeta)
	if err != nil {
		t.Fatalf("ReplaceChartMeta: %v", err)
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

	foundValuesYAML := false
	tr := tar.NewReader(rc)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("reading tar: %v", err)
		}
		if header.Name == "my-chart/values.yaml" {
			foundValuesYAML = true
		}
	}

	if !foundValuesYAML {
		t.Error("values.yaml not found in patched chart")
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

func createTestChartWithChangelog(t *testing.T, chartName, valuesYAML, changelog string) v1.Image {
	t.Helper()

	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)

	files := []struct {
		name, data string
	}{
		{chartName + "/Chart.yaml", "apiVersion: v2\nname: " + chartName + "\nversion: 1.0.0\n"},
		{chartName + "/values.yaml", valuesYAML},
		{chartName + "/CHANGELOG.md", changelog},
	}
	for _, f := range files {
		if err := tw.WriteHeader(&tar.Header{
			Name: f.name,
			Size: int64(len(f.data)),
			Mode: 0644,
		}); err != nil {
			t.Fatalf("writing %s header: %v", f.name, err)
		}
		if _, err := tw.Write([]byte(f.data)); err != nil {
			t.Fatalf("writing %s: %v", f.name, err)
		}
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
