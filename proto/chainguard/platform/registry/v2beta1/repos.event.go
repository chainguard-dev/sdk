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
	_ events.Eventable  = (*DeleteRepoRequest)(nil)
	_ events.Extendable = (*DeleteRepoRequest)(nil)
)

func (x *DeleteRepoRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *DeleteRepoRequest) CloudEventsSubject() string {
	return x.GetUid()
}
