/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const (
	issuerA = "https://issuer.a.test"
	issuerB = "https://issuer.b.test"
	audAPI  = "https://api.a.test"
	audReg  = "cgr.a.test"
	audAPK  = "apk.a.test"
)

// isolateCache points the token cache at a fresh directory and returns the
// cache root (<UserCacheDir>/chainguard).
func isolateCache(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	t.Setenv("XDG_CACHE_HOME", dir)
	t.Setenv("LocalAppData", dir)
	base, err := os.UserCacheDir()
	if err != nil {
		t.Fatal(err)
	}
	return filepath.Join(base, parentDir)
}

// claimsOf are the fields familyToken encodes. Each token gets its own kid, as
// the issuer's do.
type claimsOf struct {
	iss, sub string
	iat, exp time.Time
	auds     []string
}

func familyToken(t *testing.T, c claimsOf) string {
	t.Helper()
	b, err := json.Marshal(map[string]any{
		"kid": rand.Text(), "iss": c.iss, "sub": c.sub, "aud": c.auds,
		"iat": c.iat.Unix(), "exp": c.exp.Unix(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return base64.StdEncoding.EncodeToString(b)
}

// login returns the claims of a fresh login by sub for auds.
func login(iss, sub string, auds ...string) claimsOf {
	now := time.Now().Truncate(time.Second)
	return claimsOf{iss: iss, sub: sub, iat: now, exp: now.Add(20 * time.Hour), auds: auds}
}

// next returns the claims of c's successor: the same family, a later iat.
func (c claimsOf) next() claimsOf {
	c.iat = c.iat.Add(time.Minute)
	return c
}

func savePerAudience(t *testing.T, tok, aud string) {
	t.Helper()
	if err := Save([]byte(tok), KindRefresh, aud); err != nil {
		t.Fatal(err)
	}
}

func mustSaveFamily(t *testing.T, tok, iss string, auds ...string) {
	t.Helper()
	if err := SaveFamily([]byte(tok), auds, WithIssuer(iss)); err != nil {
		t.Fatalf("SaveFamily() = %v", err)
	}
}

func loadString(t *testing.T, aud string, opts ...Option) (string, error) {
	t.Helper()
	b, err := Load(KindRefresh, aud, opts...)
	return string(b), err
}

// familyView is the store's content in a comparable form: each family's token
// and the audiences it covers.
type familyView map[string][]string

func viewFamilies(t *testing.T) familyView {
	t.Helper()
	set, err := loadFamilies(false)
	if err != nil {
		t.Fatal(err)
	}
	v := familyView{}
	for _, f := range set.all {
		v[f.Token] = f.Audiences
	}
	return v
}

// diffFamilies compares the store with want, ignoring audience order.
func diffFamilies(t *testing.T, want familyView) string {
	t.Helper()
	return cmp.Diff(want, viewFamilies(t), cmpopts.SortSlices(strings.Compare))
}

func perAudienceExists(t *testing.T, aud string) bool {
	t.Helper()
	p, err := Path(KindRefresh, aud)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(p)
	return err == nil
}

// declining is an ExchangeFunc that exchanges nothing, so a Refresh with it
// only takes the lock and adopts.
func declining(context.Context, string) (string, error) { return "", nil }

func TestFamily_IssuerScopedSelection(t *testing.T) {
	isolateCache(t)
	a := familyToken(t, login(issuerA, "alice", audAPI, audReg))
	b := familyToken(t, login(issuerB, "alice", audReg))
	mustSaveFamily(t, a, issuerA, audAPI, audReg)
	mustSaveFamily(t, b, issuerB, audReg)

	for _, tc := range []struct {
		name string
		aud  string
		opts []Option
		want string
	}{
		{name: "issuer A, API", aud: audAPI, opts: []Option{WithIssuer(issuerA)}, want: a},
		{name: "issuer A, registry", aud: audReg, opts: []Option{WithIssuer(issuerA)}, want: a},
		{name: "issuer B, registry", aud: audReg, opts: []Option{WithIssuer(issuerB)}, want: b},
		{name: "issuer B, API", aud: audAPI, opts: []Option{WithIssuer(issuerB)}},
		{name: "no issuer sees the copy another issuer's login kept", aud: audReg, want: a},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := loadString(t, tc.aud, tc.opts...)
			if tc.want == "" {
				if !errors.Is(err, fs.ErrNotExist) {
					t.Errorf("Load() = %q, %v; want fs.ErrNotExist", got, err)
				}
				if life := RemainingLife(KindRefresh, tc.aud, 0, tc.opts...); life != 0 {
					t.Errorf("RemainingLife() = %v, want 0", life)
				}
				return
			}
			if err != nil || got != tc.want {
				t.Errorf("Load() = %q, %v; want %q", got, err, tc.want)
			}
			if life := RemainingLife(KindRefresh, tc.aud, 0, tc.opts...); life <= 0 {
				t.Errorf("RemainingLife() = %v, want > 0", life)
			}
		})
	}
}

func TestFamily_NoIssuerSeesPerAudienceOnly(t *testing.T) {
	isolateCache(t)
	perAud := familyToken(t, login(issuerA, "alice", audAPI))
	fam := familyToken(t, login(issuerA, "bob", audAPI))
	savePerAudience(t, perAud, audAPI)
	// Written directly so adoption cannot fold perAud in first.
	set := &familySet{dirty: map[*family]struct{}{}}
	set.cover(newFamily(issuerA, fam), []string{audAPI})
	if err := set.flush(); err != nil {
		t.Fatal(err)
	}

	if got, err := loadString(t, audAPI); err != nil || got != perAud {
		t.Errorf("Load() without issuer = %q, %v; want the per-audience token", got, err)
	}
	if got, err := loadString(t, audAPI, WithIssuer(issuerA)); err != nil || got != fam {
		t.Errorf("Load() with issuer = %q, %v; want the family token", got, err)
	}
	if err := Refresh(t.Context(), audAPI, declining); !errors.Is(err, errNoIssuer) {
		t.Errorf("Refresh() without issuer = %v, want errNoIssuer", err)
	}
	if err := SaveFamily([]byte(fam), []string{audAPI}); !errors.Is(err, errNoIssuer) {
		t.Errorf("SaveFamily() without issuer = %v, want errNoIssuer", err)
	}
}

func TestFamily_PerAudienceFallback(t *testing.T) {
	for _, tc := range []struct {
		name  string
		token string
		want  bool
	}{
		{name: "same issuer", token: familyToken(t, login(issuerA, "alice", audAPI)), want: true},
		{name: "other issuer", token: familyToken(t, login(issuerB, "alice", audAPI))},
		{name: "opaque", token: "zendesk-" + rand.Text()},
	} {
		t.Run(tc.name, func(t *testing.T) {
			isolateCache(t)
			savePerAudience(t, tc.token, audAPI)
			got, err := loadString(t, audAPI, WithIssuer(issuerA))
			if tc.want {
				if err != nil || got != tc.token {
					t.Errorf("Load() = %q, %v; want the per-audience token", got, err)
				}
				return
			}
			if !errors.Is(err, fs.ErrNotExist) {
				t.Errorf("Load() = %q, %v; want fs.ErrNotExist", got, err)
			}
		})
	}
}

func TestSaveFamily(t *testing.T) {
	t.Run("Save with an issuer is a one-audience family", func(t *testing.T) {
		isolateCache(t)
		tok := familyToken(t, login(issuerA, "alice", audAPI))
		if err := Save([]byte(tok), KindRefresh, audAPI, WithIssuer(issuerA)); err != nil {
			t.Fatal(err)
		}
		if diff := diffFamilies(t, familyView{tok: {audAPI}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
		if got, err := loadString(t, audAPI); err != nil || got != tok {
			t.Errorf("issuer-less Load() = %q, %v; want the family's token", got, err)
		}
	})

	t.Run("an issuer-less reader finds the token in every covered audience", func(t *testing.T) {
		base := isolateCache(t)
		tok := familyToken(t, login(issuerA, "alice", audAPI, audReg, audAPK))
		mustSaveFamily(t, tok, issuerA, audAPI, audReg, audAPK)
		for _, aud := range []string{audAPI, audReg, audAPK} {
			if got, err := loadString(t, aud); err != nil || got != tok {
				t.Errorf("Load(%s) without issuer = %q, %v; want the family's token", aud, got, err)
			}
			if runtime.GOOS == "windows" {
				continue
			}
			fi, err := os.Stat(filepath.Join(base, audienceDir(aud), string(KindRefresh)))
			if err != nil {
				t.Fatal(err)
			}
			if got := fi.Mode().Perm(); got != 0o600 {
				t.Errorf("per-audience copy for %s: mode = %v, want 0600", aud, got)
			}
		}
	})

	t.Run("coverage is exclusive within an issuer", func(t *testing.T) {
		isolateCache(t)
		first := familyToken(t, login(issuerA, "alice", audAPI, audReg, audAPK))
		second := familyToken(t, login(issuerA, "alice", audReg))
		other := familyToken(t, login(issuerB, "alice", audReg))
		mustSaveFamily(t, first, issuerA, audAPI, audReg, audAPK)
		mustSaveFamily(t, other, issuerB, audReg)
		mustSaveFamily(t, second, issuerA, audReg)
		want := familyView{first: {audAPI, audAPK}, second: {audReg}, other: {audReg}}
		if diff := diffFamilies(t, want); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
	})

	// An issuer-less reader sees only the per-audience file, so it holds the
	// latest login of that issuer. Another issuer's live session there is
	// kept: it may be that issuer's only refresh token, and an issuer-less
	// reader configured for it would otherwise present this login's token.
	t.Run("replaces its own issuer's per-audience copies", func(t *testing.T) {
		isolateCache(t)
		old := familyToken(t, login(issuerA, "alice", audAPI, audReg))
		foreign := familyToken(t, login(issuerB, "alice", audReg))
		expired := familyToken(t, claimsOf{iss: issuerB, sub: "alice", iat: time.Now().Add(-2 * time.Hour), exp: time.Now().Add(-time.Hour), auds: []string{audAPK}})
		savePerAudience(t, old, audAPI)
		savePerAudience(t, foreign, audReg)
		savePerAudience(t, expired, audAPK)
		fresh := familyToken(t, login(issuerA, "bob", audAPI, audReg, audAPK))
		mustSaveFamily(t, fresh, issuerA, audAPI, audReg, audAPK)
		if diff := diffFamilies(t, familyView{fresh: {audAPI, audReg, audAPK}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
		for aud, want := range map[string]string{audAPI: fresh, audReg: foreign, audAPK: fresh} {
			if got, err := loadString(t, aud); err != nil || got != want {
				t.Errorf("Load(%s) without issuer = %q, %v; want %q", aud, got, err, want)
			}
		}
		if got, err := loadString(t, audReg, WithIssuer(issuerB)); err != nil || got != foreign {
			t.Errorf("Load(%s, issuer B) = %q, %v; want issuer B's session kept", audReg, got, err)
		}
	})

	t.Run("overwrites another issuer's copy in the primary audience", func(t *testing.T) {
		isolateCache(t)
		foreign := familyToken(t, login(issuerB, "alice", audAPI, audReg))
		savePerAudience(t, foreign, audAPI)
		savePerAudience(t, foreign, audReg)
		fresh := familyToken(t, login(issuerA, "bob", audAPI, audReg))
		if err := SaveFamily([]byte(fresh), []string{audAPI, audReg}, WithIssuer(issuerA), WithPrimaryAudience(audAPI)); err != nil {
			t.Fatalf("SaveFamily() = %v", err)
		}
		for aud, want := range map[string]string{audAPI: fresh, audReg: foreign} {
			if got, err := loadString(t, aud); err != nil || got != want {
				t.Errorf("Load(%s) without issuer = %q, %v; want %q", aud, got, err, want)
			}
		}
	})

	// The family's issuer is the dialed name, so a token whose iss names
	// another URL is compared with the copies by its claim.
	t.Run("replaces a copy from the same issuer under another name", func(t *testing.T) {
		isolateCache(t)
		const internal = "http://issuer.internal.svc"
		old := familyToken(t, login(internal, "alice", audAPI))
		savePerAudience(t, old, audAPI)
		fresh := familyToken(t, login(internal, "alice", audAPI))
		mustSaveFamily(t, fresh, issuerA, audAPI)
		if got, err := loadString(t, audAPI); err != nil || got != fresh {
			t.Errorf("Load(%s) without issuer = %q, %v; want the fresh login", audAPI, got, err)
		}
	})

	// An issuer reached through a mapped domain mints tokens whose iss names
	// its cluster-internal URL, as the Kind e2e environment's does.
	t.Run("records the family under the issuer it was obtained from", func(t *testing.T) {
		isolateCache(t)
		tok := familyToken(t, login("http://issuer.internal.svc", "alice", audAPI))
		mustSaveFamily(t, tok, issuerA, audAPI)
		if got, err := loadString(t, audAPI, WithIssuer(issuerA)); err != nil || got != tok {
			t.Errorf("Load(issuer A) = %q, %v; want the family's token", got, err)
		}
		set, err := loadFamilies(false)
		if err != nil {
			t.Fatal(err)
		}
		if f := set.covering("http://issuer.internal.svc", audAPI); f != nil {
			t.Errorf("a family was recorded under the token's iss: %+v", f)
		}
	})

	// A rotation-capable family is refreshed only through the store, under its
	// lock. The copies a login would have overwritten are deleted, so the next
	// locked call cannot adopt an earlier login's copy back over the family.
	t.Run("a rotation-capable family writes no per-audience copies", func(t *testing.T) {
		isolateCache(t)
		old := familyToken(t, login(issuerA, "alice", audAPI, audReg))
		foreign := familyToken(t, login(issuerB, "alice", audReg, audAPK))
		savePerAudience(t, old, audAPI)
		savePerAudience(t, foreign, audReg)
		savePerAudience(t, foreign, audAPK)
		fresh := familyToken(t, login(issuerA, "bob", audAPI, audReg, audAPK))
		if err := SaveFamily([]byte(fresh), []string{audAPI, audReg, audAPK}, WithIssuer(issuerA), WithPrimaryAudience(audAPK), WithRotationCapable()); err != nil {
			t.Fatalf("SaveFamily() = %v", err)
		}
		if diff := diffFamilies(t, familyView{fresh: {audAPI, audReg, audAPK}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
		for _, aud := range []string{audAPI, audAPK} {
			if perAudienceExists(t, aud) {
				t.Errorf("per-audience copy for %s exists, want it deleted", aud)
			}
		}
		if got, err := loadString(t, audReg); err != nil || got != foreign {
			t.Errorf("Load(%s) without issuer = %q, %v; want issuer B's session kept", audReg, got, err)
		}
		// Adoption must not bring the deleted copy back over the family.
		if err := Refresh(t.Context(), audAPI, declining, WithIssuer(issuerA)); err != nil {
			t.Fatalf("Refresh() = %v", err)
		}
		if diff := diffFamilies(t, familyView{fresh: {audAPI, audReg, audAPK}}); diff != "" {
			t.Errorf("families after adoption (-want +got):\n%s", diff)
		}
		if got, err := loadString(t, audAPI, WithIssuer(issuerA)); err != nil || got != fresh {
			t.Errorf("Load(%s, issuer A) = %q, %v; want the family's token", audAPI, got, err)
		}
	})

	// A copy Windows will not delete, because a lock-free reader has it open,
	// is emptied instead. The login succeeds and adoption finds nothing to take
	// the audience back with.
	t.Run("a rotation-capable family empties a copy it cannot delete", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("a read-only directory mode does not block removal on Windows")
		}
		if os.Geteuid() == 0 {
			t.Skip("running as root ignores the read-only mode this test depends on")
		}
		base := isolateCache(t)
		old := familyToken(t, login(issuerA, "alice", audAPI))
		savePerAudience(t, old, audAPI)
		dir := filepath.Join(base, audienceDir(audAPI))
		if err := os.Chmod(dir, 0o500); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

		fresh := familyToken(t, login(issuerA, "bob", audAPI))
		if err := SaveFamily([]byte(fresh), []string{audAPI}, WithIssuer(issuerA), WithRotationCapable()); err != nil {
			t.Fatalf("SaveFamily() = %v", err)
		}
		if err := os.Chmod(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := Refresh(t.Context(), audAPI, declining, WithIssuer(issuerA)); err != nil {
			t.Fatalf("Refresh() = %v", err)
		}
		if diff := diffFamilies(t, familyView{fresh: {audAPI}}); diff != "" {
			t.Errorf("families after adoption (-want +got):\n%s", diff)
		}
	})

	t.Run("rejects a token that is not a Chainguard refresh token", func(t *testing.T) {
		isolateCache(t)
		if err := SaveFamily([]byte("opaque-"+rand.Text()), []string{audAPI}, WithIssuer(issuerA)); err == nil {
			t.Error("SaveFamily() = nil, want an error")
		}
	})
}

func TestAdoption(t *testing.T) {
	alice := login(issuerA, "alice", audAPI, audReg, audAPK)
	aliceTok := familyToken(t, alice)
	aliceNext := familyToken(t, alice.next())
	bobTok := familyToken(t, login(issuerA, "bob", audAPI, audReg, audAPK))
	foreignTok := familyToken(t, login(issuerB, "alice", audAPI, audReg))
	expiredTok := familyToken(t, claimsOf{iss: issuerA, sub: "alice", iat: time.Now().Add(-2 * time.Hour), exp: time.Now().Add(-time.Hour), auds: []string{audAPK}})

	for _, tc := range []struct {
		name string
		// family is a family written before adoption, as token -> audiences.
		family familyView
		// copies are per-audience files, as audience -> token.
		copies map[string]string
		want   familyView
		// gone are the audiences whose per-audience file must be removed.
		gone []string
	}{{
		name:   "N copies of one login collapse into one family",
		copies: map[string]string{audAPI: aliceTok, audReg: aliceTok, audAPK: aliceTok},
		want:   familyView{aliceTok: {audAPI, audAPK, audReg}},
	}, {
		name:   "coverage only where a copy was found",
		copies: map[string]string{audAPI: aliceTok, audReg: aliceTok},
		want:   familyView{aliceTok: {audAPI, audReg}},
	}, {
		name:   "a copy of a family's token adds its audience",
		family: familyView{aliceTok: {audAPI}},
		copies: map[string]string{audReg: aliceTok},
		want:   familyView{aliceTok: {audAPI, audReg}},
	}, {
		name:   "two logins are grouped by token",
		copies: map[string]string{audAPI: aliceTok, audReg: bobTok, audAPK: bobTok},
		want:   familyView{aliceTok: {audAPI}, bobTok: {audAPK, audReg}},
	}, {
		name:   "an older writer's login takes its audience from the family",
		family: familyView{aliceTok: {audAPI, audReg}},
		copies: map[string]string{audReg: bobTok},
		want:   familyView{aliceTok: {audAPI}, bobTok: {audReg}},
	}, {
		name:   "a superseded parent is deleted, not adopted",
		family: familyView{aliceNext: {audAPI}},
		copies: map[string]string{audReg: aliceTok},
		want:   familyView{aliceNext: {audAPI}},
		gone:   []string{audReg},
	}, {
		name:   "another issuer's copies are left alone",
		copies: map[string]string{audAPI: foreignTok},
		want:   familyView{},
	}, {
		name:   "an expired copy is not adopted",
		copies: map[string]string{audAPK: expiredTok},
		want:   familyView{},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			base := isolateCache(t)
			set := &familySet{dirty: map[*family]struct{}{}}
			for tok, auds := range tc.family {
				set.cover(newFamily(issuerA, tok), auds)
			}
			if err := set.flush(); err != nil {
				t.Fatal(err)
			}
			for aud, tok := range tc.copies {
				savePerAudience(t, tok, aud)
			}
			// cg zendesk's files share the cache and must survive.
			zdDir := filepath.Join(base, audienceDir("https://unit.zendesk.com"))
			if err := os.MkdirAll(zdDir, 0o700); err != nil {
				t.Fatal(err)
			}
			for _, name := range []string{string(KindRefresh), "refresh.lock"} {
				if err := os.WriteFile(filepath.Join(zdDir, name), []byte("zendesk"), 0o600); err != nil {
					t.Fatal(err)
				}
			}

			if err := Refresh(t.Context(), audAPI, declining, WithIssuer(issuerA)); err != nil && !errors.Is(err, fs.ErrNotExist) {
				t.Fatalf("Refresh() = %v", err)
			}

			if diff := diffFamilies(t, tc.want); diff != "" {
				t.Errorf("families (-want +got):\n%s", diff)
			}
			for aud := range tc.copies {
				if got, want := perAudienceExists(t, aud), !slices.Contains(tc.gone, aud); got != want {
					t.Errorf("per-audience file for %s exists = %t, want %t", aud, got, want)
				}
			}
			for _, name := range []string{string(KindRefresh), "refresh.lock"} {
				if b, err := os.ReadFile(filepath.Join(zdDir, name)); err != nil || string(b) != "zendesk" {
					t.Errorf("zendesk %s: got %q, %v; want it untouched", name, b, err)
				}
			}
		})
	}
}

// TestAdoption_CompatibilityCopies pins that the per-audience copies SaveFamily
// writes are read as the family's own: a second family for the same token
// would split coverage, and rotating one would leave the other holding the
// superseded token.
func TestAdoption_CompatibilityCopies(t *testing.T) {
	isolateCache(t)
	tok := familyToken(t, login(issuerA, "alice", audAPI, audReg, audAPK))
	mustSaveFamily(t, tok, issuerA, audAPI, audReg, audAPK)
	before, err := loadFamilies(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(before.all) != 1 {
		t.Fatalf("families after SaveFamily = %d, want 1", len(before.all))
	}

	if err := Refresh(t.Context(), audReg, declining, WithIssuer(issuerA)); err != nil {
		t.Fatalf("Refresh() = %v", err)
	}

	after, err := loadFamilies(false)
	if err != nil {
		t.Fatal(err)
	}
	type view struct {
		Name      string
		Token     string
		Audiences []string
	}
	got := make([]view, 0, len(after.all))
	for _, f := range after.all {
		got = append(got, view{Name: f.name, Token: f.Token, Audiences: f.Audiences})
	}
	want := []view{{Name: before.all[0].name, Token: tok, Audiences: []string{audAPI, audAPK, audReg}}}
	if diff := cmp.Diff(want, got, cmpopts.SortSlices(strings.Compare)); diff != "" {
		t.Errorf("families after adoption (-want +got):\n%s", diff)
	}
	for _, aud := range []string{audAPI, audReg, audAPK} {
		if got, err := loadString(t, aud); err != nil || got != tok {
			t.Errorf("Load(%s) without issuer after adoption = %q, %v; want the copy kept", aud, got, err)
		}
	}
}

// TestAdoption_IssuerAlias pins that a compatibility copy of a token held by
// a family under one name for an issuer is not adopted by a caller using
// another name for it: rotating either family would leave the other holding
// the superseded parent, and presenting that revokes the family.
func TestAdoption_IssuerAlias(t *testing.T) {
	isolateCache(t)
	const canonical = "https://issuer.canonical.test"
	c := login(canonical, "alice", audAPI, audReg)
	tok, successor := familyToken(t, c), familyToken(t, c.next())
	mustSaveFamily(t, tok, issuerA, audAPI, audReg)

	// Readers agree with Refresh: a positive life Refresh cannot act on
	// sends a caller into a retry loop.
	for _, aud := range []string{audAPI, audReg} {
		if got, err := loadString(t, aud, WithIssuer(canonical)); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Load(%s, canonical) = %q, %v; want fs.ErrNotExist", aud, got, err)
		}
		if life := RemainingLife(KindRefresh, aud, 0, WithIssuer(canonical)); life != 0 {
			t.Errorf("RemainingLife(%s, canonical) = %v, want 0", aud, life)
		}
		if got, err := loadString(t, aud, WithIssuer(issuerA)); err != nil || got != tok {
			t.Errorf("Load(%s, issuer A) = %q, %v; want the family's token", aud, got, err)
		}
	}
	if err := Refresh(t.Context(), audAPI, declining, WithIssuer(canonical)); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Refresh(canonical) = %v, want fs.ErrNotExist", err)
	}
	if diff := diffFamilies(t, familyView{tok: {audAPI, audReg}}); diff != "" {
		t.Errorf("families after the canonical name's adoption (-want +got):\n%s", diff)
	}

	rotate := func(_ context.Context, presented string) (string, error) {
		if presented != tok {
			t.Errorf("presented %q, want the family's token", presented)
		}
		return successor, nil
	}
	if err := Refresh(t.Context(), audReg, rotate, WithIssuer(issuerA)); err != nil {
		t.Fatalf("Refresh(issuer A) = %v", err)
	}
	if diff := diffFamilies(t, familyView{successor: {audAPI, audReg}}); diff != "" {
		t.Errorf("families after rotation (-want +got):\n%s", diff)
	}
	if err := Refresh(t.Context(), audAPI, declining, WithIssuer(canonical)); !errors.Is(err, fs.ErrNotExist) {
		t.Errorf("Refresh(canonical) after rotation = %v, want fs.ErrNotExist", err)
	}
}

func TestRefresh(t *testing.T) {
	session := login(issuerA, "alice", audAPI, audReg, audAPK)

	t.Run("a successor replaces the family's token and every copy of the parent", func(t *testing.T) {
		isolateCache(t)
		presented := familyToken(t, session)
		foreign := familyToken(t, login(issuerB, "alice", audAPK))
		for _, aud := range []string{audAPI, audReg} {
			savePerAudience(t, presented, aud)
		}
		savePerAudience(t, foreign, audAPK)
		successor := familyToken(t, session.next())

		var got []string
		err := Refresh(t.Context(), audReg, func(_ context.Context, tok string) (string, error) {
			got = append(got, tok)
			return successor, nil
		}, WithIssuer(issuerA))
		if err != nil {
			t.Fatalf("Refresh() = %v", err)
		}
		if diff := cmp.Diff([]string{presented}, got); diff != "" {
			t.Errorf("presented (-want +got):\n%s", diff)
		}
		if diff := diffFamilies(t, familyView{successor: {audAPI, audReg}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
		for _, aud := range []string{audAPI, audReg} {
			if perAudienceExists(t, aud) {
				t.Errorf("a copy of the superseded parent survived under %s", aud)
			}
			if tok, err := loadString(t, aud, WithIssuer(issuerA)); err != nil || tok != successor {
				t.Errorf("Load(%s) = %q, %v; want the successor", aud, tok, err)
			}
		}
		if tok, err := loadString(t, audAPK); err != nil || tok != foreign {
			t.Errorf("another issuer's token: got %q, %v; want it untouched", tok, err)
		}
	})

	// A sliding session's successor carries a later exp than its parent.
	t.Run("a successor with a later expiry replaces its parent", func(t *testing.T) {
		isolateCache(t)
		presented := familyToken(t, session)
		mustSaveFamily(t, presented, issuerA, audAPI, audReg)
		slid := session.next()
		slid.exp = session.exp.Add(29 * 24 * time.Hour)
		successor := familyToken(t, slid)

		if err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
			return successor, nil
		}, WithIssuer(issuerA)); err != nil {
			t.Fatalf("Refresh() = %v", err)
		}
		if diff := diffFamilies(t, familyView{successor: {audAPI, audReg}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
		for _, aud := range []string{audAPI, audReg} {
			if tok, err := loadString(t, aud, WithIssuer(issuerA)); err != nil || tok != successor {
				t.Errorf("Load(%s) = %q, %v; want the successor", aud, tok, err)
			}
			if life := RemainingLife(KindRefresh, aud, 0, WithIssuer(issuerA)); life <= 20*time.Hour {
				t.Errorf("RemainingLife(%s) = %v, want the successor's later expiry", aud, life)
			}
		}
	})

	// A successor written to a per-audience copy would be presented by an
	// issuer-less reader after the family moved on, which is replay.
	t.Run("rotation leaves no per-audience copy of either token", func(t *testing.T) {
		base := isolateCache(t)
		presented := familyToken(t, session)
		mustSaveFamily(t, presented, issuerA, audAPI, audReg, audAPK)
		successor := familyToken(t, session.next())
		if err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
			return successor, nil
		}, WithIssuer(issuerA)); err != nil {
			t.Fatalf("Refresh() = %v", err)
		}
		if diff := diffFamilies(t, familyView{successor: {audAPI, audReg, audAPK}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
		for _, aud := range []string{audAPI, audReg, audAPK} {
			b, err := os.ReadFile(filepath.Join(base, audienceDir(aud), string(KindRefresh)))
			switch {
			case errors.Is(err, fs.ErrNotExist):
			case err != nil:
				t.Fatal(err)
			case string(b) == presented:
				t.Errorf("per-audience copy for %s holds the superseded token", aud)
			case string(b) == successor:
				t.Errorf("per-audience copy for %s holds the successor", aud)
			}
		}
	})

	for _, tc := range []struct {
		name      string
		successor func(presented string) string
	}{
		{name: "echo", successor: func(p string) string { return p }},
		{name: "echo with surrounding space", successor: func(p string) string { return " " + p + "\n" }},
		{name: "empty", successor: func(string) string { return "" }},
	} {
		t.Run(tc.name+" writes nothing", func(t *testing.T) {
			isolateCache(t)
			presented := familyToken(t, session)
			mustSaveFamily(t, presented, issuerA, audAPI, audReg)
			before := viewFamilies(t)
			if err := Refresh(t.Context(), audAPI, func(_ context.Context, tok string) (string, error) {
				return tc.successor(tok), nil
			}, WithIssuer(issuerA)); err != nil {
				t.Fatalf("Refresh() = %v", err)
			}
			if diff := diffFamilies(t, before); diff != "" {
				t.Errorf("families (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("an exchange error is returned and writes nothing", func(t *testing.T) {
		isolateCache(t)
		presented := familyToken(t, session)
		mustSaveFamily(t, presented, issuerA, audAPI)
		want := errors.New("rejected")
		err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
			return familyToken(t, session.next()), want
		}, WithIssuer(issuerA))
		if !errors.Is(err, want) {
			t.Errorf("Refresh() = %v, want %v", err, want)
		}
		if diff := diffFamilies(t, familyView{presented: {audAPI}}); diff != "" {
			t.Errorf("families (-want +got):\n%s", diff)
		}
	})

	t.Run("nothing to refresh", func(t *testing.T) {
		isolateCache(t)
		called := false
		err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
			called = true
			return "", nil
		}, WithIssuer(issuerA))
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Refresh() = %v, want fs.ErrNotExist", err)
		}
		if called {
			t.Error("exchange was called with no token")
		}
	})
}

// TestRefresh_SuccessorNotStored pins the persistent write failure: the
// family is discarded, so the superseded parent is never presented again, and
// the error says the successor was lost.
func TestRefresh_SuccessorNotStored(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a read-only directory mode does not block writes on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root ignores the read-only mode this test depends on")
	}
	base := isolateCache(t)
	persistBackoff = 0
	t.Cleanup(func() { persistBackoff = 50 * time.Millisecond })

	session := login(issuerA, "alice", audAPI, audReg)
	presented := familyToken(t, session)
	mustSaveFamily(t, presented, issuerA, audAPI, audReg)

	// A read-only families directory refuses the staging file and the
	// remove, so the discard has to fall back to truncating the file.
	dir := filepath.Join(base, familiesDir)
	if err := os.Chmod(dir, 0o500); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })

	err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
		return familyToken(t, session.next()), nil
	}, WithIssuer(issuerA))
	if !errors.Is(err, ErrSuccessorNotStored) {
		t.Fatalf("Refresh() = %v, want ErrSuccessorNotStored", err)
	}

	for _, aud := range []string{audAPI, audReg} {
		if tok, err := loadString(t, aud, WithIssuer(issuerA)); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Load(%s) = %q, %v; want the parent gone", aud, tok, err)
		}
	}
	var presentedAgain []string
	err = Refresh(t.Context(), audReg, func(_ context.Context, tok string) (string, error) {
		presentedAgain = append(presentedAgain, tok)
		return "", nil
	}, WithIssuer(issuerA))
	if !errors.Is(err, fs.ErrNotExist) || len(presentedAgain) != 0 {
		t.Errorf("second Refresh() = %v, presenting %q; want fs.ErrNotExist and nothing presented", err, presentedAgain)
	}
}

