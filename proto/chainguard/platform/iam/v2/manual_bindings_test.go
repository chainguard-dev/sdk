/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"slices"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"chainguard.dev/sdk/events"
	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
	"chainguard.dev/sdk/uidp"
)

func Test_IdentityProviders_ManualBindings_Annotations(t *testing.T) {
	for _, tt := range []struct {
		method string
		want   []capabilities.Capability
		mcp    func(*cgannotations.MCP) bool
		desc   string
	}{{
		method: "ListManualBindings",
		want:   []capabilities.Capability{capabilities.Capability_CAP_IAM_ROLE_BINDINGS_LIST, capabilities.Capability_CAP_IAM_SCIM_USERS_LIST},
		// The response is a customer's employee directory, so the list is
		// withheld from MCP as the SCIM user and group lists are.
		mcp:  func(m *cgannotations.MCP) bool { return m.GetReadOnly() && m.GetHidden() },
		desc: "read-only and hidden from MCP",
	}, {
		method: "RemoveManualBindings",
		want: []capabilities.Capability{
			capabilities.Capability_CAP_IAM_ROLE_BINDINGS_DELETE,
			// The dry run lists what would be removed, so it needs the list too.
			capabilities.Capability_CAP_IAM_ROLE_BINDINGS_LIST,
			capabilities.Capability_CAP_IAM_SCIM_MANAGE,
		},
		// A removal is never offered as a tool, and is not read-only.
		mcp:  func(m *cgannotations.MCP) bool { return m.GetHidden() && !m.GetReadOnly() },
		desc: "hidden from MCP and not read-only",
	}} {
		t.Run(tt.method, func(t *testing.T) {
			opts := methodOptions(t, identityProvidersService(), tt.method)
			iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
			if !ok || iam.GetEnabled() == nil {
				t.Fatalf("%s missing IAM rules", tt.method)
			}
			if diff := cmp.Diff(tt.want, iam.GetEnabled().GetCapabilities()); diff != "" {
				t.Errorf("capabilities (-want, +got):\n%s", diff)
			}
			mcp, ok := proto.GetExtension(opts, cgannotations.E_Mcp).(*cgannotations.MCP)
			if !ok || mcp == nil || !tt.mcp(mcp) {
				t.Errorf("mcp annotation: got = %v, want %s", mcp, tt.desc)
			}
		})
	}
}

func Test_IdentityProviders_RemoveManualBindings_EventAnnotations(t *testing.T) {
	runAnnotationTests(t, identityProvidersService(), []annotationTest{
		{"RemoveManualBindings", "dev.chainguard.api.iam.identity_providers.manual_bindings.removed.v1", []string{"group"}},
	})
}

// The delete event a removal publishes for each binding is the one deleting it
// through the single-binding RPC publishes, so the two cannot drift apart.
func Test_RemoveManualBindingsResponse_DeleteEventMatchesDeleteRoleBinding(t *testing.T) {
	sd := File_chainguard_platform_iam_v2_role_bindings_proto.Services().ByName("RoleBindingsService")
	if sd == nil {
		t.Fatal("RoleBindingsService not found")
	}
	want := getEventAttributes(t, sd, "DeleteRoleBinding")
	resp := &RemoveManualBindingsResponse{Removals: []*ManualBindingRemoval{{
		Binding: &ManualBinding{Uid: "a"},
		Outcome: ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REMOVED,
	}}}
	extras := resp.CloudEventsAdditional()
	if len(extras) != 1 {
		t.Fatalf("additional events: got = %d, want = 1", len(extras))
	}
	got := extras[0]
	if got.Type != want.GetType() || !slices.Equal(got.Extensions, want.GetExtensions()) || got.Audience != events.AudienceCustomer {
		t.Errorf("additional event = {%q %v %q}, want the DeleteRoleBinding event {%q %v %q}", got.Type, got.Extensions, got.Audience, want.GetType(), want.GetExtensions(), events.AudienceCustomer)
	}
}

