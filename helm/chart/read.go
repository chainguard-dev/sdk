/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"chainguard.dev/sdk/helm/images"
	"gopkg.in/yaml.v3"
)

// HelmMetadataFilename is the standard Helm chart metadata filename.
const HelmMetadataFilename = "Chart.yaml"

// Meta holds metadata from Chart.yaml.
type Meta struct {
	APIVersion   string       `yaml:"apiVersion" json:"apiVersion"`
	Name         string       `yaml:"name" json:"name"`
	Version      string       `yaml:"version" json:"version"`
	AppVersion   string       `yaml:"appVersion,omitempty" json:"appVersion,omitempty"`
	Description  string       `yaml:"description,omitempty" json:"description,omitempty"`
	Dependencies []Dependency `yaml:"dependencies,omitempty" json:"dependencies,omitempty"`
}

// Dependency represents a dependency entry in Chart.yaml.
type Dependency struct {
	Name  string `yaml:"name" json:"name"`
	Alias string `yaml:"alias,omitempty" json:"alias,omitempty"`
}

// Key returns the Helm values namespace key for this dependency —
// the alias if set, otherwise the name.
func (d Dependency) Key() string {
	if d.Alias != "" {
		return d.Alias
	}
	return d.Name
}

// Chart holds metadata and image mappings extracted from a chart filesystem.
type Chart struct {
	Meta      *Meta
	Mapping   *images.Mapping
	Subcharts map[string]*Chart // dep.Key() -> subchart
}

// dirEntry holds Chart.yaml and cg.json for a single directory in the tar.
type dirEntry struct {
	meta    *Meta
	mapping *images.Mapping
}

// Read extracts Chart.yaml metadata and cg.json image mappings from a chart
// filesystem. Subcharts are discovered from Chart.yaml dependencies and
// populated in Chart.Subcharts recursively.
func Read(fsys fs.FS) (*Chart, error) {
	dirs := map[string]*dirEntry{}

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		dir := filepath.Dir(path)
		switch d.Name() {
		case HelmMetadataFilename, images.ChainguardChartMetadataFilename:
		default:
			return nil
		}

		f, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("opening %s: %w", path, err)
		}
		defer f.Close()

		if dirs[dir] == nil {
			dirs[dir] = &dirEntry{}
		}

		switch d.Name() {
		case HelmMetadataFilename:
			meta := &Meta{}
			if err := yaml.NewDecoder(f).Decode(meta); err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}
			dirs[dir].meta = meta

		case images.ChainguardChartMetadataFilename:
			m, err := images.Parse(f)
			if err != nil {
				return fmt.Errorf("parsing %s: %w", path, err)
			}
			dirs[dir].mapping = m
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Find the root chart (shallowest directory by path depth).
	var rootDir string
	rootDepth := -1
	for dir := range dirs {
		depth := strings.Count(dir, string(filepath.Separator))
		if rootDepth < 0 || depth < rootDepth {
			rootDir = dir
			rootDepth = depth
		}
	}
	root := dirs[rootDir]
	if root == nil || root.meta == nil {
		return nil, fmt.Errorf("Chart.yaml not found") //nolint:staticcheck // Chart.yaml is the actual filename
	}
	if root.mapping == nil {
		return nil, fmt.Errorf("%s not found", images.ChainguardChartMetadataFilename)
	}

	chart := &Chart{Meta: root.meta, Mapping: root.mapping}
	buildSubchartTree(rootDir, root, dirs, chart)

	return chart, nil
}

// buildSubchartTree populates chart.Subcharts by matching Chart.yaml
// dependencies to subdirectories.
func buildSubchartTree(parentDir string, parent *dirEntry, dirs map[string]*dirEntry, chart *Chart) {
	if parent.meta == nil || parent.mapping == nil {
		return
	}
	for _, dep := range parent.meta.Dependencies {
		// Helm stores subcharts in charts/{name}/ regardless of alias.
		subDir := filepath.Join(parentDir, "charts", dep.Name)
		sub, ok := dirs[subDir]
		if !ok || sub.mapping == nil {
			continue
		}
		subChart := &Chart{Meta: sub.meta, Mapping: sub.mapping}
		buildSubchartTree(subDir, sub, dirs, subChart)

		if chart.Subcharts == nil {
			chart.Subcharts = make(map[string]*Chart)
		}
		// Use the alias (if set) as the map key since Helm uses it
		// as the values namespace for the subchart.
		chart.Subcharts[dep.Key()] = subChart
	}
}
