/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	matcher "chainguard.dev/sdk/proto/platform/matcher/v1"
)

var _ matcher.Clients = (*MockImageMatcherClients)(nil)

type MockImageMatcherClients struct {
	OnClose error

	ImageMatcherClient MockImageMatcherClient
}

func (m MockImageMatcherClients) ImageMatcher() matcher.ImageMatcherClient {
	return &m.ImageMatcherClient
}

func (m MockImageMatcherClients) Close() error {
	return m.OnClose
}

var _ matcher.ImageMatcherClient = (*MockImageMatcherClient)(nil)

type MockImageMatcherClient struct {
	OnMatchImage []MatchImageOnMatch
}

type MatchImageOnMatch struct {
	Given    *matcher.MatchImageRequest
	Response *matcher.MatchedImages
	Error    error
}

func (m MockImageMatcherClient) MatchImage(_ context.Context, given *matcher.MatchImageRequest, _ ...grpc.CallOption) (*matcher.MatchedImages, error) { //nolint: revive
	for _, o := range m.OnMatchImage {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Response, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
