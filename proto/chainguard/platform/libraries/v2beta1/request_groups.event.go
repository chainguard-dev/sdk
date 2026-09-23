/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import "chainguard.dev/sdk/events"

// Every mutating RPC in RequestGroupsService carries a
// (chainguard.annotations.events) annotation, so each type those RPCs emit has to
// be able to describe itself as a CloudEvent. Create, Submit, Update and
// RefreshCoverage all return RequestGroup; RemoveItems and RestoreItems return
// their own response messages; both deletes return Empty, so their request types
// carry the event instead.
var (
	_ events.Eventable  = (*RequestGroup)(nil)
	_ events.Extendable = (*RequestGroup)(nil)
	_ events.Eventable  = (*RemoveItemsResponse)(nil)
	_ events.Extendable = (*RemoveItemsResponse)(nil)
	_ events.Eventable  = (*RestoreItemsResponse)(nil)
	_ events.Extendable = (*RestoreItemsResponse)(nil)
	_ events.Eventable  = (*DeleteRequestGroupRequest)(nil)
	_ events.Extendable = (*DeleteRequestGroupRequest)(nil)
	_ events.Eventable  = (*DeleteSubmittedRequestGroupRequest)(nil)
	_ events.Extendable = (*DeleteSubmittedRequestGroupRequest)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
//
// The group extension is read from parent_id rather than derived from the uid.
// RequestGroup already carries the owning organization as a field, so using it is
// both exact and self-describing; deriving it would restate an assumption about
// the uid's shape that the message does not need to make.
func (x *RequestGroup) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetParentId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *RequestGroup) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
func (x *RemoveItemsResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetParentId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
//
// The subject is the request group, not the individual items. A subscriber reads
// subject as one resource identifier, so a list of item uids would both break that
// expectation and risk the attribute size limit on a bulk edit. Which items moved
// is in the event body.
func (x *RemoveItemsResponse) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
func (x *RestoreItemsResponse) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetParentId(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
// As with RemoveItemsResponse, the group is the subject rather than the items.
func (x *RestoreItemsResponse) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
//
// DeleteRequestGroup returns Empty, so the event is built from the request. The
// organization is read from parent, the field that carries the IAM scope, rather
// than derived from the uid: uid is a flat UUID here, not a Chainguard ID, so it
// has no parent to take. DeleteGroupRequest and DeleteRoleBindingRequest derive
// theirs because their uid is a UIDP.
func (x *DeleteRequestGroupRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetParent(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteRequestGroupRequest) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
//
// Identical to the customer delete's, because both name the same group in the
// same two fields, and both emit the same event type.
func (x *DeleteSubmittedRequestGroupRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return x.GetParent(), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteSubmittedRequestGroupRequest) CloudEventsSubject() string {
	return x.GetUid()
}
