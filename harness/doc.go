/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package harness holds shared, dependency-light helpers for the Chainguard
// harness hunt API that both the mono server and the chainctl client link:
// HuntName, the deterministic hunt name, and ParseReportPublicKey, the check
// every report key must pass. See HuntName for what the name is derived from
// and why the two callers must hash identically.
package harness
