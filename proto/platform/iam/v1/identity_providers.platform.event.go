/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"chainguard.dev/sdk/uidp"
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
func (x *IdentityProvider) CloudEventsRedact() interface{} {
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
			},
		}
	default:
		// no redaction
		idp.Configuration = x.Configuration
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
