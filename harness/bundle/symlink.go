/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle

import (
	"path"
	"strings"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const maxSymlinkHops = 40

type entryTree struct {
	seen map[string]struct{}
	// Folded path to what the tree needs that path to be.
	claims   map[string]claim
	symlinks map[string]string
	// Original symlink paths, in arrival order. checkLinks walks these rather
	// than the map's folded keys, which would fold a second time.
	linkPaths []string
}

// claim is what one entry requires of one path: the kind it has to be -- TypeDir
// for a directory, whether an entry declared it or a descendant implied it,
// otherwise the entry's own type -- and the spelling the entry used for it.
// A host that folds gives a path one kind and one spelling, so two entries that
// disagree on either describe a tree it cannot materialize.
type claim struct {
	kind EntryType
	path string
}

// fold keys the index the way a normalization- and case-insensitive extractor
// resolves names: decompose to NFD, then map every rune to the minimum of its
// SimpleFold orbit. Decomposing is what makes this canonical in one pass --
// NFC recomposes into characters like "ᾴ" that have no simple case mapping, so
// a precomposed name and its decomposed twin would settle on different keys no
// matter how often they were folded.
// Windows and macOS open "a" through an entry named "A", so indexing by exact
// case would let "A -> ." and "a/B -> ../outside" look unrelated here and
// alias on disk. Folding can only merge paths, so it rejects more than it
// admits: a bundle carrying both "A" and "a" cannot be extracted faithfully on
// those hosts anyway.
// Case folding rather than lowercasing, because those hosts fold: "Σ", "σ" and
// "ς" are one name on disk, while ToLower keeps "ς" apart from the other two.
// Mapping is per rune and no decomposition crosses "/", so folding a whole
// path and folding each component give the same key -- the index and the chain
// walk cannot drift apart.
func fold(p string) string {
	return strings.Map(func(r rune) rune {
		// The smallest rune in the SimpleFold orbit is one key for every
		// spelling of the letter, and the orbit is a cycle through them all.
		lo := r
		for f := unicode.SimpleFold(r); f != r; f = unicode.SimpleFold(f) {
			lo = min(lo, f)
		}
		return lo
	}, norm.NFD.String(p))
}

// Every path an entry names, and every path it implies, must have one kind and
// one spelling. Keeping the index faithful to the paths an extractor would
// create is what stops entries aliasing each other through a link or a folded
// name, and what keeps tree_sha256 a pin on the tree that actually lands.
func (t *entryTree) add(e Entry) error {
	// Two paths that fold together address one object on a case-insensitive
	// host, so the exact-duplicate rule does not cover them: extraction either
	// fails or the later entry replaces the earlier, and the tree that lands
	// no longer matches tree_sha256. Indexing them would also drop one of two
	// colliding symlinks from the chain walk.
	if t.seen == nil {
		t.seen = make(map[string]struct{})
		t.claims = make(map[string]claim)
	}
	f := fold(e.Path)
	if _, ok := t.seen[f]; ok {
		return rejectf(ReasonDuplicatePath, e.Path, "path collides with an earlier entry once case is folded")
	}
	t.seen[f] = struct{}{}
	// One rule over the whole path rather than one per failure mode: every
	// ancestor has to be a directory, the entry has to be its own type, and all
	// of them have to agree on how the path is spelled. A folded key that two
	// entries claim differently is an alias whichever of them arrives first, so
	// "A/b" then "a" is the same collision as "A" then "a/b" -- only the sort
	// order differs, which the archive chooses.
	for p, want := e.Path, e.Type; p != "."; p, want = path.Dir(p), TypeDir {
		key := fold(p)
		prev, ok := t.claims[key]
		switch {
		case !ok:
			t.claims[key] = claim{kind: want, path: p}
		case (prev.kind == TypeDir) != (want == TypeDir):
			kind := prev.kind
			if want != TypeDir {
				kind = want
			}
			return rejectf(ReasonUnsafePath, e.Path, "%q is both a directory and a %s once case is folded", p, kind)
		case prev.path != p:
			// Only directories reach here: a second entry at one folded path is
			// already a duplicate, and a file under a file is already a kind
			// conflict. Both spellings would extract into one directory, so the
			// paths tree_sha256 covers are not the paths that would land.
			return rejectf(ReasonDuplicatePath, e.Path, "directory %q is spelled %q elsewhere, and a host that folds can keep only one", p, prev.path)
		}
	}
	if e.Type == TypeSymlink {
		if t.symlinks == nil {
			t.symlinks = make(map[string]string)
		}
		t.symlinks[f] = e.LinkTarget
		t.linkPaths = append(t.linkPaths, e.Path)
	}
	return nil
}

// All links must be indexed before checking chains: a target can name a link
// that sorts after it in the archive.
func (t *entryTree) checkLinks() error {
	for _, from := range t.linkPaths {
		if err := t.resolve(from); err != nil {
			return err
		}
	}
	return nil
}

// resolve carries the folded path as bytes rather than re-joining a component
// slice at every step. Each component triggers a lookup, so joining per step
// would cost O(components²) copying on a chain a bundle controls the width of.
// Indexing a map by string(key) does not allocate.
func (t *entryTree) resolve(from string) error {
	pending := strings.Split(from, "/")
	var key []byte
	var lens []int // key's length before each component, for popping ".."
	var hops int
	for len(pending) > 0 {
		part := pending[0]
		pending = pending[1:]
		switch part {
		case "", ".":
			continue
		case "..":
			if len(lens) == 0 {
				return rejectf(ReasonSymlinkEscape, from, "symlink chain resolves outside the tree")
			}
			key = key[:lens[len(lens)-1]]
			lens = lens[:len(lens)-1]
			continue
		}
		lens = append(lens, len(key))
		if len(key) > 0 {
			key = append(key, '/')
		}
		key = append(key, fold(part)...)
		target, ok := t.symlinks[string(key)]
		if !ok {
			continue
		}
		hops++
		if hops > maxSymlinkHops {
			return rejectf(ReasonSymlinkEscape, from, "symlink chain exceeds %d links", maxSymlinkHops)
		}
		// Resolve the target from the link's parent before examining any
		// remaining components. Cleaning first would apply ".." too early.
		key = key[:lens[len(lens)-1]]
		lens = lens[:len(lens)-1]
		pending = append(strings.Split(target, "/"), pending...)
	}
	return nil
}
