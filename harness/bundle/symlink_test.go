/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle_test

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"chainguard.dev/sdk/harness/bundle"
)

//	SECURITY FINDING: Transitive symlink escape in checkLinkTarget {
//	  CWE: CWE-22
//	  Description: Lexical checks collapse .. before following preceding symlinks.
//	  Risk: HIGH
//	  Confidence: 10
//	  Attacker: A caller supplying a source bundle.
//	  Entry: Writer.Write and Validate
//	  Flow: tar Linkname -> checkLinkTarget -> accepted Entry -> filesystem resolution
//	  Exploit: a/a -> .. and a/b -> a/../outside make a/b resolve outside the tree.
//	  Impact: Accepted links can redirect filesystem access outside staging.
//	}
func TestSymlinkChains(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree []rawEntry
		want bundle.Reason
	}{
		{
			name: "backward escape",
			tree: []rawEntry{link("a/a", tar.TypeSymlink, ".."), link("a/b", tar.TypeSymlink, "a/../outside")},
			want: bundle.ReasonSymlinkEscape,
		},
		{
			name: "forward escape",
			tree: []rawEntry{link("a/a", tar.TypeSymlink, "z/../outside"), link("a/z", tar.TypeSymlink, "..")},
			want: bundle.ReasonSymlinkEscape,
		},
		{
			name: "symlink parent aliases an entry",
			tree: []rawEntry{link("a", tar.TypeSymlink, "."), link("a/b", tar.TypeSymlink, "..")},
			want: bundle.ReasonUnsafePath,
		},
		{
			name: "file parent",
			tree: []rawEntry{file("a", 0o644, ""), file("a/b", 0o644, "payload")},
			want: bundle.ReasonUnsafePath,
		},
		{
			name: "self cycle",
			tree: []rawEntry{link("a", tar.TypeSymlink, "a")},
			want: bundle.ReasonSymlinkEscape,
		},
		{
			name: "mutual cycle",
			tree: []rawEntry{link("a", tar.TypeSymlink, "b"), link("b", tar.TypeSymlink, "a")},
			want: bundle.ReasonSymlinkEscape,
		},
		{
			name: "contained forward chain",
			tree: []rawEntry{link("a", tar.TypeSymlink, "b/link"), link("b/link", tar.TypeSymlink, "../file"), file("file", 0o644, "safe")},
		},
		{
			name: "contained parent traversal",
			tree: []rawEntry{link("a/link", tar.TypeSymlink, "../file"), file("file", 0o644, "safe")},
		},
		{
			name: "repeated link without a cycle",
			tree: []rawEntry{link("a/up", tar.TypeSymlink, ".."), file("file", 0o644, "safe"), link("link", tar.TypeSymlink, "a/up/a/up/file")},
		},
		{
			name: "dangling inside the tree",
			tree: []rawEntry{link("a/link", tar.TypeSymlink, "../missing")},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkTree(t, tc.tree, tc.want)
		})
	}
}

func TestSymlinkChainDepth(t *testing.T) {
	for _, tc := range []struct {
		depth int
		want  bundle.Reason
	}{{depth: 40}, {depth: 41, want: bundle.ReasonSymlinkEscape}} {
		t.Run(fmt.Sprint(tc.depth), func(t *testing.T) {
			tree := make([]rawEntry, 0, tc.depth)
			for i := range tc.depth {
				tree = append(tree, link(fmt.Sprintf("link%02d", i), tar.TypeSymlink, fmt.Sprintf("link%02d", i+1)))
			}
			checkTree(t, tree, tc.want)
		})
	}
}

// The reader fixture bypasses Writer so a writer rejection cannot hide a
// missing check on the server path. Callers provide entries in path order.
func checkTree(t *testing.T, tree []rawEntry, want bundle.Reason) {
	t.Helper()
	es := make([]bundle.Entry, 0, len(tree))
	records := make([][4]string, 0, len(tree))
	m := goldenManifest()
	m.Format = bundle.Format
	for _, raw := range tree {
		e := bundle.Entry{Path: raw.h.Name, Mode: fs.FileMode(raw.h.Mode)}
		var payload string
		switch raw.h.Typeflag {
		case tar.TypeReg:
			e.Type, e.Size, e.Open = bundle.TypeFile, int64(len(raw.body)), opener(raw.body)
			payload = sha256hex([]byte(raw.body))
			m.Files++
			m.Bytes += e.Size
		case tar.TypeDir:
			e.Type, e.Path = bundle.TypeDir, strings.TrimSuffix(raw.h.Name, "/")
		case tar.TypeSymlink:
			e.Type, e.LinkTarget = bundle.TypeSymlink, raw.h.Linkname
			payload = e.LinkTarget
		default:
			t.Fatalf("unsupported test entry type %q", raw.h.Typeflag)
		}
		es = append(es, e)
		records = append(records, [4]string{e.Path, fmt.Sprintf("%04o", e.Mode), string(e.Type), payload})
	}
	m.TreeSHA256 = specTreeHash(records...)
	t.Run("writer", func(t *testing.T) {
		got, err := (&bundle.Writer{}).Write(t.Context(), io.Discard, goldenManifest(), es)
		if want != "" {
			if !isReason(err, want) || got != nil {
				t.Errorf("Write: got = %v, %v, want = nil, reason %q", got, err, want)
			}
			return
		}
		if err != nil {
			t.Fatalf("Write: %v", err)
		}
		if diff := cmp.Diff(&m, got); diff != "" {
			t.Errorf("manifest (-want +got):\n%s", diff)
		}
	})
	t.Run("reader", func(t *testing.T) {
		got, err := bundle.Validate(t.Context(), bytes.NewReader(rawBundle(t, &m, entries(tree...))))
		if want != "" {
			if !isReason(err, want) || got != nil {
				t.Errorf("Validate: got = %v, %v, want = nil, reason %q", got, err, want)
			}
			return
		}
		if err != nil {
			t.Fatalf("Validate: %v", err)
		}
		if diff := cmp.Diff(&m, got); diff != "" {
			t.Errorf("manifest (-want +got):\n%s", diff)
		}
	})
}

