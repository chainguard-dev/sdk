/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Spec-conformance harness for containers/openspec/changes/
// add-overlays-platform-api/specs/custom-assembly/overlay-api/spec.md.
// Test_Conformance_GatewayRouteOrder pins gateway dispatch: the
// /registry/v2beta1/overlayBindings/... routes are a distinct literal
// prefix from /registry/v2beta1/overlays/..., so the overlay routes'
// {uid=**}/{parent=**} wildcards can never swallow them regardless of
// registration order.
// Exercises the invariants deliverable at the SDK layer: the RPC surface
// (no update RPCs), the packages-only payload enforced by proto shape,
// the no-inline-content binding shape, and gateway dispatch.
package v2beta1

import (
	"context"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/emptypb"
)

func overlaysService(t *testing.T) protoreflect.ServiceDescriptor {
	t.Helper()
	sd := File_chainguard_platform_registry_v2beta1_overlays_proto.Services().ByName("OverlaysService")
	if sd == nil {
		t.Fatal("OverlaysService not found")
	}
	return sd
}

func overlayBindingsService(t *testing.T) protoreflect.ServiceDescriptor {
	t.Helper()
	sd := File_chainguard_platform_registry_v2beta1_overlay_bindings_proto.Services().ByName("OverlayBindingsService")
	if sd == nil {
		t.Fatal("OverlayBindingsService not found")
	}
	return sd
}

// Requirement: Overlay lifecycle / Binding lifecycle — each service's
// RPC set is exactly the four CRUD-without-update methods; no update RPC
// exists on either surface.
func Test_Conformance_RPCSurface(t *testing.T) {
	for _, tc := range []struct {
		sd   protoreflect.ServiceDescriptor
		want []string
	}{
		{overlaysService(t), []string{"CreateOverlay", "GetOverlay", "ListOverlays", "DeleteOverlay"}},
		{overlayBindingsService(t), []string{"CreateOverlayBinding", "GetOverlayBinding", "ListOverlayBindings", "DeleteOverlayBinding"}},
	} {
		t.Run(string(tc.sd.Name()), func(t *testing.T) {
			got := make([]string, 0, tc.sd.Methods().Len())
			for i := range tc.sd.Methods().Len() {
				name := string(tc.sd.Methods().Get(i).Name())
				if strings.Contains(name, "Update") {
					t.Errorf("update RPC %q exists; the milestone workflow is delete-and-recreate", name)
				}
				got = append(got, name)
			}
			want := slices.Clone(tc.want)
			slices.Sort(want)
			slices.Sort(got)
			if !slices.Equal(want, got) {
				t.Errorf("RPC surface: got = %v, want = %v", got, want)
			}
		})
	}
}

// Requirement: Packages-only payload by schema — no message reachable from
// any Overlays or OverlayBindings RPC references ImageContents,
// CustomOverlay, or any other content-model message. Certificates,
// accounts, environment, repositories, and keyring are unrepresentable,
// not merely rejected.
func Test_Conformance_PackagesOnlyByShape(t *testing.T) {
	reachable := map[protoreflect.FullName]bool{}
	var walk func(md protoreflect.MessageDescriptor)
	walk = func(md protoreflect.MessageDescriptor) {
		if reachable[md.FullName()] {
			return
		}
		reachable[md.FullName()] = true
		for i := range md.Fields().Len() {
			fd := md.Fields().Get(i)
			if fd.Kind() == protoreflect.MessageKind || fd.Kind() == protoreflect.GroupKind {
				walk(fd.Message())
			}
		}
	}
	for _, sd := range []protoreflect.ServiceDescriptor{overlaysService(t), overlayBindingsService(t)} {
		for i := range sd.Methods().Len() {
			walk(sd.Methods().Get(i).Input())
			walk(sd.Methods().Get(i).Output())
		}
	}

	for name := range reachable {
		s := string(name)
		if strings.Contains(s, "ImageContents") || strings.Contains(s, "CustomOverlay") || strings.Contains(s, "apko") {
			t.Errorf("content-model message %q is reachable from the Overlays surface", s)
		}
	}

	// The Overlay payload is exactly: uid, name, packages, referenced_by.
	overlay := (&Overlay{}).ProtoReflect().Descriptor()
	fields := make([]string, 0, overlay.Fields().Len())
	for i := range overlay.Fields().Len() {
		fields = append(fields, string(overlay.Fields().Get(i).Name()))
	}
	slices.Sort(fields)
	if want := []string{"name", "packages", "referenced_by", "uid"}; !slices.Equal(fields, want) {
		t.Errorf("Overlay fields: got = %v, want = %v", fields, want)
	}
	if fd := overlay.Fields().ByName("packages"); !fd.IsList() || fd.Kind() != protoreflect.StringKind {
		t.Error("Overlay.packages must be a repeated string")
	}
}

