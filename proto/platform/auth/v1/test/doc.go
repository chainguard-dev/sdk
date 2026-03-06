/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package test provides mock implementations of the auth service for testing.

# Overview

This package contains mock implementations of auth.AuthClient and auth.AuthServer
that can be used in unit tests to simulate authentication behavior without requiring
a real auth service. The mocks support configurable responses for all auth operations.

# Features

  - MockAuthClient: A configurable mock implementation of auth.AuthClient
  - MockAuthServer: A mock server that wraps MockAuthClient for gRPC testing
  - Flexible response configuration using slice-based matchers
  - Context-aware validation matching via FromContextFn
  - Protocol buffer comparison using protocmp for request matching

# Usage

The mock client is configured by populating slices of response configurations.
Each operation (Validate, Register, GetHeadlessSession, Delete) has a corresponding
slice that defines the expected inputs and desired outputs.

Basic mock setup:

	mock := test.MockAuthClient{
		OnValidate: []test.AuthOnValidate{{
			Given: func(ctx context.Context) bool { return true },
			Validate: &auth.WhoAmI{
				Identity: &auth.Identity{Id: "user-123"},
			},
		}},
	}

	whoami, err := mock.Validate(ctx, &emptypb.Empty{})

Request matching:

	mock := test.MockAuthClient{
		OnRegister: []test.AuthOnRegister{{
			Given: &auth.RegistrationRequest{
				Email: "user@example.com",
			},
			Created: &auth.Session{
				Token: "session-token",
			},
		}},
	}

	session, err := mock.Register(ctx, &auth.RegistrationRequest{
		Email: "user@example.com",
	})

Error simulation:

	mock := test.MockAuthClient{
		OnDelete: []test.AuthOnDelete{{
			Given: &auth.DeletionRequest{Id: "session-123"},
			Error: errors.New("deletion failed"),
		}},
	}

	_, err := mock.Delete(ctx, &auth.DeletionRequest{Id: "session-123"})
	// err will be "deletion failed"

# Integration Patterns

For testing gRPC servers, use MockAuthServer:

	server := test.MockAuthServer{
		Client: test.MockAuthClient{
			OnValidate: []test.AuthOnValidate{{
				Given: func(ctx context.Context) bool { return true },
				Validate: &auth.WhoAmI{
					Identity: &auth.Identity{Id: "user-123"},
				},
			}},
		},
	}

	// Register with grpc.Server
	auth.RegisterAuthServer(grpcServer, server)

The mock uses protocol buffer comparison (protocmp.Transform()) for matching
requests, ensuring that semantically equivalent messages match even if they
differ in internal representation.
*/
package test
