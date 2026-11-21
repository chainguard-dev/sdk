/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"strings"

	"chainguard.dev/sdk/events"
)

// The DeleteUploadedVersion event's body is its response, which names the
// version as it is stored. The request supplies the group and the ecosystem,
// which the response does not repeat, as the image delete's request does.
var (
	_ events.ResponseBodied = (*DeleteUploadedVersionRequest)(nil)
	_ events.Eventable      = (*deletedUploadedVersion)(nil)
	_ events.Extendable     = (*deletedUploadedVersion)(nil)
	_ events.Suppressible   = (*deletedUploadedVersion)(nil)
)

// CloudEventsBody makes the uploaded_version.deleted.v1 body the
// DeleteUploadedVersionResponse.
func (x *DeleteUploadedVersionRequest) CloudEventsBody(resp any) events.Eventable {
	r, ok := resp.(*DeleteUploadedVersionResponse)
	if !ok || r == nil {
		return nil
	}
	return deletedUploadedVersion{DeleteUploadedVersionResponse: r, req: x}
}

// deletedUploadedVersion is a DeleteUploadedVersionResponse with the group and
// ecosystem of the request that produced it. It marshals as the response.
type deletedUploadedVersion struct {
	*DeleteUploadedVersionResponse
	req *DeleteUploadedVersionRequest
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
//
// The group is parent, the organization that uploaded the version, which the
// method requires to be an organization.
func (d deletedUploadedVersion) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return d.req.GetParent(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
//
// An uploaded version has no uid, so the subject is ecosystem/name@version.
// The name and version are taken from the response. The response holds the
// key as the ecosystem's clients spell it. The subject therefore matches the
// version's other events whatever spelling the caller sent.
//
// The ecosystem comes first because one name and version can exist in more
// than one ecosystem.
func (d deletedUploadedVersion) CloudEventsSubject() string {
	eco := strings.ToLower(strings.TrimPrefix(d.req.GetEcosystem().String(), "STORAGE_USAGE_ECOSYSTEM_"))
	return eco + "/" + d.GetName() + "@" + d.GetVersion()
}

// CloudEventsSuppress skips the event for a version that was already
// deleted, which the call did not change.
func (d deletedUploadedVersion) CloudEventsSuppress() bool {
	return d.GetAlreadyDeleted()
}
