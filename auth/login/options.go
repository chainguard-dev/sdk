/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package login

import (
	"errors"
	"fmt"

	"chainguard.dev/sdk/uidp"
)

type config struct {
	// Audience is the intended audience for the token, if different than what is
	// configured for the oauth handler.
	Audience []string

	// URL of Chainguard Enforce OIDC Issuer. Defaults to https://issuer.enforce.dev
	Issuer string

	// Identity is the exact UIDP of an assumable identity to authenticate as.
	Identity string

	// UIDP of specific customer identity provider to log in with
	IDP string

	// OrgName is the name of an organization with custom identity provider configured to use for authentication
	OrgName string

	// Optional invite code to consume on login
	InviteCode string

	// ClientID is the ID of the oauth2 provider
	ClientID string

	// Auth0Connection is the connection parameter sent to Auth0
	// to preselect the social connection.
	//
	// See docs for details: https://auth0.com/docs/api/authentication#social
	Auth0Connection string

	// SkipRegistration tells the issuer not to attempt to register
	// if the account is not found.
	SkipRegistration bool

	// IncludeUpstreamToken tells the issuer to include the encrypted upstream token
	// in the Chainguard token
	IncludeUpstreamToken bool

	// CreateRefreshToken tells the issuer to create a refresh token
	CreateRefreshToken bool

	// SkipBrowser avoids opening a browser window for login, just print out the url
	SkipBrowser bool
}

const defaultIssuer = `https://issuer.enforce.dev`

func newDefaultConfig() *config {
	return &config{
		Issuer: defaultIssuer,
	}
}

func newConfigFromOptions(opts ...Option) (*config, error) {
	conf := newDefaultConfig()
	for _, o := range opts {
		o(conf)
	}
	if err := conf.valid(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *config) valid() error {
	if c.ClientID != "" && (c.IDP != "" || c.OrgName != "") {
		return errors.New("must specify one of client id or custom identity provider")
	}
	if c.IDP != "" && c.OrgName != "" {
		return errors.New("must specify one of identity provider id or organization name")
	}
	if c.ClientID == "" && c.IDP == "" && c.OrgName == "" {
		return errors.New("must select one of client id, custom identity provider and organization name")
	}

	switch {
	case c.IDP != "":
		if !uidp.Valid(c.IDP) {
			return errors.New("invalid ID for identity provider")
		}
		if c.Auth0Connection != "" {
			return errors.New("identity provider ID and Auth0 connection are mutually exclusive")
		}
		return nil

	case c.OrgName != "":
		if c.Auth0Connection != "" {
			return errors.New("organization name and Auth0 connection are mutually exclusive")
		}

		verified, err := orgCheck(c.OrgName, c.Issuer)
		if err != nil {
			return fmt.Errorf("error checking if organization is verified: %w", err)
		}
		if !verified {
			return errors.New("organization is not verified must use identity provider ID to log in")
		}

		return nil

	default:
		return nil
	}
}

type Option func(opt *config)

func WithAudience(aud []string) Option {
	return func(c *config) {
		c.Audience = aud
	}
}

func WithClientID(id string) Option {
	return func(c *config) {
		c.ClientID = id
	}
}
func WithIssuer(issuer string) Option {
	return func(c *config) {
		c.Issuer = issuer
	}
}

func WithIdentity(identity string) Option {
	return func(c *config) {
		c.Identity = identity
	}
}

func WithIdentityProvider(idp string) Option {
	return func(c *config) {
		c.IDP = idp
	}
}

func WithOrgName(org string) Option {
	return func(c *config) {
		c.OrgName = org
	}
}

func WithInviteCode(inviteCode string) Option {
	return func(c *config) {
		c.InviteCode = inviteCode
	}
}

func WithAuth0Connection(conn string) Option {
	return func(c *config) {
		c.Auth0Connection = conn
	}
}

func WithSkipRegistration() Option {
	return func(c *config) {
		c.SkipRegistration = true
	}
}

func WithIncludeUpstreamToken() Option {
	return func(c *config) {
		c.IncludeUpstreamToken = true
	}
}

func WithCreateRefreshToken() Option {
	return func(c *config) {
		c.CreateRefreshToken = true
	}
}

func WithSkipBrowser() Option {
	return func(c *config) {
		c.SkipBrowser = true
	}
}
