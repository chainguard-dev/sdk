/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test_test

import (
	"context"
	"errors"
	"fmt"

	"google.golang.org/protobuf/types/known/emptypb"

	auth "chainguard.dev/sdk/proto/platform/auth/v1"
	"chainguard.dev/sdk/proto/platform/auth/v1/test"
)

func ExampleMockAuthClient_Validate() {
	mock := test.MockAuthClient{
		OnValidate: []test.AuthOnValidate{{
			Given: func(_ context.Context) bool { return true },
			Validate: &auth.WhoAmI{
				Subject: "user-123",
			},
		}},
	}

	whoami, err := mock.Validate(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("identity: %s\n", whoami.Subject)
	// Output: identity: user-123
}

func ExampleMockAuthClient_Validate_contextMatching() {
	type contextKey string
	const userKey contextKey = "user"

	mock := test.MockAuthClient{
		OnValidate: []test.AuthOnValidate{{
			Given: func(ctx context.Context) bool {
				return ctx.Value(userKey) == "admin"
			},
			Validate: &auth.WhoAmI{
				Subject: "admin-456",
			},
		}, {
			Given: func(ctx context.Context) bool {
				return ctx.Value(userKey) == "user"
			},
			Validate: &auth.WhoAmI{
				Subject: "user-789",
			},
		}},
	}

	ctx := context.WithValue(context.Background(), userKey, "admin")
	whoami, err := mock.Validate(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("identity: %s\n", whoami.Subject)
	// Output: identity: admin-456
}

func ExampleMockAuthClient_Register() {
	mock := test.MockAuthClient{
		OnRegister: []test.AuthOnRegister{{
			Given: &auth.RegistrationRequest{
				Code: "invite-code-123",
			},
			Created: &auth.Session{
				Identity: "user-identity-abc",
				Group:    "group-123",
			},
		}},
	}

	session, err := mock.Register(context.Background(), &auth.RegistrationRequest{
		Code: "invite-code-123",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("identity: %s\n", session.Identity)
	// Output: identity: user-identity-abc
}

func ExampleMockAuthClient_Register_withContextCheck() {
	type contextKey string
	const authKey contextKey = "authorized"

	mock := test.MockAuthClient{
		OnRegister: []test.AuthOnRegister{{
			Given: &auth.RegistrationRequest{
				Code: "invite-code-456",
			},
			CheckContext: func(ctx context.Context) bool {
				return ctx.Value(authKey) == true
			},
			Created: &auth.Session{
				Identity: "user-identity-xyz",
				Group:    "group-456",
			},
		}},
	}

	ctx := context.WithValue(context.Background(), authKey, true)
	session, err := mock.Register(ctx, &auth.RegistrationRequest{
		Code: "invite-code-456",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("identity: %s\n", session.Identity)
	// Output: identity: user-identity-xyz
}

func ExampleMockAuthClient_GetHeadlessSession() {
	mock := test.MockAuthClient{
		OnGetHeadlessSession: []test.AuthOnGetHeadlessSession{{
			Given: &auth.GetHeadlessSessionRequest{
				Code: "headless-code-123",
			},
			Found: &auth.HeadlessSession{
				EcdhPublicKey:    []byte("public-key-data"),
				EncryptedIdtoken: []byte("encrypted-token-data"),
			},
		}},
	}

	session, err := mock.GetHeadlessSession(context.Background(), &auth.GetHeadlessSessionRequest{
		Code: "headless-code-123",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("session has public key: %v\n", len(session.EcdhPublicKey) > 0)
	// Output: session has public key: true
}

func ExampleMockAuthClient_Delete() {
	mock := test.MockAuthClient{
		OnDelete: []test.AuthOnDelete{{
			Given: &auth.DeletionRequest{
				Id: "session-456",
			},
		}},
	}

	_, err := mock.Delete(context.Background(), &auth.DeletionRequest{
		Id: "session-456",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println("deleted successfully")
	// Output: deleted successfully
}

func ExampleMockAuthClient_Delete_error() {
	mock := test.MockAuthClient{
		OnDelete: []test.AuthOnDelete{{
			Given: &auth.DeletionRequest{
				Id: "session-789",
			},
			Error: errors.New("session not found"),
		}},
	}

	_, err := mock.Delete(context.Background(), &auth.DeletionRequest{
		Id: "session-789",
	})
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// Output: error: session not found
}

func ExampleMockAuthServer() {
	server := test.MockAuthServer{
		Client: test.MockAuthClient{
			OnValidate: []test.AuthOnValidate{{
				Given: func(_ context.Context) bool { return true },
				Validate: &auth.WhoAmI{
					Subject: "server-user-123",
				},
			}},
		},
	}

	whoami, err := server.Validate(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Printf("identity: %s\n", whoami.Subject)
	// Output: identity: server-user-123
}
