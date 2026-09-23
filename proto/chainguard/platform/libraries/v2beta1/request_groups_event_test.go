/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"testing"

	"chainguard.dev/sdk/events"
	v2beta1 "chainguard.dev/sdk/proto/chainguard/platform/libraries/v2beta1"
)

// Both delete RPCs return Empty, so their request messages carry the event. The
// group extension is what routes it to an organization, and it is the one field
// neither message returns as part of a normal response -- nothing else would
// notice it being wrong.
//
// It was wrong: the extension used to be derived from uid with uidp.Parent, which
// answers "/" for a flat UUID, so every delete event named the root rather than
// the owning organization.
func Test_RequestGroup_Delete_EventInterfaces(t *testing.T) {
	const (
		org = "720909c9f7df8751f32a25ff2af1ef4c8f92dc8b"
		uid = "2f8d0c2e-6d8c-4c1a-9a4a-1a6b0a6b3f21"
	)

	for _, tc := range []struct {
		name string
		req  interface {
			events.Eventable
			events.Extendable
		}
	}{{
		name: "customer delete",
		req:  &v2beta1.DeleteRequestGroupRequest{Uid: uid, Parent: org},
	}, {
		name: "admin delete of a submitted group",
		req:  &v2beta1.DeleteSubmittedRequestGroupRequest{Uid: uid, Parent: org},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := tc.req.CloudEventsExtension("group")
			if !ok {
				t.Fatal("CloudEventsExtension(\"group\"): got ok = false, want true")
			}
			if got != org {
				t.Errorf("CloudEventsExtension(%q): got = %q, want = %q", "group", got, org)
			}

			if _, ok := tc.req.CloudEventsExtension("nonesuch"); ok {
				t.Error("CloudEventsExtension(\"nonesuch\"): got ok = true, want false")
			}

			if got := tc.req.CloudEventsSubject(); got != uid {
				t.Errorf("CloudEventsSubject(): got = %q, want = %q", got, uid)
			}
		})
	}
}
