/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"math/rand/v2"
	"strconv"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/klauspost/compress/zstd"

	"chainguard.dev/sdk/harness/bundle"
)

// The tree the determinism goldens are taken over. It covers all three entry
// types, an executable file, a nested path and a symlink.
func goldenEntries() []bundle.Entry {
	return []bundle.Entry{
		fileEntry("cmd/run.sh", 0o700, "#!/bin/sh\nexec /opt/fenugreek --listen :9173\n"),
		fileEntry("README.md", 0o640, "# fenugreek\n\nA spice rack for cautious people.\n"),
		dirEntry("cmd", 0o750),
		symlinkEntry("docs/entry.md", "../README.md"),
		fileEntry("docs/.keep", 0o644, ""),
	}
}

const (
	// goldenTreeSHA256 pins tree_sha256 over goldenEntries. Both ends and any
	// reimplementation depend on this value byte for byte, so a change to the
	// record layout must show up here and be a coordinated one.
	goldenTreeSHA256 = "42c8725595c50bb82c6cc918ccef4891903989c161195cc01371f500e4d238e1"

	// goldenTarSHA256 pins the decompressed tar a pack of goldenEntries and
	// goldenManifest produces. The tar layer is what the format makes
	// reproducible, and this covers everything that depends on — entry order,
	// headers, zeroed metadata, the manifest's own bytes. The compressed bundle
	// is deliberately not pinned: zstd's output moves with the encoder's
	// version, so a golden over it would pin the dependency rather than the
	// format. This value is a given toolchain's; archive/tar is free to frame
	// the same tree differently in a later Go.
	goldenTarSHA256 = "580fa3ca403b996ee02d674863115a9021c58a245fb6f08fb941034c0230bcc8"
)

func goldenManifest() bundle.Manifest {
	return bundle.Manifest{
		Kind:   bundle.KindGit,
		Commit: "4f3a1c9e2b7d8a60513ce4f2a9b8d7c6e5f40312",
		Client: "chainctl v0.1.2",
	}
}

// TestWriteIsDeterministic is the golden-byte test: one tree always packs to
// one tar, whoever packs it and in whatever order intake enumerated it. The
// compressed bytes are compared within this one build, where the encoder is
// fixed; only the tar is pinned to a golden, because that is the layer the
// format makes reproducible.
func TestWriteIsDeterministic(t *testing.T) {
	ctx := t.Context()
	first, _ := pack(t, goldenManifest(), goldenEntries())

	// Re-packing the same tree, with the entries handed over in a different
	// order, produces the same bytes.
	shuffled := goldenEntries()
	shuffled[0], shuffled[3] = shuffled[3], shuffled[0]
	second, m := pack(t, goldenManifest(), shuffled)
	if !bytes.Equal(first, second) {
		t.Errorf("packing the same tree twice differed: %d bytes vs %d", len(first), len(second))
	}

	if m.TreeSHA256 != goldenTreeSHA256 {
		t.Errorf("tree_sha256: got = %s, want = %s", m.TreeSHA256, goldenTreeSHA256)
	}
	if got := sha256hex(untar(t, first)); got != goldenTarSHA256 {
		t.Errorf("tar bytes: got = %s, want = %s", got, goldenTarSHA256)
	}
	if _, err := bundle.Validate(ctx, bytes.NewReader(first)); err != nil {
		t.Errorf("Validate: %v", err)
	}
}

// TestTreeHashMatchesSpec recomputes tree_sha256 straight from the definition
// in the package doc, independently of the package's own implementation. A
// reimplementation on another end has only that text to go on.
func TestTreeHashMatchesSpec(t *testing.T) {
	_, m := pack(t, goldenManifest(), goldenEntries())

	// Sorted by path in byte order, with the manifest entry excluded.
	want := specTreeHash(
		[4]string{"README.md", "0644", "file", sha256hex([]byte("# fenugreek\n\nA spice rack for cautious people.\n"))},
		[4]string{"cmd", "0755", "dir", ""},
		[4]string{"cmd/run.sh", "0755", "file", sha256hex([]byte("#!/bin/sh\nexec /opt/fenugreek --listen :9173\n"))},
		[4]string{"docs/.keep", "0644", "file", sha256hex(nil)},
		[4]string{"docs/entry.md", "0777", "symlink", "../README.md"},
	)
	if want != m.TreeSHA256 {
		t.Errorf("tree_sha256: got = %s, want (from the spec) = %s", m.TreeSHA256, want)
	}

	// "A bundle with no entries hashes the domain field alone."
	_, empty := pack(t, goldenManifest(), nil)
	if want := specTreeHash(); want != empty.TreeSHA256 {
		t.Errorf("tree_sha256 of an empty tree: got = %s, want (from the spec) = %s", empty.TreeSHA256, want)
	}
}

