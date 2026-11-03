/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
	"chainguard.dev/sdk/uidp"
)

func identityProvidersService() protoreflect.ServiceDescriptor {
	return File_chainguard_platform_iam_v2_identity_providers_proto.Services().ByName("IdentityProvidersService")
}

// The reset unlinks every provisioned user and removes their SCIM-managed role
// bindings, so it requires both the update and the SCIM management capability,
// and it is never an MCP tool.
func Test_IdentityProviders_ResetCorrelationRule_Annotations(t *testing.T) {
	opts := methodOptions(t, identityProvidersService(), "ResetCorrelationRule")

	iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
	if !ok || iam.GetEnabled() == nil {
		t.Fatal("ResetCorrelationRule missing IAM rules")
	}
	want := []capabilities.Capability{
		capabilities.Capability_CAP_IAM_IDENTITY_PROVIDERS_UPDATE,
		capabilities.Capability_CAP_IAM_SCIM_MANAGE,
	}
	if diff := cmp.Diff(want, iam.GetEnabled().GetCapabilities()); diff != "" {
		t.Errorf("ResetCorrelationRule capabilities (-want, +got):\n%s", diff)
	}

	mcp, ok := proto.GetExtension(opts, cgannotations.E_Mcp).(*cgannotations.MCP)
	if !ok || mcp == nil {
		t.Fatal("ResetCorrelationRule missing MCP annotation")
	}
	if !mcp.GetHidden() {
		t.Error("ResetCorrelationRule mcp.hidden: got = false, want = true")
	}

	runAnnotationTests(t, identityProvidersService(), []annotationTest{
		{"ResetCorrelationRule", "dev.chainguard.api.iam.identity_providers.correlation_rule.reset.v1", []string{"group"}},
	})
}

// The reset's event is its audit record: it keeps the outcome and the rule
// change, and redacts the identity provider as the identity provider's own
// events do.
func Test_IdentityProviders_ResetCorrelationRuleResponse_Event(t *testing.T) {
	idpUID := uidp.NewUIDP(uidp.NewUIDP("")).String()
	resp := &ResetCorrelationRuleResponse{
		IdentityProvider: &IdentityProvider{
			Uid:  idpUID,
			Name: "entra",
			Configuration: &IdentityProvider_Oidc{Oidc: &IdentityProvider_OIDC{
				Issuer:          "https://login.microsoftonline.com/tenant/v2.0",
				ClientId:        "client-id",
				ClientSecret:    "must-not-appear",
				CorrelationRule: IdentityProvider_OIDC_CORRELATION_RULE_OID_EQUALS_EXTERNAL_ID,
			}},
		},
		PreviousCorrelationRule: IdentityProvider_OIDC_CORRELATION_RULE_SUB_EQUALS_EXTERNAL_ID,
		Changed:                 true,
		UsersUnlinked:           12,
		BindingsDropped:         30,
		UnlinkComplete:          true,
	}

	if got := resp.CloudEventsSubject(); got != idpUID {
		t.Errorf("CloudEventsSubject(): got = %q, want = %q", got, idpUID)
	}
	if got, ok := resp.CloudEventsExtension("group"); !ok || got != uidp.Parent(idpUID) {
		t.Errorf("CloudEventsExtension(group): got = (%q, %t), want = (%q, true)", got, ok, uidp.Parent(idpUID))
	}
	// Delivery is detached from the request, so a caller that gave up while the
	// reset finished does not cancel its audit event.
	if !resp.CloudEventsAsync() {
		t.Error("CloudEventsAsync(): got = false, want = true")
	}

	want := &ResetCorrelationRuleResponse{
		IdentityProvider: &IdentityProvider{
			Uid:  idpUID,
			Name: "entra",
			Configuration: &IdentityProvider_Oidc{Oidc: &IdentityProvider_OIDC{
				Issuer:          "https://login.microsoftonline.com/tenant/v2.0",
				ClientId:        "client-id",
				CorrelationRule: IdentityProvider_OIDC_CORRELATION_RULE_OID_EQUALS_EXTERNAL_ID,
			}},
		},
		PreviousCorrelationRule: IdentityProvider_OIDC_CORRELATION_RULE_SUB_EQUALS_EXTERNAL_ID,
		Changed:                 true,
		UsersUnlinked:           12,
		BindingsDropped:         30,
		UnlinkComplete:          true,
	}
	if diff := cmp.Diff(want, resp.CloudEventsRedact(), protocmp.Transform()); diff != "" {
		t.Errorf("CloudEventsRedact() (-want, +got):\n%s", diff)
	}
	if resp.GetIdentityProvider().GetOidc().GetClientSecret() != "must-not-appear" {
		t.Error("CloudEventsRedact() modified the response it redacted")
	}
}
