/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"testing"

	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
	iamv2b "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	iamv1 "chainguard.dev/sdk/proto/platform/iam/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrTermsNotAccepted_Roundtrip(t *testing.T) {
	err := iamv2b.ErrTermsNotAccepted([]terms.Document{
		{ID: "guardener-tos.v1", Label: "Terms of Service", URL: "https://example.com"},
	})
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected IsTermsNotAccepted to return true")
	}
	if len(docs) != 1 || docs[0].ID != "guardener-tos.v1" {
		t.Fatalf("unexpected docs: %+v", docs)
	}
	if docs[0].Label != "Terms of Service" {
		t.Errorf("got label %q, want %q", docs[0].Label, "Terms of Service")
	}
}

func TestIsTermsNotAccepted_V1Error(t *testing.T) {
	err := iamv1.ErrTermsNotAccepted([]iamv1.TermsDocument{
		{ID: "sfdpa.v1", Label: "Data Privacy Agreement", URL: "https://example.com/dpa"},
	})
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected IsTermsNotAccepted to return true for v1 error")
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
	if docs[0].ID != "sfdpa.v1" || docs[0].Label != "Data Privacy Agreement" {
		t.Errorf("unexpected doc: %+v", docs[0])
	}
}

func TestIsTermsNotAccepted_MultipleDocuments(t *testing.T) {
	err := iamv2b.ErrTermsNotAccepted([]terms.Document{
		{ID: "guardener-tos.v1", Label: "Terms of Service", URL: "https://example.com/tos"},
		{ID: "sfdpa.v1", Label: "Data Privacy Agreement", URL: "https://example.com/dpa"},
	})
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected IsTermsNotAccepted to return true")
	}
	if len(docs) != 2 {
		t.Fatalf("expected 2 docs, got %d", len(docs))
	}
}

func TestIsTermsNotAccepted_EmptyMissing(t *testing.T) {
	err := iamv2b.ErrTermsNotAccepted(nil)
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected true even with empty missing list")
	}
	if len(docs) != 0 {
		t.Errorf("expected 0 docs, got %d", len(docs))
	}
}

func TestIsTermsNotAccepted_FieldPreservation(t *testing.T) {
	want := terms.Document{
		ID:    "custom-doc.v2",
		Label: "Custom Document with Special Ch@rs & <entities>",
		URL:   "https://example.com/legal?doc=custom&version=2",
	}
	err := iamv2b.ErrTermsNotAccepted([]terms.Document{want})
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected IsTermsNotAccepted to return true")
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}
	if docs[0] != want {
		t.Errorf("field values not preserved:\ngot:  %+v\nwant: %+v", docs[0], want)
	}
}

func TestIsTermsNotAccepted_V1AndV2Beta1Interop(t *testing.T) {
	doc := terms.Document{ID: "guardener-tos.v1", Label: "ToS", URL: "https://example.com"}

	v1Err := iamv1.ErrTermsNotAccepted([]iamv1.TermsDocument{
		{ID: doc.ID, Label: doc.Label, URL: doc.URL},
	})
	v2Err := iamv2b.ErrTermsNotAccepted([]terms.Document{doc})

	ok1, docs1 := terms.IsTermsNotAccepted(v1Err)
	ok2, docs2 := terms.IsTermsNotAccepted(v2Err)

	if !ok1 || !ok2 {
		t.Fatalf("both should return true: v1=%v, v2beta1=%v", ok1, ok2)
	}
	if len(docs1) != 1 || len(docs2) != 1 {
		t.Fatalf("both should return 1 doc: v1=%d, v2beta1=%d", len(docs1), len(docs2))
	}
	if docs1[0] != docs2[0] {
		t.Errorf("extracted docs should match:\nv1:      %+v\nv2beta1: %+v", docs1[0], docs2[0])
	}
}

func TestIsTermsNotAccepted_WrongCode(t *testing.T) {
	ok, _ := terms.IsTermsNotAccepted(status.Error(codes.NotFound, "not found"))
	if ok {
		t.Fatal("expected false for NotFound error")
	}
}
