/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package aws provides utilities for generating and verifying Chainguard
// tokens backed by AWS identity credentials via the AWS STS
// GetCallerIdentity API.
//
// VerifyToken treats a token as proof of an AWS identity only when the
// audience and identity binding headers — along with the host and date
// headers it relies on — are covered by the request's SigV4 signature, and
// when the request targets the expected STS host with a fresh, non-future
// timestamp. This prevents a GetCallerIdentity request signed for another
// relying party from being replayed against this verifier.
package aws
