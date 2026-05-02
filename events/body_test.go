/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events_test

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"

	"chainguard.dev/sdk/events"
	iamv2 "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	iamv1 "chainguard.dev/sdk/proto/platform/iam/v1"
)

// buildEvent encodes an Occurrence onto a CloudEvent the same way a
// publisher does: the body is protojson-marshaled and stashed as a
// json.RawMessage, then the whole envelope is encoded with encoding/json.
func buildEvent(t *testing.T, body proto.Message, actor *events.Actor, setBodyType bool) cloudevents.Event {
	t.Helper()
	occ := &events.Occurrence{Actor: actor}
	if body != nil {
		b, err := protojson.Marshal(body)
		if err != nil {
			t.Fatalf("protojson.Marshal: %v", err)
		}
		occ.Body = json.RawMessage(b)
	}
	event := cloudevents.NewEvent()
	event.SetType("dev.chainguard.api.iam.group.created.v1")
	event.SetSource("https://api.test")
	event.SetSubject("abc123")
	if setBodyType && body != nil {
		event.SetExtension(events.BodyTypeKey, string(body.ProtoReflect().Descriptor().FullName()))
	}
	if err := event.SetData(cloudevents.ApplicationJSON, occ); err != nil {
		t.Fatalf("event.SetData: %v", err)
	}
	return event
}

// eventWithBodyType builds an event with an arbitrary bodytype string,
// useful for exercising error paths that don't require a real proto type.
func eventWithBodyType(t *testing.T, bodyType string) cloudevents.Event {
	t.Helper()
	event := cloudevents.NewEvent()
	event.SetType("dev.chainguard.api.iam.group.created.v1")
	event.SetSource("https://api.test")
	event.SetSubject("abc123")
	event.SetExtension(events.BodyTypeKey, bodyType)
	if err := event.SetData(cloudevents.ApplicationJSON, &events.Occurrence{Body: json.RawMessage(`{}`)}); err != nil {
		t.Fatalf("event.SetData: %v", err)
	}
	return event
}

