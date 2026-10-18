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
	_ events.Eventable  = (*CreateUploadResponse)(nil)
	_ events.Extendable = (*CreateUploadResponse)(nil)
	_ events.Redactable = (*CreateUploadResponse)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *CreateUploadResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUpload().GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *CreateUploadResponse) CloudEventsSubject() string {
	return x.GetUpload().GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// The event carries only the upload's id: the upload URL is a bearer
// credential, and the digests have no place in a customer audit record.
func (x *CreateUploadResponse) CloudEventsRedact() any {
	return &CreateUploadResponse{
		Upload: &Upload{Id: x.GetUpload().GetId()},
	}
}
