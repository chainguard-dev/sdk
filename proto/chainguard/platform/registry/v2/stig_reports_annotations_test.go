/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Test_StigReports_Annotations pins the IAM capabilities shared by both STIG
// report reads and keeps the raw download out of MCP tool discovery. A dropped
// or swapped capability silently changes who can read reports.
func Test_StigReports_Annotations(t *testing.T) {
	sd := File_chainguard_platform_registry_v2_stig_reports_proto.Services().ByName("StigReportsService")
	if sd == nil {
		t.Fatal("StigReportsService not found")
	}
	wantCapabilities := []capabilities.Capability{
		capabilities.Capability_CAP_IAM_GROUPS_LIST,
		capabilities.Capability_CAP_REPO_LIST,
		capabilities.Capability_CAP_MANIFEST_LIST,
		capabilities.Capability_CAP_TAG_LIST,
		capabilities.Capability_CAP_REFERRERS_LIST,
		capabilities.Capability_CAP_REPO_BLOBS_GET,
	}
	methods := []struct {
		name       protoreflect.Name
		wantHidden bool
	}{
		{name: "GetStigReport"},
		{name: "DownloadStigReport", wantHidden: true},
	}
	for _, method := range methods {
		t.Run(string(method.name), func(t *testing.T) {
			md := sd.Methods().ByName(method.name)
			if md == nil {
				t.Fatalf("method %s not found", method.name)
			}
			opts, ok := md.Options().(*descriptorpb.MethodOptions)
			if !ok || opts == nil {
				t.Fatalf("method %s has no options", method.name)
			}
			iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
			if !ok || iam.GetEnabled() == nil {
				t.Fatalf("method %s missing IAM rules", method.name)
			}
			if diff := cmp.Diff(wantCapabilities, iam.GetEnabled().GetCapabilities()); diff != "" {
				t.Errorf("%s capabilities (-want, +got):\n%s", method.name, diff)
			}
			if iam.GetEnabled().GetUnscoped() {
				t.Errorf("%s must be a scoped read, got unscoped", method.name)
			}
			mcp, ok := proto.GetExtension(opts, cgannotations.E_Mcp).(*cgannotations.MCP)
			if !ok {
				t.Fatalf("method %s missing MCP rules", method.name)
			}
			if got := mcp.GetHidden(); got != method.wantHidden {
				t.Errorf("%s MCP hidden = %t, want %t", method.name, got, method.wantHidden)
			}
		})
	}
}
