/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

// CloudEventsExtension implements chainguard.dev/sdk/pkg/events/Extendable.CloudEventsExtension
func (x *IdentityMetadata) CloudEventsExtension(key string) (string, bool) { //nolint: revive
	return "", false
}

// CloudEventsSubject implements chainguard.dev/sdk/pkg/events/Eventable.CloudEventsSubject.
func (x *IdentityMetadata) CloudEventsSubject() string {
	return ""
}
