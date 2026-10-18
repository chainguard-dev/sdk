/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"chainguard.dev/sdk/events"
	cgannotations "chainguard.dev/sdk/proto/annotations"
	"chainguard.dev/sdk/uidp"
)

// eventAnnotation returns the chainguard.annotations.events option on a method.
func eventAnnotation(t *testing.T, fd protoreflect.FileDescriptor, service, method protoreflect.Name) *cgannotations.EventAttributes {
	t.Helper()
	md := fd.Services().ByName(service).Methods().ByName(method)
	if md == nil {
		t.Fatalf("%s.%s not found", service, method)
	}
	opts, ok := md.Options().(*descriptorpb.MethodOptions)
	if !ok || opts == nil {
		t.Fatalf("%s.%s has no options", service, method)
	}
	ea, ok := proto.GetExtension(opts, cgannotations.E_Events).(*cgannotations.EventAttributes)
	if !ok || ea == nil {
		t.Fatalf("%s.%s is missing chainguard.annotations.events", service, method)
	}
	return ea
}

// encodeRedacted mirrors the api-impl interceptor, which encodes
// CloudEventsRedact's result with protojson.
func encodeRedacted(t *testing.T, body any) ([]byte, map[string]any) {
	t.Helper()
	redactable, ok := body.(events.Redactable)
	if !ok {
		t.Fatalf("%T does not implement events.Redactable", body)
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
	return encoded, got
}

func TestCreateUploadEventAnnotation(t *testing.T) {
	got := eventAnnotation(t, File_chainguard_platform_harness_v1alpha1_uploads_proto, "Uploads", "CreateUpload")
	want := &cgannotations.EventAttributes{
		Type:       "dev.chainguard.api.harness.uploads.created.v1",
		Extensions: []string{"group"},
		Audience:   cgannotations.EventAttributes_CUSTOMER,
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("events annotation (-want, +got):\n%s", diff)
	}
}

func TestCreateUploadResponseEvent(t *testing.T) {
	org := uidp.NewUIDP("")
	uploadID := org.NewChild().String()
	uploadURL := "https://storage.googleapis.com/" + uidp.NewUID().String()
	digest := "sha256:" + uidp.NewUID().String()
	resp := &CreateUploadResponse{
		Upload: &Upload{
			Id:           uploadID,
			Status:       UploadStatus_UPLOAD_STATUS_PENDING,
			EnvelopeSize: 4096,
			BundleSize:   2048,
			BundleSha256: digest,
			CreateTime:   timestamppb.New(time.Now()),
			UpdateTime:   timestamppb.New(time.Now()),
		},
		UploadUrl: &UploadUrl{
			Url:        uploadURL,
			Headers:    map[string]string{"x-goog-if-generation-match": "0"},
			ExpireTime: timestamppb.New(time.Now().Add(6 * time.Hour)),
		},
	}

	before := proto.Clone(resp)

	var body any = resp
	eventable, ok := body.(events.Eventable)
	if !ok {
		t.Fatal("CreateUploadResponse does not implement events.Eventable")
	}
	if got := eventable.CloudEventsSubject(); got != uploadID {
		t.Errorf("subject: got = %q, want = %q", got, uploadID)
	}
	extendable, ok := body.(events.Extendable)
	if !ok {
		t.Fatal("CreateUploadResponse does not implement events.Extendable")
	}
	if got, ok := extendable.CloudEventsExtension("group"); !ok || got != org.String() {
		t.Errorf("group extension: got = (%q, %v), want = (%q, true)", got, ok, org)
	}
	if _, ok := extendable.CloudEventsExtension("unknown"); ok {
		t.Error("CloudEventsExtension(unknown): got = true, want = false")
	}

	encoded, got := encodeRedacted(t, body)
	want := map[string]any{
		"upload": map[string]any{"id": uploadID},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("encoded body (-want, +got):\n%s", diff)
	}

	for _, forbidden := range []string{"artifact_uuid", "artifactUuid", uploadURL, digest} {
		if strings.Contains(string(encoded), forbidden) {
			t.Errorf("encoded body %s contains %q", encoded, forbidden)
		}
	}

	// Redaction must copy, not strip the caller's response in place.
	if diff := cmp.Diff(before, resp, protocmp.Transform()); diff != "" {
		t.Errorf("CloudEventsRedact mutated the response (-before, +after):\n%s", diff)
	}
}
