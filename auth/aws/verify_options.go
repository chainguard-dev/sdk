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

	time   func() time.Time
	stsURL string
}

func newDefaultConfig() *verifyConf {
	return &verifyConf{
		allowedAudiences: sets.New("https://issuer.enforce.dev"),
		stsURL:           "https://sts.amazonaws.com/?Action=GetCallerIdentity&Version=2011-06-15",
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
