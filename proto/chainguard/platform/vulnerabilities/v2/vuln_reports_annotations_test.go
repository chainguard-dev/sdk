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
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

func vulnReportsService(t *testing.T) protoreflect.ServiceDescriptor {
	t.Helper()
	sd := File_chainguard_platform_vulnerabilities_v2_vuln_reports_proto.Services().ByName("VulnReportsService")
	if sd == nil {
		t.Fatal("VulnReportsService not found")
	}
	return sd
}

// Test_VulnReports_Annotations pins ListVulnCountReports' IAM shape at the
// proto-descriptor level: the v1 endpoint serves with CAP_VULN_REPORT_LIST
// unscoped, and the v2 endpoint must keep both or callers lose access on
// migration.
func Test_VulnReports_Annotations(t *testing.T) {
	md := vulnReportsService(t).Methods().ByName("ListVulnCountReports")
	if md == nil {
		t.Fatal("method ListVulnCountReports not found")
	}
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatal("method ListVulnCountReports has no options")
	}
	iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
	if !ok || iam.GetEnabled() == nil {
		t.Fatal("method ListVulnCountReports missing IAM rules")
	}
	if got := iam.GetEnabled().GetCapabilities(); !slices.Contains(got, capabilities.Capability_CAP_VULN_REPORT_LIST) {
		t.Errorf("capabilities = %v, missing CAP_VULN_REPORT_LIST", got)
	}
	if !iam.GetEnabled().GetUnscoped() {
		t.Error("ListVulnCountReports must be unscoped, matching the v1 endpoint")
	}
}