// TestRefresh_SuccessorNotStoredAdopted is the discard with an adopted cache:
// the per-audience copy adoption left on disk must not outlive the family, or
// the next locked call, once writes recover, adopts it as a new family and
// presents the superseded parent.
func TestRefresh_SuccessorNotStoredAdopted(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("a read-only directory mode does not block writes on Windows")
	}
	if os.Geteuid() == 0 {
		t.Skip("running as root ignores the read-only mode this test depends on")
	}
	base := isolateCache(t)
	persistBackoff = 0
	t.Cleanup(func() { persistBackoff = 50 * time.Millisecond })

	session := login(issuerA, "alice", audAPI, audReg)
	presented := familyToken(t, session)
	mustSaveFamily(t, presented, issuerA, audAPI)
	savePerAudience(t, presented, audReg)
	decline := func(context.Context, string) (string, error) { return "", nil }
	if err := Refresh(t.Context(), audReg, decline, WithIssuer(issuerA)); err != nil {
		t.Fatalf("adopting Refresh() = %v", err)
	}

	// The audience directory refuses the remove but not the truncate.
	readOnly := []string{filepath.Join(base, familiesDir), filepath.Join(base, audienceDir(audReg))}
	for _, dir := range readOnly {
		if err := os.Chmod(dir, 0o500); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
	}
	err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
		return familyToken(t, session.next()), nil
	}, WithIssuer(issuerA))
	if !errors.Is(err, ErrSuccessorNotStored) {
		t.Fatalf("Refresh() = %v, want ErrSuccessorNotStored", err)
	}
	for _, dir := range readOnly {
		if err := os.Chmod(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}

	var presentedAgain []string
	err = Refresh(t.Context(), audReg, func(_ context.Context, tok string) (string, error) {
		presentedAgain = append(presentedAgain, tok)
		return "", nil
	}, WithIssuer(issuerA))
	if !errors.Is(err, fs.ErrNotExist) || len(presentedAgain) != 0 {
		t.Errorf("Refresh() after recovery = %v, presenting %q; want fs.ErrNotExist and nothing presented", err, presentedAgain)
	}
}

