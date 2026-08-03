/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"testing"
	"time"

	v1 "chainguard.dev/sdk/proto/platform/iam/v1"
	"chainguard.dev/sdk/uidp"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestScimTokenEventRedaction proves the reveal-once plaintext token never
// survives into an audit event, that the issuance events are delivered
// asynchronously, and that lifecycle metadata is preserved. It mirrors the
// v2beta1 SCIM event tests.
func TestScimTokenEventRedaction(t *testing.T) {
	const (
		id    = "abc123/def456"
		token = "must-not-appear"
	)

	t.Run("generate strips token, keeps metadata", func(t *testing.T) {
		expiry := timestamppb.New(time.Unix(1893456000, 0))
		resp := &v1.GenerateScimTokenResponse{
			Token:              token,
			IdentityProviderId: id,
			EndpointUrl:        "https://scim.example.dev/scim/v2/" + id,
			ExpireTime:         expiry,
			Etag:               "etag1",
		}
		redacted := resp.CloudEventsRedact().(*v1.GenerateScimTokenResponse)
		if redacted.GetToken() != "" {
			t.Errorf("redacted token: got = %q, want stripped", redacted.GetToken())
		}
		if got, want := redacted.GetIdentityProviderId(), id; got != want {
			t.Errorf("identity_provider_id: got = %q, want = %q", got, want)
		}
		if got, want := redacted.GetEtag(), "etag1"; got != want {
			t.Errorf("etag: got = %q, want = %q", got, want)
		}
		if got, want := redacted.GetEndpointUrl(), resp.EndpointUrl; got != want {
			t.Errorf("endpoint_url: got = %q, want = %q", got, want)
		}
		if got := redacted.GetExpireTime(); got == nil || !got.AsTime().Equal(expiry.AsTime()) {
			t.Errorf("expire_time: got = %v, want = %v", got, expiry)
		}
		if !resp.CloudEventsAsync() {
			t.Error("CloudEventsAsync(): got = false, want = true")
		}
	})

	t.Run("regenerate strips token, keeps metadata", func(t *testing.T) {
		expiry := timestamppb.New(time.Unix(1893456000, 0))
		prevExpiry := timestamppb.New(time.Unix(1800000000, 0))
		resp := &v1.RegenerateScimTokenResponse{
			Token:                   token,
			IdentityProviderId:      id,
			EndpointUrl:             "https://scim.example.dev/scim/v2/" + id,
			ExpireTime:              expiry,
			Etag:                    "etag2",
			PreviousTokenExpireTime: prevExpiry,
			RequestedOverlap:        durationpb.New(time.Hour),
		}
		redacted := resp.CloudEventsRedact().(*v1.RegenerateScimTokenResponse)
		if redacted.GetToken() != "" {
			t.Errorf("redacted token: got = %q, want stripped", redacted.GetToken())
		}
		if got, want := redacted.GetIdentityProviderId(), id; got != want {
			t.Errorf("identity_provider_id: got = %q, want = %q", got, want)
		}
		if got, want := redacted.GetEndpointUrl(), resp.EndpointUrl; got != want {
			t.Errorf("endpoint_url: got = %q, want = %q", got, want)
		}
		if got := redacted.GetExpireTime(); got == nil || !got.AsTime().Equal(expiry.AsTime()) {
			t.Errorf("expire_time: got = %v, want = %v", got, expiry)
		}
		if got, want := redacted.GetEtag(), "etag2"; got != want {
			t.Errorf("etag: got = %q, want = %q", got, want)
		}
		if got := redacted.GetRequestedOverlap(); got == nil || got.AsDuration() != time.Hour {
			t.Errorf("requested_overlap: got = %v, want = 1h", got)
		}
		if got := redacted.GetPreviousTokenExpireTime(); got == nil || !got.AsTime().Equal(prevExpiry.AsTime()) {
			t.Errorf("previous_token_expire_time: got = %v, want = %v", got, prevExpiry)
		}
		if !resp.CloudEventsAsync() {
			t.Error("CloudEventsAsync(): got = false, want = true")
		}
	})

	// Every SCIM lifecycle response reports its identity provider as the event
	// subject and derives the "group" extension from the provider's parent
	// scope.
	t.Run("subject and group extension", func(t *testing.T) {
		type scimEvent interface {
			CloudEventsSubject() string
			CloudEventsExtension(string) (string, bool)
		}
		cases := []struct {
			name string
			ev   scimEvent
		}{
			{"generate", &v1.GenerateScimTokenResponse{IdentityProviderId: id}},
			{"regenerate", &v1.RegenerateScimTokenResponse{IdentityProviderId: id}},
			{"revoke", &v1.RevokeScimTokenResponse{IdentityProviderId: id}},
			{"setEnabled", &v1.SetScimEnabledResponse{IdentityProviderId: id}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				if got, want := tc.ev.CloudEventsSubject(), id; got != want {
					t.Errorf("subject: got = %q, want = %q", got, want)
				}
				got, ok := tc.ev.CloudEventsExtension("group")
				if want := uidp.Parent(id); !ok || got != want {
					t.Errorf("group extension: got = (%q, %v), want = (%q, true)", got, ok, want)
				}
				if _, ok := tc.ev.CloudEventsExtension("bogus"); ok {
					t.Error("unknown extension key: got ok = true, want = false")
				}
			})
		}
	})

	// The IdentityProvider read event carries SCIM status but must still strip
	// the OIDC client secret, and the SCIM sub-message is deep-copied so a
	// future credential field cannot silently ride through.
	t.Run("identity provider redaction preserves scim, strips client secret", func(t *testing.T) {
		idp := &v1.IdentityProvider{
			Id:   id,
			Name: "corp-okta",
			Configuration: &v1.IdentityProvider_Oidc{
				Oidc: &v1.IdentityProvider_OIDC{
					Issuer:       "https://issuer.example.com",
					ClientId:     "client-id",
					ClientSecret: "super-secret",
					GroupsClaim:  "groups",
					PkceEnabled:  true,
				},
			},
			Scim: &v1.IdentityProvider_SCIM{
				Enabled:         true,
				EndpointUrl:     "https://scim.example.dev/scim/v2/" + id,
				CredentialState: v1.IdentityProvider_SCIM_CREDENTIAL_STATE_LIVE,
				Etag:            "etag3",
			},
		}
		redacted := idp.CloudEventsRedact().(*v1.IdentityProvider)
		if got := redacted.GetOidc().GetClientSecret(); got != "" {
			t.Errorf("client_secret: got = %q, want stripped", got)
		}
		if got, want := redacted.GetOidc().GetClientId(), "client-id"; got != want {
			t.Errorf("client_id: got = %q, want = %q", got, want)
		}
		if got, want := redacted.GetOidc().GetGroupsClaim(), "groups"; got != want {
			t.Errorf("groups_claim: got = %q, want = %q", got, want)
		}
		// pkce_enabled is config, not a credential, and the only event-visible
		// signal of a confidential <-> public (PKCE-only) client transition.
		if !redacted.GetOidc().GetPkceEnabled() {
			t.Error("pkce_enabled: got = false, want preserved")
		}
		scim := redacted.GetScim()
		if scim == nil {
			t.Fatal("redacted IdentityProvider dropped SCIM status")
		}
		if !scim.GetEnabled() {
			t.Error("redacted SCIM: enabled = false, want = true")
		}
		if got, want := scim.GetCredentialState(), v1.IdentityProvider_SCIM_CREDENTIAL_STATE_LIVE; got != want {
			t.Errorf("scim credential_state: got = %v, want = %v", got, want)
		}
		if got, want := scim.GetEtag(), "etag3"; got != want {
			t.Errorf("scim etag: got = %q, want = %q", got, want)
		}
	})

	// An identity provider with no SCIM status must redact to a nil SCIM
	// sub-message, not an empty one, so "no information" stays distinct from an
	// explicit state.
	t.Run("identity provider without scim stays nil", func(t *testing.T) {
		idp := &v1.IdentityProvider{
			Id:   id,
			Name: "corp-okta",
			Configuration: &v1.IdentityProvider_Oidc{
				Oidc: &v1.IdentityProvider_OIDC{
					Issuer:       "https://issuer.example.com",
					ClientId:     "client-id",
					ClientSecret: "super-secret",
				},
			},
		}
		redacted := idp.CloudEventsRedact().(*v1.IdentityProvider)
		if got := redacted.GetScim(); got != nil {
			t.Errorf("scim: got = %v, want = nil", got)
		}
		if got := redacted.GetOidc().GetClientSecret(); got != "" {
			t.Errorf("client_secret: got = %q, want stripped", got)
		}
	})
}
