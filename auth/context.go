/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"context"

	"google.golang.org/grpc/credentials"
)

type contextAuth struct {
	secure bool
}

// The key to associate token with context.
type tokenKey struct{}

var _ credentials.PerRPCCredentials = (*contextAuth)(nil)

// GetRequestMetadata implements credentials.PerRPCCredentials
func (ta *contextAuth) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) { //nolint: revive
	return map[string]string{
		"Authorization": GetToken(ctx),
	}, nil
}

// RequireTransportSecurity implements credentials.PerRPCCredentials
func (ta *contextAuth) RequireTransportSecurity() bool {
	return ta.secure
}

// WithToken associates the token with the returned context.
func WithToken(ctx context.Context, authz string) context.Context {
	return context.WithValue(ctx, tokenKey{}, authz)
}

// GetToken fetches the token from the context.
func GetToken(ctx context.Context) string {
	v := ctx.Value(tokenKey{})
	if v == nil {
		return ""
	}
	return v.(string)
}
