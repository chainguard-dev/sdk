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
	_ events.Eventable  = (*Overlay)(nil)
	_ events.Extendable = (*Overlay)(nil)
	_ events.Eventable  = (*DeleteOverlayRequest)(nil)
	_ events.Extendable = (*DeleteOverlayRequest)(nil)
)

func (x *Overlay) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *Overlay) CloudEventsSubject() string {
	return x.GetUid()
}

func (x *DeleteOverlayRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *DeleteOverlayRequest) CloudEventsSubject() string {
	return x.GetUid()
}
