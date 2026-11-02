/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle_test

import (
	"archive/tar"
	"testing"

	"chainguard.dev/sdk/harness/bundle"
)

// The manifest never enters the tree index, so nothing but the reserved-path
// rule stops an entry landing on it, or in its directory, once a host that
// folds has extracted the bundle. Each tree goes through both the writer and
// Validate, and none has an entry at the manifest's exact path, so an exact
// comparison with ManifestPath refuses none of them.
func TestManifestDirectoryIsReserved(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree []rawEntry
		want bundle.Reason
	}{{
		name: "the manifest's directory as a file",
		tree: []rawEntry{file(".chainguard-harness", 0o644, "")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "the manifest's directory as a symlink",
		tree: []rawEntry{link(".chainguard-harness", tar.TypeSymlink, "docs")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "the manifest's directory as a file, by case",
		tree: []rawEntry{file(".Chainguard-Harness", 0o644, "")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "the manifest's directory declared under another case",
		tree: []rawEntry{dir(".Chainguard-Harness")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a file under the manifest's directory spelled another case",
		tree: []rawEntry{file(".Chainguard-Harness/other", 0o644, "x")},
		want: bundle.ReasonReservedPath,
	}, {
		// U+017F LATIN SMALL LETTER LONG S folds with "s".
		name: "a file under the manifest's directory spelled with a folding rune",
		tree: []rawEntry{file(".chainguard-harneſs/other", 0o644, "x")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "the manifest by case",
		tree: []rawEntry{file(".chainguard-harness/MANIFEST.json", 0o644, "{}")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a path under the manifest",
		tree: []rawEntry{file(".chainguard-harness/manifest.json/x", 0o644, "x")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a symlink to the manifest",
		tree: []rawEntry{link("m", tar.TypeSymlink, ".chainguard-harness/manifest.json")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a symlink to the manifest by case",
		tree: []rawEntry{link("m", tar.TypeSymlink, ".CHAINGUARD-HARNESS/Manifest.json")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a symlink under the manifest",
		tree: []rawEntry{link("m", tar.TypeSymlink, ".chainguard-harness/manifest.json/x")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a symlink to the manifest's directory",
		tree: []rawEntry{link("m", tar.TypeSymlink, ".chainguard-harness")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a symlink beside the manifest to it",
		tree: []rawEntry{dir(".chainguard-harness"), link(".chainguard-harness/latest", tar.TypeSymlink, "manifest.json")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "a symlink chain to the manifest",
		tree: []rawEntry{link("a", tar.TypeSymlink, "b/manifest.json"), link("b", tar.TypeSymlink, "c"), link("c", tar.TypeSymlink, ".chainguard-harness")},
		want: bundle.ReasonReservedPath,
	}, {
		name: "the manifest's directory declared",
		tree: []rawEntry{dir(".chainguard-harness")},
	}, {
		name: "other files in the manifest's directory",
		tree: []rawEntry{dir(".chainguard-harness"), file(".chainguard-harness/config.yaml", 0o644, "x"), file(".chainguard-harness/manifest.json.bak", 0o644, "x")},
	}, {
		name: "a file in the manifest's directory, which it implies",
		tree: []rawEntry{file(".chainguard-harness/config.yaml", 0o644, "x")},
	}, {
		name: "a manifest under a directory that only starts with the reserved name",
		tree: []rawEntry{file(".chainguard-harness-old/manifest.json", 0o644, "{}")},
	}, {
		name: "a symlink to another file in the manifest's directory",
		tree: []rawEntry{link("cfg", tar.TypeSymlink, ".chainguard-harness/config.yaml"), link("m", tar.TypeSymlink, ".chainguard-harness/manifest.json.bak")},
	}, {
		name: "a symlink through the manifest's directory and out",
		tree: []rawEntry{link("s", tar.TypeSymlink, ".chainguard-harness/../src")},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			checkTree(t, tc.tree, tc.want)
		})
	}
}

// Intake drops what Normalize rejects, so the path rules have to hold for one
// entry on its own. A symlink's target is a whole-tree question and passes.
func TestNormalizeRefusesReservedPaths(t *testing.T) {
	for _, tc := range []struct {
		name  string
		entry bundle.Entry
		want  bundle.Reason
	}{{
		name:  "the manifest's directory as a file",
		entry: fileEntry(".chainguard-harness", 0o644, ""),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "the manifest's directory as a symlink",
		entry: symlinkEntry(".chainguard-harness", "docs"),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "the manifest's directory under another case",
		entry: dirEntry(".CHAINGUARD-HARNESS", 0o755),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "a file under the manifest's directory spelled another case",
		entry: fileEntry(".Chainguard-Harness/other", 0o644, "x"),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "the manifest by case",
		entry: fileEntry(".chainguard-harness/Manifest.JSON", 0o644, "{}"),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "the manifest's directory",
		entry: dirEntry(".chainguard-harness", 0o755),
	}, {
		name:  "another file in the manifest's directory",
		entry: fileEntry(".chainguard-harness/config.yaml", 0o644, "x"),
	}, {
		name:  "a symlink to the manifest",
		entry: symlinkEntry("m", ".chainguard-harness/manifest.json"),
	}} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := bundle.Normalize(tc.entry)
			if tc.want == "" {
				if err != nil {
					t.Errorf("Normalize: %v", err)
				}
				return
			}
			if !isReason(err, tc.want) {
				t.Errorf("Normalize: got = %v, want reason %q", err, tc.want)
			}
		})
	}
}
