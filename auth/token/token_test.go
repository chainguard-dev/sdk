/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gopkg.in/square/go-jose.v2/jwt"

	"chainguard.dev/go-oidctest/pkg/oidctest"
)

func testToken(t *testing.T, audience, subject string, issuedAt time.Time, validDur time.Duration) string {
	t.Helper()

	signer, iss := oidctest.NewIssuer(t)

	tok, err := jwt.Signed(signer).Claims(jwt.Claims{
		Issuer:   iss,
		IssuedAt: jwt.NewNumericDate(issuedAt),
		Expiry:   jwt.NewNumericDate(issuedAt.Add(validDur)),
		Subject:  subject,
		Audience: jwt.Audience{audience},
	}).CompactSerialize()
	if err != nil {
		t.Fatalf("CompactSerialize() = %v", err)
	}

	return tok
}

func testRefreshToken(t *testing.T, issuer, audience, subject string, issuedAt time.Time, validDur time.Duration) string {
	t.Helper()

	msg := jwt.Claims{
		Issuer:   issuer,
		Subject:  subject,
		Audience: jwt.Audience{audience},
		IssuedAt: jwt.NewNumericDate(issuedAt),
		Expiry:   jwt.NewNumericDate(issuedAt.Add(validDur)),
	}

	bs, _ := json.Marshal(msg)
	code := base64.StdEncoding.EncodeToString(bs)

	return code
}

