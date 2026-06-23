/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1

import (
	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrTermsNotAccepted returns a gRPC FailedPrecondition error carrying the
// list of legal documents that the group still needs to accept, including their
// display metadata (label and URL).
func ErrTermsNotAccepted(missing []terms.Document) error {
	proto := make([]*MissingDocument, len(missing))
	for i, d := range missing {
		proto[i] = &MissingDocument{Id: d.ID, Label: d.Label, Url: d.URL}
	}
	st, err := status.New(codes.FailedPrecondition, "terms not accepted").
		WithDetails(&TermsNotAcceptedDetail{Missing: proto})
	if err != nil {
		return status.Error(codes.FailedPrecondition, "terms not accepted")
	}
	return st.Err()
}
