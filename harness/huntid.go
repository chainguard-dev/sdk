/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package harness

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// huntNameVersion is the scheme version of the hunt-name hash. It is mixed into
// the hash input so a future change to the derivation (adding, dropping or
// reordering a dimension) yields distinct names and cannot collide with v1
// ones. Bump it only alongside a coordinated change on the chainctl client.
const huntNameVersion = "v1"

// HuntNamePrefix is the prefix of every hunt name. A hunt name is the
// content-addressed hash namespaced under the submitting group:
// "hunt/{group-uidp}/{hunthash}".
const HuntNamePrefix = "hunt/"

// HuntName computes the deterministic name of a hunt from its three content
// dimensions, namespaced under the submitting group. It is content-addressed,
// so identical inputs always yield the same name — the property CreateHunt's
// idempotency rests on. The hunt table carries a unique index on this name and
// CreateHunt is insert-if-absent, so resubmitting an identical bundle returns
// the existing hunt instead of starting a second hours-long run.
//
//   - orgID:        the submitting group's UIDP (authorized against the caller).
//     It namespaces the hunt and is what GroupFromHuntName recovers for the
//     cross-tenant read check.
//   - actorID:      the submitting identity's subject. Without it one member's
//     create would silently return another member's hunt, handing them a
//     status and rejection reason they did not produce; with it each submitter
//     gets a hunt of their own to track.
//   - bundleSHA256: the sha256 of the plaintext source bundle, as
//     "sha256:<64-hex>". This is the content dimension: a changed tree is a
//     different hunt, an unchanged tree is the same one.
//
// The ciphertext envelope is deliberately not a dimension. It is sealed under a
// fresh per-upload key, so the same tree seals to different bytes every time and
// a resubmit would never match.
//
// The returned name's hash segment is a hex sha256 over a version tag and the
// three fields, each length-prefixed so no field boundary is ambiguous (so an
// actorID ending in the digest's prefix cannot alias a different split).
//
// IMPORTANT: this MUST hash identically to the chainctl client's local
// recomputation. Any change here is a coordinated cross-component change; bump
// huntNameVersion rather than silently altering the input encoding.
func HuntName(orgID, actorID, bundleSHA256 string) string {
	var b strings.Builder
	writeField(&b, huntNameVersion)
	writeField(&b, orgID)
	writeField(&b, actorID)
	writeField(&b, bundleSHA256)
	sum := sha256.Sum256([]byte(b.String()))
	return HuntNamePrefix + orgID + "/" + hex.EncodeToString(sum[:])
}

// HuntRetryName derives the next attempt from a failed hunt's name. Transport
// retries therefore return the same attempt; a further retry must name the
// newly failed hunt, so one failed hunt has exactly one successor. The caller
// verifies ownership and failure before creating it.
func HuntRetryName(failedHunt string) (string, error) {
	group, err := GroupFromHuntName(failedHunt)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256([]byte("hunt-retry-v1\n" + failedHunt))
	return HuntNamePrefix + group + "/" + hex.EncodeToString(sum[:]), nil
}

// GroupFromHuntName extracts the group UIDP embedded in a hunt name. The name is
// "hunt/{group}/{hunthash}"; the group UIDP may itself contain "/"
// (hierarchical UIDPs), so the hash is the final path segment and the group is
// everything before it. It is the inverse of HuntName and the boundary the
// api-impl cross-tenant read check uses to confirm a caller's group owns a hunt.
func GroupFromHuntName(name string) (string, error) {
	rest, ok := strings.CutPrefix(name, HuntNamePrefix)
	if !ok {
		return "", fmt.Errorf("invalid hunt name %q: missing %q prefix", name, HuntNamePrefix)
	}
	i := strings.LastIndex(rest, "/")
	if i <= 0 || i == len(rest)-1 {
		return "", fmt.Errorf("invalid hunt name %q: missing group/hash segments", name)
	}
	return rest[:i], nil
}

// writeField appends a length-prefixed field to b so field boundaries are
// unambiguous: two distinct field tuples can never produce the same byte stream.
func writeField(b *strings.Builder, s string) {
	// A newline-delimited "len:value" framing: len fixes the boundary, and the
	// separators keep the framing human-inspectable for debugging.
	b.WriteString(strconv.Itoa(len(s)))
	b.WriteByte(':')
	b.WriteString(s)
	b.WriteByte('\n')
}
