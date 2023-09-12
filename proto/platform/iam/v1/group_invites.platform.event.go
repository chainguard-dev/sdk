/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/pkg/uidp"
)

// CloudEventsExtension implements chainguard.dev/sdk/pkg/events/Extendable.CloudEventsExtension
func (x *GroupInvite) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/pkg/events/Eventable.CloudEventsSubject.
func (x *GroupInvite) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/pkg/events/Redactable.CloudEventsRedact.
func (x *GroupInvite) CloudEventsRedact() interface{} {
	return &GroupInvite{
		Id:         x.Id,
		Expiration: x.Expiration,
	}
}

// CloudEventsExtension implements chainguard.dev/sdk/pkg/events/Extendable.CloudEventsExtension
func (x *DeleteGroupInviteRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/pkg/events/Eventable.CloudEventsSubject.
func (x *DeleteGroupInviteRequest) CloudEventsSubject() string {
	return x.GetId()
}
