/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"strings"

	"chainguard.dev/sdk/uidp"
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *RoleBinding) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *RoleBinding) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *RoleBindingBatch) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		// All role bindings created in a batch have the same parent.
		if len(x.GetRoleBindings()) > 0 {
			return uidp.Parent(x.GetRoleBindings()[0].GetId()), true
		}
		// We shouldn't get here since a successful create should include at least one binding.
		return "", false
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *RoleBindingBatch) CloudEventsSubject() string {
	var sb strings.Builder
	for _, rb := range x.GetRoleBindings() {
		sb.WriteString(rb.GetId() + ",")
	}
	return strings.TrimSuffix(sb.String(), ",")
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteRoleBindingRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteRoleBindingRequest) CloudEventsSubject() string {
	return x.GetId()
}
