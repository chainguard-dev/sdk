/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts

import (
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"chainguard.dev/sdk/auth/token"
)

const (
	// supersededErrorDomain and supersededErrorReason identify the
	// google.rpc.ErrorInfo the issuer attaches when it refuses a refresh token
	// it has already superseded. They are the issuer's wire contract; its own
	// constants live in a package outside this module.
	supersededErrorDomain = "oidc.chainguard.dev"
	supersededErrorReason = "REFRESH_TOKEN_SUPERSEDED"
)

// AsSupersededRefreshToken marks err as token.ErrPresentedTokenSuperseded when
// the issuer refused the presented refresh token because it had already been
// superseded, and returns err unchanged otherwise. A nil error stays nil.
//
// Every caller that runs an Exchanger's Refresh inside token.Refresh should
// pass the exchange error through this, so that token.Refresh discards the
// refused token instead of leaving it to be presented again: the next
// presentation lands outside the issuer's grace window, which is replay, and
// the issuer revokes the family. The cache is shared between chainguard
// clients, so one that keeps a refused token hands the replay to the others.
//
// It lives here rather than in auth/token because recognizing the refusal
// means reading gRPC status details, which that package deliberately does not
// do, and it is exported rather than per-caller because the classification has
// to agree everywhere -- the reason is a wire contract, not a local choice.
//
// The reason is matched, not the status code. UNAUTHENTICATED is what every
// rejected token comes back as, and NotFound is what a scoped exchange returns
// when it resolves to no capabilities -- there the presented token is still the
// live head, so discarding on the code would destroy a working session.
func AsSupersededRefreshToken(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Unauthenticated {
		return err
	}
	for _, d := range st.Details() {
		info, ok := d.(*errdetails.ErrorInfo)
		if !ok {
			continue
		}
		if info.GetDomain() == supersededErrorDomain && info.GetReason() == supersededErrorReason {
			return fmt.Errorf("%w: %w", token.ErrPresentedTokenSuperseded, err)
		}
	}
	return err
}
