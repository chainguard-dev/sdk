/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle

import (
	"cmp"
	"fmt"
	"hash"
	"io"
	"io/fs"
	"path"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	// Format is the value of a manifest's format field, and the name of this
	// bundle format.
	Format = "harness-source-v1"

	// ManifestPath is the reserved path of the manifest, the bundle's first
	// entry. A source tree that contains this path is rejected.
	ManifestPath = ".chainguard-harness/manifest.json"

	// treeHashDomain is the first field of the tree hash input, pinning the
	// record layout documented in the package doc. Changing the layout means
	// changing this tag so v1 hashes cannot collide with the new scheme's.
	treeHashDomain = "harness-source-v1-tree"

	// maxManifestBytes bounds the manifest entry before it is read, so a bundle
	// cannot hand the JSON decoder an arbitrarily large first entry. The
	// manifest's fields are a fixed, short list; a real one is a few hundred
	// bytes.
	maxManifestBytes = 64 << 10
)

// Kind is the source kind a bundle was built from. It records what chainctl
// found and how it packed it, not the flag the customer passed: "--workspace"
// packs a git checkout one way and a plain directory another. The bundle never
// records where the source lived.
type Kind string

// The source kinds a client packs today.
const (
	KindGit    Kind = "git"
	KindFolder Kind = "folder"
	KindFiles  Kind = "files"
)

// KindSource is reserved. "chainctl harness scan" does not yet take a
// repository as a source, so no writer emits it and validation rejects it; the
// name is spoken for so that taking the work up later adds a kind rather than
// moving the format onto a name an old bundle already used for something else.
// It is "source" because that is what the deferred flag is called, and what the
// sandbox's Session.environment already calls a repository and revision.
const KindSource Kind = "source"

// EntryType is one of the three permitted entry types. Its value is the token
// the tree hash covers, so these strings are part of the format.
type EntryType string

// The permitted entry types.
const (
	TypeFile    EntryType = "file"
	TypeDir     EntryType = "dir"
	TypeSymlink EntryType = "symlink"
)

// The normalized modes. A bundle carries no other mode: the writer normalizes
// to these and the reader rejects anything else, so a setuid bit or a
// group-writable file never survives a round trip.
const (
	modeFile fs.FileMode = 0o644
	modeExec fs.FileMode = 0o755
	modeDir  fs.FileMode = 0o755
	modeLink fs.FileMode = 0o777
)

// Entry is one normalized bundle entry. The writer takes entries from intake
// and normalizes Mode; the reader yields them with Mode already normalized.
type Entry struct {
	// Path is the entry's slash-separated path relative to the bundle root,
	// with no leading or trailing slash.
	Path string

	// Type is the entry's type.
	Type EntryType

	// Mode is the entry's permissions. On the way in only the execute bits of a
	// regular file are read, to pick 0755 over 0644; everything else is
	// normalized away. On the way out it is the normalized mode.
	Mode fs.FileMode

	// Size is a regular file's content length in bytes, and 0 for other types.
	Size int64

	// LinkTarget is a symlink's relative target, and empty for other types.
	LinkTarget string

	// Open opens a regular file's content, and is nil for other types. The
	// writer calls it twice — once to digest the content for tree_sha256 and
	// once to write it — and rejects the bundle if the two reads disagree, so a
	// source that changes underneath the pack is caught rather than packed
	// inconsistently.
	Open func() (io.ReadCloser, error)
}

