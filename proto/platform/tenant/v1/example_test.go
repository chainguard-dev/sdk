/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"fmt"

	tenantv1 "chainguard.dev/sdk/proto/platform/tenant/v1"
)

// ExampleSbom2_Source demonstrates the Sbom2_Source enum values.
func ExampleSbom2_Source() {
	fmt.Println(tenantv1.Sbom2_UNKNOWN)
	fmt.Println(tenantv1.Sbom2_INGESTED)
	fmt.Println(tenantv1.Sbom2_GENERATED)
	// Output:
	// UNKNOWN
	// INGESTED
	// GENERATED
}
