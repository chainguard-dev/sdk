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
	"slices"
	"strings"
	"time"

	"github.com/chainguard-dev/clog"

	"chainguard.dev/sdk/auth"
)

// A refresh-token family is every token descended by rotation from one login.
// The store keeps one file per family under familiesDir, holding the family's
// current token and the audiences it serves, so one login's token is stored
// once rather than copied into every audience directory, where a rotation
// would leave each other copy holding a superseded token. The design is
// docs/customer-platform/chainctl-credential-store.md in chainguard-dev/mono.
//
// Both names start with a dot so that walks over the audience directories --
// DeleteAll's, ListAudiences', and cg's ClearPlatformTokens -- never take them
// for an audience.
const (
	familiesDir = ".refresh-families"

	// lockFileName is the store lock. It never holds token bytes and is never
	// removed: unlinking a lock file another process holds lets the next
	// process lock a new file while the holder still believes it excludes
	// everyone. Its name differs from the refresh.lock cg zendesk keeps in its
	// audience directory.
	lockFileName = ".refresh-families.lock"
)

var (
	// ErrLockUnavailable is returned when the store lock could not be taken:
	// another process held it past the timeout, the context ended, or the
	// filesystem does not support locking. Nothing was exchanged. A caller
	// should fail with a retryable error rather than serve a cached access
	// token, which may be for another scope, issuer or principal, and must
	// not exchange the refresh token without the lock.
	ErrLockUnavailable = errors.New("token: the refresh-token store lock is unavailable")

	// ErrSuccessorNotStored is returned by Refresh when the exchange succeeded
	// but the successor it returned could not be written. The family has been
	// discarded, so the superseded token is never presented again, and the
	// access token the exchange minted is still valid; the next refresh needs
	// a new login.
	ErrSuccessorNotStored = errors.New("token: the rotated refresh token could not be stored")

	// ErrPresentedTokenSuperseded is what an ExchangeFunc returns, or wraps, to
	// report that the issuer refused the presented token because it has already
	// been superseded. Refresh then discards the family, which is the whole
	// point of the signal: the token is dead and nothing revives it, so keeping
	// it only buys a presentation after the issuer's grace window, which is
	// replay and revokes the family.
	//
	// Recognizing the refusal is the caller's job rather than this package's.
	// It arrives as a detail of whatever protocol the exchange speaks, and this
	// package deliberately knows nothing about that.
	ErrPresentedTokenSuperseded = errors.New("token: the presented refresh token was superseded")

	errNoIssuer = errors.New("token: refresh-token families require WithIssuer")
	errAlias    = errors.New("token: refresh-token families do not support WithAlias")
)

// ExchangeFunc presents refreshToken to the issuer and returns the refresh
// token the issuer answered with. Returning "" or refreshToken itself stores
// nothing, which is how a caller that found it no longer needs an exchange
// declines one. It runs with the store lock held, so it must not call back
// into Refresh, SaveFamily or a KindRefresh Save or Delete.
type ExchangeFunc func(ctx context.Context, refreshToken string) (successor string, err error)

