/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"chainguard.dev/sdk/events"
)

var (
	_ events.Eventable  = (*Entitlement)(nil)
	_ events.Extendable = (*Entitlement)(nil)
	_ events.Eventable  = (*DeleteEntitlementRequest)(nil)
	_ events.Extendable = (*DeleteEntitlementRequest)(nil)
)

// CloudEventsExtension implements events.Extendable. The entitlement's id is
// the entitled group's UIDP, so the "group" extension returns it directly.
func (x *Entitlement) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements events.Eventable.
func (x *Entitlement) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements events.Extendable. The id is the group UIDP.
func (x *DeleteEntitlementRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements events.Eventable.
func (x *DeleteEntitlementRequest) CloudEventsSubject() string {
	return x.GetId()
}
