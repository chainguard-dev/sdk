/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"fmt"
	"testing"

	"chainguard.dev/sdk/proto/chainguard/platform/iam/terms"
	iamv2b "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1"
	iamv2btest "chainguard.dev/sdk/proto/chainguard/platform/iam/v2beta1/test"
)

func TestCheckTermsAcceptance_AllAccepted(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{
				TermsAcceptances: []*iamv2b.TermsAcceptance{
					{DocId: "guardener-tos.v1"},
					{DocId: "sfdpa.v1"},
				},
			},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1", "sfdpa.v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckTermsAcceptance_SomeMissing(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{
				TermsAcceptances: []*iamv2b.TermsAcceptance{
					{DocId: "guardener-tos.v1"},
				},
			},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1", "sfdpa.v1"})
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected terms not accepted error")
	}
	if len(docs) != 1 || docs[0].ID != "sfdpa.v1" {
		t.Fatalf("unexpected missing docs: %+v", docs)
	}
	if docs[0].Label != "Data Privacy Agreement" {
		t.Errorf("got label %q, want %q", docs[0].Label, "Data Privacy Agreement")
	}
}

func TestCheckTermsAcceptance_NoneAccepted(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1"})
	ok, docs := terms.IsTermsNotAccepted(err)
	if !ok {
		t.Fatal("expected terms not accepted error")
	}
	if len(docs) != 1 || docs[0].ID != "guardener-tos.v1" {
		t.Fatalf("unexpected missing docs: %+v", docs)
	}
}

func TestCheckTermsAcceptance_ListError(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Error: fmt.Errorf("connection refused"),
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1"})
	if err == nil {
		t.Fatal("expected error when ListTermsAcceptances fails")
	}
	ok, _ := terms.IsTermsNotAccepted(err)
	if ok {
		t.Error("ListTermsAcceptances failure should not be a terms-not-accepted error")
	}
}

func TestCheckTermsAcceptance_ExtraAccepted(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{
				TermsAcceptances: []*iamv2b.TermsAcceptance{
					{DocId: "guardener-tos.v1"},
					{DocId: "sfdpa.v1"},
					{DocId: "extra-doc.v1"},
				},
			},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckTermsAcceptance_EmptyRequired(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, nil)
	if err != nil {
		t.Fatalf("empty required list should pass: %v", err)
	}
}

func TestCheckTermsAcceptance_MultiPage(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{
				TermsAcceptances: []*iamv2b.TermsAcceptance{
					{DocId: "guardener-tos.v1"},
				},
				NextPageToken: "page2",
			},
		}, {
			Resp: &iamv2b.ListTermsAcceptancesResponse{
				TermsAcceptances: []*iamv2b.TermsAcceptance{
					{DocId: "sfdpa.v1"},
				},
			},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1", "sfdpa.v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckTermsAcceptance_EarlyExit(t *testing.T) {
	mock := &iamv2btest.MockTermsServiceClient{
		OnListTermsAcceptances: []iamv2btest.TermsOnListTermsAcceptances{{
			Resp: &iamv2b.ListTermsAcceptancesResponse{
				TermsAcceptances: []*iamv2b.TermsAcceptance{
					{DocId: "guardener-tos.v1"},
				},
				NextPageToken: "page2",
			},
		}},
	}
	err := iamv2b.CheckTermsAcceptance(t.Context(), "group-1", mock, []string{"guardener-tos.v1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDocumentMetadata_KnownDocs(t *testing.T) {
	tests := []struct {
		id    string
		label string
	}{
		{"guardener-tos.v1", "Terms of Service"},
		{"sfdpa.v1", "Data Privacy Agreement"},
		{"skills-tos.v1", "Skills Registry Terms of Service"},
	}
	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			doc := terms.DocumentMetadata(tt.id)
			if doc.ID != tt.id {
				t.Errorf("ID = %q, want %q", doc.ID, tt.id)
			}
			if doc.Label != tt.label {
				t.Errorf("Label = %q, want %q", doc.Label, tt.label)
			}
			if doc.URL == "" {
				t.Error("URL should not be empty for known document")
			}
		})
	}
}

func TestDocumentMetadata_UnknownDoc(t *testing.T) {
	doc := terms.DocumentMetadata("unknown-doc.v99")
	if doc.ID != "unknown-doc.v99" {
		t.Errorf("ID = %q, want %q", doc.ID, "unknown-doc.v99")
	}
	if doc.Label != "" {
		t.Errorf("Label = %q, want empty for unknown doc", doc.Label)
	}
	if doc.URL != "" {
		t.Errorf("URL = %q, want empty for unknown doc", doc.URL)
	}
}
