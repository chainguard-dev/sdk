/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	"chainguard.dev/sdk/helm/images"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"gopkg.in/yaml.v3"
)

// MediaType is the OCI media type for Helm chart content layers.
const MediaType types.MediaType = "application/vnd.cncf.helm.chart.content.v1.tar+gzip"

// maxFileSize is the maximum size of any single file in a Helm chart (10 MB).
const maxFileSize = 10 * 1024 * 1024

// PathMatcher defines a function type for matching file paths in a chart tarball.
// Used to identify specific files like values.yaml or Chart.yaml while iterating through tarball contents.
type PathMatcher func(string) bool

// ReadValues extracts the top-level values.yaml from a Helm chart OCI artifact.
// Returns nil, nil if the chart has no values.yaml.
func ReadValues(chart v1.Image) ([]byte, error) {
	return readChartFile(chart, isTopLevelValuesYAML)
}

// ReadChartMeta extracts and parses the top-level Chart.yaml metadata from a Helm chart OCI artifact.
func ReadChartMeta(chart v1.Image) (*Meta, error) {
	chartYAML, err := readChartFile(chart, isTopLevelChartYAML)
	if err != nil {
		return nil, err
	}
	if chartYAML == nil {
		return nil, fmt.Errorf("chart has no Chart.yaml")
	}

	var meta Meta
	if err := yaml.Unmarshal(chartYAML, &meta); err != nil {
		return nil, fmt.Errorf("parsing Chart.yaml: %w", err)
	}

	return &meta, nil
}

// ReadValues extracts the top-level values.yaml from a Helm chart OCI artifact.
// Returns nil, nil if the chart has no values.yaml.
func readChartFile(chart v1.Image, pathMatcher PathMatcher) ([]byte, error) {
	layer, err := getChartLayer(chart)
	if err != nil {
		return nil, err
	}

	rc, err := layer.Uncompressed()
	if err != nil {
		return nil, fmt.Errorf("uncompressing layer: %w", err)
	}
	defer rc.Close()

	tr := tar.NewReader(rc)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar: %w", err)
		}

		if pathMatcher(header.Name) {
			content, err := io.ReadAll(tr)
			if err != nil {
				return nil, fmt.Errorf("reading values.yaml: %w", err)
			}
			return content, nil
		}
	}

	return nil, nil
}

// ReplaceValues returns a new Helm chart OCI artifact with the top-level
// values.yaml transformed by applying the given mapping template with the
// provided image refs.
func ReplaceValues(chart v1.Image, m *images.Mapping, refs map[string]string) (v1.Image, error) {
	values, err := ReadValues(chart)
	if err != nil {
		return nil, fmt.Errorf("reading values: %w", err)
	}
	if values == nil {
		return nil, fmt.Errorf("chart has no values.yaml")
	}

	newValues, err := m.Resolve(refs, bytes.NewReader(values))
	if err != nil {
		return nil, fmt.Errorf("resolving values: %w", err)
	}

	return rewriteChart(chart, newValues, isTopLevelValuesYAML)
}

// ReplaceChartMeta returns a new Helm chart OCI artifact with the top-level
// Chart.yaml replaced with the provided metadata.
func ReplaceChartMeta(chart v1.Image, meta *Meta) (v1.Image, error) {
	chartYAML, err := yaml.Marshal(meta)
	if err != nil {
		return nil, fmt.Errorf("marshaling Chart.yaml: %w", err)
	}

	return rewriteChart(chart, chartYAML, isTopLevelChartYAML)
}

