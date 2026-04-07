/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrTermsNotAccepted returns a gRPC FailedPrecondition error carrying the
// list of legal documents that the group still needs to accept, including their
// display metadata (label and URL).
func ErrTermsNotAccepted(missing []TermsDocument) error {
	proto := make([]*MissingDocument, len(missing))
	for i, d := range missing {
		proto[i] = &MissingDocument{Id: d.ID, Label: d.Label, Url: d.URL}
	}
	st, err := status.New(codes.FailedPrecondition, "terms not accepted").
		WithDetails(&TermsNotAcceptedDetail{Missing: proto})
	if err != nil {
		// Fallback if WithDetails fails (should not happen in practice).
		return status.Error(codes.FailedPrecondition, "terms not accepted")
	}
	return st.Err()
}

// IsTermsNotAccepted reports whether err signals that required legal documents
// have not been accepted. If it does, it also returns the missing documents with
// their display metadata.
func IsTermsNotAccepted(err error) (bool, []TermsDocument) {
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.FailedPrecondition {
		return false, nil
	}
	for _, detail := range st.Details() {
		if d, ok := detail.(*TermsNotAcceptedDetail); ok {
			docs := make([]TermsDocument, len(d.Missing))
			for i, m := range d.Missing {
				docs[i] = TermsDocument{ID: m.GetId(), Label: m.GetLabel(), URL: m.GetUrl()}
			}
			return true, docs
		}
	}
	return false, nil
}
