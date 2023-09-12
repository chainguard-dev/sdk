/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

// CloudEventsExtension implements chainguard.dev/sdk/pkg/events/Extendable.CloudEventsExtension
func (x *Session) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetGroup(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/pkg/events/Eventable.CloudEventsSubject.
func (x *Session) CloudEventsSubject() string {
	return ""
}
