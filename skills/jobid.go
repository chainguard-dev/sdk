/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package skills

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
)

// hardenJobIDVersion is the scheme version of the job-id hash. It is mixed into
// the hash input so a future change to the id scheme (e.g. adding or reordering
// a dimension) yields distinct ids and cannot collide with v1 ids. Bump it only
// alongside a coordinated change on the chainctl client and skillchain2.
const hardenJobIDVersion = "v1"

// HardenNamePrefix is the prefix of every harden operation name. The operation
// name is the job id, namespaced under the target group:
// "harden/{group-uidp}/{jobhash}".
const HardenNamePrefix = "harden/"

// HardenJobID computes the deterministic harden job id from its four content
// dimensions. It is content-addressed, so identical inputs always yield the same
// id — the property the operation store relies on for idempotent submit (a
// re-submit of the same skill into the same group by the same user is a no-op).
//
//   - orgID:         the target group UIDP (authorized against the caller).
//   - userID:        the submitting user's identity subject, so each user gets
//     their own tracked job.
//   - skillName:     the skill's catalog name.
//   - contentDigest: the sha256 digest of the uploaded uploads.cgr.dev artifact.
//
// Two users hardening the same bundle run separate jobs AND land separate
// outputs: each user's hardened artifact is published under a per-user path,
// skills.cgr.dev/<group>/<user>/<skill_name>@sha256:<hardened-digest>, so the two
// never overwrite or race the same digest. The user-less
// skills.cgr.dev/<group>/<skill_name> path is reserved for promotion (blessing a
// hardened skill to the org level). userID being a dimension of both this id and
// the output path is what keeps that per-user isolation intact.
//
// The returned value is a hex sha256 over a version tag and the four fields,
// each length-prefixed so no field boundary is ambiguous (e.g. so a userID
// ending in the skillName's prefix cannot alias a different split).
//
// IMPORTANT: this MUST hash identically to the chainctl client's local
// recomputation and to skillchain2's job key. Any change here is a coordinated
// cross-repo change; bump hardenJobIDVersion rather than silently altering the
// input encoding.
func HardenJobID(orgID, userID, skillName, contentDigest string) string {
	var b strings.Builder
	writeField(&b, hardenJobIDVersion)
	writeField(&b, orgID)
	writeField(&b, userID)
	writeField(&b, skillName)
	writeField(&b, contentDigest)
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

// HardenOperationName returns the full operation name for a harden job: the job
// id namespaced under the target group. Both the store record and the id handed
// to skillchain2 use this name.
func HardenOperationName(orgID, userID, skillName, contentDigest string) string {
	return HardenNamePrefix + orgID + "/" + HardenJobID(orgID, userID, skillName, contentDigest)
}

// GroupFromHardenName extracts the group UIDP embedded in a harden operation
// name. The name is "harden/{group}/{jobhash}"; the group UIDP may itself
// contain "/" (hierarchical UIDPs), so the job hash is the final path segment
// and the group is everything before it. It is the inverse of
// HardenOperationName and the boundary the api-impl cross-tenant read check uses
// to confirm a caller's group owns an operation.
func GroupFromHardenName(name string) (string, error) {
	rest, ok := strings.CutPrefix(name, HardenNamePrefix)
	if !ok {
		return "", fmt.Errorf("invalid harden operation name %q: missing %q prefix", name, HardenNamePrefix)
	}
	i := strings.LastIndex(rest, "/")
	if i <= 0 || i == len(rest)-1 {
		return "", fmt.Errorf("invalid harden operation name %q: missing group/id segments", name)
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