// Refresh loads the refresh token from the issuer named by WithIssuer that
// covers audience, passes it to exchange, and stores the successor exchange
// returns, all while holding a lock shared with every other process using the
// cache. The token is loaded after the lock is taken, so a process that waited
// on another's rotation presents the successor, never the parent it would have
// read before waiting.
//
// It returns exchange's error unchanged, ErrLockUnavailable without calling
// exchange, an error wrapping fs.ErrNotExist when no token covers audience, or
// ErrSuccessorNotStored after a successful exchange whose result could not be
// written.
//
// The one exchange error it acts on is ErrPresentedTokenSuperseded: the family
// and the per-audience copies of its token are discarded before the error is
// returned, so the refused token is never presented again. The error still
// wraps what exchange returned, so a caller can re-authenticate on it as
// before; what it can no longer do is retry with the same token.
func Refresh(ctx context.Context, audience string, exchange ExchangeFunc, opts ...Option) error {
	t := newToken(opts...)
	if err := t.familyCheck(); err != nil {
		return err
	}
	unlock, err := lockStore(ctx)
	if err != nil {
		return err
	}
	defer unlock()

	set, err := t.adoptedFamilies(audience)
	if err != nil {
		return err
	}
	f := set.covering(t.issuer, audience)
	if f == nil {
		return fmt.Errorf("no refresh token from %s for %s: %w", t.issuer, audience, fs.ErrNotExist)
	}

	presented := f.Token
	successor, err := exchange(ctx, presented)
	if err != nil {
		if errors.Is(err, ErrPresentedTokenSuperseded) {
			// Discard here rather than leaving it to whoever handles the
			// error. A caller with no terminal cannot re-authenticate, and
			// re-authenticating is not the same as discarding: nothing else on
			// a rejection path deletes a stored refresh token, so the token
			// would survive to be presented again on the next run, after the
			// grace window, as replay. The copies go with it for the reason
			// persist gives -- a copy whose family is gone is adopted as a new
			// family and presented again.
			return errors.Join(err, removeFamily(f.name), removeCopiesOf(presented))
		}
		return err
	}
	successor = strings.TrimSpace(successor)
	if successor == "" || successor == presented {
		return nil
	}
	return persist(ctx, f.name, presented, successor)
}

// SaveFamily stores tok, a refresh token obtained from the issuer named by
// WithIssuer, as that issuer's family serving audiences. Those audiences are
// removed from every other family of that issuer, so each audience resolves
// to exactly one token, and tok is written to each one's per-audience file
// unless that file holds a live token from another issuer and the audience is
// not the one named by WithPrimaryAudience. With WithRotationCapable, each of
// those per-audience files is deleted instead of written.
//
// The family is recorded under the issuer the caller names, which is where
// the token is presented, not under the token's iss claim: an issuer reached
// through another name, such as a cluster-internal issuer behind a mapped
// domain, mints tokens whose iss is a URL the client never dials.
func SaveFamily(tok []byte, audiences []string, opts ...Option) error {
	t := newToken(opts...)
	return t.saveFamily(string(tok), audiences)
}

func (t token) saveFamily(tok string, audiences []string) error {
	if err := t.familyCheck(); err != nil {
		return err
	}
	tok = strings.TrimSpace(tok)
	claims, ok := parseRefresh(tok)
	if !ok {
		return errors.New("token: not a Chainguard refresh token")
	}
	auds := slices.DeleteFunc(slices.Clone(audiences), func(a string) bool { return a == "" })
	slices.Sort(auds)
	auds = slices.Compact(auds)
	if len(auds) == 0 {
		return errors.New("token: a refresh-token family needs at least one audience")
	}

	unlock, err := lockStore(context.Background())
	if err != nil {
		return err
	}
	defer unlock()

	set, err := t.adoptedFamilies("")
	if err != nil {
		return err
	}
	f := set.holding(t.issuer, tok)
	if f == nil {
		f = newFamily(t.issuer, tok)
	}
	set.cover(f, auds)
	if err := set.flush(); err != nil {
		return err
	}
	// The per-audience copies are for readers that pass no issuer and so
	// never see a family: terraform-provider-chainguard and older SDK
	// callers. They can go once those readers move to Refresh. Adoption
	// reads a copy of the family's token as that family's, and overwriting
	// also replaces a copy of an earlier token, which adoption would
	// otherwise take as a competing login. Rotation deletes the copies and
	// never writes a successor into them: an issuer-less reader that kept
	// one would present it after the family moved on, which is replay.
	//
	// A live copy minted by another issuer is kept: it may be that issuer's
	// only session for a hostname both serve, and an issuer-less reader
	// configured for it would otherwise present this token to it. Issuers are
	// compared by iss claim, since the family's issuer may be another name
	// for the same one. The audience named by WithPrimaryAudience is written
	// regardless: the user logged in to replace the session there.
	//
	// A rotation-capable family gets no copies (see WithRotationCapable). The
	// copies it would have overwritten are deleted rather than left: one
	// holding an earlier login's token would be adopted by the next locked
	// call and take its audience back from this family.
	for _, aud := range auds {
		path, err := t.path(KindRefresh, aud)
		if err != nil {
			return err
		}
		if aud != t.primary && heldByOtherIssuer(path, claims.Issuer) {
			continue
		}
		if t.rotationCapable {
			if err := removeOrTruncate(path); err != nil {
				return fmt.Errorf("removing refresh token %s: %w", path, err)
			}
			continue
		}
		if err := writeFileAtomic(path, []byte(tok)); err != nil {
			return fmt.Errorf("writing refresh token %s: %w", path, err)
		}
	}
	return nil
}

