/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"chainguard.dev/sdk/events"
	cgannotations "chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/uidp"
)

func TestCreateHuntEventAnnotation(t *testing.T) {
	md := File_chainguard_platform_harness_v1alpha1_hunts_proto.Services().ByName("Hunts").Methods().ByName("CreateHunt")
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatal("CreateHunt has no options")
	}
	ea, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
	if !ok || ea == nil {
		t.Fatal("CreateHunt is missing chainguard.annotations.events")
	}
	if got, want := ea.GetType(), "dev.chainguard.api.harness.hunts.created.v1"; got != want {
		t.Errorf("type: got = %q, want = %q", got, want)
	}
	if got, want := ea.GetAudience(), cgannotations.EventAttributes_CUSTOMER; got != want {
		t.Errorf("audience: got = %v, want = %v", got, want)
	}
	if got, want := ea.GetExtensions(), []string{"group"}; !slices.Equal(got, want) {
		t.Errorf("extensions: got = %v, want = %v", got, want)
	}
}

func TestCreateHuntResponseEvent(t *testing.T) {
	org := uidp.NewUIDP("")
	huntID := org.NewChild().String()
	uploadURL := "https://storage.googleapis.com/" + uidp.NewUID().String()
	digest := "sha256:" + uidp.NewUID().String()
	resp := &CreateHuntResponse{
		Hunt: &Hunt{
			Id:           huntID,
			Name:         "hunt/" + org.String() + "/" + uidp.NewUID().String(),
			Status:       HuntStatus_HUNT_STATUS_PENDING_UPLOAD,
			EnvelopeSize: 4096,
			BundleSize:   2048,
			BundleSha256: digest,
			CreateTime:   timestamppb.New(time.Now()),
			UpdateTime:   timestamppb.New(time.Now()),
		},
		Upload: &HuntUpload{
			Url:        uploadURL,
			Headers:    map[string]string{"x-goog-if-generation-match": "0"},
			ExpireTime: timestamppb.New(time.Now().Add(6 * time.Hour)),
		},
		Created: true,
	}

	var body any = resp
	eventable, ok := body.(events.Eventable)
	if !ok {
		t.Fatal("CreateHuntResponse does not implement events.Eventable")
	}
	if got := eventable.CloudEventsSubject(); got != huntID {
		t.Errorf("subject: got = %q, want = %q", got, huntID)
	}
	extendable, ok := body.(events.Extendable)
	if !ok {
		t.Fatal("CreateHuntResponse does not implement events.Extendable")
	}
	if got, ok := extendable.CloudEventsExtension("group"); !ok || got != org.String() {
		t.Errorf("group extension: got = (%q, %v), want = (%q, true)", got, ok, org)
	}
	if _, ok := extendable.CloudEventsExtension("unknown"); ok {
		t.Error("CloudEventsExtension(unknown): got = true, want = false")
	}

	// Mirror the api-impl interceptor: it encodes CloudEventsRedact's result with protojson.
	redactable, ok := body.(events.Redactable)
	if !ok {
		t.Fatal("CreateHuntResponse does not implement events.Redactable")
	}
	redacted, ok := redactable.CloudEventsRedact().(proto.Message)
	if !ok {
		t.Fatal("CloudEventsRedact did not return a proto.Message")
	}
	encoded, err := protojson.Marshal(redacted)
	if err != nil {
		t.Fatalf("protojson.Marshal: %v", err)
	}

	var got map[string]any
	if err := json.Unmarshal(encoded, &got); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	want := map[string]any{
		"hunt":    map[string]any{"id": huntID},
		"created": true,
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("encoded body (-want, +got):\n%s", diff)
	}

	for _, forbidden := range []string{"artifact_uuid", "artifactUuid", uploadURL, digest, resp.GetHunt().GetName()} {
		if strings.Contains(string(encoded), forbidden) {
			t.Errorf("encoded body %s contains %q", encoded, forbidden)
		}
	}

	// Redaction must copy, not strip the caller's response in place.
	if resp.GetUpload().GetUrl() != uploadURL {
		t.Error("CloudEventsRedact mutated the response")
	}
}
