/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"encoding/json"
	"maps"
	"slices"
	"strings"
	"testing"
	"time"

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
		Withdrawn:     timestamppb.Now(),
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
				SinkLocator: []*argosv1.SinkLocator{{
					Class:            "example.pkg.Sink",
					Method:           "method(str)",
					FileLine:         "example/pkg/sink.py:1-2",
					ObservedVersions: []string{"1.2.2"},
				}},
				DefectKind: "incorrect-control",
			},
		}},
		References: []*argosv1.Reference{{
			Type: argosv1.Reference_FIX,
			Url:  "https://github.com/django/django/commit/cf694e6852b0da7799f8b53f1fb2f7d20cf17534",
		}},
		Upstream: []string{"CGP-DDDD-EEEE-FFFF"},
	}
	b, err := protojson.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got := string(b)
	for _, want := range []string{
		`"affected"`, `"ranges"`, `"introduced"`, `"fixed"`, `"withdrawn"`,
		// Multi-word OSV keys are the ones protojson would lowerCamelCase
		// without their json_name pins — assert the exact spec form.
		`"schema_version"`, `"last_affected"`, `"database_specific"`,
		`"cwe_ids"`, `"sink_locator"`, `"class"`, `"file_line"`, `"defect_kind"`,
		`"observed_versions"`,
		// Range type must render the bare OSV enum form ("ECOSYSTEM"), not a
		// proto-prefixed value name.
		`"ECOSYSTEM"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("got = %s, want it to contain %s", got, want)
		}
	}
	for _, reject := range []string{
		"schemaVersion", "lastAffected", "databaseSpecific",
		"cweIds", "sinkLocator", "fileLine", "defectKind",
		"observedVersions",
		"RANGE_TYPE_", "REFERENCE_TYPE_",
	} {
		if strings.Contains(got, reject) {
			t.Errorf("non-OSV-spec token %q leaked into customer OSV JSON: %s", reject, got)
		}
	}

	// sink_locator is repeated: the served JSON shape is an array of entry
	// objects, never a bare object. Assert structurally (protojson output
	// spacing is deliberately non-deterministic, so no substring check).
	var top struct {
		Affected []struct {
			DatabaseSpecific struct {
				SinkLocator []struct {
					ObservedVersions []string `json:"observed_versions"`
				} `json:"sink_locator"`
			} `json:"database_specific"`
		} `json:"affected"`
	}
	if err := json.Unmarshal(b, &top); err != nil {
		t.Fatalf("sink_locator is not an array of objects: %v in %s", err, got)
	}
	if n := len(top.Affected[0].DatabaseSpecific.SinkLocator); n != 1 {
		t.Errorf("sink_locator entries: got = %d, want = 1: %s", n, got)
	}

	// references must render the OSV spec's exact shape: an array of
	// {"type": "<bare spec name>", "url": "..."} objects, with the type as
	// the spec's literal string (osv-schema §references-field). An exact
	// equality on "FIX" (not a substring check) pins the served vocabulary
	// to the spec, so renaming the enum value in BOTH protos still fails
	// here. Structural, like sink_locator above, because protojson spacing
	// is non-deterministic.
	var refs struct {
		References []struct {
			Type string `json:"type"`
			URL  string `json:"url"`
		} `json:"references"`
	}
	if err := json.Unmarshal(b, &refs); err != nil {
		t.Fatalf("references is not an array of objects: %v in %s", err, got)
	}
	if n := len(refs.References); n != 1 {
		t.Fatalf("references entries: got = %d, want = 1: %s", n, got)
	}
	if refs.References[0].Type != "FIX" {
		t.Errorf("reference type: got = %q, want the bare OSV spec form %q: %s", refs.References[0].Type, "FIX", got)
	}
	if want := "https://github.com/django/django/commit/cf694e6852b0da7799f8b53f1fb2f7d20cf17534"; refs.References[0].URL != want {
		t.Errorf("reference url: got = %q, want = %q", refs.References[0].URL, want)
	}

	// upstream renders as OSV 1.7's array of record ids.
	var up struct {
		Upstream []string `json:"upstream"`
	}
	if err := json.Unmarshal(b, &up); err != nil {
		t.Fatalf("upstream is not an array of strings: %v in %s", err, got)
	}
	if want := []string{"CGP-DDDD-EEEE-FFFF"}; !slices.Equal(up.Upstream, want) {
		t.Errorf("upstream: got = %v, want = %v: %s", up.Upstream, want, got)
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

// TestAssessmentJSONUsesOpenVEXTokens pins database_specific.assessment's keys
// to the OpenVEX field names, which a VEX consumer copies, and superseded_by
// to its snake_case key.
func TestAssessmentJSONUsesOpenVEXTokens(t *testing.T) {
	assessedAt := time.Date(2026, 11, 20, 18, 0, 0, 0, time.UTC)
	rec := &argosv1.OSVRecord{
		Id:        "CGP-AAAA-BBBB-CCCC",
		Withdrawn: timestamppb.New(assessedAt),
		DatabaseSpecific: &argosv1.TopLevelDatabaseSpecific{
			Assessment: &argosv1.Assessment{
				Status:          "not_affected",
				Justification:   "vulnerable_code_cannot_be_controlled_by_adversary",
				ImpactStatement: "reachable only from the admin console",
				AssessedAt:      timestamppb.New(assessedAt),
				SupersededBy:    "CGP-DDDD-EEEE-FFFF",
			},
		},
	}
	b, err := protojson.Marshal(rec)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var got struct {
		DatabaseSpecific struct {
			Assessment map[string]string `json:"assessment"`
		} `json:"database_specific"`
	}
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal %s: %v", b, err)
	}
	want := map[string]string{
		"status":           "not_affected",
		"justification":    "vulnerable_code_cannot_be_controlled_by_adversary",
		"impact_statement": "reachable only from the admin console",
		"assessed_at":      "2026-11-20T18:00:00Z",
		"superseded_by":    "CGP-DDDD-EEEE-FFFF",
	}
	if !maps.Equal(got.DatabaseSpecific.Assessment, want) {
		t.Errorf("assessment: got = %v, want = %v", got.DatabaseSpecific.Assessment, want)
	}
}
