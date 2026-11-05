/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/uidp"
)

var (
	_ events.Eventable      = (*DeleteImageRequest)(nil)
	_ events.Extendable     = (*DeleteImageRequest)(nil)
	_ events.Suppressible   = (*DeleteImageRequest)(nil)
	_ events.ResponseBodied = (*DeleteImageRequest)(nil)
	_ events.Extendable     = (*deletedImage)(nil)
	_ events.Suppressible   = (*deletedImage)(nil)
)

func (x *DeleteImageRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetRepoUid()), true
	default:
		return "", false
	}
}

func (x *DeleteImageRequest) CloudEventsSubject() string {
	return x.GetRepoUid()
}

// CloudEventsSuppress skips the deletion event for a validate-only preview,
// which deletes nothing.
func (x *DeleteImageRequest) CloudEventsSuppress() bool {
	return x.GetValidateOnly()
}

// CloudEventsBody makes the image.deleted.v1 body the DeleteImageResponse: a
// deletion selected by tag names its digest, the other tags and the manifests
// it removed only there.
func (x *DeleteImageRequest) CloudEventsBody(resp any) events.Eventable {
	r, ok := resp.(*DeleteImageResponse)
	if !ok || r == nil {
		return nil
	}
	return deletedImage{DeleteImageResponse: r, req: x}
}

// deletedImage is a DeleteImageResponse with the subject, group and
// suppression of the request that produced it. It marshals as the response.
type deletedImage struct {
	*DeleteImageResponse
	req *DeleteImageRequest
}

func (d deletedImage) CloudEventsSubject() string { return d.req.CloudEventsSubject() }

func (d deletedImage) CloudEventsExtension(key string) (string, bool) {
	return d.req.CloudEventsExtension(key)
}

func (d deletedImage) CloudEventsSuppress() bool { return d.req.CloudEventsSuppress() }
