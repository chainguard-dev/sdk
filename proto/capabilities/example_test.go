/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities_test

import (
	"fmt"

	"chainguard.dev/sdk/proto/capabilities"
)

// ExampleNames demonstrates listing all known capability names.
func ExampleNames() {
	names := capabilities.Names()
	fmt.Println(len(names) > 0)
	// Output:
	// true
}

// ExampleParse demonstrates parsing a capability by name.
func ExampleParse() {
	// Capability_UNKNOWN is the zero value and is not in the name map.
	capability, err := capabilities.Parse("unknown-capability-name")
	fmt.Println(err)
	fmt.Println(capability == capabilities.Capability_UNKNOWN)
	// Output:
	// <nil>
	// true
}
