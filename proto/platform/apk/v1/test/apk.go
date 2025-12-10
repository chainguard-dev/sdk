/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	apk "chainguard.dev/sdk/proto/platform/apk/v1"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
)

var _ apk.Clients = (*MockAPKClients)(nil)

type MockAPKClients struct {
	OnClose error

	APKClient MockAPKClient
}

func (m MockAPKClients) APK() apk.APKClient {
	return &m.APKClient
}

func (m MockAPKClients) Close() error {
	return m.OnClose
}

var _ apk.APKClient = (*MockAPKClient)(nil)

type MockAPKClient struct {
	apk.APKClient

	OnListAPKs []APKsOnList
}

type APKsOnList struct {
	Given *apk.APKFilter
	List  *apk.APKList
	Error error
}

func (m MockAPKClient) ListAPKs(_ context.Context, given *apk.APKFilter, _ ...grpc.CallOption) (*apk.APKList, error) {
	for _, o := range m.OnListAPKs {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
