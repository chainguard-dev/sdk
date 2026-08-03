/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/events"
	"chainguard.dev/sdk/uidp"
)

var (
	_ events.Eventable      = (*IdentityProvider)(nil)
	_ events.Extendable     = (*IdentityProvider)(nil)
	_ events.Redactable     = (*IdentityProvider)(nil)
	_ events.Eventable      = (*DeleteIdentityProviderRequest)(nil)
	_ events.Extendable     = (*DeleteIdentityProviderRequest)(nil)
	_ events.Eventable      = (*GenerateScimTokenResponse)(nil)
	_ events.Extendable     = (*GenerateScimTokenResponse)(nil)
	_ events.Redactable     = (*GenerateScimTokenResponse)(nil)
	_ events.AsyncEventable = (*GenerateScimTokenResponse)(nil)
	_ events.Eventable      = (*RegenerateScimTokenResponse)(nil)
	_ events.Extendable     = (*RegenerateScimTokenResponse)(nil)
	_ events.Redactable     = (*RegenerateScimTokenResponse)(nil)
	_ events.AsyncEventable = (*RegenerateScimTokenResponse)(nil)
	_ events.Eventable      = (*RevokeScimTokenResponse)(nil)
	_ events.Extendable     = (*RevokeScimTokenResponse)(nil)
	_ events.Eventable      = (*SetScimEnabledResponse)(nil)
	_ events.Extendable     = (*SetScimEnabledResponse)(nil)
)

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *IdentityProvider) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *IdentityProvider) CloudEventsSubject() string {
	return x.GetId()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
func (x *IdentityProvider) CloudEventsRedact() any {
	idp := &IdentityProvider{
		Id:          x.Id,
		Name:        x.Name,
		Description: x.Description,
		DefaultRole: x.DefaultRole,
	}

	switch cfg := x.Configuration.(type) {
	case *IdentityProvider_Oidc:
		// redact OIDC configuration
		idp.Configuration = &IdentityProvider_Oidc{
			Oidc: &IdentityProvider_OIDC{
				Issuer:   cfg.Oidc.Issuer,
				ClientId: cfg.Oidc.ClientId,
				// ClientSecret is redacted.
				AdditionalScopes: cfg.Oidc.AdditionalScopes,
				// GroupsClaim must survive redaction so the audit event records
				// changes to the trusted group-mapping claim.
				GroupsClaim: cfg.Oidc.GroupsClaim,
				// PkceEnabled is configuration, not a credential: with the
				// client_secret redacted, it is the only signal in the event of
				// a confidential <-> public (PKCE-only) client transition.
				PkceEnabled: cfg.Oidc.PkceEnabled,
			},
		}
	default:
		// no redaction
		idp.Configuration = x.Configuration
	}

	// Deep-copy SCIM status field-by-field rather than sharing the pointer, so
	// that any future credential fields (e.g. bearer_token) are not silently
	// carried through redaction unless explicitly copied here.
	if scim := x.GetScim(); scim != nil {
		idp.Scim = &IdentityProvider_SCIM{
			Enabled:                 scim.Enabled,
			EndpointUrl:             scim.EndpointUrl,
			TokenExpireTime:         scim.TokenExpireTime,
			PreviousTokenExpireTime: scim.PreviousTokenExpireTime,
			CredentialState:         scim.CredentialState,
			Etag:                    scim.Etag,
		}
	}

	return idp
}

// CloudEventsExtension implements chainguard.dev/sdk/events/Extendable.CloudEventsExtension
func (x *DeleteIdentityProviderRequest) CloudEventsExtension(key string) (string, bool) {
	switch key {
	case "group":
		return uidp.Parent(x.GetId()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteIdentityProviderRequest) CloudEventsSubject() string {
	return x.GetId()
}

func scimEventExtension(id, key string) (string, bool) {
	if key == "group" {
		return uidp.Parent(id), true
	}
	return "", false
}

func (x *GenerateScimTokenResponse) CloudEventsSubject() string { return x.GetIdentityProviderId() }
func (x *GenerateScimTokenResponse) CloudEventsAsync() bool     { return true }
func (x *GenerateScimTokenResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderId(), key)
}

// CloudEventsRedact omits the reveal-once token from the audit event; only the
// non-secret result fields are carried.
func (x *GenerateScimTokenResponse) CloudEventsRedact() any {
	return &GenerateScimTokenResponse{
		IdentityProviderId: x.IdentityProviderId,
		EndpointUrl:        x.EndpointUrl,
		ExpireTime:         x.ExpireTime,
		Etag:               x.Etag,
	}
}

func (x *RegenerateScimTokenResponse) CloudEventsSubject() string {
	return x.GetIdentityProviderId()
}
func (x *RegenerateScimTokenResponse) CloudEventsAsync() bool { return true }
func (x *RegenerateScimTokenResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderId(), key)
}

// CloudEventsRedact omits the reveal-once token from the audit event; only the
// non-secret result fields are carried.
func (x *RegenerateScimTokenResponse) CloudEventsRedact() any {
	return &RegenerateScimTokenResponse{
		IdentityProviderId:      x.IdentityProviderId,
		EndpointUrl:             x.EndpointUrl,
		ExpireTime:              x.ExpireTime,
		Etag:                    x.Etag,
		PreviousTokenExpireTime: x.PreviousTokenExpireTime,
		RequestedOverlap:        x.RequestedOverlap,
	}
}

func (x *RevokeScimTokenResponse) CloudEventsSubject() string { return x.GetIdentityProviderId() }
func (x *RevokeScimTokenResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderId(), key)
}

func (x *SetScimEnabledResponse) CloudEventsSubject() string { return x.GetIdentityProviderId() }
func (x *SetScimEnabledResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderId(), key)
}
