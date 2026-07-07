/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws

import (
	"errors"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
)

type verifyConf struct {
	allowedAudiences sets.Set[string]
	identity         string

	time             func() time.Time
	stsURL           string
	expectedSTSHost  string
	allowedClockSkew time.Duration
}

const (
	defaultAudience        = "https://issuer.enforce.dev"
	defaultExpectedSTSHost = "sts.amazonaws.com"
	// stsQuery is the GetCallerIdentity query string shared by the default STS
	// URL and any host-derived one, so the endpoint we connect to and the host
	// the signature must cover cannot silently drift.
	stsQuery      = "/?Action=GetCallerIdentity&Version=2011-06-15"
	defaultSTSURL = "https://" + defaultExpectedSTSHost + stsQuery
)

func newDefaultConfig() *verifyConf {
	return &verifyConf{
		allowedAudiences: sets.New(defaultAudience),
		stsURL:           defaultSTSURL,
		expectedSTSHost:  defaultExpectedSTSHost,
		allowedClockSkew: defaultAllowedClockSkew,
		time:             time.Now,
	}
}

func newConfigFromOptions(opts ...VerifyOption) (*verifyConf, error) {
	conf := newDefaultConfig()
	for _, o := range opts {
		o(conf)
	}
	if err := conf.valid(); err != nil {
		return nil, err
	}
	return conf, nil
}

func (c *verifyConf) valid() error {
	if c.identity == "" {
		return errors.New("must specify assume identity to verify")
	}
	if c.allowedAudiences.Len() == 0 {
		return errors.New("must specify audience to verify")
	}
	if c.stsURL == "" {
		return errors.New("must specify AWS STS endpoint url")
	}
	if c.expectedSTSHost == "" {
		return errors.New("must specify expected AWS STS host")
	}
	if c.allowedClockSkew < 0 {
		return errors.New("allowed clock skew must not be negative")
	}
	return nil
}

type VerifyOption func(*verifyConf)

func WithAudience(aud sets.Set[string]) VerifyOption {
	return func(c *verifyConf) {
		c.allowedAudiences = aud
	}
}

func WithIdentity(id string) VerifyOption {
	return func(c *verifyConf) {
		c.identity = id
	}
}

// WithExpectedSTSHost sets the STS host a token's GetCallerIdentity request must
// be addressed to: both the host required in the signed request (compared
// against the request's Host header) and the endpoint VerifyToken forwards to.
// It defaults to "sts.amazonaws.com"; override it to verify identities whose
// tokens are signed against a regional, FIPS, or GovCloud/China-partition STS
// endpoint (all of which use https://<host>/?Action=GetCallerIdentity&Version=2011-06-15).
//
// The value must be a lowercase, port-less hostname with no trailing dot,
// matching the host as it appears in the SigV4 canonical request. It assumes the
// forwarded endpoint host equals the signed host; deployments where they differ
// (PrivateLink/VPC-interface endpoints, egress proxies) are not yet supported.
func WithExpectedSTSHost(host string) VerifyOption {
	return func(c *verifyConf) {
		c.expectedSTSHost = host
		c.stsURL = "https://" + host + stsQuery
	}
}

// WithAllowedClockSkew sets how far a token's X-Amz-Date may be ahead of the
// verifier's clock before it is rejected as future-dated. It defaults to
// defaultAllowedClockSkew (5 minutes, matching AWS SigV4's tolerance). Increase
// it for environments with looser clock synchronization.
func WithAllowedClockSkew(skew time.Duration) VerifyOption {
	return func(c *verifyConf) {
		c.allowedClockSkew = skew
	}
}

// Unexported testing options

func withSTSURL(url string) VerifyOption {
	return func(c *verifyConf) {
		c.stsURL = url
	}
}

func withTimestamp(t time.Time) VerifyOption {
	return func(c *verifyConf) {
		c.time = func() time.Time {
			return t
		}
	}
}