// TestRefresh_FamilyChangedDuringExchange pins what Refresh stores when the
// family it read is changed by a writer that did not hold the store lock, as
// on a filesystem that does not enforce it. The change lands inside the
// exchange, after the load and before the successor is written.
func TestRefresh_FamilyChangedDuringExchange(t *testing.T) {
	session := login(issuerA, "alice", audAPI, audReg)
	for _, tc := range []struct {
		name string
		// change rewrites or removes f on disk and returns the token the
		// family must hold afterwards, or "" when it must be gone.
		change func(t *testing.T, f *family) string
	}{{
		name: "another write is kept and the successor dropped",
		change: func(t *testing.T, f *family) string {
			f.Token = familyToken(t, session.next().next())
			if err := writeFamily(f); err != nil {
				t.Fatal(err)
			}
			return f.Token
		},
	}, {
		name: "a logout is not undone",
		change: func(t *testing.T, f *family) string {
			if err := removeFamily(f.name); err != nil {
				t.Fatal(err)
			}
			return ""
		},
	}} {
		t.Run(tc.name, func(t *testing.T) {
			isolateCache(t)
			presented := familyToken(t, session)
			mustSaveFamily(t, presented, issuerA, audAPI, audReg)
			successor := familyToken(t, session.next())

			var after string
			err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
				set, err := loadFamilies(false)
				if err != nil {
					t.Fatal(err)
				}
				f := set.covering(issuerA, audAPI)
				if f == nil {
					t.Fatal("no family covers the audience being refreshed")
				}
				after = tc.change(t, f)
				return successor, nil
			}, WithIssuer(issuerA))
			if err != nil {
				t.Errorf("Refresh() = %v, want nil", err)
			}

			want := familyView{}
			if after != "" {
				want[after] = []string{audAPI, audReg}
			}
			if diff := diffFamilies(t, want); diff != "" {
				t.Errorf("families (-want +got):\n%s", diff)
			}

			var next []string
			err = Refresh(t.Context(), audReg, func(_ context.Context, tok string) (string, error) {
				next = append(next, tok)
				return "", nil
			}, WithIssuer(issuerA))
			if after == "" {
				if !errors.Is(err, fs.ErrNotExist) || len(next) != 0 {
					t.Errorf("next Refresh() = %v, presenting %q; want fs.ErrNotExist and nothing presented", err, next)
				}
				return
			}
			if err != nil {
				t.Fatalf("next Refresh() = %v", err)
			}
			if diff := cmp.Diff([]string{after}, next); diff != "" {
				t.Errorf("next presented (-want +got):\n%s", diff)
			}
		})
	}
}