// Manifest describes a bundle, and is its first entry at [ManifestPath].
//
// These fields are the complete list of metadata a bundle carries, and none of
// them resolves to a customer: no URL, remote, author or local path, and no
// customer-supplied free text. Commit and TreeSHA256 pin the tree without them.
type Manifest struct {
	// Format is always [Format]. The writer sets it.
	Format string `json:"format"`

	// Kind is the source kind the bundle was built from.
	Kind Kind `json:"kind"`

	// Commit is the 40-hex base commit the tree was packed against, on
	// [KindGit] only, and empty otherwise. A folder and a file list have no
	// commit to name.
	Commit string `json:"commit,omitempty"`

	// Dirty records that the checkout's working tree differed from Commit when
	// it was packed. It is false on every other kind.
	Dirty bool `json:"dirty"`

	// IncludeIgnored records that the checkout's ignore rules were disregarded
	// and files .gitignore excludes were packed anyway. Without it a bundle
	// holding .env and build output is indistinguishable from one that
	// respected those rules: both are [KindGit] at the same commit.
	//
	// It is false on every other kind, not because they exclude anything but
	// because they have no ignore rules to disregard; the kind is already what
	// says everything on disk was packed.
	IncludeIgnored bool `json:"include_ignored"`

	// Files is the number of regular-file entries. The writer sets it.
	Files int `json:"files"`

	// Bytes is the total content size of those files. The writer sets it.
	// Neither Files nor Bytes counts the manifest, directories or symlinks.
	Bytes int64 `json:"bytes"`

	// TreeSHA256 is the hex tree hash over the entries, defined exactly in the
	// package doc. The writer sets it; the reader recomputes it from the
	// entries it streamed and rejects a bundle that disagrees.
	TreeSHA256 string `json:"tree_sha256"`

	// Client identifies the packer, e.g. "chainctl v0.1.2".
	Client string `json:"client"`
}

// Limits bounds a bundle. One value applies on both sides: the client for an
// early, specific error, the server because the client's checks are not a
// control.
//
// Callers do not pick limits — [DefaultLimits] is the policy. A zero field
// takes its default, which is what lets a test drive one guard at a small size
// without restating the rest.
type Limits struct {
	// MaxCompressedBytes bounds the bundle as it arrives, before decompression.
	MaxCompressedBytes int64

	// MaxDecodedBytes bounds the whole decoded tar stream, tar headers and
	// skipped payloads included, not just retained file content.
	MaxDecodedBytes int64

	// MaxFileBytes bounds one regular file.
	MaxFileBytes int64

	// MaxEntries bounds the number of tar entries, the manifest included.
	MaxEntries int

	// MaxPathBytes and MaxPathDepth bound one path's length in bytes and its
	// number of slash-separated segments.
	MaxPathBytes int
	MaxPathDepth int

	// MaxRatio bounds decoded bytes divided by compressed bytes read.
	MaxRatio int64

	// RatioFloorBytes is the decoded size below which MaxRatio does not apply.
	// The ratio is measured over the stream so far, and a tar's leading bytes
	// are its most compressible: headers are mostly padding, so a tree whose
	// first paths are small or empty files inflates far above what the whole
	// bundle does. Without a floor those bundles would be refused for what they
	// start with rather than for what they are. The floor is also all the
	// decoded output MaxRatio gives a bomb for free, so it is sized in
	// megabytes; MaxDecodedBytes is what bounds a bomb that pays for its
	// inflation.
	RatioFloorBytes int64

	// MaxTrailingBytes bounds the zero padding a bundle may carry after the
	// tar's end-of-archive marker. Every one of those bytes must be zero, so
	// what is left of the channel is how many of them there are; the allowance
	// caps that at a handful of bits, and the bundle's own SHA-256 covers them.
	// It is sized to accept a GNU-style blocked writer: GNU tar pads its output
	// to a whole 20-block record by default, and this leaves room for the
	// 64-block records a writer may be configured with.
	MaxTrailingBytes int64
}

// DefaultLimits is the v1 policy, sized for large monorepos. It is the single
// place these values change; no call site sets its own.
func DefaultLimits() Limits {
	return Limits{
		MaxCompressedBytes: 512 << 20,
		MaxDecodedBytes:    2 << 30,
		MaxFileBytes:       256 << 20,
		MaxEntries:         250_000,
		MaxPathBytes:       4096,
		MaxPathDepth:       64,
		MaxRatio:           100,
		RatioFloorBytes:    8 << 20,
		MaxTrailingBytes:   32 << 10,
	}
}

