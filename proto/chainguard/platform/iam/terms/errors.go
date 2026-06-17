/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorDetail is the interface that each version's TermsNotAcceptedDetail
// proto type implements. Adding support for a new API version requires only
// implementing this method on the new version's proto type — no changes to
// this package.
type ErrorDetail interface {
	MissingTermsDocs() []Document
}

// IsTermsNotAccepted reports whether err signals that required legal documents
// have not been accepted. If it does, it also returns the missing documents
// with their display metadata.
//
// This function is version-agnostic: it checks for any gRPC status detail that
// satisfies the ErrorDetail interface, regardless of which API version
// produced the error.
func IsTermsNotAccepted(err error) (bool, []Document) {
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.FailedPrecondition {
		return false, nil
	}
	for _, detail := range st.Details() {
		if d, ok := detail.(ErrorDetail); ok {
			return true, d.MissingTermsDocs()
		}
	}
	return false, nil
}
