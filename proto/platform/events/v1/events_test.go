/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"testing"
)

func TestIdentityMetadata_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		want   string
		wantOk bool
	}{{
		name:   "any key returns false",
		key:    "group",
		want:   "",
		wantOk: false,
	}, {
		name:   "unknown key returns false",
		key:    "unknown",
		want:   "",
		wantOk: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metadata := &IdentityMetadata{}
			got, ok := metadata.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestIdentityMetadata_CloudEventsSubject(t *testing.T) {
	metadata := &IdentityMetadata{}
	got := metadata.CloudEventsSubject()
	if got != "" {
		t.Errorf("CloudEventsSubject() = %q, want = %q", got, "")
	}
}

func TestSubscription_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{{
		name:   "group extension returns parent UIDP",
		id:     "abc123/def456/subscription789",
		key:    "group",
		want:   "abc123/def456",
		wantOk: true,
	}, {
		name:   "unknown extension returns false",
		id:     "abc123/def456/subscription789",
		key:    "unknown",
		want:   "",
		wantOk: false,
	}, {
		name:   "group extension with empty id",
		id:     "",
		key:    "group",
		want:   "/",
		wantOk: true,
	}, {
		name:   "group extension with single segment id",
		id:     "abc123",
		key:    "group",
		want:   "/",
		wantOk: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{Id: tt.id}
			got, ok := sub.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestSubscription_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{{
		name: "returns subscription id",
		id:   "abc123/def456/subscription789",
		want: "abc123/def456/subscription789",
	}, {
		name: "empty id",
		id:   "",
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := &Subscription{Id: tt.id}
			got := sub.CloudEventsSubject()
			if got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeleteSubscriptionRequest_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{{
		name:   "group extension returns parent UIDP",
		id:     "abc123/def456/subscription789",
		key:    "group",
		want:   "abc123/def456",
		wantOk: true,
	}, {
		name:   "unknown extension returns false",
		id:     "abc123/def456/subscription789",
		key:    "unknown",
		want:   "",
		wantOk: false,
	}, {
		name:   "group extension with empty id",
		id:     "",
		key:    "group",
		want:   "/",
		wantOk: true,
	}, {
		name:   "group extension with single segment id",
		id:     "abc123",
		key:    "group",
		want:   "/",
		wantOk: true,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeleteSubscriptionRequest{Id: tt.id}
			got, ok := req.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeleteSubscriptionRequest_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{{
		name: "returns request id",
		id:   "abc123/def456/subscription789",
		want: "abc123/def456/subscription789",
	}, {
		name: "empty id",
		id:   "",
		want: "",
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeleteSubscriptionRequest{Id: tt.id}
			got := req.CloudEventsSubject()
			if got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeleteSubscriptionRequest_CloudEventsRedact(t *testing.T) {
	req := &DeleteSubscriptionRequest{Id: "abc123/def456/subscription789"}
	got := req.CloudEventsRedact()
	if got != nil {
		t.Errorf("CloudEventsRedact() = %v, want = nil", got)
	}
}
