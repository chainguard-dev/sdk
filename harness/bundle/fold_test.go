/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package bundle

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/text/cases"
	"golang.org/x/text/unicode/norm"
)

// Greek Extended is where precomposed letters, iota subscript and case
// interact, so every rune in it must key the same as its decomposed, fully
// case-folded spelling. Full folding is only the oracle here: it maps "ß" to
// "ss", which the filesystems this index models keep apart, so the index
// itself folds per rune.
func TestGreekExtendedKeysMatchFullCaseFold(t *testing.T) {
	folder := cases.Fold()
	for r := rune(0x1F00); r <= 0x1FFF; r++ {
		full := folder.String(norm.NFD.String(string(r)))
		if got, want := fold(string(r)), fold(full); got != want {
			t.Errorf("U+%04X: fold(%s) = %s, fold(full case fold %s) = %s",
				r, codepoints(string(r)), codepoints(got), codepoints(full), codepoints(want))
		}
	}
}

// add keys whole paths and resolve keys one component at a time. "/" survives
// NFD as a starter that no decomposition crosses, so the two agree; pin that
// rather than trusting it, since the two drifting apart is a silent escape.
func TestFoldDistributesOverPathSeparator(t *testing.T) {
	for _, p := range []string{
		"a/ᾴ/outside", "ΆΙ/a", "a/́b", "Å/Å", "a/ΑΙ/ᾼ/Σ/ς",
		"ᾼ/ᾳ/ΑΙ", "ͅ/Ι/ι/ι", "a/b/c", "", "/", "a//b",
	} {
		parts := strings.Split(p, "/")
		for i, part := range parts {
			parts[i] = fold(part)
		}
		if got, want := fold(p), strings.Join(parts, "/"); got != want {
			t.Errorf("fold(%q) = %s, per component = %s", p, codepoints(got), codepoints(want))
		}
	}
}

func codepoints(s string) string {
	var b strings.Builder
	for _, r := range s {
		fmt.Fprintf(&b, "U+%04X ", r)
	}
	return strings.TrimSpace(b.String())
}