// Only a removal that happened is an audit record: a dry run publishes nothing,
// and the record carries identifiers, classes and outcomes, not display names.
func Test_IdentityProviders_RemoveManualBindingsResponse_Event(t *testing.T) {
	idpUID := uidp.NewUIDP(uidp.NewUIDP("")).String()
	resp := &RemoveManualBindingsResponse{
		IdentityProviderUid: idpUID,
		PeopleLosingAccess:  1,
		MoreRemain:          true,
		PlanToken:           "token",
		Removals: []*ManualBindingRemoval{{
			IdentityUid: "identity-uid",
			Binding: &ManualBinding{
				Uid: "binding-uid", RoleUid: "role-uid", RoleName: "viewer", GroupUid: "group-uid", GroupName: "Engineering",
				BindingClass:       ManualBindingClass_MANUAL_BINDING_CLASS_NOT_COVERED,
				SupersedingSources: []*ScimBindingSource{{RoleBindingUid: "scim-binding", ScimGroupName: "Platform"}},
			},
			Outcome: ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REMOVED,
		}},
	}

	if got := resp.CloudEventsSubject(); got != idpUID {
		t.Errorf("subject: got = %q, want = %q", got, idpUID)
	}
	if got, ok := resp.CloudEventsExtension("group"); !ok || got != uidp.Parent(idpUID) {
		t.Errorf("group extension: got = %q (%v), want = %q", got, ok, uidp.Parent(idpUID))
	}
	if _, ok := resp.CloudEventsExtension("other"); ok {
		t.Error("an unknown extension is supported, want it reported as unsupported")
	}
	if resp.CloudEventsSuppress() {
		t.Error("an applied removal is suppressed, want published")
	}
	if !(&RemoveManualBindingsResponse{DryRun: true}).CloudEventsSuppress() {
		t.Error("a dry run is published, want suppressed")
	}

	want := &RemoveManualBindingsResponse{
		IdentityProviderUid: idpUID,
		PeopleLosingAccess:  1,
		MoreRemain:          true,
		Removals: []*ManualBindingRemoval{{
			IdentityUid: "identity-uid",
			Binding: &ManualBinding{
				Uid: "binding-uid", RoleUid: "role-uid", GroupUid: "group-uid",
				BindingClass: ManualBindingClass_MANUAL_BINDING_CLASS_NOT_COVERED,
			},
			Outcome: ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REMOVED,
		}},
	}
	if diff := cmp.Diff(want, resp.CloudEventsRedact(), protocmp.Transform()); diff != "" {
		t.Errorf("redacted body (-want +got):\n%s", diff)
	}
}

// Only a removed binding adds a role binding delete event, in the order the
// call handled them, and the group extension names the binding's group.
func Test_RemoveManualBindingsResponse_AdditionalEvents(t *testing.T) {
	const group = "aaaabbbbccccddddeeeeffff0000111122223333"
	removal := func(uid string, outcome ManualBindingRemovalOutcome) *ManualBindingRemoval {
		return &ManualBindingRemoval{Binding: &ManualBinding{Uid: group + "/" + uid}, Outcome: outcome}
	}
	resp := &RemoveManualBindingsResponse{Removals: []*ManualBindingRemoval{
		removal("a", ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REMOVED),
		removal("b", ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_ALREADY_GONE),
		removal("c", ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REFUSED),
		removal("d", ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_PLANNED),
		removal("e", ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REMOVED),
	}}
	extras := resp.CloudEventsAdditional()
	if len(extras) != 2 {
		t.Fatalf("additional events: got = %d, want = 2 (only the removed bindings)", len(extras))
	}
	for i, uid := range []string{group + "/a", group + "/e"} {
		if got := extras[i].Body.CloudEventsSubject(); got != uid {
			t.Errorf("subject %d: got = %q, want = %q", i, got, uid)
		}
		if ext, ok := extras[i].Body.(events.Extendable).CloudEventsExtension("group"); !ok || ext != group {
			t.Errorf("group extension %d: got = %q, %v, want = %q", i, ext, ok, group)
		}
	}
	if extras := (&RemoveManualBindingsResponse{}).CloudEventsAdditional(); len(extras) != 0 {
		t.Errorf("a response with no removals adds %d events, want none", len(extras))
	}
}
