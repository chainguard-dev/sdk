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
	_ events.Eventable  = (*Repo)(nil)
	_ events.Extendable = (*Repo)(nil)
	_ events.Eventable  = (*DeleteRepoRequest)(nil)
	_ events.Extendable = (*DeleteRepoRequest)(nil)

	_ events.Eventable  = (*RepoReadme)(nil)
	_ events.Extendable = (*RepoReadme)(nil)
)

func (x *Repo) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *Repo) CloudEventsSubject() string {
	return x.GetUid()
}

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

func (x *RepoReadme) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

func (x *RepoReadme) CloudEventsSubject() string {
	return x.GetUid()
}
