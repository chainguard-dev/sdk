/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"slices"
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func TestBinding_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "group extension returns parent UIDP",
			id:     "abc123/def456/binding789",
			key:    "group",
			want:   "abc123/def456",
			wantOk: true,
		},
		{
			name:   "unknown extension returns false",
			id:     "abc123/def456/binding789",
			key:    "unknown",
			want:   "",
			wantOk: false,
		},
		{
			name:   "group extension with empty id",
			id:     "",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
		{
			name:   "group extension with single segment id",
			id:     "abc123",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binding{Id: tt.id}
			got, ok := b.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestBinding_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "non-empty id",
			id:   "abc123/def456/binding789",
			want: "abc123/def456/binding789",
		},
		{
			name: "empty id",
			id:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binding{Id: tt.id}
			if got := b.CloudEventsSubject(); got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeleteBindingRequest_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "group extension returns parent UIDP",
			id:     "abc123/def456/binding789",
			key:    "group",
			want:   "abc123/def456",
			wantOk: true,
		},
		{
			name:   "unknown extension returns false",
			id:     "abc123/def456/binding789",
			key:    "unknown",
			want:   "",
			wantOk: false,
		},
		{
			name:   "group extension with empty id",
			id:     "",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
		{
			name:   "group extension with single segment id",
			id:     "abc123",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeleteBindingRequest{Id: tt.id}
			got, ok := req.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeleteBindingRequest_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "non-empty id",
			id:   "abc123/def456/binding789",
			want: "abc123/def456/binding789",
		},
		{
			name: "empty id",
			id:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeleteBindingRequest{Id: tt.id}
			if got := req.CloudEventsSubject(); got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestPolicy_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "group extension returns parent UIDP",
			id:     "abc123/def456/policy789",
			key:    "group",
			want:   "abc123/def456",
			wantOk: true,
		},
		{
			name:   "unknown extension returns false",
			id:     "abc123/def456/policy789",
			key:    "unknown",
			want:   "",
			wantOk: false,
		},
		{
			name:   "group extension with empty id",
			id:     "",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
		{
			name:   "group extension with single segment id",
			id:     "abc123",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Policy{Id: tt.id}
			got, ok := p.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestPolicy_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "non-empty id",
			id:   "abc123/def456/policy789",
			want: "abc123/def456/policy789",
		},
		{
			name: "empty id",
			id:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Policy{Id: tt.id}
			if got := p.CloudEventsSubject(); got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeletePolicyRequest_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "group extension returns parent UIDP",
			id:     "abc123/def456/policy789",
			key:    "group",
			want:   "abc123/def456",
			wantOk: true,
		},
		{
			name:   "unknown extension returns false",
			id:     "abc123/def456/policy789",
			key:    "unknown",
			want:   "",
			wantOk: false,
		},
		{
			name:   "group extension with empty id",
			id:     "",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
		{
			name:   "group extension with single segment id",
			id:     "abc123",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeletePolicyRequest{Id: tt.id}
			got, ok := req.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeletePolicyRequest_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "non-empty id",
			id:   "abc123/def456/policy789",
			want: "abc123/def456/policy789",
		},
		{
			name: "empty id",
			id:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeletePolicyRequest{Id: tt.id}
			if got := req.CloudEventsSubject(); got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

// annotationTest defines a test case for verifying proto event annotations.
type annotationTest struct {
	method   string
	wantType string
	wantExts []string
}

// getEventAttributes extracts EventAttributes from a proto method descriptor,
// failing the test if the annotation is missing.
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

// TestPoliciesEventAnnotations locks in the proto-level event annotations on
// the Policies service: each mutating RPC declares the right type string,
// CUSTOMER audience, and ["group"] extension. A regression that drops or
// renames one of these annotations would silently stop emitting audit
// events; this catches it at build time.
func TestPoliciesEventAnnotations(t *testing.T) {
	sd := File_policies_platform_proto.Services().ByName("Policies")
	if sd == nil {
		t.Fatal("Policies service not found")
	}
	runAnnotationTests(t, sd, []annotationTest{{
		method:   "CreatePolicy",
		wantType: "dev.chainguard.api.policies.policies.created.v1",
		wantExts: []string{"group"},
	}, {
		method:   "UpdatePolicy",
		wantType: "dev.chainguard.api.policies.policies.updated.v1",
		wantExts: []string{"group"},
	}, {
		method:   "DeletePolicy",
		wantType: "dev.chainguard.api.policies.policies.deleted.v1",
		wantExts: []string{"group"},
	}})
}

// TestBindingsEventAnnotations is the matching regression guard for the
// Bindings service. The annotations were added in PR #41069; this test
// pins them in place so a future rename or removal is caught here.
func TestBindingsEventAnnotations(t *testing.T) {
	sd := File_policies_platform_proto.Services().ByName("Bindings")
	if sd == nil {
		t.Fatal("Bindings service not found")
	}
	runAnnotationTests(t, sd, []annotationTest{{
		method:   "CreateBinding",
		wantType: "dev.chainguard.api.policies.bindings.created.v1",
		wantExts: []string{"group"},
	}, {
		method:   "UpdateBinding",
		wantType: "dev.chainguard.api.policies.bindings.updated.v1",
		wantExts: []string{"group"},
	}, {
		method:   "DeleteBinding",
		wantType: "dev.chainguard.api.policies.bindings.deleted.v1",
		wantExts: []string{"group"},
	}})
}
