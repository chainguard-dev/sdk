/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package sts provides OIDC token exchange functionality for Chainguard services.

# Overview

The sts package implements secure token exchange (STS) for converting third-party
OIDC tokens into Chainguard-issued tokens. It supports both standard gRPC-based
exchanges and HTTP/1.1 downgrade mode for compatibility with environments that
require HTTP/1.1.

# Features

  - OIDC token exchange with Chainguard issuers
  - Refresh token support for long-lived sessions
  - Configurable scopes and capabilities
  - Identity and identity provider selection
  - HTTP/1.1 downgrade mode for compatibility
  - oauth2.TokenSource integration for seamless authentication

# Basic Usage

The simplest way to exchange a token is using the ExchangePair function:

	ctx := context.Background()
	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken)
	if err != nil {
		log.Fatal(err)
	}
	// Use tokenPair.AccessToken for authenticated requests

# Exchanger Interface

For more control, create an Exchanger instance with New():

	exchanger := sts.New(issuer, audience)
	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

# Refresh Tokens

Exchange refresh tokens for new access tokens:

	accessToken, refreshToken, err := exchanger.Refresh(ctx, oldRefreshToken)
	if err != nil {
		log.Fatal(err)
	}

# Customization Options

Configure the exchanger with options:

	exchanger := sts.New(issuer, audience,
		sts.WithScope("read", "write"),
		sts.WithCapabilities("groups.list"),
		sts.WithIdentity(identityUID),
		sts.WithUserAgent("my-app/1.0"),
	)

# HTTP/1.1 Downgrade Mode

For environments requiring HTTP/1.1:

	exchanger := sts.NewHTTP1DowngradeExchanger(issuer, audience)
	tokenPair, err := exchanger.Exchange(ctx, idToken)

Or use the WithHTTP1Downgrade option:

	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken,
		sts.WithHTTP1Downgrade(),
	)

# oauth2 Integration

Wrap an existing oauth2.TokenSource to automatically exchange tokens:

	baseTokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: thirdPartyToken,
	})
	exchanger := sts.New(issuer, audience)
	chainguardTokenSource := sts.NewContextTokenSource(ctx, baseTokenSource, exchanger)

	// Use with any oauth2-compatible client
	client := oauth2.NewClient(ctx, chainguardTokenSource)

# Integration Patterns

The package integrates with standard Go authentication patterns:

  - Use with google.golang.org/grpc/credentials/oauth for gRPC authentication
  - Use with golang.org/x/oauth2 for HTTP client authentication
  - Combine with token caching for improved performance
  - Chain with other TokenSource implementations for multi-stage authentication

# Error Handling

All functions return descriptive errors that can be inspected and wrapped:

	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken)
	if err != nil {
		return fmt.Errorf("exchanging token: %w", err)
	}

# Thread Safety

Exchanger instances are safe for concurrent use. Multiple goroutines can call
Exchange and Refresh on the same Exchanger instance simultaneously.
*/
package sts