// heldByOtherIssuer reports whether the per-audience file at path holds an
// unexpired Chainguard refresh token whose iss claim is not issuer.
func heldByOtherIssuer(path, issuer string) bool {
	b, err := os.ReadFile(path) //nolint:gosec // G304: reading the user's own token cache
	if err != nil {
		return false
	}
	c, ok := parseRefresh(string(b))
	if !ok || c.Issuer == issuer {
		return false
	}
	return c.Expiry == 0 || time.Unix(c.Expiry, 0).After(time.Now())
}

// usesFamilies reports whether a call for kind goes through the family store
// rather than the per-audience file.
func (t token) usesFamilies(kind Kind) bool {
	return kind == KindRefresh && t.issuer != "" && t.alias == ""
}

func (t token) familyCheck() error {
	switch {
	case t.issuer == "":
		return errNoIssuer
	case t.alias != "":
		return errAlias
	}
	return nil
}

// loadSelected returns the token of the issuer's family covering audience,
// else the per-audience file if it holds a token from that issuer that no
// family under another name for the issuer holds. Refresh's adoption refuses
// such a copy, so reporting it here would have callers retry a refresh that
// can never find it. It takes no lock: every file it reads is replaced by
// rename, never rewritten.
func (t token) loadSelected(audience string) ([]byte, error) {
	set, err := loadFamilies(false)
	if err != nil {
		return nil, err
	}
	if f := set.covering(t.issuer, audience); f != nil {
		return []byte(f.Token), nil
	}
	b, err := t.load(KindRefresh, audience)
	if err != nil {
		return nil, err
	}
	tok := strings.TrimSpace(string(b))
	if c, ok := parseRefresh(tok); !ok || c.Issuer != t.issuer || set.heldByAlias(t.issuer, tok) {
		return nil, fmt.Errorf("reading token file: no refresh token from %s for %s: %w", t.issuer, audience, fs.ErrNotExist)
	}
	return b, nil
}

// uncover removes audience from every family's coverage, or only from the
// families of t.issuer when one is set.
func (t token) uncover(audience string) error {
	dir, err := familiesPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	// Removing never presents a token, so a lock that cannot be had does not
	// stop a logout. Without it, though, a family is never rewritten: a
	// rewrite from this snapshot could put back the parent of a successor a
	// lock holder stored after the snapshot, and presenting that parent would
	// revoke the family. So each family that covers audience is removed
	// whole instead, at the cost of its other audiences.
	unlock, err := lockStore(context.Background())
	locked := err == nil
	if locked {
		defer unlock()
	}
	set, err := loadFamilies(locked)
	if err != nil {
		return err
	}
	for _, f := range set.all {
		if t.issuer != "" && f.Issuer != t.issuer {
			continue
		}
		if !locked && slices.Contains(f.Audiences, audience) {
			f.Audiences = nil
			set.dirty[f] = struct{}{}
			continue
		}
		set.remove(f, []string{audience})
	}
	return set.flush()
}

