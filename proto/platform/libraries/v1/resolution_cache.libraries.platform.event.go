/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/events"
)

var (
	_ events.Eventable  = (*ZapResolutionCacheSummary)(nil)
	_ events.Extendable = (*ZapResolutionCacheSummary)(nil)
	_ events.Eventable  = (*ResolutionCacheOptOut)(nil)
	_ events.Extendable = (*ResolutionCacheOptOut)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
// The resolution cache is an org-level resource, so the group extension is the
// org itself. The event is deliberately coordinate-free: it never carries
// package names or versions.
func (x *ZapResolutionCacheSummary) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *ZapResolutionCacheSummary) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
func (x *ResolutionCacheOptOut) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *ResolutionCacheOptOut) CloudEventsSubject() string {
	return x.GetId()
}
