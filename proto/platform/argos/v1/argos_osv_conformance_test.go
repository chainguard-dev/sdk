/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"

	argosv1 "chainguard.dev/sdk/proto/platform/argos/v1"
)

func TestOSVRecordJSONConformsToOSVSpec(t *testing.T) {
	rec := &argosv1.OSVRecord{
		Id:            "CGP-AAAA-BBBB-CCCC",
		SchemaVersion: "1.6.0",
		Summary:       "example",
		Modified:      timestamppb.Now(),
		Affected: []*argosv1.Affected{{
			Package: &argosv1.Package{Ecosystem: "PyPI", Name: "example"},
			Ranges: []*argosv1.Range{{
				Type: argosv1.Range_RANGE_TYPE_ECOSYSTEM,
				Events: []*argosv1.Event{
					{Event: &argosv1.Event_Introduced{Introduced: "0"}},
					{Event: &argosv1.Event_Fixed{Fixed: "1.2.3"}},
				},
			}},
		}},
	}
	b, err := protojson.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	for _, want := range []string{`"affected"`, `"ranges"`, `"introduced"`, `"fixed"`} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %s in %s", want, got)
		}
	}
}
