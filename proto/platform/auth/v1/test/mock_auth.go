/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"
	"google.golang.org/protobuf/types/known/emptypb"

	auth "chainguard.dev/sdk/proto/platform/auth/v1"
)

var _ auth.AuthClient = (*MockAuthClient)(nil)

type MockAuthClient struct {
	OnValidate []AuthOnValidate
	OnRegister []AuthOnRegister
}

type FromContextFn func(context.Context) bool

type AuthOnValidate struct {
	Given    FromContextFn
	Validate *auth.WhoAmI
	Error    error
}

type AuthOnRegister struct {
	Given   *auth.RegistrationRequest
	Created *auth.Session
	Error   error
}

func (m MockAuthClient) Validate(ctx context.Context, _ *emptypb.Empty, _ ...grpc.CallOption) (*auth.WhoAmI, error) {
	for _, o := range m.OnValidate {
		if o.Given(ctx) {
			return o.Validate, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for context: %v", ctx)
}

func (m MockAuthClient) Register(_ context.Context, given *auth.RegistrationRequest, _ ...grpc.CallOption) (*auth.Session, error) {
	for _, o := range m.OnRegister {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Created, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

// --- Server ---

type MockAuthServer struct {
	auth.UnimplementedAuthServer
	Client MockAuthClient
}

func (m MockAuthServer) Validate(ctx context.Context, empty *emptypb.Empty) (*auth.WhoAmI, error) {
	return m.Client.Validate(ctx, empty)
}

func (m MockAuthServer) Register(ctx context.Context, req *auth.RegistrationRequest) (*auth.Session, error) {
	return m.Client.Register(ctx, req)
}
