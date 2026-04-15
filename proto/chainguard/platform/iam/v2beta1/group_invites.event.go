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
	_ events.Eventable  = (*GroupInvite)(nil)
	_ events.Extendable = (*GroupInvite)(nil)
	_ events.Redactable = (*GroupInvite)(nil)
	_ events.Eventable  = (*DeleteGroupInviteRequest)(nil)
	_ events.Extendable = (*DeleteGroupInviteRequest)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *GroupInvite) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *GroupInvite) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// Strips code and key_id from the event payload.
func (x *GroupInvite) CloudEventsRedact() any {
	return &GroupInvite{
		Uid:            x.Uid,
		ExpirationTime: x.ExpirationTime,
	}
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteGroupInviteRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteGroupInviteRequest) CloudEventsSubject() string {
	return x.GetUid()
}
