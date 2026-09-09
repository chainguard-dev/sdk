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
	_ events.Eventable  = (*OverlayBinding)(nil)
	_ events.Extendable = (*OverlayBinding)(nil)
	_ events.Eventable  = (*DeleteOverlayBindingRequest)(nil)
	_ events.Extendable = (*DeleteOverlayBindingRequest)(nil)
)

func (x *OverlayBinding) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *OverlayBinding) CloudEventsSubject() string {
	return x.GetUid()
}

func (x *DeleteOverlayBindingRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Root(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *DeleteOverlayBindingRequest) CloudEventsSubject() string {
	return x.GetUid()
}
