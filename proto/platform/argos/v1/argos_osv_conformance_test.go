/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"strings"
	"testing"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
				Type: argosv1.Range_ECOSYSTEM,
				Events: []*argosv1.Event{
					{Event: &argosv1.Event_Introduced{Introduced: "0"}},
					{Event: &argosv1.Event_Fixed{Fixed: "1.2.3"}},
					{Event: &argosv1.Event_LastAffected{LastAffected: "1.2.2"}},
				},
			}},
			DatabaseSpecific: &argosv1.DatabaseSpecific{
				CweIds: []string{"CWE-0000"},
				SinkLocator: &argosv1.SinkLocator{
					Class:    "example.pkg.Sink",
					Method:   "method(str)",
					FileLine: "example/pkg/sink.py:1-2",
				},
				DefectKind: "incorrect-control",
			},
		}},
	}
	b, err := protojson.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	for _, want := range []string{
		`"affected"`, `"ranges"`, `"introduced"`, `"fixed"`,
		// Multi-word OSV keys are the ones protojson would lowerCamelCase
		// without their json_name pins — assert the exact spec form.
		`"schema_version"`, `"last_affected"`, `"database_specific"`,
		`"cwe_ids"`, `"sink_locator"`, `"class"`, `"file_line"`, `"defect_kind"`,
		// Range type must render the bare OSV enum form ("ECOSYSTEM"), not a
		// proto-prefixed value name.
		`"ECOSYSTEM"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("expected %s in %s", want, got)
		}
	}
	for _, reject := range []string{
		"schemaVersion", "lastAffected", "databaseSpecific",
		"cweIds", "sinkLocator", "fileLine", "defectKind",
		"RANGE_TYPE_",
	} {
		if strings.Contains(got, reject) {
			t.Errorf("non-OSV-spec token %q leaked into customer OSV JSON: %s", reject, got)
		}
	}
}

func TestOSVPaginationJSONConformsToOSVSpec(t *testing.T) {
	for _, tc := range []struct {
		name   string
		msg    proto.Message
		want   string
		reject string
	}{
		{"OSVQueryRequest", &argosv1.OSVQueryRequest{PageToken: "tok"}, `"page_token"`, "pageToken"},
		{"OSVQueryResponse", &argosv1.OSVQueryResponse{NextPageToken: "tok"}, `"next_page_token"`, "nextPageToken"},
		{"OSVQueryBatchResult", &argosv1.OSVQueryBatchResult{NextPageToken: "tok"}, `"next_page_token"`, "nextPageToken"},
	} {
		b, err := protojson.Marshal(tc.msg)
		if err != nil {
			t.Fatalf("%s: marshal: %v", tc.name, err)
		}
		got := string(b)
		if !strings.Contains(got, tc.want) {
			t.Errorf("%s: expected %s in %s", tc.name, tc.want, got)
		}
		if strings.Contains(got, tc.reject) {
			t.Errorf("%s: lowerCamelCase key %q leaked into customer OSV JSON (missing json_name pin): %s", tc.name, tc.reject, got)
		}
	}
}
