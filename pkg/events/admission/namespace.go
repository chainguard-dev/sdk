/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package admission

import (
	"chainguard.dev/sdk/pkg/uidp"
)

// NamespaceEventType is the cloudevents event type when a namespace's scope
// changes to include or exclude a particular admission controller.
const NamespaceEventType = "dev.chainguard.admission.namespace.v1"

// Change is an enumeration of the types of changes we will emit events for.
type Change string

const (
	// CreatedChange is emitted the first time we see a resource
	CreatedChange Change = "created"

	// UpdatedChange is emitted when we see a resource change.
	UpdatedChange Change = "updated"
)

// EnforcerState is an enumeration of the possible states that the enforcer may
// be in.
type EnforcerState string

const (
	// EnabledEnforcerState is emitted when enforcement is enabled.
	EnabledEnforcerState EnforcerState = "enabled"

	// DisabledEnforcerState is emitted when enforcement is disabled.
	DisabledEnforcerState EnforcerState = "disabled"
)

// NamespaceBody is the body of the Chainguard event Occurrence when the event
// type is NamespaceEventType.
type NamespaceBody struct {
	// Name is the name of the namespace as it appears within the user's cluster
	// e.g. kube-system
	Name string `json:"name"`

	// ID is the UIDP of the Namespace (whose parent is the Cluster UIDP)
	ID uidp.UIDP `json:"id"`

	// Change holds the type of change to the namespace we have observed.
	Change Change `json:"change"`

	// EnforcerState holds the state that policy enforcement is in for a
	// particular namespace.
	EnforcerState EnforcerState `json:"enforcer_state"`
}