// DeleteFamilies removes every refresh-token family and the per-audience
// copies of each family's token, leaving other per-audience token files and
// the store's lock file in place. Every family holds a
// Chainguard platform refresh token, so a caller clearing platform
// credentials by walking audience directories needs this as well: a family
// can cover audiences that have no directory.
func DeleteFamilies() error {
	dir, err := familiesPath()
	if err != nil {
		return err
	}
	if _, err := os.Stat(dir); errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	// Taken so a refresh in flight cannot write its successor back after the
	// family file is removed. Removing never presents a token, so a lock that
	// cannot be had does not stop the clear.
	if unlock, err := lockStore(context.Background()); err == nil {
		defer unlock()
	}
	// A copy left behind would be adopted as a new family by the next locked
	// call, undoing the clear.
	set, err := loadFamilies(false)
	if err != nil {
		return err
	}
	var errs []error
	for _, f := range set.all {
		errs = append(errs, removeCopiesOf(f.Token))
	}
	root, err := os.OpenRoot(dir)
	if err != nil {
		return fmt.Errorf("opening %s: %w", dir, err)
	}
	defer root.Close()
	ents, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return fmt.Errorf("reading %s: %w", dir, err)
	}
	for _, e := range ents {
		switch {
		case e.IsDir():
		case isTempFile(e.Name()):
			if err := root.Remove(e.Name()); err != nil && !errors.Is(err, fs.ErrNotExist) {
				errs = append(errs, fmt.Errorf("removing %s: %w", e.Name(), err))
			}
		default:
			errs = append(errs, removeFamily(e.Name()))
		}
	}
	return errors.Join(errs...)
}

// family is one family file. name is the file's name, a random local
// identifier: nothing in a refresh token is stable across a rotation, so the
// family keeps its name only because each successor is written into the
// same file.
type family struct {
	Issuer    string    `json:"iss"`
	Audiences []string  `json:"aud"`
	Token     string    `json:"token"`
	Created   time.Time `json:"created"`

	name string
}

func newFamily(issuer, tok string) *family {
	return &family{Issuer: issuer, Token: tok, Created: time.Now().UTC(), name: rand.Text()}
}

// familySet is the store's families in memory. Changes are recorded in dirty
// and written by flush.
type familySet struct {
	all   []*family
	dirty map[*family]struct{}
}

func familiesPath() (string, error) {
	base, err := cacheFilePath("")
	if err != nil {
		return "", err
	}
	return filepath.Join(base, familiesDir), nil
}

// loadFamilies reads every family file, newest first. A file that does not
// parse is skipped, and with prune, which only a lock holder may ask for, it
// is also removed: it is what a discard that had to truncate leaves behind.
func loadFamilies(prune bool) (*familySet, error) {
	set := &familySet{dirty: map[*family]struct{}{}}
	dir, err := familiesPath()
	if err != nil {
		return nil, err
	}
	root, err := os.OpenRoot(dir)
	if errors.Is(err, fs.ErrNotExist) {
		return set, nil
	}
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", dir, err)
	}
	defer root.Close()
	ents, err := fs.ReadDir(root.FS(), ".")
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", dir, err)
	}
	for _, e := range ents {
		if e.IsDir() || isTempFile(e.Name()) {
			continue
		}
		f := &family{name: e.Name()}
		b, err := root.ReadFile(e.Name())
		if err == nil {
			err = json.Unmarshal(b, f)
		}
		if err != nil || f.Issuer == "" || f.Token == "" {
			if prune {
				_ = root.Remove(e.Name())
			}
			continue
		}
		set.all = append(set.all, f)
	}
	slices.SortFunc(set.all, func(a, b *family) int {
		if c := b.Created.Compare(a.Created); c != 0 {
			return c
		}
		return strings.Compare(a.name, b.name)
	})
	return set, nil
}

