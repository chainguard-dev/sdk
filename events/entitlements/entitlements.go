/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package entitlements defines the public CloudEvent emitted by the Centralized
// Entitlements Service when an entitlement or its bound resources change. It
// follows the same convention as the other public event packages (see
// events/registry, events/apk): a plain payload struct plus the cloudevents
// event-type string, wrapped by events.Occurrence at emit time.
package entitlements

import (
	"encoding/json"

	"chainguard.dev/sdk/events"
)

// ChangedEventType is the cloudevents event type emitted whenever an
// entitlement (or one of its bound resources) is created, updated, or removed.
const ChangedEventType = "dev.chainguard.entitlement.changed.v1"

// The Domain* constants are the canonical wire tokens for the
// ChangeEvent.Domain field.
const (
	// DomainImage is the IMAGE entitlement domain. Resourceful: binds
	// individual resources.
	DomainImage = "IMAGE"
	// DomainLibrary is the LIBRARY entitlement domain.
	DomainLibrary = "LIBRARY"
	// DomainPackage is the PACKAGE entitlement domain.
	DomainPackage = "PACKAGE"
	// DomainHelmChart is the HELM_CHART entitlement domain. Resourceful: binds
	// individual resources.
	DomainHelmChart = "HELM_CHART"
	// DomainFeature is the FEATURE entitlement domain.
	DomainFeature = "FEATURE"
)

// Operation is the kind of change that produced a ChangeEvent. It is a closed
// set corresponding to the four entitlement mutation RPCs. New operations may be
// added in a minor version, so consumers should treat an unrecognized value as a
// signal to re-fetch authoritative state (or log-and-skip) — not to silently
// discard the change.
type Operation string

const (
	// OperationWrite is a create or update of an entitlement (entitlement grain).
	OperationWrite Operation = "WRITE"
	// OperationDelete is the deletion of an entitlement (entitlement grain).
	OperationDelete Operation = "DELETE"
	// OperationAddResource binds one or more resources to an entitlement
	// (resource grain).
	OperationAddResource Operation = "ADD_RESOURCE"
	// OperationRemoveResource unbinds one or more resources from an entitlement
	// (resource grain).
	OperationRemoveResource Operation = "REMOVE_RESOURCE"
)

// ChangeEvent is the events.Occurrence body for ChangedEventType.
//
// Domain, OrgUIDP, EntitlementUIDP, and Operation are populated for every
// operation. Which of the remaining fields are populated depends on Operation;
// switch on Operation and read only the fields listed for that grain:
//
//   - WRITE:           After is set; Before is set iff the write updated an
//     existing entitlement (Before nil ⇒ the write created it). Resources empty.
//   - DELETE:          Before is set; After is nil. Resources empty.
//   - ADD_RESOURCE:    Resources holds the rows actually bound. Before/After nil.
//   - REMOVE_RESOURCE: Resources holds the rows actually removed. Before/After nil.
//
// Before and After hold the changed entitlement as JSON, interpreted per
// Domain (each entitlement domain has its own public entitlement shape). They
// are kept opaque — rather than a single typed message — so one event type can
// serve every domain. The change timestamp lives on the CloudEvent envelope
// (its time attribute), not in the body.
type ChangeEvent struct {
	// Domain is the entitlement domain that changed; one of the Domain*
	// constants (DomainImage, DomainLibrary, DomainPackage, DomainHelmChart,
	// DomainFeature).
	Domain string `json:"domain"`

	// OrgUIDP is the organization (UIDP) the entitlement belongs to.
	OrgUIDP string `json:"org_uidp"`

	// EntitlementUIDP is the UIDP of the entitlement that changed — the stable
	// identity of the subject across every Operation. OrgUIDP + Domain is not
	// unique to an entitlement, so subscribers should identify the entitlement
	// from this field rather than parsing Before/After. It is set for all
	// operations, both entitlement and resource grains.
	EntitlementUIDP string `json:"entitlement_uidp"`

	// Operation is the kind of change; see the grain rules on the type doc.
	Operation Operation `json:"operation"`

	// Source is the system of record for the entitlement (e.g. "SALESFORCE"),
	// when known.
	Source string `json:"source,omitempty"`

	// Actor is the identity the change is attributed to (recorded with the change
	// itself), when known. This is distinct from the emit-time caller carried on
	// the events.Occurrence envelope: for machine-driven syncs (e.g. Salesforce
	// or GitHub-YAML) the envelope carries the service-account credential while
	// this field carries the originator the change is attributed to.
	Actor string `json:"actor,omitempty"`

	// ChangeReason is the free-text reason recorded with the change, when
	// provided.
	ChangeReason string `json:"change_reason,omitempty"`

	// Before is the entitlement prior to the change (the domain's public
	// entitlement shape as JSON), or nil when there was no prior state (creates
	// and resource operations).
	Before json.RawMessage `json:"before,omitempty"`

	// After is the entitlement after the change, or nil for deletes and
	// resource operations.
	After json.RawMessage `json:"after,omitempty"`

	// Resources are the resource rows that changed, set only for ADD_RESOURCE
	// and REMOVE_RESOURCE. It reflects the rows that actually changed (absorbed
	// duplicate adds and already-absent removes are excluded).
	Resources []ResourceChange `json:"resources,omitempty"`
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension.
// It exposes the event's organization (events.GroupKey), entitlement domain
// (events.EntitlementDomainKey), and operation (events.EntitlementOperationKey)
// as CloudEvent context attributes, so subscribers can filter on org, domain,
// and operation without decoding the body.
func (e ChangeEvent) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case events.GroupKey:
		return e.OrgUIDP, true
	case events.EntitlementDomainKey:
		return e.Domain, true
	case events.EntitlementOperationKey:
		return string(e.Operation), true
	default:
		return "", false
	}
}

// ResourceChange identifies a single resource bound to or removed from an
// entitlement. Emitted only for resourceful domains (e.g. IMAGE, HELM_CHART).
// A resource is identified by its Resource (name) and Tier within the
// entitlement, which is itself identified by the event's top-level
// EntitlementUIDP. There is deliberately no per-resource id: no
// customer-meaningful id exists at emit time, so the resource name is the
// identity carried on the wire.
type ResourceChange struct {
	// Resource is the resource's own name — e.g. the repository name for
	// IMAGE. It is the resource's identity within the entitlement, and lets a
	// REMOVE_RESOURCE consumer act on the change without a lookup against a row
	// that may already be gone by the time it processes the event.
	Resource string `json:"resource,omitempty"`

	// Tier is the catalog tier of the resource (e.g. "APPLICATION", "BASE"),
	// when applicable to the domain.
	Tier string `json:"tier,omitempty"`
}
