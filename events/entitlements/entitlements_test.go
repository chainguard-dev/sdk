/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package entitlements_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/events/entitlements"
)

// ChangeEvent must satisfy events.Extendable as a value, since the events
// framework probes the Occurrence body for extensions.
var _ events.Extendable = entitlements.ChangeEvent{}

func TestChangedEventType(t *testing.T) {
	if got := entitlements.ChangedEventType; got != "dev.chainguard.entitlement.changed.v1" {
		t.Errorf("ChangedEventType = %q, want dev.chainguard.entitlement.changed.v1", got)
	}
}

// TestChangeEvent_JSONRoundTrip pins the wire shape: marshaling, unmarshaling,
// and re-marshaling a ChangeEvent must produce identical JSON for each
// operation grain (write-update, delete, add-resource). This is the
// encode/decode round-trip property for the event body, checked across the
// grain variants rather than one hand-picked value.
func TestChangeEvent_JSONRoundTrip(t *testing.T) {
	for _, tc := range []struct {
		name string
		ev   entitlements.ChangeEvent
	}{
		{
			name: "write-update (before+after)",
			ev: entitlements.ChangeEvent{
				Domain:          "IMAGE",
				OrgUIDP:         "abcd1234/ef567890",
				EntitlementUIDP: "abcd1234/ef567890/aabbccdd",
				Operation:       entitlements.OperationWrite,
				Source:          "SALESFORCE",
				Actor:           "user/42",
				ChangeReason:    "contract renewal",
				Before:          json.RawMessage(`{"max_quota":5}`),
				After:           json.RawMessage(`{"max_quota":10}`),
			},
		},
		{
			name: "delete (before only, no after)",
			ev: entitlements.ChangeEvent{
				Domain:          "IMAGE",
				OrgUIDP:         "abcd1234/ef567890",
				EntitlementUIDP: "abcd1234/ef567890/aabbccdd",
				Operation:       entitlements.OperationDelete,
				Actor:           "user/42",
				Before:          json.RawMessage(`{"max_quota":10}`),
			},
		},
		{
			name: "add-resource (resources, no entitlement snapshot)",
			ev: entitlements.ChangeEvent{
				Domain:          "IMAGE",
				OrgUIDP:         "abcd1234/ef567890",
				EntitlementUIDP: "abcd1234/ef567890/aabbccdd",
				Operation:       entitlements.OperationAddResource,
				Resources: []entitlements.ResourceChange{
					{Resource: "repo-nginx", Tier: "APPLICATION"},
					{Resource: "repo-go", Tier: "BASE"},
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			b1, err := json.Marshal(tc.ev)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}
			var got entitlements.ChangeEvent
			if err := json.Unmarshal(b1, &got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}
			b2, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("re-marshal: %v", err)
			}
			if !bytes.Equal(b1, b2) {
				t.Errorf("round-trip changed the JSON:\n before=%s\n after =%s", b1, b2)
			}
			if got.Operation != tc.ev.Operation {
				t.Errorf("operation not preserved: got %q want %q", got.Operation, tc.ev.Operation)
			}
		})
	}
}

// The TestChangeEvent_GoldenJSONKeys_* pair pins the exact wire keys of a
// ChangeEvent (and its nested ResourceChange). A struct-tag rename — e.g.
// org_uidp→orgUidp or resource→resourceName — changes the marshaled key and fails here.
// The round-trip test can't catch that because it renames the encode and decode
// side together, so a renamed tag still round-trips cleanly. json.Marshal emits
// struct fields in declaration order deterministically, so a literal-string
// compare of the whole document is a stable assertion on the key set. Two
// fixtures are needed because no single event can legally populate every field:
// the grain rules (TestChangeEvent_GrainInvariants) forbid Before/After and
// Resources on the same event, so the WRITE fixture pins the entitlement-grain
// keys and the ADD_RESOURCE fixture pins the resource-grain keys. Each fixture
// is a well-formed event of its grain, safe to feed through any validation the
// grain rules grow into.

// TestChangeEvent_GoldenJSONKeys_Write pins the entitlement-grain keys: domain,
// org_uidp, entitlement_uidp, operation, source, actor, change_reason, before,
// after.
func TestChangeEvent_GoldenJSONKeys_Write(t *testing.T) {
	ev := entitlements.ChangeEvent{
		Domain:          "IMAGE",
		OrgUIDP:         "abcd1234/ef567890",
		EntitlementUIDP: "abcd1234/ef567890/aabbccdd",
		Operation:       entitlements.OperationWrite,
		Source:          "SALESFORCE",
		Actor:           "user/42",
		ChangeReason:    "contract renewal",
		Before:          json.RawMessage(`{"max_quota":5}`),
		After:           json.RawMessage(`{"max_quota":10}`),
	}

	const want = `{"domain":"IMAGE","org_uidp":"abcd1234/ef567890","entitlement_uidp":"abcd1234/ef567890/aabbccdd","operation":"WRITE","source":"SALESFORCE","actor":"user/42","change_reason":"contract renewal","before":{"max_quota":5},"after":{"max_quota":10}}`

	got, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(got) != want {
		t.Errorf("ChangeEvent JSON keys drifted:\n got=%s\nwant=%s", got, want)
	}
}

