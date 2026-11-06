/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"slices"
	"testing"
	"time"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/durationpb"

	"chainguard.dev/sdk/events"
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
		// Starting a validation is part of configuring the policy (owner
		// tier), and reading a session's status is held to the same
		// capability: its diagnostics name the person who signed in.
		"StartValidation":     capabilities.Capability_CAP_IAM_AUTH_POLICIES_UPDATE,
		"GetValidationStatus": capabilities.Capability_CAP_IAM_AUTH_POLICIES_UPDATE,
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

// TestGetValidationStatusBinding pins the route: the session is addressed as
// a resource under the policy, so its id travels as a path segment. The path
// offers no secrecy a query string lacks (request logs keep both); the pin is
// about the route's shape, which the generated gateway and the SDK clients
// build on.
func TestGetValidationStatusBinding(t *testing.T) {
	rule, ok := proto.GetExtension(methodOptions(t, authPoliciesService(t), "GetValidationStatus"), annotations.E_Http).(*annotations.HttpRule)
	if !ok || rule == nil {
		t.Fatal("GetValidationStatus missing an HTTP binding")
	}
	const want = "/iam/v2/authPolicies/{uid}/validations/{session_id}"
	if got := rule.GetGet(); got != want {
		t.Errorf("GetValidationStatus GET path = %q, want %q", got, want)
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
	// Both mutations: flipping sso_only can lock members out, and a
	// validation drives a real sign-in a person must complete.
	for _, method := range []string{"UpdateAuthPolicy", "StartValidation"} {
		mcp, ok := proto.GetExtension(methodOptions(t, authPoliciesService(t), method), cgannotations.E_Mcp).(*cgannotations.MCP)
		if !ok || mcp == nil {
			t.Fatalf("%s missing MCP annotation", method)
		}
		if !mcp.GetHidden() {
			t.Errorf("%s mcp.hidden = false, want true", method)
		}
	}
}

func TestAuthPoliciesEventAnnotations(t *testing.T) {
	runAnnotationTests(t, authPoliciesService(t), []annotationTest{
		{"UpdateAuthPolicy", "dev.chainguard.api.iam.auth_policies.updated.v1", []string{"group"}},
		{"StartValidation", "dev.chainguard.api.iam.auth_policies.validation.started.v1", []string{"group"}},
	})
	// Reads emit nothing.
	for _, method := range []string{"GetAuthPolicy", "GetValidationStatus"} {
		if opts := methodOptions(t, authPoliciesService(t), method); proto.HasExtension(opts, cgannotations.E_Events) {
			t.Errorf("%s carries an events annotation; only the mutations should", method)
		}
	}
}

// TestValidationEventRedaction pins what the started event carries: the
// organization, kind, designee and state, and never the open_url or the
// session id, which together are the sign-in's only credential.
func TestValidationEventRedaction(t *testing.T) {
	const orgUID = "0123456789abcdef0123456789abcdef01234567"
	v := &Validation{
		Id: "secret-session-id", Uid: orgUID, Kind: ValidationKind_VALIDATION_KIND_BREAK_GLASS,
		IdentityUidp: "designated-break-glass-identity",
		State:        ValidationState_VALIDATION_STATE_PENDING, OpenUrl: "https://issuer/validate?session=secret-session-id",
	}
	if got := v.CloudEventsSubject(); got != orgUID {
		t.Errorf("Validation.CloudEventsSubject() = %q, want %q", got, orgUID)
	}
	if got, ok := v.CloudEventsExtension("group"); !ok || got != orgUID {
		t.Errorf("Validation.CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, orgUID)
	}
	redacted, ok := v.CloudEventsRedact().(*Validation)
	if !ok {
		t.Fatalf("CloudEventsRedact() = %T, want *Validation", v.CloudEventsRedact())
	}
	if redacted.GetOpenUrl() != "" || redacted.GetId() != "" {
		t.Errorf("redacted event carries open_url=%q id=%q, want neither", redacted.GetOpenUrl(), redacted.GetId())
	}
	// The event's audience is every subscriber of the organization, not only
	// the owners who may read the policy's break-glass designations, so the
	// designee is not part of it either.
	if redacted.GetIdentityUidp() != "" {
		t.Errorf("redacted event carries identity_uidp=%q, want none", redacted.GetIdentityUidp())
	}
	if redacted.GetUid() != orgUID || redacted.GetKind() != v.GetKind() || redacted.GetState() != v.GetState() {
		t.Errorf("redacted event lost uid/kind/state: %v", redacted)
	}
}

// TestAuthPolicyEventRedaction pins what the updated event carries: the
// policy's uid and its writable fields, sso_only and session_duration, and
// nothing else. The UpdateAuthPolicy response answers a caller who may enable
// sso_only with the validation view (which designees have a current
// break-glass validation), the readiness verdict and the entitlement, and the
// event interceptor publishes that response as the CUSTOMER event body; all
// three are held to that caller on the read side, so the event must not
// carry them to every subscriber of the organization. The writable fields
// stay: the event (and the previous state the interceptor runs through the
// same redaction) is the audit record of what an update changed.
func TestAuthPolicyEventRedaction(t *testing.T) {
	const orgUID = "0123456789abcdef0123456789abcdef01234567"
	p := &AuthPolicy{
		Uid:             orgUID,
		SsoOnly:         proto.Bool(true),
		SessionDuration: durationpb.New(time.Hour),
		Validations: &ValidationStatus{
			IdpValidationCurrent: true,
			BreakGlassDesignees: []*BreakGlassValidationStatus{
				{IdentityUidp: "designated-break-glass-identity", ValidationCurrent: true},
			},
		},
		Readiness:           &SsoOnlyReadiness{IdpValidatedCurrent: true, BreakGlassReady: true},
		EnforcementEntitled: proto.Bool(true),
	}
	r, ok := any(p).(events.Redactable)
	if !ok {
		t.Fatal("AuthPolicy does not implement events.Redactable: the validation view reaches the updated event body")
	}
	redacted, ok := r.CloudEventsRedact().(*AuthPolicy)
	if !ok {
		t.Fatalf("CloudEventsRedact() = %T, want *AuthPolicy", r.CloudEventsRedact())
	}
	if want := (&AuthPolicy{Uid: orgUID, SsoOnly: proto.Bool(true), SessionDuration: durationpb.New(time.Hour)}); !proto.Equal(redacted, want) {
		t.Errorf("redacted event = %v, want exactly uid, sso_only and session_duration: %v", redacted, want)
	}
	// The response the caller receives is not the event body: the redaction
	// copies, so the owner still sees the view, and the body shares no
	// storage with it.
	if p.GetValidations() == nil || len(p.GetValidations().GetBreakGlassDesignees()) != 1 || p.GetReadiness() == nil {
		t.Errorf("CloudEventsRedact() mutated the response: %v", p)
	}
	*redacted.SsoOnly = false
	if !p.GetSsoOnly() {
		t.Error("CloudEventsRedact() shares sso_only storage with the response")
	}
	if redacted.SessionDuration == nil {
		t.Fatalf("redacted event dropped session_duration: %v", redacted)
	}
	redacted.SessionDuration.Seconds = 0
	if p.GetSessionDuration().AsDuration() != time.Hour {
		t.Error("CloudEventsRedact() shares session_duration storage with the response")
	}
	// An unset session_duration stays unset: "Chainguard's default" must not
	// become an explicit zero, which only Chainguard sets.
	if r := (&AuthPolicy{Uid: orgUID}).CloudEventsRedact().(*AuthPolicy); r.SessionDuration != nil {
		t.Errorf("redacted unset session_duration = %v, want nil", r.SessionDuration)
	}
	// An unset sso_only stays unset: redaction must not turn "never
	// configured" into an explicit false.
	if r := (&AuthPolicy{Uid: orgUID}).CloudEventsRedact().(*AuthPolicy); r.SsoOnly != nil {
		t.Errorf("redacted unset sso_only = %v, want nil", *r.SsoOnly)
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
