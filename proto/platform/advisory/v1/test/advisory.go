/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	advisory "chainguard.dev/sdk/proto/platform/advisory/v1"
)

var _ advisory.Clients = (*MockSecurityAdvisoryClients)(nil)

type MockSecurityAdvisoryClients struct {
	OnClose error

	SecurityAdvisoryClient MockSecurityAdvisoryClient
}

func (m MockSecurityAdvisoryClients) SecurityAdvisory() advisory.SecurityAdvisoryClient {
	return &m.SecurityAdvisoryClient
}

func (m MockSecurityAdvisoryClients) Close() error {
	return m.OnClose
}

var _ advisory.SecurityAdvisoryClient = (*MockSecurityAdvisoryClient)(nil)

type MockSecurityAdvisoryClient struct {
	OnListDocuments []DocumentsOnList
}

type DocumentsOnList struct {
	Given *advisory.DocumentFilter
	List  *advisory.DocumentList
	Error error
}

func (m MockSecurityAdvisoryClient) ListDocuments(_ context.Context, given *advisory.DocumentFilter, _ ...grpc.CallOption) (*advisory.DocumentList, error) { //nolint: revive
	for _, o := range m.OnListDocuments {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
