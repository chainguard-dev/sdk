/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package headless

import (
	"strings"
	"testing"
)

func TestVerifyCode(t *testing.T) {
	for _, tt := range []struct {
		name    string
		code    string
		wantErr string
	}{
		{
			name:    "invalid public key",
			code:    "yolo",
			wantErr: "invalid public key",
		}, {
			name:    "not base64",
			code:    "y@l@",
			wantErr: "illegal base64",
		}, {
			name: "valid code",
			code: "uFK0TyUaZCMB6QBGC-HVKGJKqI4-U9q5rV9k8I_YUQ4",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			err := VerifyCode(tt.code)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error, see none")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestNewSession(t *testing.T) {
	for _, tt := range []struct {
		name    string
		code    Code
		wantErr string
	}{{
		name:    "invalid public key",
		code:    "yolo",
		wantErr: "invalid public key",
	}, {
		name:    "not base64",
		code:    "y@l@",
		wantErr: "illegal base64",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.code.NewSession([]byte("idtoken"))
			if err == nil && tt.wantErr != "" {
				t.Fatal("expected error, see none")
			}
			if err != nil && tt.wantErr == "" {
				t.Fatalf("unexpected error: %v", err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncrypt(t *testing.T) {
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
			_, err := symmetricEncrypt([]byte("token"), tt.key)
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
