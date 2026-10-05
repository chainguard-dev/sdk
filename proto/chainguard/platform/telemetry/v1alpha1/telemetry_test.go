/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	telemetry "chainguard.dev/sdk/proto/chainguard/platform/telemetry/v1alpha1"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestListImagesSkipHTTPBinding(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name       string
		query      string
		wantSkip   int32
		wantStatus int
	}{
		{"omitted", "", 0, http.StatusOK},
		{"zero", "&skip=0", 0, http.StatusOK},
		{"page jump", "&skip=250", 250, http.StatusOK},
		{"not a number", "&skip=invalid", 0, http.StatusBadRequest},
		{"overflow", "&skip=2147483648", 0, http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server := &requestRecorder{}
			mux := runtime.NewServeMux()
			if err := telemetry.RegisterTelemetryServiceHandlerServer(t.Context(), mux, server); err != nil {
				t.Fatal(err)
			}
			response := httptest.NewRecorder()
			request := httptest.NewRequestWithContext(t.Context(), http.MethodGet,
				"/telemetry/v1alpha1/root/group/images?page_size=25&page_token=opaque&source_repository=github.com/example/app&categories=IMAGE_CATEGORY_CHAINGUARD"+tc.query, nil)
			mux.ServeHTTP(response, request)
			if got := response.Code; got != tc.wantStatus {
				t.Fatalf("HTTP status: got = %d, want = %d; body = %s", got, tc.wantStatus, response.Body)
			}
			if tc.wantStatus != http.StatusOK {
				if server.request != nil {
					t.Error("handler request: got a request, want none after query decoding failure")
				}
				return
			}
			want := &telemetry.ListImagesRequest{
				Group: "root/group", Categories: []telemetry.ImageCategory{telemetry.ImageCategory_IMAGE_CATEGORY_CHAINGUARD},
				SourceRepository: "github.com/example/app", PageSize: 25, PageToken: "opaque", Skip: tc.wantSkip,
			}
			if diff := cmp.Diff(want, server.request, protocmp.Transform()); diff != "" {
				t.Errorf("request (-want +got):\n%s", diff)
			}
		})
	}
}

// requestRecorder observes the generated HTTP binding; it implements no paging
// or backend validation policy.
type requestRecorder struct {
	telemetry.UnimplementedTelemetryServiceServer
	request *telemetry.ListImagesRequest
}

func (s *requestRecorder) ListImages(_ context.Context, request *telemetry.ListImagesRequest) (*telemetry.ListImagesResponse, error) {
	s.request = request
	return &telemetry.ListImagesResponse{}, nil
}
