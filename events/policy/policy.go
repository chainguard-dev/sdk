/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package policy

import cloudevents "github.com/cloudevents/sdk-go/v2"

// ChangedEventType is the cloudevents event type for validation state change for policy.
const ChangedEventType = "dev.chainguard.policy.validation.changed.v1"

const (
	// NewChange is for new policy state.
	NewChange = "new"
	// DegradedChange says the policy was passing and now is failing.
	DegradedChange = "degraded"
	// ImprovedChange says the policy was failing and now is passing.
	ImprovedChange = "improved"
)

// ImagePolicyRecord is policy states for an image in a cluster.
type ImagePolicyRecord struct {
	// ClusterID identifies the specific cluster the Request pertains to.
	ClusterID string `json:"cluster_id,omitempty"`
	// ImageID that this ExistenceRecord belongs to.
	ImageID string `json:"image_id,omitempty"`
	// LastSeen is the last time we've seen this image_id anywhere on this cluster.
	LastSeen *cloudevents.Timestamp `json:"last_seen,omitempty"`
	// Policies are a map of policy name to policy state that apply to this image.
	Policies map[string]*State `json:"policies,omitempty"`
}

// State is the state of a policy and how it has changed.
type State struct {
	// LastChecked is the time the information was last updated.
	LastChecked *cloudevents.Timestamp `json:"last_checked,omitempty"`
	// Valid is if the image passes the policy.
	Valid bool `json:"valid"`
	// Diagnostic holds any diagnostic messages surfaced during the evaluation
	// of this policy.
	Diagnostic string `json:"diagnostic,omitempty"`
	// Change is the kind of change we have seen for this image between checks.
	// Can be [Empty, "new", "degraded", "improved"]
	Change string `json:"change,omitempty"`
}
