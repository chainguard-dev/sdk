/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle_test

import (
	"archive/tar"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"chainguard.dev/sdk/harness/bundle"
)

//	SECURITY FINDING: Windows drive paths in checkPath and checkLinkTarget {
//	  CWE: CWE-22
//	  Description: POSIX checks accept drive-qualified paths with Windows semantics.
//	  Risk: MEDIUM
//	  Confidence: 10
//	  Attacker: A caller supplying a source bundle.
//	  Entry: Normalize, Writer.Write and Validate
//	  Flow: Entry.Path or tar Linkname -> POSIX path checks -> Windows consumer
//	  Exploit: C:/outside addresses a drive root instead of the bundle tree.
//	  Impact: Accepted paths or links can address files outside staging on Windows.
//	}
func TestWindowsDrivePathsRejected(t *testing.T) {
	for _, p := range []string{"C:/outside", "c:/outside", "Z:/outside", "z:/outside", "C:outside", "c:outside", "C:", "c:"} {
		t.Run(fmt.Sprintf("%q", p), func(t *testing.T) {
			checkPathForms(t, p, bundle.ReasonUnsafePath, bundle.ReasonSymlinkEscape)
		})
	}
}

//	SECURITY FINDING: Windows path aliases in checkPath and checkLinkTarget {
//	  CWE: CWE-22
//	  Description: Trailing spaces and dots allow Windows to interpret names differently.
//	  Risk: MEDIUM
//	  Confidence: 9
//	  Attacker: A caller supplying a source bundle.
//	  Entry: Normalize, Writer.Write and Validate
//	  Flow: Entry.Path or tar Linkname -> literal component checks -> Windows consumer
//	  Exploit: dir./payload may address dir/payload despite having a distinct bundle path.
//	  Impact: Filesystem aliases can invalidate the validated tree's containment checks.
//	}
func TestWindowsPathSuffixesRejected(t *testing.T) {
	for _, p := range []string{
		".. /payload", "a/.. /payload", ". /payload", "..  /payload",
		"dir./payload", "dir /payload", "dir. /payload", ".../payload",
		"payload.", "payload ", "payload..", "payload .", ".. ", ". ",
	} {
		t.Run(fmt.Sprintf("%q", p), func(t *testing.T) {
			checkPathForms(t, p, bundle.ReasonUnsafePath, bundle.ReasonSymlinkEscape)
		})
	}
}

func TestRelativePathSpellingsPreserved(t *testing.T) {
	for _, p := range []string{"dir name/file.txt", ".hidden", "file..txt", "notes/5hello\nworld.txt", "résumé.txt"} {
		t.Run(fmt.Sprintf("%q", p), func(t *testing.T) {
			checkPathForms(t, p, "", "")
		})
	}
}

func TestRelativeSymlinkDotSegmentsPreserved(t *testing.T) {
	for _, target := range []string{".", "..", "./file", "../file", "../dir/./file"} {
		t.Run(fmt.Sprintf("%q", target), func(t *testing.T) {
			e := symlinkEntry("dir/link", target)
			got, err := bundle.Normalize(e)
			if err != nil {
				t.Fatalf("Normalize: %v", err)
			}
			if diff := cmp.Diff(e, got); diff != "" {
				t.Errorf("entry (-want +got):\n%s", diff)
			}
			checkTree(t, []rawEntry{link(e.Path, tar.TypeSymlink, target)}, "")
		})
	}
}

func checkPathForms(t *testing.T, p string, pathReason, targetReason bundle.Reason) {
	t.Helper()
	for _, tc := range []struct {
		name  string
		entry bundle.Entry
		raw   rawEntry
		want  bundle.Reason
	}{
		{name: "path", entry: fileEntry(p, 0o644, "payload"), raw: file(p, 0o644, "payload"), want: pathReason},
		{name: "target", entry: symlinkEntry("dir/link", p), raw: link("dir/link", tar.TypeSymlink, p), want: targetReason},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("normalize", func(t *testing.T) {
				got, err := bundle.Normalize(tc.entry)
				want := tc.entry
				if tc.want != "" {
					want = bundle.Entry{}
					if !isReason(err, tc.want) {
						t.Errorf("Normalize: got = %v, want reason %q", err, tc.want)
					}
				} else if err != nil {
					t.Fatalf("Normalize: %v", err)
				}
				if diff := cmp.Diff(want, got, cmpopts.IgnoreFields(bundle.Entry{}, "Open")); diff != "" {
					t.Errorf("entry (-want +got):\n%s", diff)
				}
			})
			checkTree(t, []rawEntry{tc.raw}, tc.want)
		})
	}
}

