/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import "chainguard.dev/sdk/events"

var (
	_ events.Eventable  = (*AccountAssociation)(nil)
	_ events.Extendable = (*AccountAssociation)(nil)
	_ events.Eventable  = (*DeleteAccountAssociationRequest)(nil)
	_ events.Extendable = (*DeleteAccountAssociationRequest)(nil)
	_ events.Redactable = (*DeleteAccountAssociationRequest)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *AccountAssociation) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		// The UID of an AccountAssociation is the group UIDP it belongs to.
		return x.GetUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *AccountAssociation) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteAccountAssociationRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		// The UID of an AccountAssociation is the group UIDP it belongs to.
		return x.GetUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteAccountAssociationRequest) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
func (x *DeleteAccountAssociationRequest) CloudEventsRedact() any {
	return nil
}
