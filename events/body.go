/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events

import (
	"encoding/json"
	"errors"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

// ErrNoBodyType is returned when a CloudEvent has no BodyTypeKey extension.
// Events published before the emitter started setting the extension will hit
// this; callers that need to decode those can use DecodeInto with the proto
// type they expect.
var ErrNoBodyType = errors.New("cloudevent missing bodytype extension")

// DecodeInto populates msg from the Occurrence body embedded in event.Data.
// It handles the envelope unwrap and uses protojson.Unmarshal so that proto
// oneofs and enum string values round-trip correctly.
//
// Use DecodeInto when you already know the expected proto type and do not
// want to rely on the BodyTypeKey extension (e.g., during a rollout window
// where some events are emitted without the extension).
func DecodeInto(event cloudevents.Event, msg proto.Message) error {
	envelope, err := parseEnvelope(event)
	if err != nil {
		return err
	}
	if err := protojson.Unmarshal(envelope.Body, msg); err != nil {
		return fmt.Errorf("unmarshaling occurrence body: %w", err)
	}
	return nil
}

// Body decodes the Occurrence body into a concrete proto.Message resolved
// from the event's BodyTypeKey extension via protoregistry.GlobalTypes.
// Returns ErrNoBodyType if the extension is absent, and a descriptive error
// if the named type is not registered in the calling binary.
func Body(event cloudevents.Event) (proto.Message, error) {
	env, err := parseEnvelope(event)
	if err != nil {
		return nil, err
	}
	return decodeTyped(event, env)
}

// BodyAs is Body with a compile-time type assertion. It returns an error if
// the event's bodytype does not resolve to T — useful for handlers that want
// to reject cross-version deliveries (e.g. a v2beta1-only subscriber that
// receives a v1 event).
func BodyAs[T proto.Message](event cloudevents.Event) (T, error) {
	var zero T
	msg, err := Body(event)
	if err != nil {
		return zero, err
	}
	typed, ok := msg.(T)
	if !ok {
		return zero, fmt.Errorf("event body is %T, wanted %T", msg, zero)
	}
	return typed, nil
}

// Decode returns the full Occurrence (Actor + Body) from event.Data. The
// returned Occurrence.Body is a concrete proto.Message resolved the same way
// as Body.
func Decode(event cloudevents.Event) (*Occurrence, error) {
	env, err := parseEnvelope(event)
	if err != nil {
		return nil, err
	}
	msg, err := decodeTyped(event, env)
	if err != nil {
		return nil, err
	}
	return &Occurrence{Actor: env.Actor, Body: msg}, nil
}

func decodeTyped(event cloudevents.Event, env *envelope) (proto.Message, error) {
	name, err := bodyTypeName(event)
	if err != nil {
		return nil, err
	}
	mt, err := protoregistry.GlobalTypes.FindMessageByName(protoreflect.FullName(name))
	if err != nil {
		return nil, fmt.Errorf("resolving bodytype %q: %w", name, err)
	}
	msg := mt.New().Interface()
	if err := protojson.Unmarshal(env.Body, msg); err != nil {
		return nil, fmt.Errorf("unmarshaling occurrence body: %w", err)
	}
	return msg, nil
}

func bodyTypeName(event cloudevents.Event) (string, error) {
	raw, ok := event.Extensions()[BodyTypeKey]
	if !ok {
		return "", ErrNoBodyType
	}
	name, ok := raw.(string)
	if !ok || name == "" {
		return "", fmt.Errorf("%w: extension value is %T", ErrNoBodyType, raw)
	}
	return name, nil
}

type envelope struct {
	Actor *Actor          `json:"actor,omitempty"`
	Body  json.RawMessage `json:"body"`
}

func parseEnvelope(event cloudevents.Event) (*envelope, error) {
	var env envelope
	if err := json.Unmarshal(event.Data(), &env); err != nil {
		return nil, fmt.Errorf("parsing occurrence envelope: %w", err)
	}
	if len(env.Body) == 0 {
		return nil, errors.New("occurrence envelope has empty body")
	}
	return &env, nil
}
