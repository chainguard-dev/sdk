/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package headless

import (
	"crypto/ecdh"
	"fmt"
	"strconv"
	"strings"
	"testing"

	auth "chainguard.dev/sdk/proto/platform/auth/v1"
)

func TestTokenRoundTrip(t *testing.T) {
	for _, tt := range []struct {
		name  string
		token []byte
	}{
		{
			name:  "common token",
			token: testIDToken(4000), // 4KB tokens are relatively common
		},
		{
			name:  "short",
			token: testIDToken(100), // Not common, but we want to ensure small tokens are ok
		},
		{
			name:  "short",
			token: testIDToken(50000), // Not common, but we want to ensure large tokens are ok
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			pk, err := GenerateKeyPair()
			if err != nil {
				t.Fatal(err)
			}

			code := NewCode(pk.PublicKey())
			fmt.Println("code", code)
			session, err := code.NewSession(tt.token)
			if err != nil {
				t.Fatal(err)
			}

			decrypted, err := DecryptIDToken(session, pk)
			if err != nil {
				t.Fatal(err)
			}

			if string(tt.token) != string(decrypted) {
				t.Fatalf("expected %s, got %s", tt.token, decrypted)
			}
		})
	}
}

func TestDecryptToken_invalid(t *testing.T) {
	validPK, err := GenerateKeyPair()
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name            string
		headlessSession *auth.HeadlessSession
		privateKey      *ecdh.PrivateKey
		err             string
	}{
		{
			name: "invalid public key",
			headlessSession: &auth.HeadlessSession{
				EcdhPublicKey:    []byte("invalid"),
				EncryptedIdtoken: []byte("encrypted"),
			},
			privateKey: validPK,
			err:        "invalid public key",
		},
		{
			name: "invalid encrypted token, mangled",
			headlessSession: func() *auth.HeadlessSession {
				code := NewCode(validPK.PublicKey())
				session, err := code.NewSession([]byte("token"))
				if err != nil {
					t.Fatal(err)
				}
				mangled := "mangled_" + string(session.EncryptedIdtoken)
				session.EncryptedIdtoken = []byte(mangled)
				return session
			}(),
			privateKey: validPK,
			err:        "failed to decrypt id token",
		},
		{
			name: "invalid encrypted token, truncated",
			headlessSession: func() *auth.HeadlessSession {
				code := NewCode(validPK.PublicKey())
				session, err := code.NewSession([]byte("token"))
				if err != nil {
					t.Fatal(err)
				}
				session.EncryptedIdtoken = []byte("truncated")
				return session
			}(),
			privateKey: validPK,
			err:        "failed to decrypt id token",
		},
		{
			name: "invalid server public key",
			headlessSession: &auth.HeadlessSession{
				EcdhPublicKey:    []byte("foo"),
				EncryptedIdtoken: []byte("bar"),
			},
			privateKey: nil,
			err:        "invalid public key",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DecryptIDToken(tt.headlessSession, tt.privateKey)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.err) {
				t.Fatalf("expected error %s, got %s", tt.err, err)
			}
		})
	}
}

func testIDToken(length int) []byte {
	idtoken := ""

	// We want a pretty long token to ensure that the encryption/decryption round trip
	// does	not truncate the token.
	for i := 0; len(idtoken) < length; i++ {
		// it's ok to not use StringBuilder for testing
		idtoken += "token " + strconv.Itoa(i) + "\n"
	}
	return []byte(idtoken)
}

func TestDecrypt(t *testing.T) {
	for _, tt := range []struct {
		name    string
		key     []byte
		wantErr string
	}{{
		name:    "invalid key size",
		key:     []byte("invalid"),
		wantErr: "invalid key size",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := symmetricDecrypt([]byte("token"), tt.key)
			if err == nil && tt.wantErr != "" {
				t.Fatalf("expected error, got none")
			}
			if err != nil && tt.wantErr == "" {
				t.Fatalf("unexpected error %s", err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error %s, got %s", tt.wantErr, err)
			}
		})
	}
}
