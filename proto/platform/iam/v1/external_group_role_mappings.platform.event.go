/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/uidp"
)

var (
	_ events.Eventable  = (*BatchDeleteExternalGroupRoleMappingsResponse)(nil)
	_ events.Extendable = (*BatchDeleteExternalGroupRoleMappingsResponse)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *ExternalGroupRoleMapping) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetId()), true
	case "identityprovider":
		// Use the authoritative field; the create path enforces it equals the
		// mapping's parent IdP.
		return x.GetIdentityProviderUidp(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *ExternalGroupRoleMapping) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteExternalGroupRoleMappingRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetId()), true
	case "identityprovider":
		// The mapping UIDP is {org}/{idp}/{mapping}; its parent is the IdP UIDP.
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteExternalGroupRoleMappingRequest) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.Redact.
func (x *DeleteExternalGroupRoleMappingRequest) CloudEventsRedact() any {
	return nil
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *BatchDeleteExternalGroupRoleMappingsResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetParentId()), true
	case "identityprovider":
		return x.GetParentId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
// The identity provider (the mappings' common parent) is the subject — a single
// resource per the CloudEvents spec intent — while the deleted mapping IDs ride
// in the Occurrence body. Deriving from the echoed parent keeps the event
// well-formed even when no mappings matched.
func (x *BatchDeleteExternalGroupRoleMappingsResponse) CloudEventsSubject() string {
	return x.GetParentId()
}