// TestChangeEvent_GoldenJSONKeys_AddResource pins the resource-grain keys: the
// resources array and its nested resource, tier.
func TestChangeEvent_GoldenJSONKeys_AddResource(t *testing.T) {
	ev := entitlements.ChangeEvent{
		Domain:          "IMAGE",
		OrgUIDP:         "abcd1234/ef567890",
		EntitlementUIDP: "abcd1234/ef567890/aabbccdd",
		Operation:       entitlements.OperationAddResource,
		Resources: []entitlements.ResourceChange{
			{Resource: "repo-nginx", Tier: "APPLICATION"},
		},
	}

	const want = `{"domain":"IMAGE","org_uidp":"abcd1234/ef567890","entitlement_uidp":"abcd1234/ef567890/aabbccdd","operation":"ADD_RESOURCE","resources":[{"resource":"repo-nginx","tier":"APPLICATION"}]}`

	got, err := json.Marshal(ev)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(got) != want {
		t.Errorf("ChangeEvent JSON keys drifted:\n got=%s\nwant=%s", got, want)
	}
}

// TestResourceChange_GoldenJSONKeys pins the nested ResourceChange keys on their
// own so a tag rename there is caught even if ResourceChange is ever emitted
// outside a ChangeEvent.
func TestResourceChange_GoldenJSONKeys(t *testing.T) {
	rc := entitlements.ResourceChange{
		Resource: "repo-nginx",
		Tier:     "APPLICATION",
	}

	const want = `{"resource":"repo-nginx","tier":"APPLICATION"}`

	got, err := json.Marshal(rc)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if string(got) != want {
		t.Errorf("ResourceChange JSON keys drifted:\n got=%s\nwant=%s", got, want)
	}
}

// TestOperationValues pins the closed set of operation tokens so a rename can't
// silently change the wire contract subscribers switch on.
func TestOperationValues(t *testing.T) {
	for op, want := range map[entitlements.Operation]string{
		entitlements.OperationWrite:          "WRITE",
		entitlements.OperationDelete:         "DELETE",
		entitlements.OperationAddResource:    "ADD_RESOURCE",
		entitlements.OperationRemoveResource: "REMOVE_RESOURCE",
	} {
		if string(op) != want {
			t.Errorf("operation token = %q, want %q", op, want)
		}
	}
}

// TestChangeEvent_CloudEventsExtension pins the extension attributes a
// ChangeEvent exposes for subscriber-side filtering: the known keys must
// surface the org, domain, and operation from the event, and any other key
// must report absent — never an empty-but-present attribute.
func TestChangeEvent_CloudEventsExtension(t *testing.T) {
	ev := entitlements.ChangeEvent{
		Domain:          entitlements.DomainHelmChart,
		OrgUIDP:         "9f8e7d6c/5b4a3210",
		EntitlementUIDP: "9f8e7d6c/5b4a3210/deadbeef",
		Operation:       entitlements.OperationRemoveResource,
	}

	for _, tc := range []struct {
		name    string
		key     string
		want    string
		wantSet bool
	}{
		{"group key surfaces org", events.GroupKey, "9f8e7d6c/5b4a3210", true},
		{"domain key surfaces domain", events.EntitlementDomainKey, "HELM_CHART", true},
		{"operation key surfaces operation", events.EntitlementOperationKey, "REMOVE_RESOURCE", true},
		{"unknown key absent", "fricassee", "", false},
		{"empty key absent", "", "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, gotSet := ev.CloudEventsExtension(tc.key)
			if got != tc.want || gotSet != tc.wantSet {
				t.Errorf("CloudEventsExtension(%q): got = (%q, %t), want = (%q, %t)", tc.key, got, gotSet, tc.want, tc.wantSet)
			}
		})
	}
}

// TestDomainValues pins the canonical Domain wire tokens so a rename can't
// silently change the vocabulary emitters and subscribers share; agreement
// with the service's own domain vocabulary is asserted on the emitter side.
func TestDomainValues(t *testing.T) {
	for domain, want := range map[string]string{
		entitlements.DomainImage:     "IMAGE",
		entitlements.DomainLibrary:   "LIBRARY",
		entitlements.DomainPackage:   "PACKAGE",
		entitlements.DomainHelmChart: "HELM_CHART",
		entitlements.DomainFeature:   "FEATURE",
	} {
		if domain != want {
			t.Errorf("domain token = %q, want %q", domain, want)
		}
	}
}

