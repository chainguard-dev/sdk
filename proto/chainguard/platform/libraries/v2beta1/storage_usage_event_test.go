/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"testing"

	"chainguard.dev/sdk/events"
	v2beta1 "chainguard.dev/sdk/proto/chainguard/platform/libraries/v2beta1"
)

// Test_StorageUsage_DeleteUploadedVersion_EventInterfaces pins the delete
// event: the body is the response, the group extension is the organization,
// the subject names the ecosystem with the stored version, so the same name
// and version in two ecosystems are two subjects, and a version that was
// already deleted sends no event.
func Test_StorageUsage_DeleteUploadedVersion_EventInterfaces(t *testing.T) {
	const org = "720909c9f7df8751f32a25ff2af1ef4c8f92dc8b"
	req := &v2beta1.DeleteUploadedVersionRequest{
		Parent:    org,
		Ecosystem: v2beta1.StorageUsageEcosystem_STORAGE_USAGE_ECOSYSTEM_PYTHON,
		Name:      "Acme_Widget",
		Version:   "1.2.3",
	}
	resp := &v2beta1.DeleteUploadedVersionResponse{Name: "acme-widget", Version: "1.2.3"}

	body := req.CloudEventsBody(resp)
	if body == nil {
		t.Fatal("CloudEventsBody: got = nil, want the response")
	}
	if got, want := body.CloudEventsSubject(), "python/acme-widget@1.2.3"; got != want {
		t.Errorf("subject: got = %q, want = %q", got, want)
	}
	ext, ok := body.(events.Extendable)
	if !ok {
		t.Fatal("the body is not Extendable")
	}
	if got, ok := ext.CloudEventsExtension("group"); !ok || got != org {
		t.Errorf("group extension: got = %q, %v, want = %q, true", got, ok, org)
	}
	if _, ok := ext.CloudEventsExtension("other"); ok {
		t.Error("an unknown extension: got = true, want = false")
	}
	if s, ok := body.(events.Suppressible); !ok || s.CloudEventsSuppress() {
		t.Error("a delete that changed something: want an event")
	}

	already := req.CloudEventsBody(&v2beta1.DeleteUploadedVersionResponse{Name: "acme-widget", Version: "1.2.3", AlreadyDeleted: true})
	if s, ok := already.(events.Suppressible); !ok || !s.CloudEventsSuppress() {
		t.Error("an already deleted version: want no event")
	}

	other := &v2beta1.DeleteUploadedVersionRequest{Parent: org, Ecosystem: v2beta1.StorageUsageEcosystem_STORAGE_USAGE_ECOSYSTEM_JAVASCRIPT}
	if other.CloudEventsBody(resp).CloudEventsSubject() == body.CloudEventsSubject() {
		t.Errorf("two ecosystems share the subject %q", body.CloudEventsSubject())
	}
	if req.CloudEventsBody(nil) != nil {
		t.Error("CloudEventsBody of no response: want nil, so no event is sent")
	}
}
