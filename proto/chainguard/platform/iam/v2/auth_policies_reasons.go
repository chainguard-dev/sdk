/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2

// SsoOnlyNotReadyReason is the google.rpc.ErrorInfo reason carried by the
// FAILED_PRECONDITION that UpdateAuthPolicy returns when sso_only may not be
// enabled yet. The same status carries one google.rpc.PreconditionFailure
// violation per blocker, whose type is one of the SsoOnlyReason* codes below
// and whose description is the blocker's message.
const SsoOnlyNotReadyReason = "SSO_ONLY_NOT_READY"

// The stable codes of SsoOnlyReadiness.Reason.code and of the violations a
// refused enable carries. New codes may be added; render a reason's message
// for any code not recognised.
const (
	// SsoOnlyReasonIdentityProviderMissing: the organization has no identity
	// provider, so there is nothing to validate and nothing to sign in with.
	SsoOnlyReasonIdentityProviderMissing = "IDENTITY_PROVIDER_MISSING"
	// SsoOnlyReasonIdpValidationMissing: no passed identity-provider
	// validation is recorded.
	SsoOnlyReasonIdpValidationMissing = "IDP_VALIDATION_MISSING"
	// SsoOnlyReasonIdpValidationStale: a validation is recorded, but the
	// provider's configuration has changed since.
	SsoOnlyReasonIdpValidationStale = "IDP_VALIDATION_STALE"
	// SsoOnlyReasonBreakGlassNotDesignated: no break-glass identity is
	// designated.
	SsoOnlyReasonBreakGlassNotDesignated = "BREAK_GLASS_NOT_DESIGNATED"
	// SsoOnlyReasonBreakGlassValidationMissing: a designated identity has
	// never validated its non-SSO sign-in.
	SsoOnlyReasonBreakGlassValidationMissing = "BREAK_GLASS_VALIDATION_MISSING"
	// SsoOnlyReasonBreakGlassValidationExpired: a designated identity's
	// validation is older than the configured maximum age.
	SsoOnlyReasonBreakGlassValidationExpired = "BREAK_GLASS_VALIDATION_EXPIRED"
	// SsoOnlyReasonBreakGlassNotOwner: a validated designated identity cannot
	// sign in or holds no owner role on the organization, so it could not
	// turn sso_only back off.
	SsoOnlyReasonBreakGlassNotOwner = "BREAK_GLASS_NOT_OWNER"
	// SsoOnlyReasonOrganizationNotVerified: the organization is not verified
	// and so can have neither an identity provider nor a policy.
	SsoOnlyReasonOrganizationNotVerified = "ORGANIZATION_NOT_VERIFIED"
)
