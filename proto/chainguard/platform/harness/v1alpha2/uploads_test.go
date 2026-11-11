/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha2

import (
	stdcmp "cmp"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/proto/capabilities"
)

func methodOptions(t *testing.T, method protoreflect.Name) *descriptorpb.MethodOptions {
	t.Helper()
	md := File_chainguard_platform_harness_v1alpha2_uploads_proto.Services().ByName("Uploads").Methods().ByName(method)
	if md == nil {
		t.Fatalf("Uploads.%s not found", method)
	}
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatalf("Uploads.%s has no options", method)
	}
	return opts
}

func TestUploadsMethods(t *testing.T) {
	create, list := capabilities.Capability_CAP_HARNESS_UPLOADS_CREATE, capabilities.Capability_CAP_HARNESS_UPLOADS_LIST
	tests := []struct {
		method    protoreflect.Name
		wantCap   capabilities.Capability
		wantRoute string
	}{
		{method: "GetOrgPublicKey", wantCap: create, wantRoute: "/harness/v1alpha2/orgs/{parent_id=**}/publicKey"},
		{method: "CreateUpload", wantCap: create, wantRoute: "/harness/v1alpha2/uploads/{parent_id=**}"},
		{method: "WriteUpload", wantCap: create, wantRoute: "/harness/v1alpha2/uploads/{id=**}:writeUpload"},
		{method: "CompleteUpload", wantCap: create, wantRoute: "/harness/v1alpha2/uploads/{id=**}:completeUpload"},
		{method: "GetUpload", wantCap: list, wantRoute: "/harness/v1alpha2/uploads/{id=**}"},
		{method: "ListUploads", wantCap: list, wantRoute: "/harness/v1alpha2/uploads"},
	}
	if got, want := File_chainguard_platform_harness_v1alpha2_uploads_proto.Services().ByName("Uploads").Methods().Len(), len(tests); got != want {
		t.Errorf("Uploads methods: got = %d, want = %d", got, want)
	}
	for _, tc := range tests {
		t.Run(string(tc.method), func(t *testing.T) {
			opts := methodOptions(t, tc.method)

			iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
			if !ok || iam == nil {
				t.Fatal("no iam annotation")
			}
			if diff := cmp.Diff([]capabilities.Capability{tc.wantCap}, iam.GetEnabled().GetCapabilities()); diff != "" {
				t.Errorf("capabilities (-want, +got):\n%s", diff)
			}

			rule, ok := proto.GetExtension(opts, annotations.E_Http).(*annotations.HttpRule)
			if !ok || rule == nil {
				t.Fatal("no google.api.http annotation")
			}
			route := stdcmp.Or(rule.GetGet(), rule.GetPost())
			if route != tc.wantRoute {
				t.Errorf("route: got = %q, want = %q", route, tc.wantRoute)
			}
			// A route outside /harness/v1alpha2/ could collide with a v1alpha1
			// route on the shared gateway.
			if !strings.HasPrefix(route, "/harness/v1alpha2/") {
				t.Errorf("route %q is outside /harness/v1alpha2/", route)
			}
		})
	}
}

// Every request carries the organization the identity interceptor authorizes
// against. WriteUpload and CompleteUpload are routed by id, which does not
// carry it.
func TestUploadsRequestsHaveIAMScope(t *testing.T) {
	sd := File_chainguard_platform_harness_v1alpha2_uploads_proto.Services().ByName("Uploads")
	for i := range sd.Methods().Len() {
		md := sd.Methods().Get(i)
		t.Run(string(md.Name()), func(t *testing.T) {
			f := md.Input().Fields().ByName("parent_id")
			if f == nil {
				t.Fatalf("%s has no parent_id", md.Input().Name())
			}
			opts, ok := f.Options().(*descriptorpb.FieldOptions)
			if !ok || !proto.GetExtension(opts, cgannotations.E_IamScope).(bool) {
				t.Errorf("%s.parent_id is not the iam_scope", md.Input().Name())
			}
		})
	}
}

// WriteUpload carries sealed bytes, which an agent never has.
func TestWriteUploadHiddenFromMCP(t *testing.T) {
	mcp, ok := proto.GetExtension(methodOptions(t, "WriteUpload"), cgannotations.E_Mcp).(*cgannotations.MCP)
	if !ok || !mcp.GetHidden() {
		t.Error("Uploads.WriteUpload mcp annotation: hidden = false, want true")
	}
}

// Each CreateUpload mints a fresh upload, so nothing on the API is content-addressed.
func TestCreateUploadIsNotContentAddressed(t *testing.T) {
	mcp, ok := proto.GetExtension(methodOptions(t, "CreateUpload"), cgannotations.E_Mcp).(*cgannotations.MCP)
	if !ok || mcp == nil {
		t.Fatal("Uploads.CreateUpload has no mcp annotation")
	}
	if mcp.GetIdempotent() {
		t.Error("Uploads.CreateUpload mcp annotation: idempotent = true, want false")
	}
	fd := File_chainguard_platform_harness_v1alpha2_uploads_proto
	for _, f := range []struct{ msg, field protoreflect.Name }{
		{"Upload", "name"},
		{"CreateUploadResponse", "created"},
	} {
		if fd.Messages().ByName(f.msg).Fields().ByName(f.field) != nil {
			t.Errorf("%s.%s exists, want no content-derived field", f.msg, f.field)
		}
	}
}

// Sealed bytes reach storage only through WriteUpload, which api-impl relays.
// A signed URL would let a client write to storage around it.
func TestNoSignedUploadURL(t *testing.T) {
	fd := File_chainguard_platform_harness_v1alpha2_uploads_proto
	if fd.Messages().ByName("UploadUrl") != nil {
		t.Error("message UploadUrl exists")
	}
	if fd.Messages().ByName("CreateUploadResponse").Fields().ByName("upload_url") != nil {
		t.Error("CreateUploadResponse.upload_url exists")
	}
}
