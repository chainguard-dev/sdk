/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
)

type MockOIDCClient struct {
	OnClose error

	STSClient MockSTSClient
}

var _ oidc.Clients = (*MockOIDCClient)(nil)

func (m MockOIDCClient) STS() oidc.SecurityTokenServiceClient {
	return &m.STSClient
}

func (m MockOIDCClient) Close() error {
	return m.OnClose
}
