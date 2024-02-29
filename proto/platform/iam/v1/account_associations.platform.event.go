/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *AccountAssociations) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetGroup(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *AccountAssociations) CloudEventsSubject() string {
	return x.Group
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteAccountAssociationsRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetGroup(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteAccountAssociationsRequest) CloudEventsSubject() string {
	return x.Group
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.Redact.
func (x *DeleteAccountAssociationsRequest) CloudEventsRedact() interface{} {
	return nil
}
