/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v2beta1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/testing/protocmp"

	v2beta1 "chainguard.dev/sdk/proto/chainguard/platform/libraries/v2beta1"
)

// deleteRecorder records the DeleteUploadedVersion request the gateway built.
type deleteRecorder struct {
	v2beta1.UnimplementedStorageUsageServiceServer

	got *v2beta1.DeleteUploadedVersionRequest
}

func (r *deleteRecorder) DeleteUploadedVersion(_ context.Context, req *v2beta1.DeleteUploadedVersionRequest) (*v2beta1.DeleteUploadedVersionResponse, error) {
	r.got = req
	return &v2beta1.DeleteUploadedVersionResponse{}, nil
}

// Test_StorageUsage_DeleteUploadedVersion_QueryParameters pins how the gateway
// reads DeleteUploadedVersion's fields from the query string. Package names
// carry "/" and "@", and versions carry "+", so the request has no path
// parameters, and what reaches the RPC is what the query decodes to.
func Test_StorageUsage_DeleteUploadedVersion_QueryParameters(t *testing.T) {
	const (
		org  = "f7d31f768c62310afac57b60fb89324ef4410ce1"
		path = "/libraries/v2beta1/uploadedVersions"
	)
	for _, tt := range []struct {
		name  string
		query string
		want  *v2beta1.DeleteUploadedVersionRequest
	}{{
		name:  "a scoped npm name, unencoded",
		query: "parent=" + org + "&ecosystem=STORAGE_USAGE_ECOSYSTEM_JAVASCRIPT&name=@acme/widget&version=1.2.3",
		want: &v2beta1.DeleteUploadedVersionRequest{
			Parent: org, Ecosystem: v2beta1.StorageUsageEcosystem_STORAGE_USAGE_ECOSYSTEM_JAVASCRIPT,
			Name: "@acme/widget", Version: "1.2.3",
		},
	}, {
		name:  "a Go module path, and a percent-encoded plus",
		query: "parent=" + org + "&ecosystem=STORAGE_USAGE_ECOSYSTEM_GO&name=example.com/acme/widget&version=v2.0.0%2Bincompatible",
		want: &v2beta1.DeleteUploadedVersionRequest{
			Parent: org, Ecosystem: v2beta1.StorageUsageEcosystem_STORAGE_USAGE_ECOSYSTEM_GO,
			Name: "example.com/acme/widget", Version: "v2.0.0+incompatible",
		},
	}, {
		// The reason the field comments say to percent-encode.
		name:  "an unencoded plus reads as a space",
		query: "parent=" + org + "&ecosystem=STORAGE_USAGE_ECOSYSTEM_GO&name=example.com/acme/widget&version=v2.0.0+incompatible",
		want: &v2beta1.DeleteUploadedVersionRequest{
			Parent: org, Ecosystem: v2beta1.StorageUsageEcosystem_STORAGE_USAGE_ECOSYSTEM_GO,
			Name: "example.com/acme/widget", Version: "v2.0.0 incompatible",
		},
	}, {
		name:  "the ecosystem by number",
		query: "parent=" + org + "&ecosystem=2&name=com.example:app&version=1.0",
		want: &v2beta1.DeleteUploadedVersionRequest{
			Parent: org, Ecosystem: v2beta1.StorageUsageEcosystem_STORAGE_USAGE_ECOSYSTEM_JAVA,
			Name: "com.example:app", Version: "1.0",
		},
	}} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := &deleteRecorder{}
			mux := runtime.NewServeMux()
			if err := v2beta1.RegisterStorageUsageServiceHandlerServer(t.Context(), mux, recorder); err != nil {
				t.Fatalf("registering the handlers: %v", err)
			}
			req := httptest.NewRequest(http.MethodDelete, path+"?"+tt.query, nil)
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Fatalf("status = %d (body %s), want 200", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			if diff := cmp.Diff(tt.want, recorder.got, protocmp.Transform()); diff != "" {
				t.Errorf("request (-want, +got):\n%s", diff)
			}
		})
	}
}