// withDefaults fills every zero field from [DefaultLimits].
func (l Limits) withDefaults() Limits {
	d := DefaultLimits()
	return Limits{
		MaxCompressedBytes: cmp.Or(l.MaxCompressedBytes, d.MaxCompressedBytes),
		MaxDecodedBytes:    cmp.Or(l.MaxDecodedBytes, d.MaxDecodedBytes),
		MaxFileBytes:       cmp.Or(l.MaxFileBytes, d.MaxFileBytes),
		MaxEntries:         cmp.Or(l.MaxEntries, d.MaxEntries),
		MaxPathBytes:       cmp.Or(l.MaxPathBytes, d.MaxPathBytes),
		MaxPathDepth:       cmp.Or(l.MaxPathDepth, d.MaxPathDepth),
		MaxRatio:           cmp.Or(l.MaxRatio, d.MaxRatio),
		RatioFloorBytes:    cmp.Or(l.RatioFloorBytes, d.RatioFloorBytes),
		MaxTrailingBytes:   cmp.Or(l.MaxTrailingBytes, d.MaxTrailingBytes),
	}
}

// Reason is a stable code for why a bundle was rejected. It is what a caller
// branches on and surfaces; the message is for a human.
type Reason string

// The rejection reasons.
const (
	ReasonUnsafePath       Reason = "unsafe-path"
	ReasonPathTooLong      Reason = "path-too-long"
	ReasonTooDeep          Reason = "too-deep"
	ReasonReservedPath     Reason = "reserved-path"
	ReasonDuplicatePath    Reason = "duplicate-path"
	ReasonUnsortedEntry    Reason = "unsorted-entry"
	ReasonEntryType        Reason = "unsupported-entry-type"
	ReasonSymlinkEscape    Reason = "symlink-escape"
	ReasonMetadata         Reason = "unnormalized-metadata"
	ReasonFileTooLarge     Reason = "file-too-large"
	ReasonBundleTooLarge   Reason = "bundle-too-large"
	ReasonStreamTooLarge   Reason = "stream-too-large"
	ReasonTooManyEntries   Reason = "too-many-entries"
	ReasonDecompressRatio  Reason = "decompress-ratio"
	ReasonTrailingBytes    Reason = "trailing-bytes"
	ReasonInvalidManifest  Reason = "invalid-manifest"
	ReasonTreeHashMismatch Reason = "tree-hash-mismatch"
	ReasonMalformed        Reason = "malformed-bundle"
	ReasonSourceChanged    Reason = "source-changed"
)

// Error is a bundle rejection.
type Error struct {
	// Reason is the stable code for the rejection.
	Reason Reason
	// Path is the offending entry's path, when the rejection has one.
	Path string

	err error
}

var _ error = (*Error)(nil)

// Error renders the rejection. The path is quoted: it comes from an untrusted
// archive and reaches a terminal and a log.
func (e *Error) Error() string {
	if e.Path != "" {
		return fmt.Sprintf("%s: %q: %s", e.Reason, e.Path, e.err)
	}
	return fmt.Sprintf("%s: %s", e.Reason, e.err)
}

// Unwrap exposes the cause so a malformed archive's tar or JSON error survives.
func (e *Error) Unwrap() error { return e.err }

// rejectf builds a rejection. Pass an empty path when the rejection is about
// the bundle rather than one entry; wrap a cause with %w as with fmt.Errorf.
func rejectf(reason Reason, p, format string, args ...any) *Error {
	return &Error{Reason: reason, Path: p, err: fmt.Errorf(format, args...)}
}

// checkPath enforces the format's path rules. It rejects rather than sanitizes:
// a path that needed cleaning is a path whose writer disagreed with this one.
func checkPath(p string, lim Limits) error {
	switch {
	case p == "":
		return rejectf(ReasonUnsafePath, p, "empty path")
	case len(p) > lim.MaxPathBytes:
		return rejectf(ReasonPathTooLong, p, "path is %d bytes, limit is %d", len(p), lim.MaxPathBytes)
	case strings.ContainsRune(p, 0):
		return rejectf(ReasonUnsafePath, p, "path contains NUL")
	case !utf8.ValidString(p):
		return rejectf(ReasonUnsafePath, p, "path is not valid UTF-8")
	case strings.HasPrefix(p, "/"):
		return rejectf(ReasonUnsafePath, p, "path is absolute")
	case hasWindowsDrive(p):
		return rejectf(ReasonUnsafePath, p, "path has a Windows drive prefix")
	case strings.ContainsRune(p, '\\'):
		// A backslash is an ordinary byte to the checks below, and a separator
		// to a Windows extractor: "a\..\..\b" passes every rule here and still
		// escapes. The format refuses the byte rather than asking every
		// consumer to be the one that normalizes it.
		return rejectf(ReasonUnsafePath, p, "path contains a backslash")
	}
	segments := strings.Split(p, "/")
	if len(segments) > lim.MaxPathDepth {
		return rejectf(ReasonTooDeep, p, "path is %d segments deep, limit is %d", len(segments), lim.MaxPathDepth)
	}
	for _, s := range segments {
		switch s {
		case "":
			return rejectf(ReasonUnsafePath, p, "path has an empty segment")
		case ".", "..":
			return rejectf(ReasonUnsafePath, p, "path has a %q segment", s)
		}
		if strings.HasSuffix(s, " ") || strings.HasSuffix(s, ".") {
			return rejectf(ReasonUnsafePath, p, "path segment %q ends in a space or dot", s)
		}
		if strings.ContainsRune(s, ':') {
			return rejectf(ReasonUnsafePath, p, "path segment %q contains a colon", s)
		}
		if isWindowsDevice(s) {
			return rejectf(ReasonUnsafePath, p, "path segment %q is a reserved Windows device name", s)
		}
	}
	return nil
}

