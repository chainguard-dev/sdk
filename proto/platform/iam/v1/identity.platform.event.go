/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/uidp"
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *Identity) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *Identity) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteIdentityRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteIdentityRequest) CloudEventsSubject() string {
	return x.GetId()
}
