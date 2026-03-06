/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package annotations_test

import (
	"fmt"

	"chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/proto/capabilities"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Example demonstrates creating an IAM annotation with disabled mode.
func Example_iamDisabled() {
	iam := &annotations.IAM{
		Mode: &annotations.IAM_Disabled{
			Disabled: &emptypb.Empty{},
		},
	}
	fmt.Printf("IAM disabled: %v\n", iam.GetDisabled() != nil)
	// Output: IAM disabled: true
}

// Example demonstrates creating an IAM annotation with enabled mode and rules.
func Example_iamEnabled() {
	iam := &annotations.IAM{
		Mode: &annotations.IAM_Enabled{
			Enabled: &annotations.IAM_Rules{
				Capabilities: []capabilities.Capability{
					capabilities.Capability_CAP_IAM_GROUPS_LIST,
				},
				Unscoped: false,
			},
		},
	}
	rules := iam.GetEnabled()
	fmt.Printf("IAM enabled with %d capabilities\n", len(rules.GetCapabilities()))
	// Output: IAM enabled with 1 capabilities
}

// Example demonstrates creating an IAM Rules annotation with unscoped access.
func Example_iamRulesUnscoped() {
	rules := &annotations.IAM_Rules{
		Capabilities: []capabilities.Capability{
			capabilities.Capability_CAP_IAM_GROUPS_LIST,
			capabilities.Capability_CAP_IAM_GROUPS_CREATE,
		},
		Unscoped: true,
	}
	fmt.Printf("Unscoped: %v, Capabilities: %d\n", rules.GetUnscoped(), len(rules.GetCapabilities()))
	// Output: Unscoped: true, Capabilities: 2
}

// Example demonstrates creating EventAttributes for internal events.
func Example_eventAttributesInternal() {
	attrs := &annotations.EventAttributes{
		Type:       "dev.chainguard.api.platform.created.v1",
		Extensions: []string{"subject", "source"},
		Audience:   annotations.EventAttributes_INTERNAL,
	}
	fmt.Printf("Event type: %s, Audience: %s\n", attrs.GetType(), attrs.GetAudience())
	// Output: Event type: dev.chainguard.api.platform.created.v1, Audience: INTERNAL
}

// Example demonstrates creating EventAttributes for customer-facing events.
func Example_eventAttributesCustomer() {
	attrs := &annotations.EventAttributes{
		Type:       "dev.chainguard.api.iam.group.created.v1",
		Extensions: []string{"actor", "group"},
		Audience:   annotations.EventAttributes_CUSTOMER,
	}
	fmt.Printf("Event type: %s, Audience: %s, Extensions: %d\n",
		attrs.GetType(), attrs.GetAudience(), len(attrs.GetExtensions()))
	// Output: Event type: dev.chainguard.api.iam.group.created.v1, Audience: CUSTOMER, Extensions: 2
}