//	SECURITY FINDING: Undeclared bytes smuggled in a zstd skippable frame {
//	  CWE: CWE-20
//	  Description: The decoder silently skips 0x184D2A5x frames, so their bytes
//	               are never decoded and no entry declares them.
//	  Risk: MEDIUM
//	  Attacker: A caller supplying a source bundle.
//	  Entry: Validate and Reader.Read
//	  Flow: compressed stream -> zstd decoder -> frame skipped -> tar reader
//	  Exploit: A skippable frame before or after the real one carries up to
//	           MaxCompressedBytes that tree_sha256 and the trailer never cover.
//	  Impact: A bundle can carry content no manifest accounts for, and the
//	          extra bytes loosen the ratio guard by inflating the divisor.
//	}
func TestZstdFrameSmuggling(t *testing.T) {
	var buf bytes.Buffer
	if _, err := (&bundle.Writer{}).Write(t.Context(), &buf, goldenManifest(), goldenEntries()); err != nil {
		t.Fatalf("Write: %v", err)
	}
	good := buf.Bytes()
	payload := bytes.Repeat([]byte("SMUGGLED"), 1024)

	skippable := func(b []byte) []byte {
		f := make([]byte, 8, 8+len(b))
		binary.LittleEndian.PutUint32(f[0:4], 0x184D2A50)
		binary.LittleEndian.PutUint32(f[4:8], uint32(len(b)))
		return append(f, b...)
	}
	second := func() []byte { return append(append([]byte{}, good...), good...) }

	for _, tc := range []struct {
		name string
		data []byte
	}{
		{"leading skippable frame", append(skippable(payload), good...)},
		{"trailing skippable frame", append(append([]byte{}, good...), skippable(payload)...)},
		{"a second standard frame", second()},
		{"trailing garbage", append(append([]byte{}, good...), []byte("tail")...)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := bundle.Validate(t.Context(), bytes.NewReader(tc.data)); !isReason(err, bundle.ReasonMalformed) {
				t.Errorf("Validate: got = %v, want reason %q", err, bundle.ReasonMalformed)
			}
		})
	}

	// The guard must not reject what the writer actually emits.
	if _, err := bundle.Validate(t.Context(), bytes.NewReader(good)); err != nil {
		t.Errorf("Validate rejected a bundle the writer produced: %v", err)
	}
}

// Windows resolves a reserved device name anywhere and ignores the extension,
// and a colon selects an NTFS alternate data stream, so either makes a
// validated entry land somewhere other than where it says.
func TestWindowsDeviceAndStreamNamesRejected(t *testing.T) {
	for _, p := range []string{
		"NUL", "nul.txt", "a/CON", "a/com1.log", "PRN.tar.gz", "aux",
		"notes.txt:hidden", "a/b:stream",
		"COM¹", "com².txt", "a/LPT³", "lpt¹.log",
		"CONIN$", "conout$.txt", "a/ConIn$",
	} {
		t.Run(p, func(t *testing.T) {
			checkPathForms(t, p, bundle.ReasonUnsafePath, bundle.ReasonSymlinkEscape)
		})
	}
	// Names that merely start with a device name stay legal.
	for _, p := range []string{
		"console.txt", "a/nullable", "comic.png", "auxiliary/x",
		"conin.txt", "a/conout", "com¹x", "lpt²ab/x",
	} {
		t.Run("legal/"+p, func(t *testing.T) {
			checkPathForms(t, p, "", "")
		})
	}
}

