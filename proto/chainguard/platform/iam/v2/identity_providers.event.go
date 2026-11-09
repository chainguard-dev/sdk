/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

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
	_ events.Eventable      = (*ResetCorrelationRuleResponse)(nil)
	_ events.Extendable     = (*ResetCorrelationRuleResponse)(nil)
	_ events.Redactable     = (*ResetCorrelationRuleResponse)(nil)
	_ events.AsyncEventable = (*ResetCorrelationRuleResponse)(nil)

	_ events.Eventable           = (*RemoveManualBindingsResponse)(nil)
	_ events.Extendable          = (*RemoveManualBindingsResponse)(nil)
	_ events.Redactable          = (*RemoveManualBindingsResponse)(nil)
	_ events.Suppressible        = (*RemoveManualBindingsResponse)(nil)
	_ events.AdditionalEventable = (*RemoveManualBindingsResponse)(nil)
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
				// correlation_rule is configuration, not a credential: the audit
				// event must record how logins bind to SCIM-provisioned users.
				CorrelationRule: cfg.Oidc.CorrelationRule,
				// token_endpoint_auth_method is configuration, not a credential:
				// the audit event must record how the client authenticates at the
				// token endpoint.
				TokenEndpointAuthMethod: cfg.Oidc.TokenEndpointAuthMethod,
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

func (x *ResetCorrelationRuleResponse) CloudEventsSubject() string {
	return x.GetIdentityProvider().GetUid()
}

// CloudEventsAsync delivers the reset's audit event detached from the request,
// so a caller that gave up while the reset finished does not cancel it.
func (x *ResetCorrelationRuleResponse) CloudEventsAsync() bool { return true }
func (x *ResetCorrelationRuleResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProvider().GetUid(), key)
}

// CloudEventsRedact keeps the reset's outcome, which is the audit record, and
// redacts the identity provider as its own events do.
func (x *ResetCorrelationRuleResponse) CloudEventsRedact() any {
	r := &ResetCorrelationRuleResponse{
		PreviousCorrelationRule: x.PreviousCorrelationRule,
		Changed:                 x.Changed,
		UsersUnlinked:           x.UsersUnlinked,
		BindingsDropped:         x.BindingsDropped,
		UnlinkComplete:          x.UnlinkComplete,
	}
	if x.IdentityProvider != nil {
		r.IdentityProvider = x.IdentityProvider.CloudEventsRedact().(*IdentityProvider)
	}
	return r
}

func (x *RemoveManualBindingsResponse) CloudEventsSubject() string {
	return x.GetIdentityProviderUid()
}

func (x *RemoveManualBindingsResponse) CloudEventsExtension(key string) (string, bool) {
	return scimEventExtension(x.GetIdentityProviderUid(), key)
}

// CloudEventsSuppress publishes nothing for a dry run, which removed nothing.
func (x *RemoveManualBindingsResponse) CloudEventsSuppress() bool { return x.GetDryRun() }

// roleBindingDeletedEventType is the type of the DeleteRoleBinding event.
const roleBindingDeletedEventType = "dev.chainguard.api.iam.rolebindings.deleted.v1"

// CloudEventsAdditional publishes the role binding delete event for each binding
// this call removed, the same event deleting it one at a time publishes. A
// binding that was already gone or was refused was not removed here and gets
// none.
func (x *RemoveManualBindingsResponse) CloudEventsAdditional() []events.AdditionalEvent {
	var out []events.AdditionalEvent
	for _, removal := range x.GetRemovals() {
		if removal.GetOutcome() != ManualBindingRemovalOutcome_MANUAL_BINDING_REMOVAL_OUTCOME_REMOVED {
			continue
		}
		out = append(out, events.AdditionalEvent{
			Type:       roleBindingDeletedEventType,
			Extensions: []string{"group"},
			Audience:   events.AudienceCustomer,
			Body:       &DeleteRoleBindingRequest{Uid: removal.GetBinding().GetUid()},
		})
	}
	return out
}

// CloudEventsRedact keeps what the audit record needs: which bindings were
// removed, refused or already gone, for whom, and in what class. It drops the
// display names of the people, roles and groups, which the UIDs identify, and
// the superseding sources whole, so the record does not say which SCIM
// bindings or groups made each binding redundant.
func (x *RemoveManualBindingsResponse) CloudEventsRedact() any {
	r := &RemoveManualBindingsResponse{
		DryRun:              x.DryRun,
		PeopleLosingAccess:  x.PeopleLosingAccess,
		MoreRemain:          x.MoreRemain,
		IdentityProviderUid: x.IdentityProviderUid,
	}
	for _, removal := range x.GetRemovals() {
		b := removal.GetBinding()
		r.Removals = append(r.Removals, &ManualBindingRemoval{
			IdentityUid: removal.GetIdentityUid(),
			Binding: &ManualBinding{
				Uid:          b.GetUid(),
				RoleUid:      b.GetRoleUid(),
				GroupUid:     b.GetGroupUid(),
				BindingClass: b.GetBindingClass(),
				Exclusion:    b.GetExclusion(),
			},
			Outcome: removal.GetOutcome(),
			Reason:  removal.GetReason(),
		})
	}
	return r
}
