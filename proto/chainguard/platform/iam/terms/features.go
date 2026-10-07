/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms

import "slices"

// featureDocuments are the documents an organization must accept to use each
// entitlement feature, keyed by feature key, however the organization comes to
// be entitled: a trial, a staff grant, or a sale. Every feature has an entry.
var featureDocuments = map[string][]string{
	// Checks runs customer code under the Code Analysis Terms, and is a
	// Technology Preview as the MSLA defines one.
	"sandbox_checks": {"code-analysis-terms.v1", "msla.v1"},
	// Workspaces has its own terms, not yet published. They must be before
	// Workspaces is granted to customers outside a trial.
	"sandbox_workspaces": nil,
}

// FeatureDocuments returns the IDs of the documents an organization must
// accept to use the entitlement feature with key feature ("sandbox_checks",
// "sandbox_workspaces"), or nil when it requires none.
func FeatureDocuments(feature string) []string {
	return slices.Clone(featureDocuments[feature])
}
