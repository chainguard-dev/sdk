/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"go.uber.org/zap"
	"knative.dev/pkg/logging"
)

var (
	ErrTokenRejected                    = errors.New("token rejected by AWS STS endpoint")
	ErrTokenExpired                     = errors.New("token expired")
	ErrInvalidAudience                  = errors.New("audience header in token does not match expected audience")
	ErrInvalidIdentity                  = errors.New("identity header in token does not match expected identity")
	ErrInvalidEncoding                  = errors.New("invalid token encoding")
	ErrInvalidVerificationConfiguration = errors.New("verifcation was incorrectly configured")
)

type VerifiedClaims struct {
	UserID  string `json:"UserId"`
	Arn     string `json:"Arn"`
	Account string `json:"Account"`
}

func VerifyToken(ctx context.Context, token string, opts ...VerifyOption) (*VerifiedClaims, error) {
	conf, err := newConfigFromOptions(opts...)
	if err != nil {
		logging.FromContext(ctx).Errorw("invalid verification configuration", zap.Error(err))
		return nil, ErrInvalidVerificationConfiguration
	}

	var req *http.Request
	{
		decoded, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			return nil, ErrInvalidEncoding
		}
		r := bufio.NewReader(bytes.NewReader(decoded))
		req, err = http.ReadRequest(r)
		if err != nil {
			return nil, ErrInvalidEncoding
		}

		// NB: If RequestURI is set on an outbound http.Request the client will error. This field
		// is set because of how we serial and then parse the request so we need to unset it here.
		req.RequestURI = ""
		req.URL, err = url.Parse(conf.stsURL)
		if err != nil {
			logging.FromContext(ctx).Errorw("invalid verification configuration. invalid sts url",
				"sts_url", conf.stsURL,
				zap.Error(err))
			return nil, ErrInvalidVerificationConfiguration
		}
	}

	if got := req.Header.Get(audHeader); !conf.allowedAudiences.Has(got) {
		logging.FromContext(ctx).Warnw("verification failed with audience mismatch",
			"received", got,
		)
		return nil, ErrInvalidAudience
	}
	if got := req.Header.Get(idHeader); got != conf.identity {
		logging.FromContext(ctx).Warnw("verification failed with identity mismatch",
			"wanted", conf.identity,
			"received", got,
		)
		return nil, ErrInvalidIdentity
	}

	timestamp, err := time.Parse("20060102T150405Z", req.Header.Get("X-Amz-Date"))
	if err != nil {
		logging.FromContext(ctx).Warnw("verification failed because of a poorly formatted x-amz-date header format", zap.Error(err))
		return nil, ErrInvalidEncoding
	}
	expiry, now := timestamp.Add(15*time.Minute), conf.time()
	if expiry.Before(now) {
		// According to AWS docs
		// > The signed portions (using AWS Signatures) of requests are valid within 15 minutes of the timestamp in the request.
		// If the signature timestamp is already older than 15 minutes the token is expired and we reject it.
		// c.f https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-authenticating-requests.html
		logging.FromContext(ctx).Warnw("verification failed because of expired token")
		return nil, ErrTokenExpired
	}

	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		logging.FromContext(ctx).Warnw("verification failed because of failure to make AWS STS request", zap.Error(err))
		return nil, fmt.Errorf("failed to reach AWS STS endpoint: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body from AWS STS endpoint: %w", err)
		}
		logging.FromContext(ctx).Warnw("verification failed because it was rejected by AWS STS endpoint",
			"response_code", resp.StatusCode,
			"response", string(body),
		)
		return nil, ErrTokenRejected
	}

	var response struct {
		GetCallerIdentityResponse struct {
			GetCallerIdentityResult VerifiedClaims
		}
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		logging.FromContext(ctx).Warnw("verification failed because json parsing err in response", zap.Error(err))
		return nil, fmt.Errorf("failed to parse json from AWS STS response %w", err)
	}

	return &response.GetCallerIdentityResponse.GetCallerIdentityResult, nil
}
