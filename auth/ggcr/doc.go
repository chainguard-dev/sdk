/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package ggcr provides authentication integration between Chainguard identity tokens
and the go-containerregistry library for accessing the cgr.dev container registry.

# Overview

This package implements the authn.Keychain interface from go-containerregistry,
enabling seamless authentication to cgr.dev using Chainguard identity tokens.
It handles the token exchange flow using STS (Security Token Service) to convert
base tokens into Chainguard-specific access tokens.

# Features

  - authn.Keychain implementation for cgr.dev registry authentication
  - Automatic token exchange using Chainguard STS
  - Support for any OAuth2 token source as the base credential
  - Token reuse to minimize exchange operations
  - Registry-aware authentication (only authenticates for cgr.dev)

# Usage

The primary entry point is the Keychain function, which creates a keychain
from an identity UIDP and a base token source:

	import (
		"context"

		"chainguard.dev/sdk/auth/ggcr"
		"github.com/google/go-containerregistry/pkg/name"
		"github.com/google/go-containerregistry/pkg/v1/remote"
		"google.golang.org/api/idtoken"
	)

	func pullImage(ctx context.Context, identity string) error {
		// Create a token source using ambient GCP credentials
		ts, err := idtoken.NewTokenSource(ctx, "https://cgr.dev")
		if err != nil {
			return err
		}

		// Create the keychain for the specified identity
		kc := ggcr.Keychain(identity, ts)

		// Use the keychain to pull images
		ref, err := name.ParseReference("cgr.dev/my/image:latest")
		if err != nil {
			return err
		}
		img, err := remote.Image(ref, remote.WithAuthFromKeychain(kc))
		if err != nil {
			return err
		}

		// Work with the image...
		return nil
	}

# Integration Patterns

The package integrates with go-containerregistry's remote operations:

  - remote.Image: Pull container images
  - remote.Write: Push container images
  - remote.Index: Pull multi-platform image indexes
  - remote.Layer: Pull individual layers

All remote operations accept the keychain via remote.WithAuthFromKeychain(kc).

For advanced use cases where you already have an OAuth2 token source configured
with the complete token exchange flow, use TokenSourceKeychain directly:

	kc := ggcr.TokenSourceKeychain(myConfiguredTokenSource)

# Authentication Flow

The authentication flow works as follows:

 1. Base token source provides an initial token (e.g., GCP ID token)
 2. STS exchanges the base token for a Chainguard access token
 3. Token is cached and reused until expiration
 4. Keychain provides the token to go-containerregistry as HTTP Basic auth
 5. Registry validates the token and grants access

The keychain only provides credentials when the target registry is cgr.dev.
For other registries, it returns anonymous authentication, allowing
go-containerregistry to fall back to other authentication methods.
*/
package ggcr
