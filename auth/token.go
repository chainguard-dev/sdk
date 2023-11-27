/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"context"

	"google.golang.org/grpc/credentials"
)

type tokenAuth struct {
	token  string
	secure bool
}

var _ credentials.PerRPCCredentials = (*tokenAuth)(nil)

// GetRequestMetadata implements credentials.PerRPCCredentials
func (ta *tokenAuth) GetRequestMetadata(_ context.Context, uri ...string) (map[string]string, error) { //nolint: revive
	return map[string]string{
		"Authorization": ta.token,
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (ta *tokenAuth) RequireTransportSecurity() bool {
	return ta.secure
}