func TestSave(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		kind     Kind
		audience string
		wantPath string
	}{{
		name:     "sans audience",
		kind:     KindAccess,
		audience: "", // Intentionally blank.
		wantPath: filepath.Join(cacheDir, parentDir, string(KindAccess)),
	}, {
		name:     "with audience (sans replacement)",
		kind:     KindAccess,
		audience: "audience",
		wantPath: filepath.Join(cacheDir, parentDir, "audience", string(KindAccess)),
	}, {
		name:     "with audience (with replacement)",
		kind:     KindRefresh,
		audience: "https://audience.unit.test",
		wantPath: filepath.Join(cacheDir, parentDir, "https:--audience.unit.test", string(KindRefresh)),
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Save a token
			tokenContents := []byte("mytoken")
			if err := Save(tokenContents, test.kind, test.audience); err != nil {
				t.Fatalf("Save() unexpected error=%v", err)
			}

			// Manually check the expected file location
			_, err := os.Stat(test.wantPath)
			if err != nil {
				t.Fatalf("Expected token path returned error: %v", err)
			}

			gotContents, err := os.ReadFile(test.wantPath)
			if err != nil {
				t.Fatalf("Unexpected error reading token filepath=%s, err=%v", test.wantPath, err)
			}

			if diff := cmp.Diff(tokenContents, gotContents); diff != "" {
				t.Errorf("Reading %s returned unexpected results (-want, +got): %s", test.wantPath, diff)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		kind     Kind
		audience string
	}{{
		name:     "sans audience",
		kind:     KindAccess,
		audience: "", // Intentionally blank.
	}, {
		name:     "with audience (sans replacement)",
		kind:     KindAccess,
		audience: "audience",
	}, {
		name:     "with audience (with replacement)",
		kind:     KindRefresh,
		audience: "https://audience.unit.test",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Manually save a token
			tokenContents := []byte("mytoken")
			mutatedAud := strings.ReplaceAll(test.audience, "/", "-")
			path := filepath.Join(cacheDir, parentDir, mutatedAud)
			if err := os.MkdirAll(path, 0777); err != nil {
				t.Fatalf("Unexpected error creating temp dir: %v", err)
			}
			if err := os.WriteFile(filepath.Join(path, string(test.kind)), tokenContents, 0600); err != nil {
				t.Fatalf("Unexpected error writing test token: %v", err)
			}

			// Load the token, check its contents
			got, err := Load(test.kind, test.audience)
			if err != nil {
				t.Fatalf("Load() unexpected error=%v", err)
			}
			if string(tokenContents) != string(got) {
				t.Errorf("Load() return mismatch, want=%s, got=%s", tokenContents, got)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name     string
		audience string
		kind     Kind
		exists   bool
		wantErr  bool
	}{{
		name:     "no audience ",
		audience: "", // Intentionally blank.
		kind:     KindRefresh,
		wantErr:  false,
	}, {
		name:     "token doesn't exist",
		audience: "audience",
		exists:   false,
		kind:     KindAccess,
		wantErr:  false,
	}, {
		name:     "token does exist",
		audience: "audience",
		exists:   true,
		kind:     KindAccess,
		wantErr:  false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Manually save a token if it should exist
			if test.exists {
				tokenContents := []byte("mytoken")
				path := filepath.Join(cacheDir, parentDir, test.audience)
				if err := os.MkdirAll(path, 0777); err != nil {
					t.Fatalf("Unexpected error creating temp dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(path, string(test.kind)), tokenContents, 0600); err != nil {
					t.Fatalf("Unexpected error writing test token: %v", err)
				}
			}

			// Attempt to delete the token.
			err := Delete(test.kind, test.audience)
			if (err != nil) != test.wantErr {
				t.Fatalf("Delete() error return mismatch, want=%t got=%v", test.wantErr, err)
			}
		})
	}
}

func TestDeleteAll(t *testing.T) {
	tests := []struct {
		name           string
		auds           []string // Create sub-dirs for each of these in cacheDir
		audsWithTokens []string // Dirs from above that should also have a token (subset of auds)
		extraFiles     []string // Any extra files outside of audience dirs to include (in root of cacheDir)
		wantErr        bool
	}{{
		name:           "no tokens, no extra files/dirs",
		auds:           nil,
		audsWithTokens: nil,
		extraFiles:     nil,
		wantErr:        false,
	}, {
		name:           "all aud dirs have token, no extra files",
		auds:           []string{"audience1", "audience2"},
		audsWithTokens: []string{"audience1", "audience2"},
		extraFiles:     nil,
		wantErr:        false,
	}, {
		name:           "some empty aud dirs, no extra files",
		auds:           []string{"audience1", "audience2", "audience3", "audience4"},
		audsWithTokens: []string{"audience1", "audience2"},
		extraFiles:     nil,
		wantErr:        false,
	}, {
		name:           "all empty aud dirs, no extra files",
		auds:           []string{"audience1", "audience2", "audience3", "audience4"},
		audsWithTokens: []string{},
		extraFiles:     nil,
		wantErr:        false,
	}, {
		name:           "all aud dirs have token, extra files in root",
		auds:           []string{"audience1", "audience2"},
		audsWithTokens: []string{"audience1", "audience2"},
		extraFiles:     []string{"extra1", "extra2"},
		wantErr:        false,
	}, {
		name:           "all aud dirs have token, extra files in dirs",
		auds:           []string{"audience1", "audience2"},
		audsWithTokens: []string{"audience1", "audience2"},
		extraFiles:     []string{filepath.Join("audience1", "extra1"), filepath.Join("some-dir", "extra2")},
		wantErr:        false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Make sure parent cache directory exists
			t.Setenv("HOME", t.TempDir())
			cacheDir, err := os.UserCacheDir()
			if err != nil {
				t.Fatal(err)
			}
			cacheDir = filepath.Join(cacheDir, parentDir)

			if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
				t.Fatalf("failed to create parent dir %s: %s", cacheDir, err.Error())
			}
			// Create audience directories
			for _, aud := range test.auds {
				dir := filepath.Join(cacheDir, aud)
				if err := os.MkdirAll(dir, os.ModePerm); err != nil {
					t.Fatalf("failed to create dir %s: %s", dir, err.Error())
				}
			}
			// Populate "token" files
			for _, aud := range test.audsWithTokens {
				for _, kind := range AllKinds {
					tok := filepath.Join(cacheDir, aud, string(kind))
					if _, err := os.Create(tok); err != nil {
						t.Fatalf("failed to create fake token %s: %s", tok, err.Error())
					}
				}
			}
			// Create extra files
			for _, extra := range test.extraFiles {
				p := filepath.Join(cacheDir, extra)
				if err := os.MkdirAll(filepath.Dir(p), os.ModePerm); err != nil {
					t.Fatalf("failed to create dir %s: %s", filepath.Dir(p), err.Error())
				}
				if _, err := os.Create(p); err != nil {
					t.Fatalf("failed to create extra file %s: %s", p, err.Error())
				}
			}

			// Then burn it to the ground.
			err = DeleteAll()

			// Sift through the ashes.
			if (err != nil) != test.wantErr {
				t.Errorf("DeleteAll() error return mismatch, want=%t, got=%v", test.wantErr, err)
			}
			files, err := os.ReadDir(cacheDir)
			if err != nil {
				t.Fatalf("failed to list files in temp cache dir: %v", err)
			}
			// Only expect test.extraFiles to remain, all other files/dirs should have been removed.
			// TODO: make this more robust and count subfiles within dirs in `files`
			// This will currently fail to catch situation where one dir has >1 extra file
			if len(files) != len(test.extraFiles) {
				t.Errorf("Remaining files mismatch, want=%d, got=%d", len(test.extraFiles), len(files))
			}
		})
	}
}

func TestSaveLoadToken(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	tests := []struct {
		description string
		audience    string
	}{
		{
			description: "default audience",
			audience:    "http://api.api-system.svc",
		}, {
			description: "default audience with audience",
			audience:    "http://api.api-system.svc",
		}, {
			description: "custom audience",
			audience:    "https://abc-sigstore.enforce.dev",
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			for _, kind := range AllKinds {
				tokenContents := []byte("mytoken")
				if err := Save(tokenContents, kind, test.audience); err != nil {
					t.Fatal(err)
				}
				contents, err := Load(kind, test.audience)
				if err != nil {
					t.Fatal(err)
				}
				if string(contents) != string(tokenContents) {
					t.Fatalf("expected %s got %s", string(tokenContents), string(contents))
				}
			}
		})
	}
}

