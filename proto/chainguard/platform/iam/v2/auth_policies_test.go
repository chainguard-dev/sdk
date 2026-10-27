/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"slices"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
)

func authPoliciesService(t *testing.T) protoreflect.ServiceDescriptor {
	t.Helper()
	sd := File_chainguard_platform_iam_v2_auth_policies_proto.Services().ByName("AuthPoliciesService")
	if sd == nil {
		t.Fatal("AuthPoliciesService not found")
	}
	return sd
}

// TestAuthPolicyEventIsAsync pins that the policy's events are delivered off
// the request path (events.AsyncEventable): the sso_only toggle's audit event
// must not hold the response behind the event sink. A compile-time assertion
// only proves the method exists; this pins its value.
func TestAuthPolicyEventIsAsync(t *testing.T) {
	if !(&AuthPolicy{}).CloudEventsAsync() {
		t.Error("AuthPolicy.CloudEventsAsync() = false, want true")
	}
}

// TestAuthPoliciesCapabilities pins the authorization split: reading a policy
// needs only the list capability every role carries, while changing one needs
// the owner-only update capability.
func TestAuthPoliciesCapabilities(t *testing.T) {
	sd := authPoliciesService(t)
	for method, want := range map[string]capabilities.Capability{
		"GetAuthPolicy":    capabilities.Capability_CAP_IAM_AUTH_POLICIES_LIST,
		"UpdateAuthPolicy": capabilities.Capability_CAP_IAM_AUTH_POLICIES_UPDATE,
	} {
		iam, ok := proto.GetExtension(methodOptions(t, sd, method), cgannotations.E_Iam).(*cgannotations.IAM)
		if !ok || iam.GetEnabled() == nil {
			t.Fatalf("%s missing IAM rules", method)
		}
		if got := iam.GetEnabled().GetCapabilities(); !slices.Equal(got, []capabilities.Capability{want}) {
			t.Errorf("%s capabilities = %v, want [%v]", method, got, want)
		}
		if iam.GetEnabled().GetUnscoped() {
			t.Errorf("%s is unscoped; both RPCs must be scoped to the organization in uid", method)
		}
	}
}

// TestAuthPoliciesRolePlacement keeps the two capabilities on exactly the
// roles that hold the identity_providers capabilities they mirror: list in
// every role via ConsoleViewerCaps, update in OwnerCaps alone.
func TestAuthPoliciesRolePlacement(t *testing.T) {
	for _, tc := range []struct {
		name   string
		caps   []capabilities.Capability
		list   bool
		update bool
	}{
		{"ConsoleViewerCaps", capabilities.ConsoleViewerCaps, true, false},
		{"ViewerCaps", capabilities.ViewerCaps, true, false},
		{"EditorCaps", capabilities.EditorCaps, true, false},
		{"OwnerCaps", capabilities.OwnerCaps, true, true},
	} {
		if got := slices.Contains(tc.caps, capabilities.Capability_CAP_IAM_AUTH_POLICIES_LIST); got != tc.list {
			t.Errorf("%s contains auth_policies.list = %v, want %v", tc.name, got, tc.list)
		}
		if got := slices.Contains(tc.caps, capabilities.Capability_CAP_IAM_AUTH_POLICIES_UPDATE); got != tc.update {
			t.Errorf("%s contains auth_policies.update = %v, want %v", tc.name, got, tc.update)
		}
		// The mirror itself: wherever identity_providers.{list,update} is, so is ours.
		if got := slices.Contains(tc.caps, capabilities.Capability_CAP_IAM_IDENTITY_PROVIDERS_LIST); got != tc.list {
			t.Errorf("%s contains identity_providers.list = %v but auth_policies.list = %v", tc.name, got, tc.list)
		}
		if got := slices.Contains(tc.caps, capabilities.Capability_CAP_IAM_IDENTITY_PROVIDERS_UPDATE); got != tc.update {
			t.Errorf("%s contains identity_providers.update = %v but auth_policies.update = %v", tc.name, got, tc.update)
		}
	}
}

// TestAuthPoliciesUpdateHiddenFromMCP keeps UpdateAuthPolicy off any MCP
// surface. The backend builds its tools from v2beta1 descriptors only, so no
// RPC here is a tool today; the hidden flag is the guard for the day v2 joins
// that surface, because flipping sso_only can lock an organization's members
// out.
func TestAuthPoliciesUpdateHiddenFromMCP(t *testing.T) {
	mcp, ok := proto.GetExtension(methodOptions(t, authPoliciesService(t), "UpdateAuthPolicy"), cgannotations.E_Mcp).(*cgannotations.MCP)
	if !ok || mcp == nil {
		t.Fatal("UpdateAuthPolicy missing MCP annotation")
	}
	if !mcp.GetHidden() {
		t.Error("UpdateAuthPolicy mcp.hidden = false, want true")
	}
}

func TestAuthPoliciesEventAnnotations(t *testing.T) {
	runAnnotationTests(t, authPoliciesService(t), []annotationTest{
		{"UpdateAuthPolicy", "dev.chainguard.api.iam.auth_policies.updated.v1", []string{"group"}},
	})
	// Reads emit nothing.
	if opts := methodOptions(t, authPoliciesService(t), "GetAuthPolicy"); proto.HasExtension(opts, cgannotations.E_Events) {
		t.Error("GetAuthPolicy carries an events annotation; only the mutation should")
	}
}

func TestAuthPoliciesEventInterfaces(t *testing.T) {
	const orgUID = "0123456789abcdef0123456789abcdef01234567"

	p := &AuthPolicy{Uid: orgUID}
	if got := p.CloudEventsSubject(); got != orgUID {
		t.Errorf("AuthPolicy.CloudEventsSubject() = %q, want %q", got, orgUID)
	}
	if got, ok := p.CloudEventsExtension("group"); !ok || got != orgUID {
		t.Errorf("AuthPolicy.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, orgUID)
	}
	if _, ok := p.CloudEventsExtension("unknown"); ok {
		t.Error("AuthPolicy.CloudEventsExtension(unknown) returned true")
	}
}

// TestAuthPolicySsoOnlyPresence pins sso_only as a presence-tracked field in
// the descriptor: that is what keeps "never configured" distinct from an
// explicit false, and what lets an implied update mask carry a false.
func TestAuthPolicySsoOnlyPresence(t *testing.T) {
	fd := File_chainguard_platform_iam_v2_auth_policies_proto.Messages().ByName("AuthPolicy").Fields().ByName("sso_only")
	if fd == nil {
		t.Fatal("AuthPolicy.sso_only not found")
	}
	if !fd.HasPresence() {
		t.Error("AuthPolicy.sso_only has no presence; it must stay `optional` so unset and false remain distinguishable")
	}
}
