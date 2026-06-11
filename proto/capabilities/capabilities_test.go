/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestStringify(t *testing.T) {
	tests := []struct {
		name       string
		capability Capability
		want       string
		wantErr    error
	}{{
		name: "no requirements",
	}, {
		name:       "simple",
		capability: Capability_CAP_EVENTS_SUBSCRIPTION_DELETE,
		want:       "subscriptions.delete",
	}, {
		name:       "unknown",
		capability: 0,
		want:       "",
	}, {
		name:       "invalid",
		capability: 1,
		wantErr:    status.Error(codes.Internal, `capability has no descriptor: 1`),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, gotErr := Stringify(test.capability)

			switch {
			case (gotErr != nil) != (test.wantErr != nil):
				t.Fatalf("Stringify() = %v, %v, wanted %v, %v", got, gotErr, test.want, test.wantErr)
			case gotErr != nil && gotErr.Error() != test.wantErr.Error():
				t.Fatalf("Stringify() = %v, wanted %v", gotErr, test.wantErr)
			case gotErr == nil && got != test.want:
				t.Fatalf("Stringify() = %v, wanted %v", got, test.want)
			}
		})
	}
}

func TestStringifyAll(t *testing.T) {
	// staleCap stands in for a capability whose enum value has been removed
	// (deprecated then deleted). StringifyAll must skip it rather than failing
	// the whole list, so a stale capability on a stored role can't make that
	// role impossible to read. CUS-843.
	const staleCap = Capability(42)
	if _, err := Stringify(staleCap); err == nil {
		t.Fatalf("capability %d now has a descriptor; choose a different unassigned value for this test", staleCap)
	}

	tests := []struct {
		name string
		caps []Capability
		want []string
	}{{
		name: "empty input",
		caps: nil,
		want: []string{},
	}, {
		name: "all valid",
		caps: []Capability{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_EVENTS_SUBSCRIPTION_DELETE},
		want: []string{"groups.list", "subscriptions.delete"},
	}, {
		name: "all stale",
		caps: []Capability{staleCap},
		want: []string{},
	}, {
		name: "skips stale capability in the middle",
		caps: []Capability{
			Capability_CAP_IAM_GROUPS_LIST,
			staleCap,
			Capability_CAP_EVENTS_SUBSCRIPTION_DELETE,
		},
		want: []string{"groups.list", "subscriptions.delete"},
	}, {
		name: "skips stale capability at head",
		caps: []Capability{staleCap, Capability_CAP_IAM_GROUPS_LIST},
		want: []string{"groups.list"},
	}, {
		name: "skips stale capability at tail",
		caps: []Capability{Capability_CAP_IAM_GROUPS_LIST, staleCap},
		want: []string{"groups.list"},
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := StringifyAll(test.caps)
			if err != nil {
				t.Fatalf("StringifyAll() error: got = %v, want = nil", err)
			}
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("StringifyAll() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestStringifyAll_CUS843 replays a real stored capability set that triggered
// CUS-843: an owner-style role whose capabilities_json still includes capability
// 1601 (CAP_REGISTRY_PULL), removed from the enum in mono#23247, so it now has
// no descriptor. The datastore decodes this stored JSON straight into a
// []Capability (datastore/internal/persistence/iam/role.go), which then flows
// into StringifyAll on every role read. Before the fix the stale 1601 made
// StringifyAll fail the whole list, taking down the Console IDP settings page;
// it must now be skipped while every live capability is still returned.
func TestStringifyAll_CUS843(t *testing.T) {
	const storedCapabilitiesJSON = `[660,1601,505,623,703,1303,103,1503,1203,203,603,1605,403,1003,633,1609,903,1613,303,503,803,613,640,901,1615,670,650,1103,1703]`

	var caps []Capability
	if err := json.Unmarshal([]byte(storedCapabilitiesJSON), &caps); err != nil {
		t.Fatalf("decoding stored capabilities_json: %v", err)
	}

	// Premise: 1601 must be unresolvable for this regression to be meaningful.
	if _, err := Stringify(Capability(1601)); err == nil {
		t.Fatal("capability 1601 (CAP_REGISTRY_PULL) now has a descriptor; this regression test no longer reproduces CUS-843")
	}

	// Expected output: every capability that still resolves, in order — i.e.
	// the full stored set minus the deleted 1601.
	var want []string
	for _, c := range caps {
		if s, err := Stringify(c); err == nil {
			want = append(want, s)
		}
	}

	got, err := StringifyAll(caps)
	if err != nil {
		t.Fatalf("StringifyAll() on the stored caps: got error %v, want nil", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("StringifyAll() mismatch (-want +got):\n%s", diff)
	}
	if len(got) != len(caps)-1 {
		t.Errorf("StringifyAll() returned %d names, want %d (all %d stored caps minus the deleted 1601)", len(got), len(caps)-1, len(caps))
	}
}

func TestDeprecated(t *testing.T) {
	tests := []struct {
		name       string
		capability Capability
		want       bool
	}{{
		name:       "is deprecated",
		capability: Capability_CAP_TENANT_CLUSTERS_CREATE,
		want:       true,
	}, {
		name:       "not deprecated",
		capability: Capability_CAP_IAM_GROUPS_CREATE,
		want:       false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Deprecated(test.capability)

			if got != test.want {
				t.Errorf("Depcrecated() mismatch for %s: want=%t, got=%t", test.capability, test.want, got)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	for cap := range Capability_name {
		scap, _ := Stringify(Capability(cap))
		got, _ := Parse(scap)
		if Capability(cap) != got {
			t.Fatalf("Parse(Stringify()) = %v, wanted %v", got, Capability(cap))
		}
	}
}

func TestEncoding(t *testing.T) {
	all := make(Set, 0, len(Capability_name))
	for cap := range Capability_name {
		if cap == int32(Capability_UNKNOWN) {
			continue
		}
		all = append(all, Capability(cap))
	}
	all = SortCaps(all)

	tests := []struct {
		name string
		caps Set
	}{{
		name: "owner",
		caps: OwnerCaps,
	}, {
		name: "editor",
		caps: EditorCaps,
	}, {
		name: "viewer",
		caps: ViewerCaps,
	}, {
		name: "all",
		caps: all,
	}, {
		// SortCaps removes duplicates.
		name: "duplicates",
		caps: SortCaps(Set{Capability_CAP_IAM_GROUPS_LIST, Capability_CAP_IAM_GROUPS_LIST}),
	}}

	for _, test := range tests {
		t.Run(test.name+"-standard", func(t *testing.T) {
			raw, err := json.Marshal(test.caps)
			if err != nil {
				t.Fatalf("json.Marshal() = %v", err)
			}

			t.Logf("ENCODED: %s", raw)

			// Confirm that we decode it and get what we expect.
			got := make(Set, 0, len(test.caps))
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("json.Unmarshal() = %v", err)
			}
			if diff := cmp.Diff(got, test.caps); diff != "" {
				t.Errorf("(-got +want) = %s", diff)
			}
		})

		t.Run(test.name+"-legacy", func(t *testing.T) {
			// Remove our type alias, so that we use the legacy encoding.
			legacy := []Capability(test.caps)
			raw, err := json.Marshal(legacy)
			if err != nil {
				t.Fatalf("json.Marshal() = %v", err)
			}

			t.Logf("ENCODED: %s", raw)

			// Confirm that we decode it and get what we expect when coming from
			// the legacy encoding.
			got := make(Set, 0, len(test.caps))
			if err := json.Unmarshal(raw, &got); err != nil {
				t.Fatalf("json.Unmarshal() = %v", err)
			}
			if diff := cmp.Diff(got, test.caps); diff != "" {
				t.Errorf("(-got +want) = %s", diff)
			}
		})
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	caps := Set{
		Capability_CAP_IAM_GROUPS_LIST,

		Capability_CAP_REPO_LIST,
		Capability_CAP_MANIFEST_LIST,
		Capability_CAP_TAG_LIST,
		Capability_CAP_MANIFEST_METADATA_LIST,

		Capability_CAP_TENANT_RECORD_SIGNATURES_LIST,
		Capability_CAP_TENANT_SBOMS_LIST,
		Capability_CAP_TENANT_VULN_REPORTS_LIST,

		Capability_CAP_REPO_CREATE,
		Capability_CAP_REPO_UPDATE,
		Capability_CAP_REPO_DELETE,

		Capability_CAP_MANIFEST_CREATE,
		Capability_CAP_MANIFEST_UPDATE,
		Capability_CAP_MANIFEST_DELETE,

		Capability_CAP_TAG_CREATE,
		Capability_CAP_TAG_UPDATE,
		Capability_CAP_TAG_DELETE,

		// To create nested groups as needed on push.
		Capability_CAP_IAM_GROUPS_CREATE,
	}
	raw, err := json.Marshal(caps)
	if err != nil {
		b.Fatalf("json.Marshal() = %v", err)
	}

	for b.Loop() {
		var got Set
		if err := json.Unmarshal(raw, &got); err != nil {
			b.Fatalf("json.Unmarshal() = %v", err)
		}
	}
}
