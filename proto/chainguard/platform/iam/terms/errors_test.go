/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package terms_test

import (
	"fmt"
	"testing"

	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestIsTermsNotAccepted_NilError(t *testing.T) {
	ok, _ := terms.IsTermsNotAccepted(nil)
	if ok {
		t.Fatal("expected false for nil error")
	}
}

func TestIsTermsNotAccepted_WrongStatusCode(t *testing.T) {
	ok, _ := terms.IsTermsNotAccepted(status.Error(codes.NotFound, "not found"))
	if ok {
		t.Fatal("expected false for NotFound error")
	}
}

func TestIsTermsNotAccepted_FailedPreconditionWithoutDetail(t *testing.T) {
	ok, _ := terms.IsTermsNotAccepted(status.Error(codes.FailedPrecondition, "other precondition"))
	if ok {
		t.Fatal("expected false for FailedPrecondition without detail")
	}
}

func TestIsTermsNotAccepted_NonGRPCError(t *testing.T) {
	ok, _ := terms.IsTermsNotAccepted(fmt.Errorf("plain error"))
	if ok {
		t.Fatal("expected false for non-gRPC error")
	}
}

func TestDocumentMetadata_Known(t *testing.T) {
	doc := terms.DocumentMetadata("guardener-tos.v1")
	if doc.Label != "Terms of Service" {
		t.Errorf("Label = %q, want %q", doc.Label, "Terms of Service")
	}
	if doc.URL == "" {
		t.Error("URL should not be empty for known document")
	}
}

func TestDocumentMetadata_Unknown(t *testing.T) {
	doc := terms.DocumentMetadata("unknown-doc.v99")
	if doc.ID != "unknown-doc.v99" {
		t.Errorf("ID = %q, want %q", doc.ID, "unknown-doc.v99")
	}
	if doc.Label != "" || doc.URL != "" {
		t.Errorf("unknown doc should have empty Label/URL, got %+v", doc)
	}
}
