/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"fmt"

	telemetrytest "chainguard.dev/sdk/proto/chainguard/platform/telemetry/v1alpha1/test"
)

// ExampleMockTelemetryClients demonstrates constructing a mock telemetry
// client.
func ExampleMockTelemetryClients() {
	mock := telemetrytest.MockTelemetryClients{}
	fmt.Println(mock.Close())
	// Output:
	// <nil>
}