// TestWriteNormalizes checks what a round trip must preserve and what it must
// flatten.
func TestWriteNormalizes(t *testing.T) {
	data, m := pack(t, goldenManifest(), goldenEntries())

	if diff := cmp.Diff(goldenManifest(), *m, cmpopts.IgnoreFields(bundle.Manifest{}, "Format", "Files", "Bytes", "TreeSHA256")); diff != "" {
		t.Errorf("manifest round trip (-want +got):\n%s", diff)
	}
	if m.Format != bundle.Format {
		t.Errorf("format: got = %q, want = %q", m.Format, bundle.Format)
	}
	// Three regular files holding 92 bytes; the directory, the symlink and the
	// manifest itself do not count.
	if m.Files != 3 || m.Bytes != 92 {
		t.Errorf("counts: got = %d files / %d bytes, want = 3 / 92", m.Files, m.Bytes)
	}

	got, entries, contents := read(t, data)
	if diff := cmp.Diff(m, got); diff != "" {
		t.Errorf("manifest (-want +got):\n%s", diff)
	}
	want := []bundle.Entry{
		{Path: "README.md", Type: bundle.TypeFile, Mode: 0o644, Size: 47},
		{Path: "cmd", Type: bundle.TypeDir, Mode: 0o755},
		{Path: "cmd/run.sh", Type: bundle.TypeFile, Mode: 0o755, Size: 45},
		{Path: "docs/.keep", Type: bundle.TypeFile, Mode: 0o644},
		{Path: "docs/entry.md", Type: bundle.TypeSymlink, Mode: 0o777, LinkTarget: "../README.md"},
	}
	if diff := cmp.Diff(want, entries, cmpopts.IgnoreFields(bundle.Entry{}, "Open")); diff != "" {
		t.Errorf("entries (-want +got):\n%s", diff)
	}
	if got := contents["cmd/run.sh"]; got != "#!/bin/sh\nexec /opt/fenugreek --listen :9173\n" {
		t.Errorf("content of cmd/run.sh: got = %q", got)
	}
}

