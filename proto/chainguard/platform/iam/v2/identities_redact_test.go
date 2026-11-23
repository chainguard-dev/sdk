/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

import (
	"strings"
	"testing"

	"chainguard.dev/sdk/events"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func Test_Identity_CloudEventsRedact(t *testing.T) {
	id := &Identity{
		Uid:            "org/identity",
		Name:           "pull-token",
		CreatedBy:      "org/creator",
		CreatedByEmail: "creator@example.com",
	}
	orig := proto.Clone(id)

	var r events.Redactable = id
	got, ok := r.CloudEventsRedact().(*Identity)
	if !ok {
		t.Fatalf("CloudEventsRedact() = %T, want *Identity", r.CloudEventsRedact())
	}
	b, err := protojson.Marshal(got)
	if err != nil {
		t.Fatalf("protojson.Marshal() = %v", err)
	}
	if strings.Contains(string(b), "creator@example.com") {
		t.Errorf("event body = %s, want no creator email", b)
	}
	if got.GetCreatedBy() != "org/creator" || got.GetName() != "pull-token" {
		t.Errorf("CloudEventsRedact() = %v, want every other field kept", got)
	}
	if !proto.Equal(id, orig) {
		t.Errorf("CloudEventsRedact() changed the response to %v", id)
	}
}
