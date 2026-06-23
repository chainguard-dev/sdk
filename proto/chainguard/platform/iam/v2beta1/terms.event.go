/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import "chainguard.dev/sdk/events"

var (
	_ events.Eventable  = (*AcceptTermsResponse)(nil)
	_ events.Extendable = (*AcceptTermsResponse)(nil)
)

// CloudEventsSubject implements events.Eventable.
func (x *AcceptTermsResponse) CloudEventsSubject() string {
	return x.GetGroup()
}

// CloudEventsExtension implements events.Extendable.
func (x *AcceptTermsResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetGroup(), true
	default:
		return "", false
	}
}
