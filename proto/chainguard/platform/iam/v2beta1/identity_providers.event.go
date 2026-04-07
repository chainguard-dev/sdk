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
	_ events.Eventable  = (*IdentityProvider)(nil)
	_ events.Extendable = (*IdentityProvider)(nil)
	_ events.Redactable = (*IdentityProvider)(nil)
	_ events.Eventable  = (*DeleteIdentityProviderRequest)(nil)
	_ events.Extendable = (*DeleteIdentityProviderRequest)(nil)
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
				// ClientSecret is redacted.
			},
		}
	default:
		idp.Configuration = x.Configuration
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
