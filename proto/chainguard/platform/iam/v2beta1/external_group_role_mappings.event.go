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
	_ events.Eventable  = (*BatchDeleteExternalGroupRoleMappingsResponse)(nil)
	_ events.Extendable = (*BatchDeleteExternalGroupRoleMappingsResponse)(nil)
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

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *BatchDeleteExternalGroupRoleMappingsResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetParent()), true
	case "identityprovider":
		return x.GetParent(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
// The identity provider (the mappings' common parent) is the subject — a single
// resource per the CloudEvents spec intent — while the deleted mapping UIDs ride
// in the Occurrence body. Deriving from the echoed parent keeps the event
// well-formed even when no mappings matched.
func (x *BatchDeleteExternalGroupRoleMappingsResponse) CloudEventsSubject() string {
	return x.GetParent()
}
