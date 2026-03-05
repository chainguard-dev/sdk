/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package ggcr_test

import (
	"context"
	"fmt"
	"log"

	"chainguard.dev/sdk/auth/ggcr"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

// ExampleKeychain demonstrates how to create a keychain for authenticating
// to cgr.dev using a Chainguard identity and a base token source.
func ExampleKeychain() {
	ctx := context.Background()

	// Create a token source using ambient GCP credentials
	ts, err := idtoken.NewTokenSource(ctx, "https://cgr.dev")
	if err != nil {
		log.Fatal(err)
	}

	// Create the keychain for a specific Chainguard identity
	kc := ggcr.Keychain("my-identity-uidp", ts)

	// Use the keychain to pull an image from cgr.dev
	ref, err := name.ParseReference("cgr.dev/my/image:latest")
	if err != nil {
		log.Fatal(err)
	}
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(kc))
	if err != nil {
		log.Fatal(err)
	}

	// Get the image digest
	digest, err := img.Digest()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Pulled image with digest: %s\n", digest)
}

// ExampleKeychain_withContext demonstrates using the keychain with a context
// for timeout and cancellation support.
func ExampleKeychain_withContext() {
	ctx := context.Background()

	// Create a token source
	ts, err := idtoken.NewTokenSource(ctx, "https://cgr.dev")
	if err != nil {
		log.Fatal(err)
	}

	// Create the keychain
	kc := ggcr.Keychain("my-identity-uidp", ts)

	// Use with context for operations
	ref, err := name.ParseReference("cgr.dev/my/image:latest")
	if err != nil {
		log.Fatal(err)
	}
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(kc), remote.WithContext(ctx))
	if err != nil {
		log.Fatal(err)
	}

	// Work with the image
	_ = img
}

// ExampleTokenSourceKeychain demonstrates how to create a keychain from
// a pre-configured OAuth2 token source.
func ExampleTokenSourceKeychain() {
	// Create a custom token source (this example uses a static token,
	// but in practice you would use a properly configured OAuth2 source)
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: "my-chainguard-token",
	})

	// Create the keychain directly from the token source
	kc := ggcr.TokenSourceKeychain(ts)

	// Use the keychain with go-containerregistry operations
	ref, err := name.ParseReference("cgr.dev/my/image:latest")
	if err != nil {
		log.Fatal(err)
	}
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(kc))
	if err != nil {
		log.Fatal(err)
	}

	// Work with the image
	_ = img
}

// ExampleKeychain_pushImage demonstrates using the keychain to push an image
// to cgr.dev.
func ExampleKeychain_pushImage() {
	ctx := context.Background()

	// Create a token source
	ts, err := idtoken.NewTokenSource(ctx, "https://cgr.dev")
	if err != nil {
		log.Fatal(err)
	}

	// Create the keychain
	kc := ggcr.Keychain("my-identity-uidp", ts)

	// Pull an image from one location
	srcRef, err := name.ParseReference("cgr.dev/source/image:latest")
	if err != nil {
		log.Fatal(err)
	}
	img, err := remote.Image(srcRef, remote.WithAuthFromKeychain(kc))
	if err != nil {
		log.Fatal(err)
	}

	// Push it to another location
	dstRef, err := name.ParseReference("cgr.dev/destination/image:latest")
	if err != nil {
		log.Fatal(err)
	}
	if err := remote.Write(dstRef, img, remote.WithAuthFromKeychain(kc)); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Image pushed successfully")
}