// TestChangeEvent_GrainInvariants pins the Operation⇄field-population rules from
// the ChangeEvent doc as a mechanism, not just prose. grainViolation encodes the
// contract — "" iff an event populates exactly the fields its grain allows — and
// the test drives it in BOTH directions: every well-formed grain is accepted,
// and each illegal combination the struct can represent (DELETE-with-After,
// resource-op-with-Before, WRITE-without-After, …) is rejected. Asserting only
// the well-formed side would pass even for a producer that mis-emits, since the
// fixtures already obey the rules; the negative cases are what actually hold a
// producer to the contract. Because the emitting RPCs don't exist yet (they
// panic("unimplemented")), this is the only such guard. Iterating the closed
// Operation set on the valid side, plus an unknown-operation negative case,
// keeps the switch total: a new token without a grain rule falls to default and
// fails.
func TestChangeEvent_GrainInvariants(t *testing.T) {
	before := json.RawMessage(`{"max_quota":5}`)
	after := json.RawMessage(`{"max_quota":10}`)
	res := []entitlements.ResourceChange{{Resource: "repo-nginx", Tier: "APPLICATION"}}
	const entUIDP = "abcd1234/ef567890/aabbccdd"

	// grainViolation returns "" iff ev populates exactly the fields its Operation
	// grain permits; otherwise the first violation. An unknown operation is itself
	// a violation — the total-switch guard. EntitlementUIDP (like Domain, OrgUIDP,
	// and Operation) is grain-unrestricted and must be non-empty for every
	// operation.
	grainViolation := func(ev entitlements.ChangeEvent) string {
		if ev.EntitlementUIDP == "" {
			return "EntitlementUIDP must be set for every operation"
		}
		switch ev.Operation {
		case entitlements.OperationWrite:
			switch {
			case ev.After == nil:
				return "WRITE: After must be set"
			case len(ev.Resources) != 0:
				return "WRITE: Resources must be empty"
			}
		case entitlements.OperationDelete:
			switch {
			case ev.Before == nil:
				return "DELETE: Before must be set"
			case ev.After != nil:
				return "DELETE: After must be nil"
			case len(ev.Resources) != 0:
				return "DELETE: Resources must be empty"
			}
		case entitlements.OperationAddResource, entitlements.OperationRemoveResource:
			switch {
			case ev.Before != nil || ev.After != nil:
				return "resource op: Before and After must be nil"
			case len(ev.Resources) == 0:
				return "resource op: Resources must be non-empty"
			}
		default:
			return "unknown operation: " + string(ev.Operation)
		}
		return ""
	}

	// Well-formed events for every grain must be accepted. Iterating all four
	// operations is the totality check on the valid side: a new Operation without
	// a grain rule hits default and is reported as unknown.
	t.Run("valid", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			ev   entitlements.ChangeEvent
		}{
			{"write-create", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationWrite, After: after}},
			{"write-update", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationWrite, Before: before, After: after}},
			{"delete", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationDelete, Before: before}},
			{"add-resource", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationAddResource, Resources: res}},
			{"remove-resource", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationRemoveResource, Resources: res}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if v := grainViolation(tc.ev); v != "" {
					t.Errorf("well-formed %s rejected: %s", tc.name, v)
				}
			})
		}
	})

	// Each illegal combination the struct can represent must be rejected. Without
	// these the rules are only ever asserted against fixtures that already obey
	// them, so a mis-emitting producer would still pass CI.
	t.Run("invalid", func(t *testing.T) {
		for _, tc := range []struct {
			name string
			ev   entitlements.ChangeEvent
		}{
			{"WRITE without After", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationWrite}},
			{"WRITE with Resources", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationWrite, After: after, Resources: res}},
			{"DELETE with After", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationDelete, Before: before, After: after}},
			{"DELETE without Before", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationDelete}},
			{"resource op with Before", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationAddResource, Before: before, Resources: res}},
			{"resource op without Resources", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.OperationRemoveResource}},
			{"unknown operation", entitlements.ChangeEvent{EntitlementUIDP: entUIDP, Operation: entitlements.Operation("FRICASSEE"), After: after}},
			{"missing EntitlementUIDP", entitlements.ChangeEvent{Operation: entitlements.OperationWrite, After: after}},
		} {
			t.Run(tc.name, func(t *testing.T) {
				if grainViolation(tc.ev) == "" {
					t.Errorf("illegal event %q accepted, want rejected", tc.name)
				}
			})
		}
	})
}
