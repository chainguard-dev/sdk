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
	_ events.Eventable  = (*ExternalGroupRoleMapping)(nil)
	_ events.Extendable = (*ExternalGroupRoleMapping)(nil)
	_ events.Eventable  = (*DeleteExternalGroupRoleMappingRequest)(nil)
	_ events.Extendable = (*DeleteExternalGroupRoleMappingRequest)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *ExternalGroupRoleMapping) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	case "identityprovider":
		return x.GetIdentityProviderUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *ExternalGroupRoleMapping) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteExternalGroupRoleMappingRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	case "identityprovider":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteExternalGroupRoleMappingRequest) CloudEventsSubject() string {
	return x.GetUid()
}
