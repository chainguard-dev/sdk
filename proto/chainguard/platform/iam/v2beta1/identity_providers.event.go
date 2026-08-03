/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

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
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *IdentityProvider) CloudEventsSubject() string {
	return x.GetUid()
}

// CloudEventsRedact implements chainguard.dev/sdk/events/Redactable.CloudEventsRedact.
func (x *IdentityProvider) CloudEventsRedact() any {
	idp := &IdentityProvider{
		Uid:         x.Uid,
		Name:        x.Name,
		Description: x.Description,
		DefaultRole: x.DefaultRole,
	}

	switch cfg := x.Configuration.(type) {
	case *IdentityProvider_Oidc:
		idp.Configuration = &IdentityProvider_Oidc{
			Oidc: &IdentityProvider_OIDC{
				Issuer:           cfg.Oidc.Issuer,
				ClientId:         cfg.Oidc.ClientId,
				AdditionalScopes: cfg.Oidc.AdditionalScopes,
				GroupsClaim:      cfg.Oidc.GroupsClaim,
				// pkce_enabled is configuration, not a credential: the audit
				// event must record whether the client authenticates with a
				// secret or PKCE.
				PkceEnabled: cfg.Oidc.PkceEnabled,
				// ClientSecret is redacted.
			},
		}
	default:
		idp.Configuration = x.Configuration
	}

	// Deep-copy SCIM config field-by-field rather than sharing the pointer, so
	// that any future credential fields (e.g. bearer_token, CUS-450) are not
	// silently carried through redaction unless explicitly copied here.
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
		return uidp.Parent(x.GetUid()), true
	default:
		return "", false
	}
}

// CloudEventsSubject implements chainguard.dev/sdk/events/Eventable.CloudEventsSubject.
func (x *DeleteIdentityProviderRequest) CloudEventsSubject() string {
	return x.GetUid()
}

func scimEventExtension(uid, key string) (string, bool) {
	if key == "group" {
		return uidp.Parent(uid), true
	}
	return "", false
}

func (x *GenerateScimTokenResponse) CloudEventsSubject() string { return x.GetIdentityProviderUid() }
func (x *GenerateScimTokenResponse) CloudEventsAsync() bool     { return true }
func (x *GenerateScimTokenResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderUid(), key)
}
func (x *GenerateScimTokenResponse) CloudEventsRedact() any {
	return &GenerateScimTokenResponse{
		IdentityProviderUid: x.IdentityProviderUid,
		EndpointUrl:         x.EndpointUrl,
		ExpireTime:          x.ExpireTime,
		Etag:                x.Etag,
	}
}

func (x *RegenerateScimTokenResponse) CloudEventsSubject() string {
	return x.GetIdentityProviderUid()
}
func (x *RegenerateScimTokenResponse) CloudEventsAsync() bool { return true }
func (x *RegenerateScimTokenResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderUid(), key)
}
func (x *RegenerateScimTokenResponse) CloudEventsRedact() any {
	return &RegenerateScimTokenResponse{
		IdentityProviderUid:     x.IdentityProviderUid,
		EndpointUrl:             x.EndpointUrl,
		ExpireTime:              x.ExpireTime,
		Etag:                    x.Etag,
		PreviousTokenExpireTime: x.PreviousTokenExpireTime,
		RequestedOverlap:        x.RequestedOverlap,
	}
}

func (x *RevokeScimTokenResponse) CloudEventsSubject() string { return x.GetIdentityProviderUid() }
func (x *RevokeScimTokenResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderUid(), key)
}

func (x *SetScimEnabledResponse) CloudEventsSubject() string { return x.GetIdentityProviderUid() }
func (x *SetScimEnabledResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderUid(), key)
}
