/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"google.golang.org/protobuf/types/known/durationpb"

	"chainguard.dev/sdk/events"
)

var (
	_ events.Eventable      = (*AuthPolicy)(nil)
	_ events.Extendable     = (*AuthPolicy)(nil)
	_ events.AsyncEventable = (*AuthPolicy)(nil)
	_ events.Redactable     = (*AuthPolicy)(nil)

	_ events.Eventable  = (*Validation)(nil)
	_ events.Extendable = (*Validation)(nil)
	_ events.Redactable = (*Validation)(nil)
)

// CloudEventsAsync implements chainguard.dev/sdk/events/AsyncEventable: the
// auth_policies.updated audit event is delivered off the request path. An
// sso_only toggle is a rare administrative change whose response must not
// wait on the event sink; delivery stays best-effort, counted and logged on
// failure by the publishing interceptor.
func (x *AuthPolicy) CloudEventsAsync() bool { return true }

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *AuthPolicy) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		// The UID of an AuthPolicy is the UIDP of the root organization it governs.
		return x.GetUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *AuthPolicy) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
//
// The event body keeps the policy's address and its writable fields, sso_only
// and session_duration, and nothing else: the updated event, and the previous
// state the publishing interceptor runs through this same redaction, are the
// audit record of what an update changed. Events reach every subscriber an
// organization editor can register, while the rest of the policy (the
// validation view of which designees have a current break-glass validation,
// the readiness verdict and its reasons, the entitlement) is shown only to
// callers who may enable sso_only, so the body must not carry it. A
// subscriber that needs it reads the policy with that capability. The
// response itself is left as is: the redaction copies.
func (x *AuthPolicy) CloudEventsRedact() any {
	redacted := &AuthPolicy{Uid: x.GetUid()}
	if x.SsoOnly != nil {
		// A fresh pointer: the body must not share storage with the response,
		// and an unset sso_only must stay unset rather than become false.
		v := *x.SsoOnly
		redacted.SsoOnly = &v
	}
	if x.SessionDuration != nil {
		// Likewise a copy, and unset ("Chainguard's default") stays unset
		// rather than becoming the explicit zero only Chainguard sets.
		redacted.SessionDuration = &durationpb.Duration{
			Seconds: x.SessionDuration.GetSeconds(),
			Nanos:   x.SessionDuration.GetNanos(),
		}
	}
	return redacted
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *Validation) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		// The validation is about the organization's policy; its root UIDP is
		// the group the event scopes to.
		return x.GetUid(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *Validation) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
//
// The warehoused event records that a validation of a kind was started for an
// organization, and nothing a holder of the event could act on: not the
// open_url, and not the session id it carries, which is the only credential
// the sign-in needs. Nor the designated identity of a break-glass session:
// the designations are owner-only information on the policy, and the event
// reaches every subscriber of the organization. diagnostics are never
// present on a just-started session.
func (x *Validation) CloudEventsRedact() any {
	return &Validation{
		Uid:   x.GetUid(),
		Kind:  x.GetKind(),
		State: x.GetState(),
	}
}