func TestRemainingLife(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	now := time.Unix(0, 0)
	timeUntil = func(t time.Time) time.Duration {
		return t.Sub(now)
	}

	type ts struct {
		audience string
		subject  string
		duration time.Duration
	}

	tests := []struct {
		name     string
		token    ts
		kind     Kind
		audience string
		less     time.Duration
		want     time.Duration
	}{{
		name: "exists, 0 less, 30m remain",
		token: ts{
			audience: "audience",
			subject:  "subject",
			duration: 30 * time.Minute,
		},
		kind:     KindAccess,
		audience: "audience",
		less:     0,
		want:     30 * time.Minute,
	}, {
		name: "exists, less 10m, 20m remain",
		token: ts{
			audience: "audience",
			subject:  "subject",
			duration: 30 * time.Minute,
		},
		kind:     KindRefresh,
		audience: "audience",
		less:     10 * time.Minute,
		want:     20 * time.Minute,
	}, {
		name: "exists less 30m, 0s remain",
		token: ts{
			audience: "audience",
			subject:  "subject",
			duration: 30 * time.Minute,
		},
		kind:     KindAccess,
		audience: "audience",
		less:     30 * time.Minute,
		want:     0,
	}, {
		name: "exists less 30m, 0s remain (refresh)",
		token: ts{
			audience: "audience",
			subject:  "subject",
			duration: 30 * time.Minute,
		},
		kind:     KindRefresh,
		audience: "audience",
		less:     30 * time.Minute,
		want:     0,
	}, {
		name: "exists less 40m, 0s remain",
		token: ts{
			audience: "audience",
			subject:  "subject",
			duration: 30 * time.Minute,
		},
		kind:     KindAccess,
		audience: "audience",
		less:     40 * time.Minute,
		want:     0,
	}, {
		name: "doesn't exist, 0s returned",
		token: ts{
			audience: "audience",
			subject:  "subject",
			duration: 30 * time.Minute,
		},
		kind:     KindAccess,
		audience: "different",
		less:     0,
		want:     0,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tok string
			if test.kind == KindAccess {
				tok = testToken(t, test.token.audience, test.token.subject, now, test.token.duration)
			} else {
				tok = testRefreshToken(t, "https://issuer.unit-test.com", test.token.audience, test.token.subject, now, test.token.duration)
			}
			if err := Save([]byte(tok), test.kind, test.token.audience); err != nil {
				t.Fatalf("Save() unexpected error=%v", err)
			}

			got := RemainingLife(test.kind, test.audience, test.less)

			if got != test.want {
				t.Fatalf("RemainingLife() mismatch, want=%v, got=%v", test.want, got)
			}
		})
	}
}

func TestPath(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		kind         Kind
		audience     string
		expectedPath string
	}{
		"no audience": {
			audience:     "",
			kind:         KindAccess,
			expectedPath: filepath.Join(cacheDir, parentDir, string(KindAccess)),
		},
		"audience sans replacement": {
			audience:     "console.example.com",
			kind:         KindAccess,
			expectedPath: filepath.Join(cacheDir, parentDir, "console.example.com", string(KindAccess)),
		},
		"audience with replacement": {
			audience:     "https://console.example.com",
			kind:         KindRefresh,
			expectedPath: filepath.Join(cacheDir, parentDir, "https:--console.example.com", string(KindRefresh)),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := Path(test.kind, test.audience)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(test.expectedPath, got); diff != "" {
				t.Fatal("Path() mismatch (-want, +got):", diff)
			}
		})
	}
}