//	SECURITY FINDING: Case-insensitive alias bypasses the containment index {
//	  CWE: CWE-178
//	  Description: entryTree keys paths exactly, so "A" and "a" index separately.
//	  Risk: MEDIUM
//	  Attacker: A caller supplying a source bundle.
//	  Entry: Writer.Write and Validate
//	  Flow: tar Name -> entryTree.add -> accepted Entry -> filesystem resolution
//	  Exploit: A -> . and a/B -> ../outside make a/b/payload land outside on a
//	           case-insensitive extractor, where "a" opens "A" and "b" opens "B".
//	  Impact: Accepted entries can write outside staging on Windows and macOS.
//	}
func TestCaseInsensitiveAliases(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree []rawEntry
		want bundle.Reason
	}{
		{
			name: "symlink parent aliases a directory by case",
			tree: []rawEntry{
				link("A", tar.TypeSymlink, "."),
				link("a/B", tar.TypeSymlink, "../outside"),
				file("a/b/payload", 0o644, "x"),
			},
			want: bundle.ReasonUnsafePath,
		},
		{
			name: "file parent aliases a directory by case",
			tree: []rawEntry{file("A", 0o644, ""), file("a/b", 0o644, "payload")},
			want: bundle.ReasonUnsafePath,
		},
		{
			name: "chain resolves through a case alias",
			tree: []rawEntry{link("Up", tar.TypeSymlink, ".."), link("a/link", tar.TypeSymlink, "../up/../../outside")},
			want: bundle.ReasonSymlinkEscape,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkTree(t, tc.tree, tc.want)
		})
	}
}

// A link target is a path, so it gets the same depth ceiling an entry path has.
// Without one, a 4 KiB target carries ~2000 components and a 40-link chain
// expands to ~80000, which the decoded-size and ratio limits never see.
func TestLinkTargetDepthIsBounded(t *testing.T) {
	deep := strings.TrimSuffix(strings.Repeat("d/", 65), "/")
	checkTree(t, []rawEntry{link("a", tar.TypeSymlink, deep)}, bundle.ReasonTooDeep)
}

// The index keys a whole path while the chain walk keys one component at a
// time. Both must agree, or the lookup for "ΑΙ" misses and its chain is never
// walked even though it escapes.
func TestFoldedKeysSurviveComponentAtATimeWalk(t *testing.T) {
	checkTree(t, []rawEntry{
		link("a/a", tar.TypeSymlink, ".."),
		link("ΑΙ", tar.TypeSymlink, "a/a/../outside"),
	}, bundle.ReasonSymlinkEscape)
}

// "ΑΙ" and the case- and normalization-equivalent "ᾼ" are one name on a
// normalization-aware case-insensitive extractor, which follows the alias.
func TestNormalizationAndCaseCanonicalizeTogether(t *testing.T) {
	checkTree(t, []rawEntry{
		link("Z", tar.TypeSymlink, "a/ᾼ/../outside"),
		link("a/ΑΙ", tar.TypeSymlink, ".."),
	}, bundle.ReasonSymlinkEscape)
}

// Folding over NFC cannot reach a canonical form: "ᾴ" has no simple case
// mapping, so neither NFC nor SimpleFold ever moves it, while its decomposed
// spelling folds to the same name as "ΆΙ". Keying on NFD settles both in one
// pass. "ῄ"/"ΉΙ" and "ῴ"/"ΏΙ" split the same way.
func TestPrecomposedWithoutSimpleFoldAliases(t *testing.T) {
	for _, tc := range []struct{ name, composed, pair string }{
		{"alpha", "ᾴ", "ΆΙ"},
		{"eta", "ῄ", "ΉΙ"},
		{"omega", "ῴ", "ΏΙ"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkTree(t, []rawEntry{
				link("Z", tar.TypeSymlink, "a/"+tc.composed+"/../outside"),
				link("a/"+tc.pair, tar.TypeSymlink, ".."),
			}, bundle.ReasonSymlinkEscape)
		})
	}
}
