/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/uidp"
	"google.golang.org/protobuf/proto"
)

var (
	_ events.Eventable  = (*Identity)(nil)
	_ events.Extendable = (*Identity)(nil)
	_ events.Redactable = (*Identity)(nil)
	_ events.Eventable  = (*DeleteIdentityRequest)(nil)
	_ events.Extendable = (*DeleteIdentityRequest)(nil)
	_ events.Eventable  = (*IdentityMetadata)(nil)
	_ events.Extendable = (*IdentityMetadata)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *Identity) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *Identity) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// created_by_email is returned only by Get and List, but clearing it here
// keeps it out of the created and updated events, which customer
// subscriptions receive and recorders keep, even if a response carries it.
func (x *Identity) CloudEventsRedact() any {
	if x == nil {
		return x
	}
	c, _ := proto.Clone(x).(*Identity)
	c.CreatedByEmail = ""
	return c
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteIdentityRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteIdentityRequest) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *IdentityMetadata) CloudEventsExtension(_ string) (string, bool) {
	return "", false
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *IdentityMetadata) CloudEventsSubject() string {
	return x.GetUid()
}
