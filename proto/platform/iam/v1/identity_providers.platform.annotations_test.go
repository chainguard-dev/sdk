/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"slices"
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
	v1 "chainguard.dev/sdk/proto/platform/iam/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// TestScimLifecycleAnnotations asserts, at the proto-descriptor level, that the
// four v1 SCIM lifecycle RPCs carry the version-neutral event-type strings and
// the capability bundle CUS-1026 requires. Behavioral tests exercise redaction
// and delivery; this guards the wiring a typo would silently break — a wrong
// event-type string or a dropped capability.
func TestScimLifecycleAnnotations(t *testing.T) {
	sd := v1.File_identity_providers_platform_proto.Services().ByName("IdentityProviders")
	if sd == nil {
		t.Fatal("IdentityProviders service not found")
	}

	wantCaps := []capabilities.Capability{
		capabilities.Capability_CAP_IAM_SCIM_MANAGE,
		capabilities.Capability_CAP_IAM_IDENTITY_PROVIDERS_LIST,
	}

	for _, tc := range []struct {
		method   string
		wantType string
	}{
		{"GenerateScimToken", "dev.chainguard.api.iam.identity_providers.scim_token.generated.v1"},
		{"RegenerateScimToken", "dev.chainguard.api.iam.identity_providers.scim_token.regenerated.v1"},
		{"RevokeScimToken", "dev.chainguard.api.iam.identity_providers.scim_token.revoked.v1"},
		{"SetScimEnabled", "dev.chainguard.api.iam.identity_providers.scim_enabled.updated.v1"},
	} {
		t.Run(tc.method, func(t *testing.T) {
			md := sd.Methods().ByName(protoreflect.Name(tc.method))
			if md == nil {
				t.Fatalf("method %s not found", tc.method)
			}
			opts, ok := md.Options().(*descriptorpb.MethodOptions)
			if !ok || opts == nil {
				t.Fatalf("method %s has no options", tc.method)
			}

			ev, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
			if !ok || ev == nil {
				t.Fatalf("method %s missing events annotation", tc.method)
			}
			if got := ev.GetType(); got != tc.wantType {
				t.Errorf("event type: got = %q, want = %q", got, tc.wantType)
			}
			if got := ev.GetAudience(); got != cgannotations.EventAttributes_CUSTOMER {
				t.Errorf("audience: got = %v, want = CUSTOMER", got)
			}
			if got, want := ev.GetExtensions(), []string{"group"}; !slices.Equal(got, want) {
				t.Errorf("extensions: got = %v, want = %v", got, want)
			}

			iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
			if !ok || iam.GetEnabled() == nil {
				t.Fatalf("method %s missing IAM rules", tc.method)
			}
			gotCaps := iam.GetEnabled().GetCapabilities()
			for _, want := range wantCaps {
				if !slices.Contains(gotCaps, want) {
					t.Errorf("%s capabilities = %v, missing %v", tc.method, gotCaps, want)
				}
			}
		})
	}
}
