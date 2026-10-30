/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"chainguard.dev/sdk/events"
)

var (
	_ events.Eventable  = (*DeleteSessionRequest)(nil)
	_ events.Extendable = (*DeleteSessionRequest)(nil)
	_ events.Eventable  = (*RevokeAllSessionsResponse)(nil)
	_ events.Extendable = (*RevokeAllSessionsResponse)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
// A session belongs to an identity, not a group, so it has no extensions.
func (x *DeleteSessionRequest) CloudEventsExtension(_ string) (string, bool) {
	return "", false
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
// A delete event's body is its request, which names only the session; the
// event's actor is the identity that owned it.
func (x *DeleteSessionRequest) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
// Sessions belong to an identity, not a group, so they have no extensions.
func (x *RevokeAllSessionsResponse) CloudEventsExtension(_ string) (string, bool) {
	return "", false
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *RevokeAllSessionsResponse) CloudEventsSubject() string {
	return x.GetIdentity()
}
