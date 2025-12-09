/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	oidc "chainguard.dev/sdk/proto/platform/oidc/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ oidc.SecurityTokenServiceClient = (*MockSTSClient)(nil)

type MockSTSClient struct {
	OnExchange       []STSOnExchange
	OnExchangeTokens []STSOnExchangeTokens
	OnGetAccessToken []STSOnGetAccessToken
}

func (m MockSTSClient) Exchange(_ context.Context, given *oidc.ExchangeRequest, _ ...grpc.CallOption) (*oidc.RawToken, error) {
	for _, o := range m.OnExchange {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Exchanged, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSTSClient) ExchangeTokens(_ context.Context, given *oidc.ExchangeRequest, _ ...grpc.CallOption) (*oidc.TokenPair, error) {
	for _, o := range m.OnExchangeTokens {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Exchanged, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockSTSClient) ExchangeRefreshToken(_ context.Context, given *oidc.ExchangeRefreshTokenRequest, _ ...grpc.CallOption) (*oidc.TokenPair, error) {
	for _, o := range m.OnGetAccessToken {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Exchanged, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

type STSOnExchange struct {
	Given     *oidc.ExchangeRequest
	Exchanged *oidc.RawToken
	Error     error
}

type STSOnExchangeTokens struct {
	Given     *oidc.ExchangeRequest
	Exchanged *oidc.TokenPair
	Error     error
}

type STSOnGetAccessToken struct {
	Given     *oidc.ExchangeRefreshTokenRequest
	Exchanged *oidc.TokenPair
	Error     error
}
