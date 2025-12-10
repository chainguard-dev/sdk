/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package apk

import (
	"fmt"

	"chainguard.dev/sdk/civil"
)

const (
	// PushedEventType is the cloudevents event type for APK pushes
	PushedEventType = "dev.chainguard.apk.push.v1"

	// PulledEventType is the cloudevents event type for APK pulls
	PulledEventType = "dev.chainguard.apk.pull.v1"
)

// PushEvent describes an APK being pushed to the registry.
type PushEvent struct {
	// Repository identifies the repository being pushed
	Repository string `json:"repository"`

	// RepoID identifies the UIDP of the APK repository (group) being pushed
	RepoID string `json:"repo_id"`

	// Package holds the name of the package being pushed.
	Package string `json:"package"`

	// Origin holds the name of the origin package being pushed. For
	// subpackages, this is the name of the parent package, for main packages,
	// this is the same as Package.
	Origin string `json:"origin_package"`

	// Version holds the version of the package being pushed.
	Version string `json:"version"`

	// Architecture holds the architecture of the package being pushed.
	Architecture string `json:"architecture"`

	// Checksum holds the checksum of the package's control section as
	// it would appear in an APKINDEX entry for the package.
	Checksum string `json:"checksum"`

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

// APKPath is a convenience method for constructing the full path to the APK.
func (e PushEvent) APKPath() string {
	return fmt.Sprintf("%s/%s", e.RepoID, e.APKBasePath())
}

// APKBasePath is a convenience method for constructing the base path to the APK.
func (e PushEvent) APKBasePath() string {
	return fmt.Sprintf("%s/%s-%s.apk", e.Architecture, e.Package, e.Version)
}

type Error struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PullEvent describes an APK being pulled from the registry.
type PullEvent struct {
	// Repository identifies the repository being pulled
	Repository string `json:"repository"`

	// RepoID identifies the UIDP of the APK repository (group) being pulled
	RepoID string `json:"repo_id"`

	// Package holds the name of the package being pushed.
	Package string `json:"package"`

	// Origin holds the name of the origin package being pushed. For
	// subpackages, this is the name of the parent package, for main packages,
	// this is the same as Package.
	Origin string `json:"origin_package"`

	// Version holds the version of the package being pushed.
	Version string `json:"version"`

	// Architecture holds the architecture of the package being pushed.
	Architecture string `json:"architecture"`

	// Checksum holds the checksum of the package's control section as
	// it would appear in an APKINDEX entry for the package.
	Checksum string `json:"checksum"`

	// When holds when the pull occurred.
	When civil.DateTime `json:"when"`

	// Location holds the detected approximate location of the client who pulled.
	// For example, "ColumbusOHUS" or "Minato City13JP".
	Location string `json:"location"`

	// RemoteAddress holds the address of the client who pulled.
	RemoteAddress string `json:"remote_address"`

	// UserAgent holds the user-agent of the client who pulled.
	UserAgent string `json:"user_agent"`

	// ProxyUIDP is the UIDP of the customer associated with the apkproxy pull, if any.
	// This is only set if the pull was proxied through virtualapk.cgr.dev/<uidp>(/<hash>?)/original/<arch>/<pkg>.apk.
	ProxyUIDP string `json:"proxy_uidp,omitempty"`

	// ProxyHash is the hash of the image associated with the apkproxy pull, if any.
	// This is only set if the pull was proxied through virtualapk.cgr.dev/<uidp>/<hash>/original/<arch>/<pkg>.apk.
	ProxyHash string `json:"proxy_hash,omitempty"`
}

// APKPath is a convenience method for constructing the full path to the APK.
func (e PullEvent) APKPath() string {
	return fmt.Sprintf("%s/%s", e.RepoID, e.APKBasePath())
}

// APKBasePath is a convenience method for constructing the base path to the APK.
func (e PullEvent) APKBasePath() string {
	return fmt.Sprintf("%s/%s-%s.apk", e.Architecture, e.Package, e.Version)
}
