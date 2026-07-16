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
	_ events.Eventable  = (*Tag)(nil)
	_ events.Extendable = (*Tag)(nil)
	_ events.Eventable  = (*DeleteTagRequest)(nil)
	_ events.Extendable = (*DeleteTagRequest)(nil)
)

func (x *Tag) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *Tag) CloudEventsSubject() string {
	return x.GetUid()
}

func (x *DeleteTagRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *DeleteTagRequest) CloudEventsSubject() string {
	return x.GetUid()
}
