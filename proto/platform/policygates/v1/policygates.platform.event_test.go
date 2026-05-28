/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "testing"

func TestBinding_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "group extension returns parent UIDP",
			id:     "abc123/def456/binding789",
			key:    "group",
			want:   "abc123/def456",
			wantOk: true,
		},
		{
			name:   "unknown extension returns false",
			id:     "abc123/def456/binding789",
			key:    "unknown",
			want:   "",
			wantOk: false,
		},
		{
			name:   "group extension with empty id",
			id:     "",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
		{
			name:   "group extension with single segment id",
			id:     "abc123",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binding{Id: tt.id}
			got, ok := b.CloudEventsExtension(tt.key)
			if ok != tt.wantOk {
				t.Errorf("CloudEventsExtension() ok = %v, want = %v", ok, tt.wantOk)
			}
			if got != tt.want {
				t.Errorf("CloudEventsExtension() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestBinding_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "non-empty id",
			id:   "abc123/def456/binding789",
			want: "abc123/def456/binding789",
		},
		{
			name: "empty id",
			id:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Binding{Id: tt.id}
			if got := b.CloudEventsSubject(); got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}

func TestDeleteBindingRequest_CloudEventsExtension(t *testing.T) {
	tests := []struct {
		name   string
		id     string
		key    string
		want   string
		wantOk bool
	}{
		{
			name:   "group extension returns parent UIDP",
			id:     "abc123/def456/binding789",
			key:    "group",
			want:   "abc123/def456",
			wantOk: true,
		},
		{
			name:   "unknown extension returns false",
			id:     "abc123/def456/binding789",
			key:    "unknown",
			want:   "",
			wantOk: false,
		},
		{
			name:   "group extension with empty id",
			id:     "",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
		{
			name:   "group extension with single segment id",
			id:     "abc123",
			key:    "group",
			want:   "/",
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeleteBindingRequest{Id: tt.id}
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

func TestDeleteBindingRequest_CloudEventsSubject(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "non-empty id",
			id:   "abc123/def456/binding789",
			want: "abc123/def456/binding789",
		},
		{
			name: "empty id",
			id:   "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &DeleteBindingRequest{Id: tt.id}
			if got := req.CloudEventsSubject(); got != tt.want {
				t.Errorf("CloudEventsSubject() = %q, want = %q", got, tt.want)
			}
		})
	}
}
