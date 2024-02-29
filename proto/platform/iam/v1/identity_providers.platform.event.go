/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/uidp"
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *IdentityProvider) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *IdentityProvider) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteIdentityProviderRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteIdentityProviderRequest) CloudEventsSubject() string {
	return x.GetId()
}
