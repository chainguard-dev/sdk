/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"slices"
	"testing"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

// Test_StigReports_Annotations pins GetStigReport's IAM capabilities at the
// proto-descriptor level: the per-image registry read precedent
// (CAP_REPO_LIST + CAP_MANIFEST_METADATA_LIST). A dropped or swapped
// capability silently changes who can read STIG reports.
func Test_StigReports_Annotations(t *testing.T) {
	sd := File_chainguard_platform_registry_v2_stig_reports_proto.Services().ByName("StigReportsService")
	if sd == nil {
		t.Fatal("StigReportsService not found")
	}
	md := sd.Methods().ByName("GetStigReport")
	if md == nil {
		t.Fatal("method GetStigReport not found")
	}
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatal("method GetStigReport has no options")
	}
	iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
	if !ok || iam.GetEnabled() == nil {
		t.Fatal("method GetStigReport missing IAM rules")
	}
	got := iam.GetEnabled().GetCapabilities()
	for _, want := range []capabilities.Capability{
		capabilities.Capability_CAP_REPO_LIST,
		capabilities.Capability_CAP_MANIFEST_METADATA_LIST,
	} {
		if !slices.Contains(got, want) {
			t.Errorf("GetStigReport capabilities = %v, missing %v", got, want)
		}
	}
	if iam.GetEnabled().GetUnscoped() {
		t.Error("GetStigReport must be a scoped read, got unscoped")
	}
}