// Requirement: Binding lifecycle — one repo per binding, overlay
// referenced by a plain string (UIDP or name), exact tags only, and no
// inline-content field of any kind.
func Test_Conformance_BindingShape(t *testing.T) {
	create := (&CreateOverlayBindingRequest{}).ProtoReflect().Descriptor()
	fields := make([]string, 0, create.Fields().Len())
	for i := range create.Fields().Len() {
		fields = append(fields, string(create.Fields().Get(i).Name()))
	}
	slices.Sort(fields)
	if want := []string{"overlay", "parent", "tags"}; !slices.Equal(fields, want) {
		t.Errorf("CreateOverlayBindingRequest fields: got = %v, want = %v", fields, want)
	}
	if fd := create.Fields().ByName("overlay"); fd.Kind() != protoreflect.StringKind {
		t.Error("CreateOverlayBindingRequest.overlay must be a plain string reference — inline content is unrepresentable")
	}

	// Audit-complete shapes: binding embeds the overlay, its repo, and
	// tags; the list request narrows by scope and by overlay.
	att := (&OverlayBinding{}).ProtoReflect().Descriptor()
	attFields := make([]string, 0, att.Fields().Len())
	for i := range att.Fields().Len() {
		attFields = append(attFields, string(att.Fields().Get(i).Name()))
	}
	slices.Sort(attFields)
	if want := []string{"overlay", "repo", "tags", "uid"}; !slices.Equal(attFields, want) {
		t.Errorf("OverlayBinding fields: got = %v, want = %v", attFields, want)
	}
	list := (&ListOverlayBindingsRequest{}).ProtoReflect().Descriptor()
	for _, name := range []string{"uidp", "overlay"} {
		if list.Fields().ByName(protoreflect.Name(name)) == nil {
			t.Errorf("ListOverlayBindingsRequest missing field %q", name)
		}
	}
}

// routeRecorder records which RPC the gateway dispatched to.
type routeRecorder struct {
	UnimplementedOverlaysServiceServer
	UnimplementedOverlayBindingsServiceServer
	method string
	uid    string
}

func (r *routeRecorder) CreateOverlay(_ context.Context, req *CreateOverlayRequest) (*Overlay, error) {
	r.method, r.uid = "CreateOverlay", req.GetParent()
	return &Overlay{}, nil
}

func (r *routeRecorder) GetOverlay(_ context.Context, req *GetOverlayRequest) (*Overlay, error) {
	r.method, r.uid = "GetOverlay", req.GetUid()
	return &Overlay{}, nil
}

func (r *routeRecorder) ListOverlays(context.Context, *ListOverlaysRequest) (*ListOverlaysResponse, error) {
	r.method, r.uid = "ListOverlays", ""
	return &ListOverlaysResponse{}, nil
}

func (r *routeRecorder) DeleteOverlay(_ context.Context, req *DeleteOverlayRequest) (*emptypb.Empty, error) {
	r.method, r.uid = "DeleteOverlay", req.GetUid()
	return &emptypb.Empty{}, nil
}

func (r *routeRecorder) CreateOverlayBinding(_ context.Context, req *CreateOverlayBindingRequest) (*OverlayBinding, error) {
	r.method, r.uid = "CreateOverlayBinding", req.GetParent()
	return &OverlayBinding{}, nil
}

func (r *routeRecorder) GetOverlayBinding(_ context.Context, req *GetOverlayBindingRequest) (*OverlayBinding, error) {
	r.method, r.uid = "GetOverlayBinding", req.GetUid()
	return &OverlayBinding{}, nil
}

func (r *routeRecorder) ListOverlayBindings(context.Context, *ListOverlayBindingsRequest) (*ListOverlayBindingsResponse, error) {
	r.method, r.uid = "ListOverlayBindings", ""
	return &ListOverlayBindingsResponse{}, nil
}

func (r *routeRecorder) DeleteOverlayBinding(_ context.Context, req *DeleteOverlayBindingRequest) (*emptypb.Empty, error) {
	r.method, r.uid = "DeleteOverlayBinding", req.GetUid()
	return &emptypb.Empty{}, nil
}

// The gateway must route /overlayBindings/... to binding RPCs and
// /overlays/... to overlay RPCs; the distinct collection prefixes make
// dispatch independent of registration order.
func Test_Conformance_GatewayRouteOrder(t *testing.T) {
	rec := &routeRecorder{}
	mux := runtime.NewServeMux()
	// Registration order is the invariant: overlays first, bindings second.
	if err := RegisterOverlaysServiceHandlerServer(t.Context(), mux, rec); err != nil {
		t.Fatalf("RegisterOverlaysServiceHandlerServer: %v", err)
	}
	if err := RegisterOverlayBindingsServiceHandlerServer(t.Context(), mux, rec); err != nil {
		t.Fatalf("RegisterOverlayBindingsServiceHandlerServer: %v", err)
	}

	for _, tc := range []struct {
		verb, path, body string
		wantMethod       string
		wantUID          string
	}{
		{"GET", "/registry/v2beta1/overlayBindings", "", "ListOverlayBindings", ""},
		{"GET", "/registry/v2beta1/overlayBindings/org/repo/att", "", "GetOverlayBinding", "org/repo/att"},
		{"DELETE", "/registry/v2beta1/overlayBindings/org/repo/att", "", "DeleteOverlayBinding", "org/repo/att"},
		{"POST", "/registry/v2beta1/overlayBindings/org/repo", `{"overlay":"debug-tools","tags":["latest"]}`, "CreateOverlayBinding", "org/repo"},
		{"GET", "/registry/v2beta1/overlays", "", "ListOverlays", ""},
		{"GET", "/registry/v2beta1/overlays/org/ovl", "", "GetOverlay", "org/ovl"},
		{"DELETE", "/registry/v2beta1/overlays/org/ovl", "", "DeleteOverlay", "org/ovl"},
		{"POST", "/registry/v2beta1/overlays/org", `{"name":"n"}`, "CreateOverlay", "org"},
	} {
		t.Run(tc.verb+" "+tc.path, func(t *testing.T) {
			*rec = routeRecorder{}
			req := httptest.NewRequest(tc.verb, tc.path, strings.NewReader(tc.body))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			if rec.method != tc.wantMethod {
				t.Fatalf("routed to %q (status %d), want %q", rec.method, w.Code, tc.wantMethod)
			}
			if rec.uid != tc.wantUID {
				t.Errorf("bound uid: got = %q, want = %q", rec.uid, tc.wantUID)
			}
		})
	}
}
