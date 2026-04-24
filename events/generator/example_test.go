/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package generator_test

import (
	"fmt"

	"chainguard.dev/sdk/events/generator"
)

// ExampleGenerateVulnReportEventType demonstrates the event type constant.
func ExampleGenerateVulnReportEventType() {
	fmt.Println(generator.GenerateVulnReportEventType)
	// Output:
	// dev.chainguard.vulnreport.generate.v1
}

// ExampleGenerateVulnReportEvent demonstrates constructing a GenerateVulnReportEvent.
func ExampleGenerateVulnReportEvent() {
	evt := generator.GenerateVulnReportEvent{
		RepoID: "abc123",
		Digest: "sha256:deadbeef",
		Type:   "refresh",
	}
	fmt.Println(evt.RepoID)
	fmt.Println(evt.Type)
	// Output:
	// abc123
	// refresh
}
