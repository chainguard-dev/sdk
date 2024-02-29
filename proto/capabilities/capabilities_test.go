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
	all = sortCaps(all)

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
		name: "gulfstream",
		caps: Set{Capability_CAP_GULFSTREAM},
	}, {
		name: "all",
		caps: all,
	}, {
		// sortCaps removes duplicates.
		name: "duplicates",
		caps: sortCaps(Set{Capability_CAP_GULFSTREAM, Capability_CAP_GULFSTREAM, Capability_CAP_GULFSTREAM}),
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
