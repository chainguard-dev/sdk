/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events

// Occurrence is the CloudEvent payload for events.
type Occurrence struct {
	Actor *Actor `json:"actor,omitempty"`

	// Body is the resource that was created.
	Body any `json:"body,omitempty"`
}

// Actor is the event payload form of which identity was responsible for the
// event.
type Actor struct {
	// Subject is the identity that triggered this event.
	Subject string `json:"subject"`

	// Actor contains the name/value pairs for each of the claims that were
	// validated to assume the identity whose UIDP appears in Subject above.
	Actor map[string]string `json:"act,omitempty"`
}

// Eventable allows us to define a set of methods that allow event metadata to
// be collected.
type Eventable interface {
	// CloudEventsSubject returns the subject to use for the cloudevent.
	CloudEventsSubject() string
}

// Extendable allows us to define a generic method to return extensions based on name.
type Extendable interface {
	CloudEventsExtension(key string) (string, bool)
}

type Redactable interface {
	CloudEventsRedact() any
}
