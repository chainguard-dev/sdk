/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"errors"
	"strings"
	"testing"

	argosv1 "chainguard.dev/sdk/proto/platform/argos/v1"
)

func TestValidateRejectionReason(t *testing.T) {
	atLimit := strings.Repeat("a", argosv1.MaxRejectionReasonRunes)
	tests := []struct {
		name    string
		reason  string
		want    string
		wantErr error
	}{
		{name: "plain sentence", reason: "Taint tracking misses deserialization sinks.", want: "Taint tracking misses deserialization sinks."},
		{name: "surrounding whitespace trimmed", reason: " \tno sink coverage\n ", want: "no sink coverage"},
		{name: "surrounding non-ASCII whitespace trimmed", reason: "\u00a0\u3000no sink coverage\u0085\u2028", want: "no sink coverage"},
		{name: "inner spaces kept", reason: "two  spaces", want: "two  spaces"},
		{name: "decomposed combining mark is visible", reason: "cafe\u0301", want: "cafe\u0301"},
		{name: "CJK text", reason: "反序列化", want: "反序列化"},
		{name: "limit counts runes not bytes", reason: strings.Repeat("\u00e9", argosv1.MaxRejectionReasonRunes), want: strings.Repeat("\u00e9", argosv1.MaxRejectionReasonRunes)},
		{name: "limit counted after trimming", reason: "  " + atLimit + "\n", want: atLimit},
		{name: "one rune over the limit", reason: atLimit + "a", wantErr: argosv1.ErrRejectionReasonTooLong},
		{name: "empty", reason: "", wantErr: argosv1.ErrRejectionReasonBlank},
		{name: "whitespace only", reason: " \t\r\n ", wantErr: argosv1.ErrRejectionReasonBlank},
		{name: "non-ASCII whitespace only", reason: "\u00a0\u3000", wantErr: argosv1.ErrRejectionReasonBlank},
		{name: "invalid UTF-8", reason: "bad \xff byte", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "UTF-8-encoded surrogate half", reason: "a\xed\xa0\x80b", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "inner newline", reason: "first\nsecond", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "inner tab", reason: "first\tsecond", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "inner carriage return", reason: "first\rsecond", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "ANSI escape", reason: "\x1b[31mred", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "DEL", reason: "a\x7fb", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "C1 control", reason: "a\u0085b", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "inner no-break space", reason: "a\u00a0b", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "inner line separator", reason: "a\u2028b", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "leading zero-width space is not trimmed", reason: "\u200bab", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "bidi override", reason: "a\u202eb", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "byte order mark", reason: "a\ufeffb", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "variation selector", reason: "a\ufe0fb", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "combining grapheme joiner", reason: "a\u034fb", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "hangul filler", reason: "a\u3164b", wantErr: argosv1.ErrRejectionReasonMalformed},
		{name: "braille pattern blank", reason: "a\u2800b", wantErr: argosv1.ErrRejectionReasonMalformed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := argosv1.ValidateRejectionReason(tt.reason)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error: got = %v, want = %v", err, tt.wantErr)
			}
			// Callers return err.Error() as it is, so it must never carry the reason.
			if err != nil && err.Error() != tt.wantErr.Error() {
				t.Errorf("error text: got = %q, want = %q", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("reason: got = %q, want = %q", got, tt.want)
			}
			// cg, api-impl and the datastore each validate what the layer before
			// returned, so an accepted reason must come back unchanged.
			if err == nil {
				if again, err := argosv1.ValidateRejectionReason(got); err != nil || again != got {
					t.Errorf("revalidated: got = %q, %v, want = %q, <nil>", again, err, got)
				}
			}
		})
	}
}
