/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package aws

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

var timeNow = time.Now

const (
	audHeader = `Chainguard-Audience`
	idHeader  = `Chainguard-Identity`
)

// emptyBodySHA256 is the SHA-256 of an empty body, hex-encoded — the payload
// hash for the bodyless GetCallerIdentity request.
const emptyBodySHA256 = `e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855` //nolint:gosec

// bindingHeaders are the Chainguard headers the generator signs to bind a token
// to an audience and identity. The verifier requires the same set to be covered
// by the signature; deriving both sides from this one list prevents drift.
var bindingHeaders = []string{audHeader, idHeader}

// GenerateToken creates token using the supplied AWS credentials that can prove the user's AWS identity. Audience and identity are
// the Chainguard STS url (e.g https://issuer.enforce.dev) and the UID of the Chainguard assumable identity to assume via STS.
func GenerateToken(ctx context.Context, creds aws.Credentials, audience, identity string) (string, error) {
	req, err := http.NewRequest("POST", "https://sts.amazonaws.com", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create new HTTP request: %w", err)
	}

	req.URL.Path = "/"
	req.URL.RawQuery = "Action=GetCallerIdentity&Version=2011-06-15"
	req.Header.Add("Accept", "application/json")
	req.Header.Add(audHeader, audience)
	req.Header.Add(idHeader, identity)

	const (
		// STS service details for signing
		svc    = `sts`
		region = `us-east-1`
	)

	err = v4.NewSigner().SignHTTP(ctx, creds, req, emptyBodySHA256, svc, region, timeNow())
	if err != nil {
		return "", fmt.Errorf("failed to sign GetCallerIdentity request with AWS credentials: %w", err)
	}

	var b bytes.Buffer
	err = req.Write(&b)
	if err != nil {
		return "", fmt.Errorf("failed to serialize GetCallerIdentity HTTP request to buffer: %w", err)
	}

	return base64.URLEncoding.EncodeToString(b.Bytes()), nil
}
