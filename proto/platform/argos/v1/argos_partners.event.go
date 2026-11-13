/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "chainguard.dev/sdk/events"

var (
	_ events.Eventable  = (*UpdatePartnerCapabilitiesResponse)(nil)
	_ events.Redactable = (*UpdatePartnerCapabilitiesResponse)(nil)
	_ events.Eventable  = (*DeletePartnerCapabilitiesRequest)(nil)
	_ events.Redactable = (*DeletePartnerCapabilitiesRequest)(nil)
)

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *UpdatePartnerCapabilitiesResponse) CloudEventsSubject() string {
	return x.GetPartnerId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// The event records who changed which partner, never the capability names,
// which say what a partner can mitigate.
func (x *UpdatePartnerCapabilitiesResponse) CloudEventsRedact() any {
	return &UpdatePartnerCapabilitiesResponse{PartnerId: x.GetPartnerId()}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeletePartnerCapabilitiesRequest) CloudEventsSubject() string {
	return x.GetPartnerId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
// The event records who changed which partner, never the capability names or
// interests removed.
func (x *DeletePartnerCapabilitiesRequest) CloudEventsRedact() any {
	return &DeletePartnerCapabilitiesRequest{PartnerId: x.GetPartnerId()}
}
