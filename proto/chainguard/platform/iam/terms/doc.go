/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package terms provides version-agnostic types and error inspection for
// Chainguard Terms-of-Service acceptance.
//
// The [ErrorDetail] interface is the contract each API version fulfills by
// implementing MissingTermsDocs on its TermsNotAcceptedDetail proto type.
// [IsTermsNotAccepted] inspects gRPC errors using this interface, so it works
// with any version — past, present, or future — without modification.
package terms
