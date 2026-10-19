/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/proto/capabilities"
)

// Uploads reuses the hunt capabilities rather than taking new bits.
func TestUploadsCapabilities(t *testing.T) {
	sd := File_chainguard_platform_harness_v1alpha1_uploads_proto.Services().ByName("Uploads")
	create, list := capabilities.Capability_CAP_HARNESS_HUNTS_CREATE, capabilities.Capability_CAP_HARNESS_HUNTS_LIST
	for method, want := range map[protoreflect.Name]capabilities.Capability{
		"GetOrgPublicKey": create,
		"CreateUpload":    create,
		"CompleteUpload":  create,
		"GetUpload":       list,
		"ListUploads":     list,
	} {
		md := sd.Methods().ByName(method)
		if md == nil {
			t.Errorf("Uploads.%s not found", method)
			continue
		}
		iam, ok := proto.GetExtension(md.Options().(*descriptorpb.MethodOptions), cgannotations.E_Iam).(*cgannotations.IAM)
		if !ok || iam == nil {
			t.Errorf("Uploads.%s has no iam annotation", method)
			continue
		}
		if diff := cmp.Diff([]capabilities.Capability{want}, iam.GetEnabled().GetCapabilities()); diff != "" {
			t.Errorf("Uploads.%s capabilities (-want, +got):\n%s", method, diff)
		}
	}
}

// Each CreateUpload mints a fresh upload, so nothing on the API is content-addressed.
func TestCreateUploadIsNotContentAddressed(t *testing.T) {
	fd := File_chainguard_platform_harness_v1alpha1_uploads_proto
	md := fd.Services().ByName("Uploads").Methods().ByName("CreateUpload")
	mcp, ok := proto.GetExtension(md.Options().(*descriptorpb.MethodOptions), cgannotations.E_Mcp).(*cgannotations.MCP)
	if !ok || mcp == nil {
		t.Fatal("Uploads.CreateUpload has no mcp annotation")
	}
	if mcp.GetIdempotent() {
		t.Error("Uploads.CreateUpload mcp annotation: idempotent = true, want false")
	}
	for _, f := range []struct{ msg, field protoreflect.Name }{
		{"Upload", "name"},
		{"CreateUploadResponse", "created"},
	} {
		if fd.Messages().ByName(f.msg).Fields().ByName(f.field) != nil {
			t.Errorf("%s.%s exists, want no content-derived field", f.msg, f.field)
		}
	}
}