// Windows resolves a reserved device name in any directory and regardless of
// extension, so "NUL.txt" opens the null device rather than creating a file,
// and a colon opens an NTFS alternate data stream: "notes.txt:hidden" writes
// into a stream of notes.txt rather than to its own entry. Either makes a
// validated entry materialize somewhere other than where it says, which the
// uniqueness and containment rules cannot see. Refused for the same reason the
// backslash is -- the format decides, rather than each consumer.
// The superscript digits and the console streams are device names too, as
// Go's own internal/filepathlite records.
var winDevices = map[string]struct{}{
	"con": {}, "prn": {}, "aux": {}, "nul": {},
	"conin$": {}, "conout$": {},
	"com1": {}, "com2": {}, "com3": {}, "com4": {}, "com5": {},
	"com6": {}, "com7": {}, "com8": {}, "com9": {},
	"com¹": {}, "com²": {}, "com³": {},
	"lpt1": {}, "lpt2": {}, "lpt3": {}, "lpt4": {}, "lpt5": {},
	"lpt6": {}, "lpt7": {}, "lpt8": {}, "lpt9": {},
	"lpt¹": {}, "lpt²": {}, "lpt³": {},
}

// isWindowsDevice reports whether a path segment names a reserved device,
// ignoring any extension and case.
func isWindowsDevice(s string) bool {
	if i := strings.IndexByte(s, '.'); i >= 0 {
		s = s[:i]
	}
	_, ok := winDevices[strings.ToLower(s)]
	return ok
}

// The wire format must recognize drive prefixes on every host, including
// drive-relative forms such as C:outside that still select a Windows drive.
func hasWindowsDrive(p string) bool {
	return len(p) >= 2 && p[1] == ':' && (p[0] >= 'A' && p[0] <= 'Z' || p[0] >= 'a' && p[0] <= 'z')
}