// Two paths that fold together address one object on a case-insensitive host,
// so the exact-duplicate rule never sees them.
func TestFoldedPathCollisionsRejected(t *testing.T) {
	for _, tc := range []struct {
		name string
		tree []rawEntry
		want bundle.Reason
	}{
		{"files differing by case", []rawEntry{file("A", 0o644, "one"), file("a", 0o644, "two")}, bundle.ReasonDuplicatePath},
		// Caught a step earlier: the ancestor "a" is already the file "A".
		{"a directory colliding with a file", []rawEntry{file("A", 0o644, "x"), file("a/b", 0o644, "y")}, bundle.ReasonUnsafePath},
		{"symlinks differing by case", []rawEntry{link("L", tar.TypeSymlink, "f"), link("l", tar.TypeSymlink, "g"), file("f", 0o644, "x"), file("g", 0o644, "y")}, bundle.ReasonDuplicatePath},
		// The reverse order: the directory is implicit, so it is never an entry
		// the later file can be checked against as an ancestor. "A" sorts before
		// "a", so an archive can always present the collision this way round.
		{"a file colliding with a directory", []rawEntry{file("A/b", 0o644, "x"), file("a", 0o644, "y")}, bundle.ReasonUnsafePath},
		{"a symlink colliding with a directory", []rawEntry{link("A/B/link", tar.TypeSymlink, "../../outside"), link("a", tar.TypeSymlink, ".")}, bundle.ReasonUnsafePath},
		// Both spellings are directories, so no entry escapes -- but one physical
		// directory holds both children on a case-insensitive host, and the paths
		// that land there no longer reproduce the records tree_sha256 pins.
		// Uppercase sorts first in ASCII, so a pure case difference can only be
		// presented in this order; the implicit and explicit forms below are what
		// vary. Neither spelling is an entry of its own here.
		{"implicit directories differing by case", []rawEntry{file("A/x", 0o644, "one"), file("a/y", 0o644, "two")}, bundle.ReasonDuplicatePath},
		{"an implicit directory before an explicit one", []rawEntry{file("A/b", 0o644, "x"), dir("a")}, bundle.ReasonDuplicatePath},
		{"an explicit directory before an implicit one", []rawEntry{dir("A"), file("a/b", 0o644, "x")}, bundle.ReasonDuplicatePath},
		{"nested directories differing by case", []rawEntry{file("p/A/x", 0o644, "one"), file("p/a/y", 0o644, "two")}, bundle.ReasonDuplicatePath},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkTree(t, tc.tree, tc.want)
		})
	}
	// Folding must not merge names that are distinct on every host, and a
	// directory reached by many entries is one spelling, not a collision.
	checkTree(t, []rawEntry{
		file("README", 0o644, "x"),
		file("docs/readme.md", 0o644, "y"),
		file("docsx/y", 0o644, "z"),
		dir("src"),
		file("src/a.go", 0o644, "a"),
		file("src/b.go", 0o644, "b"),
	}, "")
}

// Case folding, not lowercasing: Greek sigma has two lowercase forms, and a
// case-insensitive filesystem treats all three as one letter. strings.ToLower
// maps "Σ" to "σ" but leaves "ς" alone, so the two spellings key apart.
func TestSigmaCaseAliasesRejected(t *testing.T) {
	// "Σ" (ce a3) sorts before "ς" (cf 82), the order the reader requires.
	for _, tc := range []struct {
		name string
		tree []rawEntry
		want bundle.Reason
	}{
		{"files differing by sigma form", []rawEntry{file("Σ", 0o644, "one"), file("ς", 0o644, "two")}, bundle.ReasonDuplicatePath},
		{"symlink parent aliases a directory by sigma form", []rawEntry{
			link("Σ", tar.TypeSymlink, "."),
			link("ς/B", tar.TypeSymlink, "../outside"),
		}, bundle.ReasonUnsafePath},
	} {
		t.Run(tc.name, func(t *testing.T) {
			checkTree(t, tc.tree, tc.want)
		})
	}
}

// APFS compares normalization-insensitively, so a composed and a decomposed
// spelling of one name are a single object on disk.
func TestUnicodeNormalizationAliasesRejected(t *testing.T) {
	const composed = "Å"    // LATIN CAPITAL LETTER A WITH RING ABOVE
	const decomposed = "Å" // A + COMBINING RING ABOVE
	// Decomposed starts with plain "A" and so sorts first; hand them over in
	// byte order, or the reader rejects for ordering before reaching the fold.
	checkTree(t, []rawEntry{
		file(decomposed, 0o644, "one"),
		file(composed, 0o644, "two"),
	}, bundle.ReasonDuplicatePath)
}

// A frame header's remainder runs to 13 bytes (window descriptor, 4-byte
// dictionary ID, 8-byte content size), more than the guard's field buffer
// holds, and a short Read splits it across walk calls.
func TestZstdFrameHeaderSplitAcrossReads(t *testing.T) {
	// Magic, then a descriptor with fcs_flag=3, single_segment=0, did_flag=3.
	data := append([]byte{0x28, 0xB5, 0x2F, 0xFD, 0xC3}, bytes.Repeat([]byte{0xAA}, 32)...)
	for chunk := 1; chunk <= 9; chunk++ {
		t.Run(fmt.Sprint(chunk), func(t *testing.T) {
			if _, err := bundle.Validate(t.Context(), &chunkedReader{data: data, n: chunk}); !isReason(err, bundle.ReasonMalformed) {
				t.Errorf("Validate: got = %v, want reason %q", err, bundle.ReasonMalformed)
			}
		})
	}
}

// chunkedReader hands out at most n bytes per Read, like a socket would.
type chunkedReader struct {
	data []byte
	n    int
}

func (r *chunkedReader) Read(p []byte) (int, error) {
	if len(r.data) == 0 {
		return 0, io.EOF
	}
	n := min(r.n, len(p), len(r.data))
	copy(p, r.data[:n])
	r.data = r.data[n:]
	return n, nil
}
