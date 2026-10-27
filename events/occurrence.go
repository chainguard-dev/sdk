/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events

import "time"

// Occurrence is the CloudEvent payload for events.
type Occurrence struct {
	Actor *Actor `json:"actor,omitempty"`

	// Body is the resource that was created.
	Body any `json:"body,omitempty"`

	// Change describes the change an update event records, beyond the state
	// Body shows: the resource before it and the client that asked for it.
	// Present only on events whose producer recorded it; absent, the event
	// says nothing about the previous state.
	Change *Change `json:"change,omitempty"`
}

// Change is the audit detail of a mutation event: what the resource looked
// like before, when the change was made and through which kind of client.
// Body remains the resource after the change.
type Change struct {
	// Previous is the resource before the change, or the part of it the
	// producer tracks, in the same encoding as Body (protojson for proto
	// resources). Absent when the change created the resource.
	Previous any `json:"previous,omitempty"`

	// At is when the change was made, as the producer's clock had it.
	At time.Time `json:"at"`

	// Client is the kind of client that requested the change, classified from
	// its user agent: "chainctl", "console", "terraform", "http" (another
	// HTTP client, such as curl or a browser outside the Console), "grpc"
	// (another gRPC client) or "unknown".
	Client string `json:"client,omitempty"`

	// ClientVersion is the client's version when its user agent carries one,
	// else "unknown".
	ClientVersion string `json:"client_version,omitempty"`
}

// Actor is the event payload form of which identity was responsible for the
// event.
type Actor struct {
	// Subject is the identity that triggered this event.
	Subject string `json:"subject"`

	// Actor contains the name/value pairs for each of the claims that were
	// validated to assume the identity whose UIDP appears in Subject above.
	Actor map[string]string `json:"act,omitempty"`

	// Chain is the ordered lineage of credential exchanges — oldest first —
	// that led to the identity in Subject, mirrored from the token's
	// act_chain claim. Actor above records only the immediate input
	// credential and is replaced on every exchange; Chain accumulates, so
	// the originating principal survives chained identity assumption. It is
	// attribution, never authorization. Empty on events whose token predates
	// the chain claim or carries none.
	Chain []ActorHop `json:"chain,omitempty"`
}

// ActorHop is one credential exchange in an actor chain. Field names and
// values mirror the hops of the token's act_chain claim.
type ActorHop struct {
	// Sub is the subject of the hop's credential — a Chainguard identity
	// UIDP when the credential was minted by the Chainguard issuer, the
	// third-party subject for an origin hop, the attributed principal for
	// asserted hops.
	Sub string `json:"sub"`
	// Kind records how the exchange this credential authorized sourced the
	// minted token's authority: "bindings" (resolved from the assumed
	// identity's role bindings) or "transferred" (a per-scope subset
	// transferred from this hop's credential).
	Kind string `json:"kind"`
	// Provenance records how this hop's subject was established: "issuer"
	// (built from the issuer's own verification of the hop's credential) or
	// "asserted" (attributed by a trusted intermediary that verified
	// something other than a Chainguard token).
	Provenance string `json:"provenance"`
	// Iss is the issuer of the hop's credential; empty for asserted hops,
	// which have no credential.
	Iss string `json:"iss,omitempty"`
	// UpstreamSub optionally names the upstream principal an asserted hop
	// was attributed from.
	UpstreamSub string `json:"upstream_sub,omitempty"`
}

// Eventable allows us to define a set of methods that allow event metadata to
// be collected.
type Eventable interface {
	// CloudEventsSubject returns the subject to use for the cloudevent.
	CloudEventsSubject() string
}

// AsyncEventable opts an Eventable response into asynchronous delivery after
// the event has been fully constructed and redacted. Use this only when waiting
// for the event sink would make delivery of the primary response materially
// less safe, such as a reveal-once credential response, or for a rare
// administrative change whose audit event must not hold its response (an
// organization's authentication policy). Delivery is best-effort either way:
// the publisher counts and logs a failed delivery, it never fails the response.
type AsyncEventable interface {
	CloudEventsAsync() bool
}

// Extendable allows us to define a generic method to return extensions based on name.
type Extendable interface {
	CloudEventsExtension(key string) (string, bool)
}

type Redactable interface {
	CloudEventsRedact() any
}
