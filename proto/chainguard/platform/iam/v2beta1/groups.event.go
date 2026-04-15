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
	_ events.Eventable  = (*Group)(nil)
	_ events.Extendable = (*Group)(nil)
	_ events.Eventable  = (*DeleteGroupRequest)(nil)
	_ events.Extendable = (*DeleteGroupRequest)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
// A Group's UID is itself a group-level UIDP, so we return it directly
// (unlike Identity/RoleBinding/etc. which call uidp.Parent()).
func (x *Group) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *Group) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteGroupRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteGroupRequest) CloudEventsSubject() string {
	return x.GetUid()
}
