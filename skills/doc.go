/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package skills holds shared, dependency-light helpers for the Chainguard
// skills harden API that both the mono server and the chainctl client link —
// chiefly HardenJobID, the deterministic harden job id. See HardenJobID for how
// the id is computed and why the two callers must hash identically.
package skills
