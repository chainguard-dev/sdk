/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package generator

import (
	"chainguard.dev/sdk/civil"
	"google.golang.org/grpc/status"
)

const GenerateVulnReportEventType = "dev.chainguard.vulnreport.generate.v1"

// GenerateVulnReportEvent is an event used to trigger the generation of scan reports on-demand
type GenerateVulnReportEvent struct {
	// RepoID identifies the UIDP of the repository used to generate the report
	RepoID string `json:"repo_id"`

	// Digest holds the digest being used to generate the report.
	// Digest will hold the sha256 content.
	Digest string `json:"digest"`

	// Type determines whether the report being generated is a refresh.
	Type string `json:"type"`

	// When holds when the generation of the report occurred.
	When civil.DateTime `json:"when"`

	// Status to represent the problem
	Status *status.Status `json:"status,omitempty"`
}
