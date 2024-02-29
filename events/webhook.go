/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package events

// WebhookCustomClaims holds the custom claims embedded in the webhook's
// OIDC authorization token.
type WebhookCustomClaims struct {
	Webhook struct {
		Digest string `json:"digest"`
	} `json:"webhook"`
}
