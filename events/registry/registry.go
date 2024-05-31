/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package registry

import "chainguard.dev/sdk/civil"

const (
	// PulledEventType is the cloudevents event type for registry pulls
	PulledEventType = "dev.chainguard.registry.pull.v1"

	// PushedEventType is the cloudevents event type for registry pushes
	PushedEventType = "dev.chainguard.registry.push.v1"

	// RepoCreatedEventType is the cloudevents event type for registry repo created
	RepoCreatedEventType = "dev.chainguard.api.platform.registry.repo.created.v1"

	// RepoUpdatedEventType is the cloudevents event type for registry repo updated
	RepoUpdatedEventType = "dev.chainguard.api.platform.registry.repo.updated.v1"
)

// PullEvent describes an item being pulled from the registry.
type PullEvent struct {
	// Repository identifies the repository being pulled
	Repository string `json:"repository"`

	// RepoID identifies the UIDP of the repository being pulled
	RepoID string `json:"repo_id"`

	// Tag holds the tag being pulled, if there is one.
	Tag string `json:"tag,omitempty"`

	// Digest holds the digest being pulled.
	// Digest will hold the sha256 of the content being fetched, whether that is
	// a blob or a manifest.
	Digest string `json:"digest"`

	// Method holds the HTTP method of the request.  For pulls, this should be
	// one of HEAD (digest resolution or existence check), or GET to actually
	// fetch the content.
	Method string `json:"method"`

	// Type determines whether the object being fetched is a manifest or blob.
	Type string `json:"type"`

	// When holds when the pull occurred.
	When civil.DateTime `json:"when"`

	// Location holds the detected approximate location of the client who pulled.
	// For example, "ColumbusOHUS" or "Minato City13JP".
	Location string `json:"location"`

	// RemoteAddress holds the address of the client who pulled.
	RemoteAddress string `json:"remote_address"`

	// UserAgent holds the user-agent of the client who pulled.
	UserAgent string `json:"user_agent"`

	Error *Error `json:"error,omitempty"`
}

// PushEvent describes an item being pushed to the registry.
type PushEvent struct {
	// Repository identifies the repository being pushed
	Repository string `json:"repository"`

	// RepoID identifies the UIDP of the repository being pushed
	RepoID string `json:"repo_id"`

	// Tag holds the tag being pushed, if there is one.
	Tag string `json:"tag,omitempty"`

	// Digest holds the digest being pushed.
	// Digest will hold the sha256 of the content being pushed, whether that is
	// a blob or a manifest.
	Digest string `json:"digest"`

	// Type determines whether the object being pushed is a manifest or blob.
	Type string `json:"type"`

	// When holds when the push occurred.
	When civil.DateTime `json:"when"`

	// Location holds the detected approximate location of the client who pushed.
	// For example, "ColumbusOHUS" or "Minato City13JP".
	Location string `json:"location"`

	// RemoteAddress holds the address of the client who pushed.
	RemoteAddress string `json:"remote_address"`

	// UserAgent holds the user-agent of the client who pushed.
	UserAgent string `json:"user_agent"`

	Error *Error `json:"error,omitempty"`
}

type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