func readFamily(name string) (*family, error) {
	path, err := cacheFilePath(filepath.Join(familiesDir, name))
	if err != nil {
		return nil, err
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f := &family{name: name}
	if err := json.Unmarshal(b, f); err != nil {
		return nil, fmt.Errorf("parsing %s: %w", path, err)
	}
	return f, nil
}

func writeFamily(f *family) error {
	path, err := cacheFilePath(filepath.Join(familiesDir, f.name))
	if err != nil {
		return err
	}
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	return writeFileAtomic(path, b)
}

// removeOrTruncate deletes the file at path, or empties it when it cannot be
// deleted. Windows refuses to remove a file another process has open, which a
// lock-free reader may briefly have, but that share mode permits truncating
// it. An empty family file parses as no family, and an empty refresh token
// file as no token.
func removeOrTruncate(path string) error {
	err := os.Remove(path)
	if err == nil || errors.Is(err, fs.ErrNotExist) {
		return nil
	}
	if terr := os.Truncate(path, 0); terr != nil {
		return errors.Join(err, terr)
	}
	return nil
}

func removeFamily(name string) error {
	path, err := cacheFilePath(filepath.Join(familiesDir, name))
	if err != nil {
		return err
	}
	if err := removeOrTruncate(path); err != nil {
		return fmt.Errorf("removing family %s: %w", path, err)
	}
	return nil
}

func (s *familySet) covering(issuer, audience string) *family {
	for _, f := range s.all {
		if f.Issuer == issuer && slices.Contains(f.Audiences, audience) {
			return f
		}
	}
	return nil
}

func (s *familySet) holding(issuer, tok string) *family {
	for _, f := range s.all {
		if f.Issuer == issuer && f.Token == tok {
			return f
		}
	}
	return nil
}

// heldByAlias reports whether a family under a configured issuer name other
// than issuer holds tok. A token whose iss claim is issuer can only be held
// there through another name for the same issuer.
func (s *familySet) heldByAlias(issuer, tok string) bool {
	return slices.ContainsFunc(s.all, func(f *family) bool { return f.Issuer != issuer && f.Token == tok })
}

// cover gives owner the audiences auds and takes them from every other family
// of its issuer. Coverage is exclusive within an issuer, so selection never
// has to choose between two families.
func (s *familySet) cover(owner *family, auds []string) {
	if !slices.Contains(s.all, owner) {
		s.all = append(s.all, owner)
		s.dirty[owner] = struct{}{}
	}
	for _, f := range s.all {
		if f != owner && f.Issuer == owner.Issuer {
			s.remove(f, auds)
		}
	}
	for _, a := range auds {
		if !slices.Contains(owner.Audiences, a) {
			owner.Audiences = append(owner.Audiences, a)
			s.dirty[owner] = struct{}{}
		}
	}
}

func (s *familySet) remove(f *family, auds []string) {
	n := len(f.Audiences)
	f.Audiences = slices.DeleteFunc(f.Audiences, func(a string) bool { return slices.Contains(auds, a) })
	if len(f.Audiences) != n {
		s.dirty[f] = struct{}{}
	}
}

// flush writes every changed family and deletes any left covering nothing.
func (s *familySet) flush() error {
	var errs []error
	for f := range s.dirty {
		if len(f.Audiences) == 0 {
			errs = append(errs, removeFamily(f.name))
			continue
		}
		errs = append(errs, writeFamily(f))
	}
	s.all = slices.DeleteFunc(s.all, func(f *family) bool { return len(f.Audiences) == 0 })
	clear(s.dirty)
	return errors.Join(errs...)
}

// refreshClaims are the claims of a Chainguard refresh token the store reads:
// the standard-base64 JSON the issuer packs, read without verification, since
// they only decide which stored token to present to the issuer that verifies
// it.
type refreshClaims struct {
	Issuer    string          `json:"iss"`
	Subject   string          `json:"sub"`
	Expiry    int64           `json:"exp"`
	IssuedAt  int64           `json:"iat"`
	Audiences json.RawMessage `json:"aud"`
}

// parseRefresh reports whether tok is a Chainguard refresh token, which is
// all the family store ever holds: anything else in the cache, such as cg
// zendesk's opaque Zendesk tokens, is never adopted or removed by it.
func parseRefresh(tok string) (refreshClaims, bool) {
	var c refreshClaims
	raw, err := base64.StdEncoding.DecodeString(strings.TrimSpace(tok))
	if err != nil || json.Unmarshal(raw, &c) != nil || c.Issuer == "" {
		return refreshClaims{}, false
	}
	c.Issuer = auth.NormalizeIssuer(c.Issuer)
	return c, true
}

// audiences returns the aud claim, which is a string or an array of strings.
func (c refreshClaims) audiences() []string {
	var one string
	if json.Unmarshal(c.Audiences, &one) == nil {
		return []string{one}
	}
	var many []string
	_ = json.Unmarshal(c.Audiences, &many)
	return many
}

// supersedes reports whether c is a later generation of the same family as
// p: a successor keeps its parent's issuer, subject and expiry and advances
// only its issued-at time.
func (c refreshClaims) supersedes(p refreshClaims) bool {
	return c.Issuer == p.Issuer && c.Subject == p.Subject && c.Expiry == p.Expiry && c.IssuedAt > p.IssuedAt
}

// audienceCopy is a refresh token found in an audience directory: the
// per-audience layout every writer used before the family store, which older
// chainctl releases, buildie and terraform-provider-chainguard still write,
// and where SaveFamily writes its compatibility copies.
type audienceCopy struct {
	path   string
	token  string
	claims refreshClaims
	// audiences are the entries of the token's aud claim whose directory
	// this is. The directory name is a lossy encoding of the audience, so the
	// claim is what recovers it.
	audiences []string
}

// audienceCopies returns the Chainguard refresh tokens in the top level of
// each audience directory. Alias subdirectories are not read.
func audienceCopies() ([]audienceCopy, error) {
	base, err := cacheFilePath("")
	if err != nil {
		return nil, err
	}
	ents, err := os.ReadDir(base)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", base, err)
	}
	var copies []audienceCopy
	for _, e := range ents {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		path := filepath.Join(base, e.Name(), string(KindRefresh))
		b, err := os.ReadFile(path) //nolint:gosec // G304: reading the user's own token cache
		if err != nil {
			continue
		}
		tok := strings.TrimSpace(string(b))
		c, ok := parseRefresh(tok)
		if !ok {
			continue
		}
		cp := audienceCopy{path: path, token: tok, claims: c}
		for _, a := range c.audiences() {
			if audienceDir(a) == e.Name() && !slices.Contains(cp.audiences, a) {
				cp.audiences = append(cp.audiences, a)
			}
		}
		copies = append(copies, cp)
	}
	return copies, nil
}

