/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/chainguard-dev/clog"
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

// TestRoundTrip is the merge-time completeness proof for the Parse name map:
// every enum value must carry a resolvable (name) option and round-trip
// through Parse.
func TestRoundTrip(t *testing.T) {
	for cap := range Capability_name {
		scap, err := Stringify(Capability(cap))
		if err != nil {
			t.Fatalf("Stringify(%d): got error %v, want nil", cap, err)
		}
		got, err := Parse(scap)
		if err != nil {
			t.Fatalf("Parse(%q): got error %v, want nil", scap, err)
		}
		if Capability(cap) != got {
			t.Fatalf("Parse(Stringify()) = %v, wanted %v", got, Capability(cap))
		}
	}
}

// TestParse_UnknownName pins the documented contract: unknown names are not an
// error; they resolve to (Capability_UNKNOWN, nil). Callers detect bad input by
// checking for Capability_UNKNOWN explicitly.
func TestParse_UnknownName(t *testing.T) {
	got, err := Parse("no.such.capability")
	if err != nil {
		t.Fatalf("Parse(unknown name) error = %v, want nil", err)
	}
	if got != Capability_UNKNOWN {
		t.Fatalf("Parse(unknown name) = %v, want Capability_UNKNOWN", got)
	}
}

func TestStringifyAllContext_WarnOncePerValue(t *testing.T) {
	// Dedicated stale values: the warn dedup is process-global, so this test
	// must not share stale values with other tests in the package.
	const staleA, staleB = Capability(424201), Capability(424202)
	for _, c := range []Capability{staleA, staleB} {
		if _, err := Stringify(c); err == nil {
			t.Fatalf("capability %d has a descriptor; pick another dedicated stale value", c)
		}
	}

	var buf bytes.Buffer
	ctx := clog.WithLogger(t.Context(), clog.New(slog.NewTextHandler(&buf, nil)))

	for range 3 {
		got, err := StringifyAllContext(ctx, []Capability{staleA, Capability_CAP_IAM_GROUPS_LIST, staleB})
		if err != nil {
			t.Fatalf("StringifyAllContext(): got error %v, want nil", err)
		}
		if diff := cmp.Diff([]string{"groups.list"}, got); diff != "" {
			t.Errorf("StringifyAllContext() mismatch (-want +got):\n%s", diff)
		}
	}

	logs := buf.String()
	if n := strings.Count(logs, "skipping capability 424201 "); n != 1 {
		t.Errorf("warned about %d %d times, want exactly 1; logs:\n%s", staleA, n, logs)
	}
	if n := strings.Count(logs, "skipping capability 424202 "); n != 1 {
		t.Errorf("warned about %d %d times, want exactly 1; logs:\n%s", staleB, n, logs)
	}
}

// TestWarnUnknownCapability_DedupLimit exercises the dedup map's boundary
// behavior: the one-time limit-reached log, continued dedup of already-seen
// values, and degrade-to-per-occurrence for novel values once the map is
// full. The process-global map is swapped out and restored so this test does
// not interfere with other tests' dedicated stale values.
//
// NOTE: no test in this package may call t.Parallel() — the swap/restore
// window here would race with any concurrent caller of
// warnUnknownCapability and silently corrupt dedup-count assertions.
func TestWarnUnknownCapability_DedupLimit(t *testing.T) {
	warnedUnknownCaps.Lock()
	saved := warnedUnknownCaps.seen
	warnedUnknownCaps.seen = make(map[Capability]struct{}, warnedUnknownCapsLimit)
	warnedUnknownCaps.Unlock()
	t.Cleanup(func() {
		warnedUnknownCaps.Lock()
		warnedUnknownCaps.seen = saved
		warnedUnknownCaps.Unlock()
	})

	var buf bytes.Buffer
	ctx := clog.WithLogger(t.Context(), clog.New(slog.NewTextHandler(&buf, nil)))

	const base = Capability(500000)
	for i := range warnedUnknownCapsLimit {
		warnUnknownCapability(ctx, base+Capability(i))
	}
	if n := strings.Count(buf.String(), "warn-dedup limit"); n != 1 {
		t.Fatalf("boundary log fired %d times while filling to the limit, want exactly 1", n)
	}

	// An already-seen value stays deduped after the map is full.
	warnUnknownCapability(ctx, base)
	if n := strings.Count(buf.String(), "skipping capability 500000 "); n != 1 {
		t.Errorf("already-seen value logged %d times, want exactly 1", n)
	}

	// A novel value past the limit degrades to logging every occurrence.
	novel := base + Capability(warnedUnknownCapsLimit)
	warnUnknownCapability(ctx, novel)
	warnUnknownCapability(ctx, novel)
	if n := strings.Count(buf.String(), "skipping capability 501024 "); n != 2 {
		t.Errorf("post-limit novel value logged %d times, want 2 (degraded, not silent)", n)
	}
	if n := strings.Count(buf.String(), "warn-dedup limit"); n != 1 {
		t.Errorf("boundary log fired %d times total, want exactly 1", n)
	}
}

func TestSetString_UnknownCap(t *testing.T) {
	const unknown = Capability(1601) // CAP_REGISTRY_PULL: reserved in the proto, removed in mono#23247
	if _, err := Stringify(unknown); err == nil {
		t.Fatal("capability 1601 has a descriptor again; pick another reserved value")
	}

	tests := []struct {
		name string
		caps Set
		want string
	}{{
		name: "unknown cap renders placeholder",
		caps: Set{Capability_CAP_IAM_GROUPS_LIST, unknown},
		want: "groups.list,unknown(cap=1601)",
	}, {
		name: "healthy set unchanged",
		caps: Set{Capability_CAP_IAM_GROUPS_LIST},
		want: "groups.list",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.caps.String(); got != test.want {
				t.Errorf("Set.String() = %q, want %q", got, test.want)
			}
		})
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
