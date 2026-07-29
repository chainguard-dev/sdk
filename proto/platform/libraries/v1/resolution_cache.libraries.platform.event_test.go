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

// TestResolutionCacheEventAnnotations pins the audit-event annotations on the
// mutating ResolutionCache methods: the events interceptor emits only for
// annotated methods, so a missing annotation silently drops the customer
// audit trail.
func TestResolutionCacheEventAnnotations(t *testing.T) {
	sd := File_resolution_cache_libraries_platform_proto.Services().ByName("ResolutionCache")
	if sd == nil {
		t.Fatal("ResolutionCache service not found in descriptors")
	}

	tests := []struct {
		method   string
		wantType string
		wantExts []string
	}{
		{"Zap", "dev.chainguard.api.libraries.resolutioncache.zapped.v1", []string{"group"}},
		{"SetOptOut", "dev.chainguard.api.libraries.resolutioncache.optout.updated.v1", []string{"group"}},
	}
	for _, tc := range tests {
		t.Run(tc.method, func(t *testing.T) {
			ext := getResolutionCacheEventAttributes(t, sd, tc.method)
			if got := ext.GetType(); got != tc.wantType {
				t.Errorf("type = %q, want %q", got, tc.wantType)
			}
			if got := ext.GetExtensions(); !slices.Equal(got, tc.wantExts) {
				t.Errorf("extensions = %v, want %v", got, tc.wantExts)
			}
			if got := ext.GetAudience(); got != cgannotations.EventAttributes_CUSTOMER {
				t.Errorf("audience = %v, want CUSTOMER", got)
			}
		})
	}

	// The read methods deliberately carry no event annotation.
	for _, method := range []string{"Status", "List"} {
		md := sd.Methods().ByName(protoreflect.Name(method))
		if md == nil {
			t.Fatalf("method %s not found", method)
		}
		opts, ok := md.Options().(*descriptorpb.MethodOptions)
		if !ok || opts == nil {
			continue
		}
		if ext, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes); ok && ext != nil {
			t.Errorf("method %s unexpectedly carries an events annotation: %v", method, ext)
		}
	}
}

func getResolutionCacheEventAttributes(t *testing.T, sd protoreflect.ServiceDescriptor, methodName string) *cgannotations.EventAttributes {
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
		t.Fatalf("method %s is missing the chainguard.annotations.events annotation", methodName)
	}
	return ext
}