// TestWriteRejects covers the rules the writer enforces on the entries intake
// hands it. Intake drops what the format cannot carry, with a warning; by the
// time entries reach the writer, anything unrepresentable is a bug and is
// refused rather than silently dropped.
func TestWriteRejects(t *testing.T) {
	for _, tc := range []struct {
		name    string
		limits  bundle.Limits
		m       bundle.Manifest
		entries []bundle.Entry
		want    bundle.Reason
	}{{
		name:    "path traversal",
		entries: []bundle.Entry{fileEntry("../escape.txt", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "traversal mid-path",
		entries: []bundle.Entry{fileEntry("a/../../b.txt", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "absolute path",
		entries: []bundle.Entry{fileEntry("/etc/shadow", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "dot segment",
		entries: []bundle.Entry{fileEntry("a/./b.txt", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "empty segment",
		entries: []bundle.Entry{fileEntry("a//b.txt", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "empty path",
		entries: []bundle.Entry{fileEntry("", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "NUL in path",
		entries: []bundle.Entry{fileEntry("a\x00b.txt", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "invalid UTF-8 path",
		entries: []bundle.Entry{fileEntry("a/\xff\xfe.txt", 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		// Legal POSIX bytes, and a traversal the moment a Windows extractor
		// splits on them.
		name:    "backslash in a path",
		entries: []bundle.Entry{fileEntry(`a\..\..\b.txt`, 0o644, "x")},
		want:    bundle.ReasonUnsafePath,
	}, {
		name:    "backslash in a link target",
		entries: []bundle.Entry{symlinkEntry("a/link", `..\..\etc\shadow`)},
		want:    bundle.ReasonSymlinkEscape,
	}, {
		name:    "duplicate path",
		entries: []bundle.Entry{fileEntry("a.txt", 0o644, "one"), fileEntry("a.txt", 0o644, "two")},
		want:    bundle.ReasonDuplicatePath,
	}, {
		name:    "reserved manifest path",
		entries: []bundle.Entry{fileEntry(bundle.ManifestPath, 0o644, "{}")},
		want:    bundle.ReasonReservedPath,
	}, {
		name:    "symlink escaping the tree",
		entries: []bundle.Entry{symlinkEntry("a/link", "../../etc/shadow")},
		want:    bundle.ReasonSymlinkEscape,
	}, {
		name:    "absolute symlink target",
		entries: []bundle.Entry{symlinkEntry("link", "/etc/shadow")},
		want:    bundle.ReasonSymlinkEscape,
	}, {
		name:    "empty symlink target",
		entries: []bundle.Entry{symlinkEntry("link", "")},
		want:    bundle.ReasonSymlinkEscape,
	}, {
		name:    "unknown entry type",
		entries: []bundle.Entry{{Path: "dev/null", Type: bundle.EntryType("device"), Mode: 0o644}},
		want:    bundle.ReasonEntryType,
	}, {
		name:    "zero entry type",
		entries: []bundle.Entry{{Path: "a.txt", Mode: 0o644}},
		want:    bundle.ReasonEntryType,
	}, {
		name:    "path too long",
		limits:  bundle.Limits{MaxPathBytes: 16},
		entries: []bundle.Entry{fileEntry("a/"+strings.Repeat("b", 32), 0o644, "x")},
		want:    bundle.ReasonPathTooLong,
	}, {
		name:    "path too deep",
		limits:  bundle.Limits{MaxPathDepth: 3},
		entries: []bundle.Entry{fileEntry("a/b/c/d.txt", 0o644, "x")},
		want:    bundle.ReasonTooDeep,
	}, {
		name:    "file too large",
		limits:  bundle.Limits{MaxFileBytes: 8},
		entries: []bundle.Entry{fileEntry("big.bin", 0o644, strings.Repeat("x", 9))},
		want:    bundle.ReasonFileTooLarge,
	}, {
		name:    "too many entries",
		limits:  bundle.Limits{MaxEntries: 3},
		entries: []bundle.Entry{fileEntry("a", 0o644, "1"), fileEntry("b", 0o644, "2"), fileEntry("c", 0o644, "3")},
		want:    bundle.ReasonTooManyEntries,
	}, {
		name:    "decoded stream too large",
		limits:  bundle.Limits{MaxDecodedBytes: 1024},
		entries: []bundle.Entry{fileEntry("big.bin", 0o644, strings.Repeat("x", 4096))},
		want:    bundle.ReasonStreamTooLarge,
	}, {
		name:    "compressed bundle too large",
		limits:  bundle.Limits{MaxCompressedBytes: 32},
		entries: goldenEntries(),
		want:    bundle.ReasonBundleTooLarge,
	}, {
		name:    "content shorter than declared",
		entries: []bundle.Entry{{Path: "a.txt", Type: bundle.TypeFile, Mode: 0o644, Size: 99, Open: opener("short")}},
		want:    bundle.ReasonSourceChanged,
	}, {
		// The writer opens a file twice, to digest it and to write it, and the
		// manifest carrying the first read's digest is already in the stream by
		// the time the second one happens. A change that keeps the length is
		// the one a length check alone would miss.
		name:    "content changed between the digest and the write",
		entries: []bundle.Entry{{Path: "a.txt", Type: bundle.TypeFile, Mode: 0o644, Size: 5, Open: rereads("alpha", "omega")}},
		want:    bundle.ReasonSourceChanged,
	}, {
		name:    "file without content",
		entries: []bundle.Entry{{Path: "a.txt", Type: bundle.TypeFile, Mode: 0o644}},
		want:    bundle.ReasonMalformed,
	}, {
		name:    "directory with content",
		entries: []bundle.Entry{{Path: "a", Type: bundle.TypeDir, Size: 4, Open: opener("oops")}},
		want:    bundle.ReasonMalformed,
	}, {
		name: "unknown kind",
		m:    bundle.Manifest{Kind: bundle.Kind("tarball"), Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		// The kind the format reserves for deferred work. A writer that learns
		// the flag before the format learns the kind is the case this refuses.
		name: "reserved source kind",
		m:    bundle.Manifest{Kind: bundle.KindSource, Commit: strings.Repeat("a", 40), Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "git kind without a base commit",
		m:    bundle.Manifest{Kind: bundle.KindGit, Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "folder kind with a commit",
		m:    bundle.Manifest{Kind: bundle.KindFolder, Commit: strings.Repeat("a", 40), Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "files kind with a commit",
		m:    bundle.Manifest{Kind: bundle.KindFiles, Commit: strings.Repeat("a", 40), Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "folder kind marked dirty",
		m:    bundle.Manifest{Kind: bundle.KindFolder, Dirty: true, Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "files kind marked dirty",
		m:    bundle.Manifest{Kind: bundle.KindFiles, Dirty: true, Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "folder kind marked include_ignored",
		m:    bundle.Manifest{Kind: bundle.KindFolder, IncludeIgnored: true, Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "files kind marked include_ignored",
		m:    bundle.Manifest{Kind: bundle.KindFiles, IncludeIgnored: true, Client: "chainctl v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "client with a control character",
		m:    bundle.Manifest{Kind: bundle.KindFolder, Client: "chainctl\x1b[31m v0.1.2"},
		want: bundle.ReasonInvalidManifest,
	}, {
		name: "no client",
		m:    bundle.Manifest{Kind: bundle.KindFolder},
		want: bundle.ReasonInvalidManifest,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			m := tc.m
			if m.Kind == "" {
				m = bundle.Manifest{Kind: bundle.KindFolder, Client: "chainctl v0.1.2"}
			}
			w := &bundle.Writer{Limits: tc.limits}
			if _, err := w.Write(t.Context(), io.Discard, m, tc.entries); !isReason(err, tc.want) {
				t.Errorf("Write: got = %v, want reason %q", err, tc.want)
			}
		})
	}
}

// TestNormalizeIsIntakesSeam covers the rule intake applies per candidate to
// decide what it drops with a warning, which is the same rule the writer
// refuses on. Without it intake would have to restate the format's rules and
// the two would drift.
func TestNormalizeIsIntakesSeam(t *testing.T) {
	for _, tc := range []struct {
		name     string
		entry    bundle.Entry
		wantMode fs.FileMode
		want     bundle.Reason
	}{{
		name:     "a file keeps only its execute bit",
		entry:    fileEntry("cmd/run.sh", 0o700, "x"),
		wantMode: 0o755,
	}, {
		name:     "a group-writable file flattens to 0644",
		entry:    fileEntry("notes.txt", 0o664, "x"),
		wantMode: 0o644,
	}, {
		name:     "a setuid file keeps no setuid bit",
		entry:    fileEntry("sbin/helper", 0o4755, "x"),
		wantMode: 0o755,
	}, {
		name:     "a directory flattens to 0755",
		entry:    dirEntry("cmd", 0o700),
		wantMode: 0o755,
	}, {
		name:     "a symlink inside the tree is kept",
		entry:    symlinkEntry("docs/entry.md", "../README.md"),
		wantMode: 0o777,
	}, {
		name:  "a symlink out of the tree is dropped",
		entry: symlinkEntry("docs/entry.md", "../../etc/shadow"),
		want:  bundle.ReasonSymlinkEscape,
	}, {
		name:  "a backslash in a name is dropped",
		entry: fileEntry(`weird\name.txt`, 0o644, "x"),
		want:  bundle.ReasonUnsafePath,
	}, {
		name:  "the reserved manifest path is dropped",
		entry: fileEntry(bundle.ManifestPath, 0o644, "{}"),
		want:  bundle.ReasonReservedPath,
	}, {
		name:  "a type the format cannot carry is dropped",
		entry: bundle.Entry{Path: "var/pipe", Type: bundle.EntryType("fifo")},
		want:  bundle.ReasonEntryType,
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := bundle.Normalize(tc.entry)
			if tc.want != "" {
				if !isReason(err, tc.want) {
					t.Fatalf("Normalize: got = %v, want reason %q", err, tc.want)
				}
				return
			}
			if err != nil {
				t.Fatalf("Normalize: %v", err)
			}
			if got.Mode != tc.wantMode {
				t.Errorf("mode: got = %04o, want = %04o", got.Mode, tc.wantMode)
			}
			if got.Path != tc.entry.Path || got.Type != tc.entry.Type {
				t.Errorf("entry: got = %q %q, want = %q %q", got.Path, got.Type, tc.entry.Path, tc.entry.Type)
			}
		})
	}
}

// TestGitKindFieldsSurviveARoundTrip pins the two kinds the fields are
// permitted on. Every other row about dirty and include_ignored is a rejection,
// so without this nothing shows they survive a round trip where they are meant
// to. include_ignored especially: it is the only thing telling a consumer that
// a bundle holding .env was a deliberate opt-in rather than the same checkout
// packed normally, so silently losing it downgrades a scan's description
// without failing anything.
func TestGitKindFieldsSurviveARoundTrip(t *testing.T) {
	m := goldenManifest()
	m.Dirty, m.IncludeIgnored = true, true
	data, written := pack(t, m, goldenEntries())
	if !written.Dirty || !written.IncludeIgnored || written.Commit != m.Commit {
		t.Errorf("written manifest: dirty = %v, include_ignored = %v, commit = %q",
			written.Dirty, written.IncludeIgnored, written.Commit)
	}
	read, _, _ := read(t, data)
	if !read.Dirty {
		t.Error("dirty did not survive the round trip")
	}
	if !read.IncludeIgnored {
		t.Error("include_ignored did not survive the round trip")
	}
}

// TestKindWireValues pins the kinds to the bytes that go on the wire. The Go
// constant can be renamed freely, but its value is what a bundle carries and
// what another end matches on, so changing one is a format change rather than
// a refactor. The reserved name is pinned for the same reason: its whole
// purpose is to still mean the same thing when a later client starts emitting
// it.
func TestKindWireValues(t *testing.T) {
	for _, tc := range []struct {
		kind bundle.Kind
		want string
	}{
		{bundle.KindGit, "git"},
		{bundle.KindFolder, "folder"},
		{bundle.KindFiles, "files"},
		{bundle.KindSource, "source"},
	} {
		if string(tc.kind) != tc.want {
			t.Errorf("wire value = %q, want %q", tc.kind, tc.want)
		}
	}
}

// TestWriteAndReadAgreeOnTheTarStream pins the two sides' accounting to each
// other. The writer's first pass only estimates the decoded stream — it has no
// manifest entry, no extended header for a long path and no end-of-archive
// marker — so a limit between the estimate and the stream's real size must
// still not let the writer emit a bundle the reader refuses.
func TestWriteAndReadAgreeOnTheTarStream(t *testing.T) {
	// goldenEntries estimate to 3584 tar bytes and pack to 5632; the limits
	// below straddle both numbers.
	for _, max := range []int64{3072, 4096, 5120, 5632, 8192} {
		t.Run(strconv.FormatInt(max, 10), func(t *testing.T) {
			lim := bundle.Limits{MaxDecodedBytes: max}
			var buf bytes.Buffer
			if _, err := (&bundle.Writer{Limits: lim}).Write(t.Context(), &buf, goldenManifest(), goldenEntries()); err != nil {
				if !isReason(err, bundle.ReasonStreamTooLarge) {
					t.Fatalf("Write: got = %v, want reason %q", err, bundle.ReasonStreamTooLarge)
				}
				return
			}
			if _, err := (&bundle.Reader{Limits: lim}).Read(t.Context(), &buf, nil); err != nil {
				t.Errorf("the writer emitted a bundle the reader refuses under the same limits: %v", err)
			}
		})
	}
}

func TestWriteAndReadAgreeOnCompressionRatio(t *testing.T) {
	noise := make([]byte, 1<<20)
	_, _ = rand.NewChaCha8([32]byte{1}).Read(noise)
	for _, tc := range []struct {
		name    string
		entries []bundle.Entry
		want    bundle.Reason
	}{
		{
			name:    "compressible below the floor",
			entries: []bundle.Entry{fileEntry("zeros.bin", 0o644, strings.Repeat("\x00", 1<<20))},
		},
		{
			name:    "compressible above the floor",
			entries: []bundle.Entry{fileEntry("zeros.bin", 0o644, strings.Repeat("\x00", 9<<20))},
			want:    bundle.ReasonDecompressRatio,
		},
		{
			name: "compressible prefix with an incompressible suffix",
			entries: []bundle.Entry{
				fileEntry("a-zeros.bin", 0o644, strings.Repeat("\x00", 9<<20)),
				fileEntry("z-noise.bin", 0o644, string(noise)),
			},
			want: bundle.ReasonDecompressRatio,
		},
		{
			name: "incompressible prefix pays for the whole stream",
			entries: []bundle.Entry{
				fileEntry("a-noise.bin", 0o644, string(noise)),
				fileEntry("z-zeros.bin", 0o644, strings.Repeat("\x00", 9<<20)),
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			m, err := (&bundle.Writer{}).Write(t.Context(), &out, goldenManifest(), tc.entries)
			if tc.want != "" {
				if !isReason(err, tc.want) || m != nil {
					t.Fatalf("Write: got = %v, %v, want = nil, reason %q", m, err, tc.want)
				}
				return
			}
			if err != nil {
				t.Fatalf("Write: %v", err)
			}
			got, err := bundle.Validate(t.Context(), bytes.NewReader(out.Bytes()))
			if err != nil {
				t.Fatalf("Validate writer output: %v", err)
			}
			if diff := cmp.Diff(m, got); diff != "" {
				t.Errorf("manifest (-want +got):\n%s", diff)
			}
		})
	}
}

func TestWriteAndReadAgreeOnDirectoryPathLimit(t *testing.T) {
	prefix := strings.Repeat(strings.Repeat("a", 127)+"/", 31)
	for _, size := range []int{4095, 4096, 4097} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			path := prefix + strings.Repeat("a", size-len(prefix))
			var out bytes.Buffer
			m, err := (&bundle.Writer{}).Write(t.Context(), &out, goldenManifest(), []bundle.Entry{dirEntry(path, 0o755)})
			if size > bundle.DefaultLimits().MaxPathBytes {
				if !isReason(err, bundle.ReasonPathTooLong) || m != nil {
					t.Fatalf("Write: got = %v, %v, want = nil, reason %q", m, err, bundle.ReasonPathTooLong)
				}
				return
			}
			if err != nil {
				t.Fatalf("Write: %v", err)
			}
			got, err := bundle.Validate(t.Context(), bytes.NewReader(out.Bytes()))
			if err != nil {
				t.Fatalf("Validate writer output: %v", err)
			}
			if diff := cmp.Diff(m, got); diff != "" {
				t.Errorf("manifest (-want +got):\n%s", diff)
			}
		})
	}
}

// TestEntriesSortByPathNotByTarName pins the ordering key. A directory is
// stored as "cmd/", which sorts after "cmd.go"; the entry path "cmd" sorts
// before it. A reader that ordered on the tar name would refuse this tree.
func TestEntriesSortByPathNotByTarName(t *testing.T) {
	data, _ := pack(t, goldenManifest(), []bundle.Entry{
		fileEntry("cmd/main.go", 0o644, "package main\n"),
		fileEntry("cmd.go", 0o644, "package widgets\n"),
		dirEntry("cmd", 0o755),
	})
	_, entries, _ := read(t, data)
	got := make([]string, 0, len(entries))
	for _, e := range entries {
		got = append(got, e.Path)
	}
	if diff := cmp.Diff([]string{"cmd", "cmd.go", "cmd/main.go"}, got); diff != "" {
		t.Errorf("entry order (-want +got):\n%s", diff)
	}
}

// TestTreeHashFramesSeparatorsInAPath pins the clause a reimplementation is
// most likely to get wrong: a newline inside a path is an ordinary byte the
// length prefix already covers, not a field boundary. The hash's other
// separator, the colon, cannot appear in a conforming path at all -- it is
// refused as an NTFS alternate-data-stream selector -- so a newline is what
// remains to test the framing with.
func TestTreeHashFramesSeparatorsInAPath(t *testing.T) {
	const path = "notes/5hello\nworld.txt"
	data, m := pack(t, goldenManifest(), []bundle.Entry{fileEntry(path, 0o644, "x")})

	if want := specTreeHash([4]string{path, "0644", "file", sha256hex([]byte("x"))}); want != m.TreeSHA256 {
		t.Errorf("tree_sha256: got = %s, want (from the spec) = %s", m.TreeSHA256, want)
	}
	if _, err := bundle.Validate(t.Context(), bytes.NewReader(data)); err != nil {
		t.Errorf("Validate: %v", err)
	}
}

func TestWriteHonoursContext(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if _, err := (&bundle.Writer{}).Write(ctx, io.Discard, goldenManifest(), goldenEntries()); !errors.Is(err, context.Canceled) {
		t.Errorf("Write: got = %v, want = %v", err, context.Canceled)
	}
}

func TestWritePreservesDestinationError(t *testing.T) {
	want := errors.New("destination failed")
	pr, pw := io.Pipe()
	t.Cleanup(func() { _ = pw.Close() })
	_ = pr.CloseWithError(want)
	m, err := (&bundle.Writer{}).Write(t.Context(), pw, goldenManifest(), goldenEntries())
	if !errors.Is(err, want) || m != nil {
		t.Errorf("Write: got = %v, %v, want = nil, %v", m, err, want)
	}
}

// pack writes a bundle and fails the test if it is rejected.
func pack(t testing.TB, m bundle.Manifest, entries []bundle.Entry) ([]byte, *bundle.Manifest) {
	t.Helper()
	var buf bytes.Buffer
	got, err := (&bundle.Writer{}).Write(t.Context(), &buf, m, entries)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}
	return buf.Bytes(), got
}

// read streams a bundle back, returning its manifest, entries and file
// contents.
func read(t *testing.T, data []byte) (*bundle.Manifest, []bundle.Entry, map[string]string) {
	t.Helper()
	var entries []bundle.Entry
	contents := make(map[string]string)
	m, err := (&bundle.Reader{}).Read(t.Context(), bytes.NewReader(data), func(e bundle.Entry, body io.Reader) error {
		entries = append(entries, e)
		b, err := io.ReadAll(body)
		if err != nil {
			return err
		}
		contents[e.Path] = string(b)
		return nil
	})
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	return m, entries, contents
}

func fileEntry(path string, mode fs.FileMode, body string) bundle.Entry {
	return bundle.Entry{
		Path: path,
		Type: bundle.TypeFile,
		Mode: mode,
		Size: int64(len(body)),
		Open: opener(body),
	}
}

func dirEntry(path string, mode fs.FileMode) bundle.Entry {
	return bundle.Entry{Path: path, Type: bundle.TypeDir, Mode: mode}
}

func symlinkEntry(path, target string) bundle.Entry {
	return bundle.Entry{Path: path, Type: bundle.TypeSymlink, Mode: 0o777, LinkTarget: target}
}

func opener(body string) func() (io.ReadCloser, error) {
	return func() (io.ReadCloser, error) { return io.NopCloser(strings.NewReader(body)), nil }
}

// rereads is a source that changes between the writer's two reads of it, the
// first returning first and every later one second.
func rereads(first, second string) func() (io.ReadCloser, error) {
	var opened int
	return func() (io.ReadCloser, error) {
		opened++
		if opened == 1 {
			return io.NopCloser(strings.NewReader(first)), nil
		}
		return io.NopCloser(strings.NewReader(second)), nil
	}
}

// isReason reports whether err is a rejection carrying want.
func isReason(err error, want bundle.Reason) bool {
	e, ok := errors.AsType[*bundle.Error](err)
	return ok && e.Reason == want
}

func sha256hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// untar decompresses a bundle to its tar bytes.
func untar(t testing.TB, data []byte) []byte {
	t.Helper()
	zr, err := zstd.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("zstd reader: %v", err)
	}
	defer zr.Close()
	raw, err := io.ReadAll(zr)
	if err != nil {
		t.Fatalf("decompressing: %v", err)
	}
	return raw
}
