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
					t.Fatal("got nil error, want error")
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
				t.Fatal("got nil error, want error")
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
				t.Fatalf("got nil error, want error")
			}
			if err != nil && tt.wantErr == "" {
				t.Fatalf("unexpected error %s", err)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("got error %s, want error containing %s", err, tt.wantErr)
			}
		})
	}
}

func TestConfirmationCode(t *testing.T) {
	for _, tt := range []struct {
		name    string
		code    Code
		want    string
		wantErr string
	}{{
		name: "valid code",
		code: "uFK0TyUaZCMB6QBGC-HVKGJKqI4-U9q5rV9k8I_YUQ4",
		want: "XRPMCQYW",
	}, {
		// Same key as above: the last base64 character has 2 unused bits.
		name: "other encoding of the same key",
		code: "uFK0TyUaZCMB6QBGC-HVKGJKqI4-U9q5rV9k8I_YUQ5",
		want: "XRPMCQYW",
	}, {
		name:    "invalid public key",
		code:    "yolo",
		wantErr: "invalid public key",
	}, {
		name:    "not base64",
		code:    "y@l@",
		wantErr: "illegal base64",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.code.ConfirmationCode()
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("got nil error, want error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("got %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfirmationCode_Unique(t *testing.T) {
	seen := make(map[string]struct{}, 5000)
	for range 5000 {
		pk, err := GenerateKeyPair()
		if err != nil {
			t.Fatal(err)
		}
		code := NewCode(pk.PublicKey())
		got, err := code.ConfirmationCode()
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := seen[got]; ok {
			t.Fatalf("duplicate confirmation code %q", got)
		}
		seen[got] = struct{}{}
	}
}

func TestMatchesConfirmationCode(t *testing.T) {
	for _, tt := range []struct {
		name  string
		code  Code
		input string
		want  bool
	}{
		{name: "exact", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: "K8013G85", want: true},
		{name: "lower case", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: "k8013g85", want: true},
		{name: "dashed and spaced", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: " K801-3G85\n", want: true},
		{name: "look-alike letters", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: "K8OI3G85", want: true},
		{name: "wrong code", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: "K8013G86", want: false},
		{name: "too short", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: "K8013G8", want: false},
		{name: "empty", code: "WeSiyKWll0A665kVZlh-N1VH_Pht07908luLaH8qQCg", input: "", want: false},
		{name: "invalid code", code: "yolo", input: "K8013G85", want: false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.MatchesConfirmationCode(tt.input); got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
