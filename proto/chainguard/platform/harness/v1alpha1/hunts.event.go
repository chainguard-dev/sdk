/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/uidp"
)

var (
	_ events.Eventable  = (*CreateHuntResponse)(nil)
	_ events.Extendable = (*CreateHuntResponse)(nil)
	_ events.Redactable = (*CreateHuntResponse)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *CreateHuntResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetHunt().GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *CreateHuntResponse) CloudEventsSubject() string {
	return x.GetHunt().GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// The event carries only the hunt's id: the upload URL is a bearer credential,
// and the name and digests have no place in a customer audit record.
func (x *CreateHuntResponse) CloudEventsRedact() any {
	return &CreateHuntResponse{
		Hunt:    &Hunt{Id: x.GetHunt().GetId()},
		Created: x.GetCreated(),
	}
}