// TestDelete_WithoutLock pins that a Delete that cannot take the lock never
// rewrites a family, since a rewrite from its snapshot could put back a
// parent a lock holder has just superseded: it removes the family whole.
func TestDelete_WithoutLock(t *testing.T) {
	isolateCache(t)
	lockTimeout = 50 * time.Millisecond
	t.Cleanup(func() { lockTimeout = 45 * time.Second })
	a := familyToken(t, login(issuerA, "alice", audAPI, audReg))
	c := familyToken(t, login(issuerA, "alice", audAPK))
	mustSaveFamily(t, a, issuerA, audAPI, audReg)
	mustSaveFamily(t, c, issuerA, audAPK)

	unlock, err := lockStore(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer unlock()

	if err := Delete(KindRefresh, audAPI, WithIssuer(issuerA)); err != nil {
		t.Fatalf("Delete() = %v", err)
	}
	if diff := diffFamilies(t, familyView{c: {audAPK}}); diff != "" {
		t.Errorf("families (-want +got):\n%s", diff)
	}
}

// TestDeleteFamilies walks the tree: every family file, staging file and
// per-audience copy of a family's token goes, and other per-audience tokens
// and the lock file stay.
func TestDeleteFamilies(t *testing.T) {
	base := isolateCache(t)
	mustSaveFamily(t, familyToken(t, login(issuerA, "alice", audAPI, audReg)), issuerA, audAPI, audReg)
	mustSaveFamily(t, familyToken(t, login(issuerB, "alice", audAPI)), issuerB, audAPI)
	savePerAudience(t, "opaque-"+rand.Text(), audAPK)
	stray := filepath.Join(base, familiesDir, stagingPrefix("x")+rand.Text())
	if err := os.WriteFile(stray, []byte("token bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := DeleteFamilies(); err != nil {
		t.Fatalf("DeleteFamilies() = %v", err)
	}

	var left []string
	if err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			left = append(left, filepath.ToSlash(strings.TrimPrefix(path, base+string(filepath.Separator))))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	want := []string{lockFileName, audienceDir(audAPK) + "/" + string(KindRefresh)}
	if diff := cmp.Diff(want, left, cmpopts.SortSlices(func(a, b string) bool { return a < b })); diff != "" {
		t.Errorf("files left (-want +got):\n%s", diff)
	}
}

func TestDelete_RemovesCoverage(t *testing.T) {
	isolateCache(t)
	a := familyToken(t, login(issuerA, "alice", audAPI, audReg))
	b := familyToken(t, login(issuerB, "alice", audAPI))
	mustSaveFamily(t, a, issuerA, audAPI, audReg)
	mustSaveFamily(t, b, issuerB, audAPI)

	steps := []struct {
		aud  string
		opts []Option
		want familyView
	}{
		{aud: audAPI, opts: []Option{WithIssuer(issuerA)}, want: familyView{a: {audReg}, b: {audAPI}}},
		{aud: audReg, opts: []Option{WithIssuer(issuerA)}, want: familyView{b: {audAPI}}},
		{aud: audAPI, want: familyView{}},
	}
	for _, s := range steps {
		if err := Delete(KindRefresh, s.aud, s.opts...); err != nil {
			t.Fatalf("Delete(%s) = %v", s.aud, err)
		}
		if diff := diffFamilies(t, s.want); diff != "" {
			t.Errorf("after Delete(%s), families (-want +got):\n%s", s.aud, diff)
		}
	}
}

// TestDeleteAll_RemovesFamilies walks the tree rather than calling the
// loader: what matters is that no token-bearing file survives, and that the
// lock file does, since removing a held lock file breaks exclusion for the
// next process.
func TestDeleteAll_RemovesFamilies(t *testing.T) {
	base := isolateCache(t)
	mustSaveFamily(t, familyToken(t, login(issuerA, "alice", audAPI)), issuerA, audAPI)
	savePerAudience(t, familyToken(t, login(issuerA, "bob", audReg)), audReg)
	stray := filepath.Join(base, familiesDir, stagingPrefix("x")+rand.Text())
	if err := os.WriteFile(stray, []byte("token bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := DeleteAll(); err != nil {
		t.Fatalf("DeleteAll() = %v", err)
	}

	var left []string
	if err := filepath.WalkDir(base, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			left = append(left, strings.TrimPrefix(path, base+string(filepath.Separator)))
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{lockFileName}, left); diff != "" {
		t.Errorf("files left (-want +got):\n%s", diff)
	}
}

// TestRefresh_DiscardsSupersededToken is the client floor for the issuer's
// grace-window refusal. A caller that cannot re-authenticate -- a credential
// helper, a --refresh-only service -- has nowhere to go when the exchange is
// refused, and nothing else on a rejection path deletes a stored refresh
// token. If the refused token survived, its next presentation would land
// outside the issuer's grace window, which is replay, and the issuer would
// revoke the whole family.
func TestRefresh_DiscardsSupersededToken(t *testing.T) {
	isolateCache(t)

	session := login(issuerA, "alice", audAPI, audReg)
	presented := familyToken(t, session)
	mustSaveFamily(t, presented, issuerA, audAPI, audReg)

	refused := fmt.Errorf("%w: rpc error: code = Unauthenticated", ErrPresentedTokenSuperseded)
	err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
		return "", refused
	}, WithIssuer(issuerA))
	if !errors.Is(err, ErrPresentedTokenSuperseded) {
		t.Fatalf("Refresh() = %v, want it to wrap ErrPresentedTokenSuperseded", err)
	}

	// Both audiences, not just the one refreshed: the family covered both, and
	// one surviving copy is a token that gets adopted and presented again.
	for _, aud := range []string{audAPI, audReg} {
		if tok, err := loadString(t, aud, WithIssuer(issuerA)); !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Load(%s) = %q, %v; want the refused token gone", aud, tok, err)
		}
		if perAudienceExists(t, aud) {
			t.Errorf("per-audience copy for %s survived the refusal", aud)
		}
	}

	// The point of the discard: nothing is left to present on the next run.
	var presentedAgain []string
	err = Refresh(t.Context(), audReg, func(_ context.Context, tok string) (string, error) {
		presentedAgain = append(presentedAgain, tok)
		return "", nil
	}, WithIssuer(issuerA))
	if !errors.Is(err, fs.ErrNotExist) || len(presentedAgain) != 0 {
		t.Errorf("second Refresh() = %v, presenting %q; want fs.ErrNotExist and nothing presented", err, presentedAgain)
	}
}

