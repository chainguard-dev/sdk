// Package ggcr provides a go-containerregistry authn.Keychain for the cgr.dev registry.
package ggcr

import (
	"context"
	"fmt"

	"chainguard.dev/sdk/sts"
	"github.com/google/go-containerregistry/pkg/authn"
	"golang.org/x/oauth2"
)

const issuer = "https://issuer.enforce.dev"

// Keychain returns an authn.Keychain used to authorize requests to the cgr.dev registry using go-containerregistry.
//
// It takes the identity UIDP to assume, and a token source to obtain the token to exchange.
//
// This can be used with google.golang.org/api/idtoken.NewTokenSource to exchange ambient GCP credentials for Chainguard tokens:
//
//	ts, err := idtoken.NewTokenSource(ctx, "https://cgr.dev")
//	kc := ggcr.Keychain("my-identity", ts)
//
// This keychain can then be used to pull images from the cgr.dev registry:
//
//	img, err := remote.Image("cgr.dev/my/image", remote.WithAuth(kc))
func Keychain(identity string, ts oauth2.TokenSource) authn.Keychain {
	return cgKeychain{identity, ts}
}

type cgKeychain struct {
	identity string
	ts       oauth2.TokenSource
}

func (k cgKeychain) Resolve(res authn.Resource) (authn.Authenticator, error) {
	if res.RegistryStr() != "cgr.dev" {
		return authn.Anonymous, nil
	}

	ctx := context.Background()
	exch := sts.New(issuer, res.RegistryStr(), sts.WithIdentity(k.identity))

	tok, err := k.ts.Token()
	if err != nil {
		return nil, fmt.Errorf("getting token: %w", err)
	}
	cgtok, err := exch.Exchange(ctx, tok.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("exchanging token: %w", err)
	}
	return &authn.Basic{
		Username: "_token",
		Password: cgtok,
	}, nil
}