func TestBody(t *testing.T) {
	for _, tt := range []struct {
		name    string
		event   cloudevents.Event
		want    proto.Message
		wantErr string
	}{{
		name:  "v1 group",
		event: buildEvent(t, &iamv1.Group{Id: "abc123", Name: "v1-group", Verified: true}, nil, true),
		want:  &iamv1.Group{Id: "abc123", Name: "v1-group", Verified: true},
	}, {
		name:  "v2beta1 group",
		event: buildEvent(t, &iamv2.Group{Uid: "abc123", Name: "v2-group", Verified: true}, nil, true),
		want:  &iamv2.Group{Uid: "abc123", Name: "v2-group", Verified: true},
	}, {
		name:    "missing bodytype extension",
		event:   buildEvent(t, &iamv2.Group{Uid: "abc123"}, nil, false),
		wantErr: "missing bodytype",
	}, {
		name:    "unknown bodytype",
		event:   eventWithBodyType(t, "chainguard.platform.nonexistent.Thing"),
		wantErr: "resolving bodytype",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := events.Body(tt.event)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("err: got = %v, wanted containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Body: %v", err)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestBody_ErrNoBodyType verifies the sentinel is reachable via errors.Is,
// in addition to the substring check already covered in TestBody.
func TestBody_ErrNoBodyType(t *testing.T) {
	event := buildEvent(t, &iamv2.Group{Uid: "abc123"}, nil, false)
	_, err := events.Body(event)
	if !errors.Is(err, events.ErrNoBodyType) {
		t.Errorf("err: got = %v, wanted errors.Is(ErrNoBodyType)", err)
	}
}

func TestBodyAs(t *testing.T) {
	for _, tt := range []struct {
		name    string
		event   cloudevents.Event
		want    *iamv2.Group
		wantErr string
	}{{
		name:  "success",
		event: buildEvent(t, &iamv2.Group{Uid: "abc123", Name: "v2-group"}, nil, true),
		want:  &iamv2.Group{Uid: "abc123", Name: "v2-group"},
	}, {
		name:    "mismatched type",
		event:   buildEvent(t, &iamv1.Group{Id: "abc123"}, nil, true),
		wantErr: "event body is ",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := events.BodyAs[*iamv2.Group](tt.event)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("err: got = %v, wanted containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("BodyAs: %v", err)
			}
			if diff := cmp.Diff(tt.want, got, protocmp.Transform()); diff != "" {
				t.Errorf("body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDecodeInto(t *testing.T) {
	// DecodeInto does not consult the bodytype extension, so these cases
	// exercise both the oneof/enum round-trip behavior and the rollout
	// path where events arrive without bodytype set.
	for _, tt := range []struct {
		name string
		body proto.Message
	}{{
		name: "oneof and enum round-trip",
		body: &iamv2.IdentityProvider{
			Uid:  "abc123/def456",
			Name: "test-idp",
			Configuration: &iamv2.IdentityProvider_Oidc{
				Oidc: &iamv2.IdentityProvider_OIDC{
					Issuer:   "https://accounts.google.com",
					ClientId: "client-id",
				},
			},
		},
	}, {
		name: "no bodytype extension",
		body: &iamv2.Group{Uid: "abc123", Name: "v2-group"},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			event := buildEvent(t, tt.body, nil, false)
			got := tt.body.ProtoReflect().New().Interface()
			if err := events.DecodeInto(event, got); err != nil {
				t.Fatalf("DecodeInto: %v", err)
			}
			if diff := cmp.Diff(tt.body, got, protocmp.Transform()); diff != "" {
				t.Errorf("body mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDecode_PreservesActor(t *testing.T) {
	wantActor := &events.Actor{
		Subject: "actor-subject",
		Actor:   map[string]string{"iss": "https://issuer.test", "sub": "user-123"},
	}
	wantBody := &iamv2.Group{Uid: "abc123", Name: "v2-group"}
	event := buildEvent(t, wantBody, wantActor, true)

	occ, err := events.Decode(event)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	if diff := cmp.Diff(wantActor, occ.Actor); diff != "" {
		t.Errorf("actor mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantBody, occ.Body, protocmp.Transform()); diff != "" {
		t.Errorf("body mismatch (-want +got):\n%s", diff)
	}
}

// TestEmptyBody verifies that all four helpers reject an envelope whose body
// field is missing. Subscribers should never see a zero-valued proto in place
// of a real one — a malformed event is an error, not a silently-empty body.
func TestEmptyBody(t *testing.T) {
	event := cloudevents.NewEvent()
	event.SetType("dev.chainguard.api.iam.group.created.v1")
	event.SetSource("https://api.test")
	event.SetSubject("abc123")
	event.SetExtension(events.BodyTypeKey, string((&iamv2.Group{}).ProtoReflect().Descriptor().FullName()))
	if err := event.SetData(cloudevents.ApplicationJSON, &events.Occurrence{}); err != nil {
		t.Fatalf("event.SetData: %v", err)
	}

	const wantErr = "empty body"
	for _, tt := range []struct {
		name string
		call func(cloudevents.Event) error
	}{{
		name: "DecodeInto",
		call: func(e cloudevents.Event) error { return events.DecodeInto(e, &iamv2.Group{}) },
	}, {
		name: "Body",
		call: func(e cloudevents.Event) error { _, err := events.Body(e); return err },
	}, {
		name: "BodyAs",
		call: func(e cloudevents.Event) error { _, err := events.BodyAs[*iamv2.Group](e); return err },
	}, {
		name: "Decode",
		call: func(e cloudevents.Event) error { _, err := events.Decode(e); return err },
	}} {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.call(event)
			if err == nil || !strings.Contains(err.Error(), wantErr) {
				t.Errorf("err: got = %v, wanted containing %q", err, wantErr)
			}
		})
	}
}
