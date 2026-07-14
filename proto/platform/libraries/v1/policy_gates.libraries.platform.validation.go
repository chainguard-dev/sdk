/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import "errors"

var (
	// ErrAllowEntryInert reports an allow-list entry that cannot affect any
	// decision: an allow entry only relaxes the cooldown and/or malware gate,
	// so one that sets neither does nothing.
	ErrAllowEntryInert = errors.New("allow entry has no effect; it must relax the cooldown gate, the malware gate, or both")

	// ErrJustificationRequired reports a malware-gate relaxation with no
	// justification, which is required for auditability.
	ErrJustificationRequired = errors.New("a justification is required to relax the malware gate")
)

// Validate reports whether the allow-list entry can meaningfully affect a
// policy decision. It is shared by chainctl and the server so the rule lives in
// one place. The wording is deliberately neutral (it names gates, not the
// bypass_*/override-* fields) so both callers can surface it verbatim.
//
// It does not validate Purl; purl syntax and normalization are enforced
// server-side.
func (e *LibraryAllowListEntry) Validate() error {
	if !e.GetBypassCooldown() && !e.GetBypassMalware() {
		return ErrAllowEntryInert
	}
	if e.GetBypassMalware() && e.GetJustification() == "" {
		return ErrJustificationRequired
	}
	return nil
}