// adoptedFamilies loads the families and folds into them the per-audience
// copies of t.issuer's tokens. It must be called with the store lock held.
//
// A copy of a family's current token adds its audience to that family. A copy
// of a token the family has since rotated past is a superseded parent and is
// deleted. Any other live copy starts a family covering only the audiences
// where a copy of it was found, which may be narrower than its aud claim: a
// login that wrote no refresh token deleted the copy for its own audience, and
// claiming that audience back would hand the next refresh someone else's
// session. Adopted copies stay on disk for the writer that put them there
// until persist deletes them after the family rotates.
//
// The audience of a copy is recovered from its aud claim. The caller's own
// audience, when given, is attributed to the copy in its directory whatever
// the claim says, since the caller has just named it.
func (t token) adoptedFamilies(audience string) (*familySet, error) {
	set, err := loadFamilies(true)
	if err != nil {
		return nil, err
	}
	copies, err := audienceCopies()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	groups := map[string][]audienceCopy{}
	var order []string
	for _, cp := range copies {
		if audience != "" && len(cp.audiences) == 0 && filepath.Base(filepath.Dir(cp.path)) == audienceDir(audience) {
			cp.audiences = []string{audience}
		}
		if cp.claims.Issuer != t.issuer || len(cp.audiences) == 0 {
			continue
		}
		if cp.claims.Expiry != 0 && time.Unix(cp.claims.Expiry, 0).Before(now) {
			continue
		}
		if _, ok := groups[cp.token]; !ok {
			order = append(order, cp.token)
		}
		groups[cp.token] = append(groups[cp.token], cp)
	}

	var errs []error
	for _, tok := range order {
		group := groups[tok]
		var auds []string
		for _, cp := range group {
			auds = append(auds, cp.audiences...)
		}
		if f := set.holding(t.issuer, tok); f != nil {
			set.cover(f, auds)
			continue
		}
		// A family under another configured issuer -- another name for the
		// same one, since the claim matched -- already holds this token. A
		// second family for it would be left holding the superseded parent
		// once the first rotates, and presenting that revokes both.
		if set.heldByAlias(t.issuer, tok) {
			continue
		}
		if set.supersededBy(t.issuer, group[0].claims) {
			for _, cp := range group {
				if err := os.Remove(cp.path); err != nil && !errors.Is(err, fs.ErrNotExist) {
					errs = append(errs, fmt.Errorf("removing superseded refresh token %s: %w", cp.path, err))
				}
			}
			continue
		}
		set.cover(newFamily(t.issuer, tok), auds)
	}
	errs = append(errs, set.flush())
	return set, errors.Join(errs...)
}