// rewriteChart returns a new Helm chart image replacing file matched by pathMatcher with the provided values content.
func rewriteChart(chart v1.Image, values []byte, pathMatcher PathMatcher) (v1.Image, error) {
	layer, err := getChartLayer(chart)
	if err != nil {
		return nil, err
	}

	origMediaType, err := layer.MediaType()
	if err != nil {
		return nil, fmt.Errorf("getting layer media type: %w", err)
	}

	rc, err := layer.Uncompressed()
	if err != nil {
		return nil, fmt.Errorf("uncompressing layer: %w", err)
	}
	defer rc.Close()

	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)

	foundPath := false
	tr := tar.NewReader(rc)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading tar: %w", err)
		}

		if pathMatcher(header.Name) {
			foundPath = true
			header.Size = int64(len(values))
			if err := tw.WriteHeader(header); err != nil {
				return nil, fmt.Errorf("writing values.yaml header: %w", err)
			}
			if _, err := tw.Write(values); err != nil {
				return nil, fmt.Errorf("writing values.yaml: %w", err)
			}
			continue
		}

		if header.Size > maxFileSize {
			return nil, fmt.Errorf("file %q exceeds maximum size (%d > %d)", header.Name, header.Size, maxFileSize)
		}
		if err := tw.WriteHeader(header); err != nil {
			return nil, fmt.Errorf("writing tar header for %q: %w", header.Name, err)
		}
		if _, err := io.CopyN(tw, tr, header.Size); err != nil {
			return nil, fmt.Errorf("copying file %q: %w", header.Name, err)
		}
	}

	if !foundPath {
		return nil, fmt.Errorf("chart has no file matching path matcher")
	}

	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("closing tar writer: %w", err)
	}

	var gzBuf bytes.Buffer
	gzw := gzip.NewWriter(&gzBuf)
	if _, err := gzw.Write(tarBuf.Bytes()); err != nil {
		return nil, fmt.Errorf("compressing tar: %w", err)
	}
	if err := gzw.Close(); err != nil {
		return nil, fmt.Errorf("closing gzip writer: %w", err)
	}

	newLayer, err := tarball.LayerFromOpener(func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(gzBuf.Bytes())), nil
	}, tarball.WithMediaType(origMediaType))
	if err != nil {
		return nil, fmt.Errorf("creating layer: %w", err)
	}

	origManifest, err := chart.Manifest()
	if err != nil {
		return nil, fmt.Errorf("getting manifest: %w", err)
	}

	base := mutate.MediaType(empty.Image, origManifest.MediaType)
	base = mutate.ConfigMediaType(base, origManifest.Config.MediaType)

	img, err := mutate.Append(base, mutate.Addendum{
		Layer:     newLayer,
		MediaType: origMediaType,
	})
	if err != nil {
		return nil, fmt.Errorf("appending layer: %w", err)
	}

	if len(origManifest.Annotations) > 0 {
		img = mutate.Annotations(img, origManifest.Annotations).(v1.Image)
	}

	return img, nil
}

func getChartLayer(chart v1.Image) (v1.Layer, error) {
	layers, err := chart.Layers()
	if err != nil {
		return nil, fmt.Errorf("getting layers: %w", err)
	}
	if len(layers) != 1 {
		return nil, fmt.Errorf("expected 1 layer, got %d", len(layers))
	}

	mt, err := layers[0].MediaType()
	if err != nil {
		return nil, fmt.Errorf("getting layer media type: %w", err)
	}
	if mt != MediaType {
		return nil, fmt.Errorf("expected layer media type %q, got %q", MediaType, mt)
	}

	return layers[0], nil
}

// isTopLevelValuesYAML a PathMatcher func that checks if a tar path is a top-level chart's
// values.yaml. Top-level values.yaml is at {chartName}/values.yaml (one directory deep).
// Subcharts are at {chartName}/charts/{subchart}/values.yaml and should be excluded.
func isTopLevelValuesYAML(path string) bool {
	if !strings.HasSuffix(path, "/values.yaml") {
		return false
	}
	parts := strings.Split(path, "/")
	return len(parts) == 2
}

// isTopLevelChartYAML a PathMatcher func that checks if a tar path is a top-level chart's
// Chart.yaml. Top-level Chart.yaml is at {chartName}/Chart.yaml (one directory deep).
// Subcharts are at {chartName}/charts/{subchart}/Chart.yaml and should be excluded.
func isTopLevelChartYAML(path string) bool {
	if !strings.HasSuffix(path, "/Chart.yaml") {
		return false
	}
	parts := strings.Split(path, "/")
	return len(parts) == 2
}
