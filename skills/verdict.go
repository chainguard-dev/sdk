/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package skills

// GateRejectedMessage is the Operation.error message a harden operation carries
// when the post-evaluation gate rejects the hardened skill. It is deliberately
// distinct from the generic hardening failure so a client can tell a verdict
// ("the gate looked at this and said no") from an outage ("we never got a
// verdict"), and retry only the latter. Both the service that produces it and
// the variance measurement that classifies it read this constant, so the two
// cannot drift apart silently.
const GateRejectedMessage = "post-evaluation gate rejected hardened skill"
