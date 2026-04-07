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
