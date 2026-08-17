/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"encoding/hex"
	"slices"
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/uidp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// getEventAttributes extracts EventAttributes from a proto method descriptor.
func getEventAttributes(t *testing.T, sd protoreflect.ServiceDescriptor, methodName string) *cgannotations.EventAttributes {
	t.Helper()
	md := sd.Methods().ByName(protoreflect.Name(methodName))
	if md == nil {
		t.Fatalf("method %s not found in service %s", methodName, sd.FullName())
	}
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatalf("method %s has no options", methodName)
	}
	ext, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
	if !ok || ext == nil {
		t.Fatalf("method %s is missing chainguard.annotations.events", methodName)
	}
	return ext
}

// annotationTest defines a test case for verifying proto event annotations.
type annotationTest struct {
	method   string
	wantType string
	wantExts []string
}

func runAnnotationTests(t *testing.T, sd protoreflect.ServiceDescriptor, tests []annotationTest) {
	t.Helper()
	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			ea := getEventAttributes(t, sd, tt.method)
			if ea.GetType() != tt.wantType {
				t.Errorf("type = %q, want %q", ea.GetType(), tt.wantType)
			}
			if ea.GetAudience() != cgannotations.EventAttributes_CUSTOMER {
				t.Errorf("audience = %v, want CUSTOMER", ea.GetAudience())
			}
			if !slices.Equal(ea.GetExtensions(), tt.wantExts) {
				t.Errorf("extensions = %v, want %v", ea.GetExtensions(), tt.wantExts)
			}
		})
	}
}

func Test_Overlays_EventAnnotations(t *testing.T) {
	runAnnotationTests(t, overlaysService(t), []annotationTest{
		{"CreateOverlay", "dev.chainguard.api.platform.registry.overlay.created.v1", []string{"group"}},
		{"DeleteOverlay", "dev.chainguard.api.platform.registry.overlay.deleted.v1", []string{"group"}},
	})
}

func Test_OverlayBindings_EventAnnotations(t *testing.T) {
	runAnnotationTests(t, overlayBindingsService(t), []annotationTest{
		{"CreateOverlayBinding", "dev.chainguard.api.platform.registry.overlay_binding.created.v1", []string{"group"}},
		{"DeleteOverlayBinding", "dev.chainguard.api.platform.registry.overlay_binding.deleted.v1", []string{"group"}},
	})
}

// assertNoEventAnnotations verifies read-only methods carry no event
// annotations.
func assertNoEventAnnotations(t *testing.T, sd protoreflect.ServiceDescriptor, methods ...string) {
	t.Helper()
	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			md := sd.Methods().ByName(protoreflect.Name(method))
			if md == nil {
				t.Fatalf("method %s not found", method)
			}
			opts, ok := md.Options().(*descriptorpb.MethodOptions)
			if !ok || opts == nil {
				return
			}
			ext, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
			if ok && ext != nil && ext.GetType() != "" {
				t.Errorf("read-only method %s carries event annotation %q", method, ext.GetType())
			}
		})
	}
}

func Test_Overlays_ReadOnlyMethodsHaveNoEvents(t *testing.T) {
	assertNoEventAnnotations(t, overlaysService(t), "GetOverlay", "ListOverlays")
}

func Test_OverlayBindings_ReadOnlyMethodsHaveNoEvents(t *testing.T) {
	assertNoEventAnnotations(t, overlayBindingsService(t), "GetOverlayBinding", "ListOverlayBindings")
}

// eventBody is the Eventable + Extendable surface the events interceptor
// consumes.
type eventBody interface {
	CloudEventsSubject() string
	CloudEventsExtension(string) (string, bool)
}

// runEventInterfaceTests verifies subject = uid and the group extension =
// the uid's root for each body.
func runEventInterfaceTests(t *testing.T, uid string, cases map[string]eventBody) {
	t.Helper()
	group := uidp.Root(uid)
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			if got := body.CloudEventsSubject(); got != uid {
				t.Errorf("CloudEventsSubject() = %q, want %q", got, uid)
			}
			if got, ok := body.CloudEventsExtension("group"); !ok || got != group {
				t.Errorf("CloudEventsExtension(group) = (%q, %v), want (%q, true)", got, ok, group)
			}
			if _, ok := body.CloudEventsExtension("unknown"); ok {
				t.Error("CloudEventsExtension(unknown) returned true")
			}
		})
	}
}

func Test_Overlays_EventInterfaces(t *testing.T) {
	uid := hex.EncodeToString([]byte("group")) + "/" + hex.EncodeToString([]byte("overlay"))
	runEventInterfaceTests(t, uid, map[string]eventBody{
		"Overlay":              &Overlay{Uid: uid},
		"DeleteOverlayRequest": &DeleteOverlayRequest{Uid: uid},
	})
}

func Test_OverlayBindings_EventInterfaces(t *testing.T) {
	uid := hex.EncodeToString([]byte("group")) + "/" + hex.EncodeToString([]byte("binding"))
	runEventInterfaceTests(t, uid, map[string]eventBody{
		"OverlayBinding":              &OverlayBinding{Uid: uid},
		"DeleteOverlayBindingRequest": &DeleteOverlayBindingRequest{Uid: uid},
	})
}
