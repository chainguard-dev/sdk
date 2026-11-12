/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"encoding/json"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
)

// RESULT_BLOCKED is an alias declared after RESULT_DENIED, so every encoding
// that writes a name keeps writing RESULT_DENIED while both names decode.
func TestResultBlockedAliasKeepsDeniedOnTheWire(t *testing.T) {
	if got, want := Result_RESULT_BLOCKED, Result(2); got != want {
		t.Errorf("Result_RESULT_BLOCKED: got = %d, want = %d", got, want)
	}
	if got, want := Result_RESULT_BLOCKED.String(), "RESULT_DENIED"; got != want {
		t.Errorf("String(): got = %q, want = %q", got, want)
	}
	if got, want := Result_name[2], "RESULT_DENIED"; got != want {
		t.Errorf("Result_name[2]: got = %q, want = %q", got, want)
	}
	for _, name := range []string{"RESULT_DENIED", "RESULT_BLOCKED"} {
		if got, want := Result_value[name], int32(2); got != want {
			t.Errorf("Result_value[%q]: got = %d, want = %d", name, got, want)
		}
	}

	b, err := protojson.Marshal(&Decision{Result: Result_RESULT_BLOCKED})
	if err != nil {
		t.Fatalf("protojson.Marshal: %v", err)
	}
	// protojson randomizes whitespace, so compare the decoded field.
	var fields map[string]any
	if err := json.Unmarshal(b, &fields); err != nil {
		t.Fatalf("json.Unmarshal(%s): %v", b, err)
	}
	if got, want := fields["result"], "RESULT_DENIED"; got != want {
		t.Errorf("JSON result: got = %v, want = %q", got, want)
	}

	var d Decision
	if err := protojson.Unmarshal([]byte(`{"result":"RESULT_BLOCKED"}`), &d); err != nil {
		t.Fatalf("protojson.Unmarshal: %v", err)
	}
	if got, want := d.GetResult(), Result_RESULT_BLOCKED; got != want {
		t.Errorf("decoded result: got = %v, want = %v", got, want)
	}
}
