/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "chainguard.dev/sdk/uidp"

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *ExternalGroupRoleMapping) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetId()), true
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