// supersededBy reports whether some family of issuer holds a later
// generation of the token whose claims are c.
func (s *familySet) supersededBy(issuer string, c refreshClaims) bool {
	for _, f := range s.all {
		if f.Issuer != issuer {
			continue
		}
		if fc, ok := parseRefresh(f.Token); ok && fc.supersedes(c) {
			return true
		}
	}
	return false
}

// persistAttempts and persistBackoff bound how long persist retries a failed
// write, which on Windows is most often a lock-free reader briefly holding
// the family file open.
var (
	persistAttempts = 3
	persistBackoff  = 50 * time.Millisecond
)

// persist replaces presented with successor in the family called name. It is
// a compare-and-set: the file is re-read and written only if it still holds
// presented, so a write is never lost to, nor clobbers, one made by a process
// the lock failed to exclude.
//
// A successor that cannot be written is fatal to the family: the issuer has
// already superseded presented, and presenting it again is replay, which
// revokes the family. So the family is discarded and ErrSuccessorNotStored
// returned, and the next refresh needs a new login.
func persist(ctx context.Context, name, presented, successor string) error {
	// presented is superseded whatever happens to the family, and a copy of
	// it with no family holding its successor would be adopted as a new
	// family and presented again. A failed removal only warns: when the
	// successor is stored, the next adoption recognizes the copy as its
	// parent and deletes it.
	dropCopies := func() {
		if err := removeCopiesOf(presented); err != nil {
			clog.WarnContextf(ctx, "removing superseded refresh-token copies: %v", err)
		}
	}
	var err error
	for attempt := range persistAttempts {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt) * persistBackoff)
		}
		var f *family
		f, err = readFamily(name)
		if errors.Is(err, fs.ErrNotExist) {
			// Removed by a logout or clear, which nothing should undo.
			dropCopies()
			return nil
		}
		if err != nil {
			continue
		}
		if f.Token != presented {
			clog.WarnContextf(ctx, "refresh-token family %s changed during an exchange that held the store lock; is locking supported on this filesystem?", name)
			dropCopies()
			return nil
		}
		f.Token = successor
		if err = writeFamily(f); err == nil {
			dropCopies()
			return nil
		}
	}
	return errors.Join(
		fmt.Errorf("%w: %w", ErrSuccessorNotStored, err),
		removeFamily(name),
		removeCopiesOf(presented),
	)
}

// removeCopiesOf deletes every per-audience copy of tok, which the caller has
// just superseded. Once its family is discarded, a copy left intact would be
// adopted as a new family by the next locked call, and its superseded token
// presented again.
func removeCopiesOf(tok string) error {
	copies, err := audienceCopies()
	if err != nil {
		return err
	}
	var errs []error
	for _, cp := range copies {
		if cp.token != tok {
			continue
		}
		if err := removeOrTruncate(cp.path); err != nil {
			errs = append(errs, fmt.Errorf("removing superseded refresh token %s: %w", cp.path, err))
		}
	}
	return errors.Join(errs...)
}
