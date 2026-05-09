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
	_ events.Eventable  = (*ArgosDocument)(nil)
	_ events.Extendable = (*ArgosDocument)(nil)
	_ events.Eventable  = (*DeleteArgosDocumentRequest)(nil)
	_ events.Extendable = (*DeleteArgosDocumentRequest)(nil)
	_ events.Redactable = (*CreateArgosDocumentRequest)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *ArgosDocument) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *ArgosDocument) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteArgosDocumentRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteArgosDocumentRequest) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// Strips the encrypted upload envelope from the request before any
// event/audit pipeline serializes it: even though it's already
// client-side-encrypted, payloads are multi-MB and have no business
// flowing through cloudevent subscribers. The Create RPC's event body
// is the response (ArgosDocument, payload-free) so the live interceptor
// path is already safe — this method makes the intent explicit and
// guards against any future code path that emits the request.
func (x *CreateArgosDocumentRequest) CloudEventsRedact() any {
	return &CreateArgosDocumentRequest{
		ParentId: x.ParentId,
		Filename: x.Filename,
	}
}
