/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/uidp"
)

var (
	_ events.Eventable  = (*RoleBinding)(nil)
	_ events.Extendable = (*RoleBinding)(nil)
	_ events.Eventable  = (*DeleteRoleBindingRequest)(nil)
	_ events.Extendable = (*DeleteRoleBindingRequest)(nil)
	_ events.Eventable  = (*BatchCreateRoleBindingsResponse)(nil)
	_ events.Extendable = (*BatchCreateRoleBindingsResponse)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *RoleBinding) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *RoleBinding) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteRoleBindingRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteRoleBindingRequest) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
// For batch operations, the "group" extension is derived from the first binding's
// parent UIDP. All bindings in a batch share the same parent group (enforced by
// core.BatchCreate validation), so any binding's parent is representative.
func (x *BatchCreateRoleBindingsResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		if len(x.GetRoleBindings()) > 0 {
			return uidp.Parent(x.GetRoleBindings()[0].GetUid()), true
		}
		return "", false
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
//
// The subject is the parent group UIDP, not a binding UID. Subscribers needing
// per-binding granularity should read the RoleBindings field of the event body.
func (x *BatchCreateRoleBindingsResponse) CloudEventsSubject() string {
	if len(x.GetRoleBindings()) > 0 {
		return uidp.Parent(x.GetRoleBindings()[0].GetUid())
	}
	return ""
}
