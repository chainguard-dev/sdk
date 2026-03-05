/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package sts_test

import (
	"context"
	"fmt"
	"log"

	"chainguard.dev/sdk/sts"
	"golang.org/x/oauth2"
)

// Example demonstrates basic token exchange.
func Example() {
	ctx := context.Background()
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
	fmt.Printf("Refresh token: %s\n", tokenPair.RefreshToken)
}

// ExampleNew demonstrates creating an Exchanger instance.
func ExampleNew() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleNew_withOptions demonstrates creating an Exchanger with options.
func ExampleNew_withOptions() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience,
		sts.WithScope("read", "write"),
		sts.WithCapabilities("groups.list", "roles.list"),
		sts.WithIdentity("my-identity-uid"),
		sts.WithUserAgent("my-app/1.0.0"),
	)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleExchanger_Exchange demonstrates exchanging a token.
func ExampleExchanger_Exchange() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	exchanger := sts.New(issuer, audience)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
	fmt.Printf("Refresh token: %s\n", tokenPair.RefreshToken)
	fmt.Printf("Expiry: %s\n", tokenPair.Expiry)
}

// ExampleExchanger_Refresh demonstrates refreshing a token.
func ExampleExchanger_Refresh() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	exchanger := sts.New(issuer, audience)

	ctx := context.Background()
	refreshToken := "refresh_token_value"

	accessToken, newRefreshToken, err := exchanger.Refresh(ctx, refreshToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("New access token: %s\n", accessToken)
	fmt.Printf("New refresh token: %s\n", newRefreshToken)
}

// ExampleExchangePair demonstrates the convenience function for token exchange.
func ExampleExchangePair() {
	ctx := context.Background()
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleExchangePair_withOptions demonstrates ExchangePair with options.
func ExampleExchangePair_withOptions() {
	ctx := context.Background()
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken,
		sts.WithScope("read"),
		sts.WithCapabilities("groups.list"),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleNewHTTP1DowngradeExchanger demonstrates creating an HTTP/1.1 exchanger.
func ExampleNewHTTP1DowngradeExchanger() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.NewHTTP1DowngradeExchanger(issuer, audience)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleNewTokenSource demonstrates wrapping a TokenSource for automatic exchange.
func ExampleNewTokenSource() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	// Base token source (e.g., from a third-party provider)
	baseTokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: "third-party-token",
	})

	// Create exchanger
	exchanger := sts.New(issuer, audience)

	// Wrap the base token source
	chainguardTokenSource := sts.NewTokenSource(baseTokenSource, exchanger)

	// Use with any oauth2-compatible client
	ctx := context.Background()
	client := oauth2.NewClient(ctx, chainguardTokenSource)

	// The client will automatically exchange tokens when needed
	_ = client
}

// ExampleNewContextTokenSource demonstrates wrapping a TokenSource with context.
func ExampleNewContextTokenSource() {
	ctx := context.Background()
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	// Base token source
	baseTokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: "third-party-token",
	})

	// Create exchanger
	exchanger := sts.New(issuer, audience)

	// Wrap with context
	chainguardTokenSource := sts.NewContextTokenSource(ctx, baseTokenSource, exchanger)

	// Use with HTTP client
	client := oauth2.NewClient(ctx, chainguardTokenSource)
	_ = client
}

// ExampleNewTokenSourceFromValues demonstrates the convenience function for creating a TokenSource.
func ExampleNewTokenSourceFromValues() {
	ctx := context.Background()
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	identity := "my-identity-uid"

	// Base token source
	baseTokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: "third-party-token",
	})

	// Create token source with common parameters
	chainguardTokenSource := sts.NewTokenSourceFromValues(
		ctx, issuer, audience, identity, baseTokenSource,
	)

	// Use with HTTP client
	client := oauth2.NewClient(ctx, chainguardTokenSource)
	_ = client
}

// ExampleWithScope demonstrates using the WithScope option.
func ExampleWithScope() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience,
		sts.WithScope("read", "write", "admin"),
	)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleWithCapabilities demonstrates using the WithCapabilities option.
func ExampleWithCapabilities() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience,
		sts.WithCapabilities("groups.list", "roles.list", "policies.read"),
	)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleWithIdentity demonstrates using the WithIdentity option.
func ExampleWithIdentity() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience,
		sts.WithIdentity("my-identity-uid"),
	)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleWithIdentityProvider demonstrates using the WithIdentityProvider option.
func ExampleWithIdentityProvider() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience,
		sts.WithIdentityProvider("my-identity-provider-uid"),
	)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleWithUserAgent demonstrates using the WithUserAgent option.
func ExampleWithUserAgent() {
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"

	exchanger := sts.New(issuer, audience,
		sts.WithUserAgent("my-application/1.0.0"),
	)

	ctx := context.Background()
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := exchanger.Exchange(ctx, idToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}

// ExampleWithHTTP1Downgrade demonstrates using the WithHTTP1Downgrade option.
func ExampleWithHTTP1Downgrade() {
	ctx := context.Background()
	issuer := "https://issuer.example.com"
	audience := "https://audience.example.com"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."

	tokenPair, err := sts.ExchangePair(ctx, issuer, audience, idToken,
		sts.WithHTTP1Downgrade(),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Access token: %s\n", tokenPair.AccessToken)
}
