/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/pkg/uidp"
)

// CloudEventsExtension implements chainguard.dev/sdk/pkg/events/Extendable.CloudEventsExtension
func (x *Cluster) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/pkg/events/Eventable.CloudEventsSubject.
func (x *Cluster) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsExtension implements chainguard.dev/sdk/pkg/events/Extendable.CloudEventsExtension
func (x *DeleteClusterRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/pkg/events/Eventable.CloudEventsSubject.
func (x *DeleteClusterRequest) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/pkg/events/Redactable.Redact.
func (x *DeleteClusterRequest) CloudEventsRedact() interface{} {
	return nil
}
