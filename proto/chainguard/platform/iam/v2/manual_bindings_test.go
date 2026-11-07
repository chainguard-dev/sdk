/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"

	cgannotations "chainguard.dev/sdk/proto/annotations"
	capabilities "chainguard.dev/sdk/proto/capabilities"
)

func Test_IdentityProviders_ManualBindings_Annotations(t *testing.T) {
	for _, tt := range []struct {
		method string
		want   []capabilities.Capability
		mcp    func(*cgannotations.MCP) bool
		desc   string
	}{{
		method: "ListManualBindings",
		want:   []capabilities.Capability{capabilities.Capability_CAP_IAM_ROLE_BINDINGS_LIST, capabilities.Capability_CAP_IAM_SCIM_USERS_LIST},
		// The response is a customer's employee directory, so the list is
		// withheld from MCP as the SCIM user and group lists are.
		mcp:  func(m *cgannotations.MCP) bool { return m.GetReadOnly() && m.GetHidden() },
		desc: "read-only and hidden from MCP",
	}} {
		t.Run(tt.method, func(t *testing.T) {
			opts := methodOptions(t, identityProvidersService(), tt.method)
			iam, ok := proto.GetExtension(opts, cgannotations.E_Iam).(*cgannotations.IAM)
			if !ok || iam.GetEnabled() == nil {
				t.Fatalf("%s missing IAM rules", tt.method)
			}
			if diff := cmp.Diff(tt.want, iam.GetEnabled().GetCapabilities()); diff != "" {
				t.Errorf("capabilities (-want, +got):\n%s", diff)
			}
			mcp, ok := proto.GetExtension(opts, cgannotations.E_Mcp).(*cgannotations.MCP)
			if !ok || mcp == nil || !tt.mcp(mcp) {
				t.Errorf("mcp annotation: got = %v, want %s", mcp, tt.desc)
			}
		})
	}
}
