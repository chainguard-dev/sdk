/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"slices"
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// capabilityCase pins one RPC's IAM capability at the proto-descriptor
// level. A dropped or swapped capability silently changes who can manage
// overlays, so the CON-2306 split — reads behind
// CAP_REGISTRY_OVERLAYS_LIST, mutations behind CAP_REGISTRY_OVERLAYS_EDIT
// — is pinned per method.
type capabilityCase struct {
	method string
	want   capabilities.Capability
}

func runCapabilityTests(t *testing.T, sd protoreflect.ServiceDescriptor, cases []capabilityCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.method, func(t *testing.T) {
			md := sd.Methods().ByName(protoreflect.Name(tc.method))
			if md == nil {
				t.Fatalf("method %s not found", tc.method)
			}
			opts, ok := md.Options().(*descriptorpb.MethodOptions)
			if !ok || opts == nil {
				t.Fatalf("method %s has no options", tc.method)
			}
			iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
			if !ok || iam.GetEnabled() == nil {
				t.Fatalf("method %s missing IAM rules", tc.method)
			}
			if got := iam.GetEnabled().GetCapabilities(); !slices.Contains(got, tc.want) {
				t.Errorf("%s capabilities = %v, missing %v", tc.method, got, tc.want)
			}
		})
	}
}

func Test_Overlays_Annotations(t *testing.T) {
	runCapabilityTests(t, overlaysService(t), []capabilityCase{
		{"CreateOverlay", capabilities.Capability_CAP_REGISTRY_OVERLAYS_EDIT},
		{"DeleteOverlay", capabilities.Capability_CAP_REGISTRY_OVERLAYS_EDIT},
		{"GetOverlay", capabilities.Capability_CAP_REGISTRY_OVERLAYS_LIST},
		{"ListOverlays", capabilities.Capability_CAP_REGISTRY_OVERLAYS_LIST},
	})
}

func Test_OverlayBindings_Annotations(t *testing.T) {
	runCapabilityTests(t, overlayBindingsService(t), []capabilityCase{
		{"CreateOverlayBinding", capabilities.Capability_CAP_REGISTRY_OVERLAYS_EDIT},
		{"DeleteOverlayBinding", capabilities.Capability_CAP_REGISTRY_OVERLAYS_EDIT},
		{"GetOverlayBinding", capabilities.Capability_CAP_REGISTRY_OVERLAYS_LIST},
		{"ListOverlayBindings", capabilities.Capability_CAP_REGISTRY_OVERLAYS_LIST},
	})
}

// Test_Overlay_RoleGrants pins the role wiring: viewers can list, registry
// editors can list and edit.
func Test_Overlay_RoleGrants(t *testing.T) {
	if !slices.Contains(capabilities.ConsoleViewerCaps, capabilities.Capability_CAP_REGISTRY_OVERLAYS_LIST) {
		t.Error("ConsoleViewerCaps missing CAP_REGISTRY_OVERLAYS_LIST")
	}
	for _, want := range []capabilities.Capability{
		capabilities.Capability_CAP_REGISTRY_OVERLAYS_LIST,
		capabilities.Capability_CAP_REGISTRY_OVERLAYS_EDIT,
	} {
		if !slices.Contains(capabilities.RegistryEditorCaps, want) {
			t.Errorf("RegistryEditorCaps missing %v", want)
		}
		if !slices.Contains(capabilities.EditorCaps, want) {
			t.Errorf("EditorCaps missing %v", want)
		}
		if !slices.Contains(capabilities.OwnerCaps, want) {
			t.Errorf("OwnerCaps missing %v", want)
		}
	}
}