// checkLinkTarget checks one target's syntax and lexical containment.
// entryTree.checkLinks also resolves intermediate links before applying "..".
func checkLinkTarget(from, target string, lim Limits) error {
	switch {
	case target == "":
		return rejectf(ReasonSymlinkEscape, from, "empty link target")
	case len(target) > lim.MaxPathBytes:
		return rejectf(ReasonPathTooLong, from, "link target is %d bytes, limit is %d", len(target), lim.MaxPathBytes)
	case strings.ContainsRune(target, 0):
		return rejectf(ReasonSymlinkEscape, from, "link target contains NUL")
	case !utf8.ValidString(target):
		return rejectf(ReasonSymlinkEscape, from, "link target is not valid UTF-8")
	case path.IsAbs(target):
		return rejectf(ReasonSymlinkEscape, from, "link target %q is absolute", target)
	case hasWindowsDrive(target):
		return rejectf(ReasonSymlinkEscape, from, "link target %q has a Windows drive prefix", target)
	case strings.ContainsRune(target, '\\'):
		// path.Join below treats a backslash as an ordinary byte; a Windows
		// extractor treats it as a separator, and the two disagree about
		// whether the target leaves the tree.
		return rejectf(ReasonSymlinkEscape, from, "link target %q contains a backslash", target)
	}
	// A target is a path, so it gets the depth ceiling a path gets. Bounding
	// only its bytes leaves ~2000 components in a 4 KiB target, and
	// entryTree.resolve substitutes a whole target for a single component, so
	// a chain multiplies that by its hop limit. None of that expansion is
	// visible to the decoded-size or ratio guards, which see only the bytes.
	segments := strings.Split(target, "/")
	if len(segments) > lim.MaxPathDepth {
		return rejectf(ReasonTooDeep, from, "link target is %d segments deep, limit is %d", len(segments), lim.MaxPathDepth)
	}
	for _, s := range segments {
		if s == "." || s == ".." {
			continue
		}
		if strings.HasSuffix(s, " ") || strings.HasSuffix(s, ".") {
			return rejectf(ReasonSymlinkEscape, from, "link target segment %q ends in a space or dot", s)
		}
		if strings.ContainsRune(s, ':') {
			return rejectf(ReasonSymlinkEscape, from, "link target segment %q contains a colon", s)
		}
		if isWindowsDevice(s) {
			return rejectf(ReasonSymlinkEscape, from, "link target segment %q is a reserved Windows device name", s)
		}
	}
	// path.Join cleans, so a target that climbs out of the tree lands on ".."
	// or a path under it.
	resolved := path.Join(path.Dir(from), target)
	if resolved == ".." || strings.HasPrefix(resolved, "../") {
		return rejectf(ReasonSymlinkEscape, from, "link target %q resolves to %q, outside the tree", target, resolved)
	}
	return nil
}

// Normalize checks one entry against the format's rules and returns it with
// its Mode normalized, or an [Error] naming the [Reason] it cannot be carried.
//
// It is the seam between intake and the writer. Intake calls it on each
// candidate and drops what it rejects, with a warning naming the reason; the
// writer applies the same rule to everything it is handed and refuses it there,
// because by then an unrepresentable entry is a bug rather than a tree the
// customer happened to have. Both share these per-entry rules.
//
// Whole-tree rules — ordering, duplicates, entry count, non-directory parents
// and symlink chains — are checked by the writer and reader. Normalize alone
// does not establish symlink containment.
func Normalize(e Entry) (Entry, error) { return normalizeEntry(e, DefaultLimits()) }

// normalizeEntry is [Normalize] under explicit limits, which only a test sets.
func normalizeEntry(e Entry, lim Limits) (Entry, error) {
	if err := checkPath(e.Path, lim); err != nil {
		return Entry{}, err
	}
	if e.Path == ManifestPath {
		return Entry{}, rejectf(ReasonReservedPath, e.Path, "source tree contains the reserved manifest path")
	}
	switch e.Type {
	case TypeFile:
		if e.Open == nil {
			return Entry{}, rejectf(ReasonMalformed, e.Path, "file entry has no content")
		}
		if e.Size < 0 || e.Size > lim.MaxFileBytes {
			return Entry{}, rejectf(ReasonFileTooLarge, e.Path, "file is %d bytes, limit is %d", e.Size, lim.MaxFileBytes)
		}
		if e.LinkTarget != "" {
			return Entry{}, rejectf(ReasonMalformed, e.Path, "file entry has a link target")
		}
		// The only thing the source mode decides is whether the file is
		// executable; everything else is normalized away.
		if e.Mode.Perm()&0o111 != 0 {
			e.Mode = modeExec
		} else {
			e.Mode = modeFile
		}
	case TypeDir:
		if e.Open != nil || e.Size != 0 || e.LinkTarget != "" {
			return Entry{}, rejectf(ReasonMalformed, e.Path, "directory entry has content or a link target")
		}
		e.Mode = modeDir
	case TypeSymlink:
		if e.Open != nil || e.Size != 0 {
			return Entry{}, rejectf(ReasonMalformed, e.Path, "symlink entry has content")
		}
		if err := checkLinkTarget(e.Path, e.LinkTarget, lim); err != nil {
			return Entry{}, err
		}
		e.Mode = modeLink
	default:
		return Entry{}, rejectf(ReasonEntryType, e.Path, "unsupported entry type %q", e.Type)
	}
	return e, nil
}

// treeHasher accumulates tree_sha256 over entries in tar order. The writer
// feeds it the tree it is about to pack; the reader feeds it the tree it
// streamed, and compares.
type treeHasher struct {
	h hash.Hash
}

