/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/chainguard-dev/clog"
	"google.golang.org/grpc/credentials"
)

// NewFromFile attempts to create a new credentials.PerRPCCredentials based on the provided file.
// Returns nil if not found.
func NewFromFile(ctx context.Context, path string, requireTransportSecurity bool) credentials.PerRPCCredentials {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		clog.FromContext(ctx).Infof("Using OIDC token from %s to authenticate requests.", path)
		return &fileAuth{
			file:   path,
			secure: requireTransportSecurity,
		}
	}
	return nil
}

// NewFromToken attempts to create a new credentials.PerRPCCredentials based on provided OIDC token.
func NewFromToken(_ context.Context, token string, requireTransportSecurity bool) credentials.PerRPCCredentials {
	return &tokenAuth{
		token:  token,
		secure: requireTransportSecurity,
	}
}

// NewFromContext creates a new credentials.PerRPCCredentials based on a token stored in context.
// This allows callers to provide a token for each RPC.
func NewFromContext(_ context.Context, requireTransportSecurity bool) credentials.PerRPCCredentials {
	return &contextAuth{
		secure: requireTransportSecurity,
	}
}

// NormalizeIssuer massages an issuer string into a canonical form, such as
// attaching a scheme when certain "special" vendors omit them.
func NormalizeIssuer(issuer string) string {
	// Similar to go-oidc, allow Google accounts to be missing scheme:
	// https://github.com/coreos/go-oidc/blob/26c5037/oidc/verify.go#L231
	if issuer == "accounts.google.com" {
		issuer = "https://accounts.google.com"
	}
	return issuer
}

type Actor struct {
	Audience string `json:"aud"`
	Issuer   string `json:"iss"`
	Subject  string `json:"sub"`
}

func decodeToken(token string) ([]byte, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("oidc: malformed jwt, expected 3 parts got %d", len(parts))
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("oidc: malformed jwt payload: %w", err)
	}
	return raw, nil
}

func ExtractActor(token string) (*Actor, error) {
	raw, err := decodeToken(token)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Actor Actor `json:"act"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, fmt.Errorf("oidc: failed to unmarshal actor: %w", err)
	}
	return &payload.Actor, nil
}

func ExtractEmail(token string) (email string, verified bool, err error) {
	raw, err := decodeToken(token)
	if err != nil {
		return "", false, err
	}

	var payload struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", false, fmt.Errorf("oidc: failed to unmarshal claims: %w", err)
	}

	return payload.Email, payload.EmailVerified, nil
}

func ExtractIssuer(token string) (string, error) {
	iss, _, err := ExtractIssuerAndSubject(token)
	return iss, err
}

func ExtractIssuerAndSubject(token string) (string, string, error) {
	raw, err := decodeToken(token)
	if err != nil {
		return "", "", err
	}

	var payload struct {
		Issuer  string `json:"iss"`
		Subject string `json:"sub"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return "", "", fmt.Errorf("oidc: failed to unmarshal claims: %w", err)
	}
	return NormalizeIssuer(payload.Issuer), payload.Subject, nil
}

func ExtractExpiry(token string) (time.Time, error) {
	raw, err := decodeToken(token)
	if err != nil {
		return time.Time{}, err
	}

	var payload struct {
		Expiry int64 `json:"exp"`
	}

	if err := json.Unmarshal(raw, &payload); err != nil {
		return time.Time{}, fmt.Errorf("oidc: failed to unmarshal claims: %w", err)
	}

	return time.Unix(payload.Expiry, 0), nil
}
