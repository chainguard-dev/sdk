/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MaxRejectionReasonRunes bounds a rejection reason, counted after trimming.
const MaxRejectionReasonRunes = 1000

// The errors ValidateRejectionReason returns. They never carry the reason,
// so a caller can log or return them as they are.
var (
	ErrRejectionReasonBlank     = errors.New("rejection reason is required and cannot be blank")
	ErrRejectionReasonTooLong   = errors.New("rejection reason exceeds 1000 characters")
	ErrRejectionReasonMalformed = errors.New("rejection reason must be valid UTF-8 with no newline, tab or invisible character")
)

// ValidateRejectionReason returns reason with surrounding whitespace trimmed,
// or one of the ErrRejectionReason* errors.
//
// A mitigation partner reads the stored reason verbatim, so a newline, tab or
// invisible rune left after trimming is refused rather than stripped:
// stripping it would store a reason that differs from the one the reviewer
// confirmed. Callers send and store the returned value, never their input.
func ValidateRejectionReason(reason string) (string, error) {
	if !utf8.ValidString(reason) {
		return "", ErrRejectionReasonMalformed
	}
	t := strings.TrimSpace(reason)
	if t == "" {
		return "", ErrRejectionReasonBlank
	}
	if utf8.RuneCountInString(t) > MaxRejectionReasonRunes {
		return "", ErrRejectionReasonTooLong
	}
	if strings.ContainsFunc(t, invisible) {
		return "", ErrRejectionReasonMalformed
	}
	return t, nil
}

// invisible ports Invisible from chainguard.dev/cg/internal/termsafe, which
// this module cannot import. Keep the two in step: cg renders a reason through
// termsafe, and a rune accepted here but dropped there would show the reviewer
// a different reason from the one the partner reads.
//
// Newline, tab and carriage return are control runes, so this refuses them.
func invisible(r rune) bool {
	switch {
	case unicode.IsControl(r):
		return true
	case r <= unicode.MaxASCII:
		return false
	case unicode.IsSpace(r):
		return true
	case unicode.Is(unicode.Cf, r):
		return true
	case unicode.Is(unicode.Other_Default_Ignorable_Code_Point, r),
		unicode.Is(unicode.Variation_Selector, r):
		return true
	}
	// BRAILLE PATTERN BLANK, in none of the tables above.
	return r == '\u2800'
}