// TestRefresh_KeepsTokenOnOtherErrors is the inverse leg, and the one that
// stops the discard from being written against the status code. Every rejected
// token comes back UNAUTHENTICATED, and a scoped exchange that resolves to no
// capabilities comes back NotFound while the presented token is still the live
// head. Discarding on anything but the superseded reason destroys a working
// session (CUS-839).
func TestRefresh_KeepsTokenOnOtherErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"scoped NotFound, parent still live", errors.New("rpc error: code = NotFound desc = no capabilities")},
		{"transient", errors.New("rpc error: code = Unavailable desc = connection refused")},
		{"rejected but not superseded", errors.New("rpc error: code = Unauthenticated desc = suspended")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			isolateCache(t)
			session := login(issuerA, "alice", audAPI, audReg)
			presented := familyToken(t, session)
			mustSaveFamily(t, presented, issuerA, audAPI, audReg)

			err := Refresh(t.Context(), audAPI, func(context.Context, string) (string, error) {
				return "", tc.err
			}, WithIssuer(issuerA))
			if !errors.Is(err, tc.err) {
				t.Fatalf("Refresh() = %v, want the exchange's own error", err)
			}

			// Still there, and still the token a later refresh presents.
			got, err := loadString(t, audAPI, WithIssuer(issuerA))
			if err != nil || got != presented {
				t.Fatalf("Load(%s) = %q, %v; want the token kept", audAPI, got, err)
			}
			var presentedAgain string
			if err := Refresh(t.Context(), audReg, func(_ context.Context, tok string) (string, error) {
				presentedAgain = tok
				return "", nil
			}, WithIssuer(issuerA)); err != nil {
				t.Fatalf("second Refresh() = %v, want nil", err)
			}
			if presentedAgain != presented {
				t.Errorf("second Refresh presented %q, want the kept token %q", presentedAgain, presented)
			}
		})
	}
}
