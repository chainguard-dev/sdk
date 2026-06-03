/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "chainguard.dev/sdk/uidp"

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
// "group" returns the binding's parent UIDP so audit consumers can filter
// events by IAM scope.
func (x *Binding) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *Binding) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
func (x *DeleteBindingRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteBindingRequest) CloudEventsSubject() string {
	return x.GetId()
}
