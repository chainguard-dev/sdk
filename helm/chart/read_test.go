/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package chart

import (
	"strings"
	"testing"
	"testing/fstest"

	"chainguard.dev/sdk/helm/images"
	"github.com/google/go-cmp/cmp"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		fsys    fstest.MapFS
		want    *Chart
		wantErr string
	}{
		{
			name: "simple chart",
			fsys: fstest.MapFS{
				"mychart/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: mychart\nversion: 1.0.0\n")},
				"mychart/cg.json":    &fstest.MapFile{Data: []byte(`{"images":{"nginx":{"values":{"image":"${registry}"}}}}`)},
			},
			want: &Chart{
				Meta: &Meta{APIVersion: "v2", Name: "mychart", Version: "1.0.0"},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"nginx": {Values: map[string]any{"image": "${registry}"}},
					},
				},
			},
		},
		{
			name: "chart with subchart",
			fsys: fstest.MapFS{
				"parent/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: parent\nversion: 1.0.0\ndependencies:\n  - name: child\n",
				)},
				"parent/cg.json":                 &fstest.MapFile{Data: []byte(`{"images":{"app":{"values":{"image":"${registry}"}}}}`)},
				"parent/charts/child/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: child\nversion: 2.0.0\n")},
				"parent/charts/child/cg.json":    &fstest.MapFile{Data: []byte(`{"images":{"db":{"values":{"image":"${registry}"}}}}`)},
			},
			want: &Chart{
				Meta: &Meta{
					APIVersion:   "v2",
					Name:         "parent",
					Version:      "1.0.0",
					Dependencies: []Dependency{{Name: "child"}},
				},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"app": {Values: map[string]any{"image": "${registry}"}},
					},
				},
				Subcharts: map[string]*Chart{
					"child": {
						Meta:    &Meta{APIVersion: "v2", Name: "child", Version: "2.0.0"},
						Mapping: &images.Mapping{Images: map[string]*images.Image{"db": {Values: map[string]any{"image": "${registry}"}}}},
					},
				},
			},
		},
		{
			name: "nested subcharts",
			fsys: fstest.MapFS{
				"root/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: root\nversion: 1.0.0\ndependencies:\n  - name: mid\n",
				)},
				"root/cg.json": &fstest.MapFile{Data: []byte(`{"images":{"top":{"values":{"k":"${registry}"}}}}`)},
				"root/charts/mid/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: mid\nversion: 1.0.0\ndependencies:\n  - name: leaf\n",
				)},
				"root/charts/mid/cg.json":                &fstest.MapFile{Data: []byte(`{"images":{"svc":{"values":{"k":"${registry}"}}}}`)},
				"root/charts/mid/charts/leaf/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: leaf\nversion: 1.0.0\n")},
				"root/charts/mid/charts/leaf/cg.json":    &fstest.MapFile{Data: []byte(`{"images":{"db":{"values":{"k":"${registry}"}}}}`)},
			},
			want: &Chart{
				Meta: &Meta{
					APIVersion:   "v2",
					Name:         "root",
					Version:      "1.0.0",
					Dependencies: []Dependency{{Name: "mid"}},
				},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"top": {Values: map[string]any{"k": "${registry}"}},
					},
				},
				Subcharts: map[string]*Chart{
					"mid": {
						Meta: &Meta{
							APIVersion:   "v2",
							Name:         "mid",
							Version:      "1.0.0",
							Dependencies: []Dependency{{Name: "leaf"}},
						},
						Mapping: &images.Mapping{
							Images: map[string]*images.Image{
								"svc": {Values: map[string]any{"k": "${registry}"}},
							},
						},
						Subcharts: map[string]*Chart{
							"leaf": {
								Meta:    &Meta{APIVersion: "v2", Name: "leaf", Version: "1.0.0"},
								Mapping: &images.Mapping{Images: map[string]*images.Image{"db": {Values: map[string]any{"k": "${registry}"}}}},
							},
						},
					},
				},
			},
		},
		{
			name: "alias uses alias as key",
			fsys: fstest.MapFS{
				"app/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: app\nversion: 1.0.0\ndependencies:\n  - name: redis-ha\n    alias: cache\n",
				)},
				"app/cg.json":                    &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
				"app/charts/redis-ha/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: redis-ha\nversion: 1.0.0\n")},
				"app/charts/redis-ha/cg.json":    &fstest.MapFile{Data: []byte(`{"images":{"redis":{"values":{"k":"${registry}"}}}}`)},
			},
			want: &Chart{
				Meta: &Meta{
					APIVersion:   "v2",
					Name:         "app",
					Version:      "1.0.0",
					Dependencies: []Dependency{{Name: "redis-ha", Alias: "cache"}},
				},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"main": {Values: map[string]any{"k": "${registry}"}},
					},
				},
				Subcharts: map[string]*Chart{
					"cache": {
						Meta:    &Meta{APIVersion: "v2", Name: "redis-ha", Version: "1.0.0"},
						Mapping: &images.Mapping{Images: map[string]*images.Image{"redis": {Values: map[string]any{"k": "${registry}"}}}},
					},
				},
			},
		},
		{
			name: "subchart without cg.json is skipped",
			fsys: fstest.MapFS{
				"app/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: app\nversion: 1.0.0\ndependencies:\n  - name: common\n",
				)},
				"app/cg.json":                  &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
				"app/charts/common/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: common\nversion: 1.0.0\ntype: library\n")},
			},
			want: &Chart{
				Meta: &Meta{
					APIVersion:   "v2",
					Name:         "app",
					Version:      "1.0.0",
					Dependencies: []Dependency{{Name: "common"}},
				},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"main": {Values: map[string]any{"k": "${registry}"}},
					},
				},
			},
		},
		{
			name: "dependency not in charts/ dir is skipped",
			fsys: fstest.MapFS{
				"app/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: app\nversion: 1.0.0\ndependencies:\n  - name: missing\n",
				)},
				"app/cg.json": &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
			},
			want: &Chart{
				Meta: &Meta{
					APIVersion:   "v2",
					Name:         "app",
					Version:      "1.0.0",
					Dependencies: []Dependency{{Name: "missing"}},
				},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"main": {Values: map[string]any{"k": "${registry}"}},
					},
				},
			},
		},
		{
			name: "extra files are ignored",
			fsys: fstest.MapFS{
				"mychart/Chart.yaml":            &fstest.MapFile{Data: []byte("apiVersion: v2\nname: mychart\nversion: 1.0.0\n")},
				"mychart/cg.json":               &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
				"mychart/values.yaml":           &fstest.MapFile{Data: []byte("key: value\n")},
				"mychart/templates/deploy.yaml": &fstest.MapFile{Data: []byte("kind: Deployment\n")},
				"mychart/README.md":             &fstest.MapFile{Data: []byte("# readme\n")},
			},
			want: &Chart{
				Meta: &Meta{APIVersion: "v2", Name: "mychart", Version: "1.0.0"},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"main": {Values: map[string]any{"k": "${registry}"}},
					},
				},
			},
		},
		{
			name: "optional meta fields",
			fsys: fstest.MapFS{
				"mychart/Chart.yaml": &fstest.MapFile{Data: []byte(
					"apiVersion: v2\nname: mychart\nversion: 1.0.0\nappVersion: 3.5.0\ndescription: A chart\n",
				)},
				"mychart/cg.json": &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
			},
			want: &Chart{
				Meta: &Meta{APIVersion: "v2", Name: "mychart", Version: "1.0.0", AppVersion: "3.5.0", Description: "A chart"},
				Mapping: &images.Mapping{
					Images: map[string]*images.Image{
						"main": {Values: map[string]any{"k": "${registry}"}},
					},
				},
			},
		},
		{
			name: "missing Chart.yaml",
			fsys: fstest.MapFS{
				"mychart/cg.json": &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
			},
			wantErr: "Chart.yaml not found",
		},
		{
			name: "missing cg.json",
			fsys: fstest.MapFS{
				"mychart/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: mychart\nversion: 1.0.0\n")},
			},
			wantErr: "not found",
		},
		{
			name: "invalid Chart.yaml",
			fsys: fstest.MapFS{
				"mychart/Chart.yaml": &fstest.MapFile{Data: []byte("{{not yaml}}")},
				"mychart/cg.json":    &fstest.MapFile{Data: []byte(`{"images":{"main":{"values":{"k":"${registry}"}}}}`)},
			},
			wantErr: "parsing",
		},
		{
			name: "invalid cg.json",
			fsys: fstest.MapFS{
				"mychart/Chart.yaml": &fstest.MapFile{Data: []byte("apiVersion: v2\nname: mychart\nversion: 1.0.0\n")},
				"mychart/cg.json":    &fstest.MapFile{Data: []byte(`{not json}`)},
			},
			wantErr: "parsing",
		},
		{
			name:    "empty filesystem",
			fsys:    fstest.MapFS{},
			wantErr: "Chart.yaml not found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Read(tc.fsys)
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
				t.Fatalf("Read: %v", err)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