func newTreeHasher(h hash.Hash) *treeHasher {
	writeField(h, treeHashDomain)
	return &treeHasher{h: h}
}

// add appends one entry's record. payload is the file's hex content digest, the
// symlink's target, or empty for a directory.
func (t *treeHasher) add(e Entry, payload string) {
	writeField(t.h, e.Path)
	writeField(t.h, fmt.Sprintf("%04o", e.Mode.Perm()))
	writeField(t.h, string(e.Type))
	writeField(t.h, payload)
}

// sum returns the accumulated tree_sha256.
func (t *treeHasher) sum() string {
	return fmt.Sprintf("%x", t.h.Sum(nil))
}

// writeField appends a length-prefixed field, the framing
// chainguard.dev/sdk/skills.HardenJobID uses: len fixes the boundary, and the
// separators keep the input human-inspectable for debugging.
func writeField(w io.Writer, s string) {
	// hash.Hash never reports a write error.
	_, _ = io.WriteString(w, strconv.Itoa(len(s)))
	_, _ = io.WriteString(w, ":")
	_, _ = io.WriteString(w, s)
	_, _ = io.WriteString(w, "\n")
}

// validate checks the manifest's own fields. The counts and the tree hash it
// declares are checked against the entries by the reader.
func (m *Manifest) validate() error {
	if err := m.checkSource(); err != nil {
		return err
	}
	if m.Files < 0 || m.Bytes < 0 {
		return rejectf(ReasonInvalidManifest, "", "negative counts: files=%d bytes=%d", m.Files, m.Bytes)
	}
	if !isHex(m.TreeSHA256, 64) {
		return rejectf(ReasonInvalidManifest, "", "tree_sha256 %q is not 64 hex digits", m.TreeSHA256)
	}
	return nil
}

// checkSource checks the fields that describe where a bundle came from, the
// ones a packer supplies rather than computes. The writer runs it before it
// digests a tree, so a manifest it could never emit costs nothing.
func (m *Manifest) checkSource() error {
	if m.Format != Format {
		return rejectf(ReasonInvalidManifest, "", "format is %q, want %q", m.Format, Format)
	}
	switch m.Kind {
	case KindGit, KindFolder, KindFiles:
	case KindSource:
		// Naming the kind as deferred rather than unknown is what a reader of a
		// bundle from a later client needs: the name is part of this format, it
		// just has no packer yet.
		return rejectf(ReasonInvalidManifest, "", "kind %q is reserved: repository sources are deferred and no client packs them yet", m.Kind)
	default:
		return rejectf(ReasonInvalidManifest, "", "unknown kind %q", m.Kind)
	}
	// Only a git checkout has a commit to be packed against, only it can differ
	// from one, and only it has ignore rules to disregard; a folder and a file
	// list have none of the three.
	if m.Kind == KindGit {
		if !isHex(m.Commit, 40) {
			return rejectf(ReasonInvalidManifest, "", "kind %q needs a 40-hex base commit, got %q", m.Kind, m.Commit)
		}
	} else {
		if m.Commit != "" {
			return rejectf(ReasonInvalidManifest, "", "kind %q carries commit %q", m.Kind, m.Commit)
		}
		if m.Dirty {
			return rejectf(ReasonInvalidManifest, "", "kind %q is dirty", m.Kind)
		}
		if m.IncludeIgnored {
			return rejectf(ReasonInvalidManifest, "", "kind %q has no ignore rules to include", m.Kind)
		}
	}
	return checkClient(m.Client)
}

// checkClient bounds the one field a packer fills in freely. It reaches the
// guest and the agent's context, so it is an ASCII name and version, not free
// text.
func checkClient(client string) error {
	if client == "" || len(client) > 64 {
		return rejectf(ReasonInvalidManifest, "", "client %q must be 1..64 bytes", client)
	}
	for _, r := range client {
		if r < ' ' || r > '~' {
			return rejectf(ReasonInvalidManifest, "", "client %q has a non-printable or non-ASCII character", client)
		}
	}
	return nil
}

// isHex reports whether s is exactly n lowercase hex digits.
func isHex(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for _, r := range s {
		if (r < '0' || r > '9') && (r < 'a' || r > 'f') {
			return false
		}
	}
	return true
}
