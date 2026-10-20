/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import "chainguard.dev/sdk/events"

var (
	_ events.Eventable  = (*AuthPolicy)(nil)
	_ events.Extendable = (*AuthPolicy)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *AuthPolicy) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		// The UID of an AuthPolicy is the UIDP of the root organization it governs.
		return x.GetUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *AuthPolicy) CloudEventsSubject() string {
	return x.GetUid()
}
