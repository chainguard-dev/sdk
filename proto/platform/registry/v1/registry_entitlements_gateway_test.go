/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	registry "chainguard.dev/sdk/proto/platform/registry/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestEntitlementImagesIncludeDeletedHTTPBinding(t *testing.T) {
	t.Parallel()
	for _, route := range []string{"catalog-images", "images"} {
		for _, tc := range []struct {
			name       string
			query      string
			want       bool
			wantStatus int
		}{
			{"omitted", "", false, http.StatusOK},
			{"false", "?include_deleted=false", false, http.StatusOK},
			{"true", "?include_deleted=true", true, http.StatusOK},
			{"not a bool", "?include_deleted=notabool", false, http.StatusBadRequest},
		} {
			t.Run(route+"/"+tc.name, func(t *testing.T) {
				t.Parallel()
				server := &entitlementImagesRecorder{}
				mux := runtime.NewServeMux()
				if err := registry.RegisterEntitlementsHandlerServer(t.Context(), mux, server); err != nil {
					t.Fatal(err)
				}
				response := httptest.NewRecorder()
				request := httptest.NewRequestWithContext(t.Context(), http.MethodGet,
					"/registry/v1/entitlements/root/entitlement/"+route+tc.query, nil)
				mux.ServeHTTP(response, request)
				if got := response.Code; got != tc.wantStatus {
					t.Fatalf("HTTP status: got = %d, want = %d; body = %s", got, tc.wantStatus, response.Body)
				}
				if tc.wantStatus != http.StatusOK {
					if server.requests[route] != nil {
						t.Error("handler request: got a request, want none after query decoding failure")
					}
					return
				}
				want := &registry.EntitlementImagesFilter{Parent: "root/entitlement", IncludeDeleted: tc.want}
				if diff := cmp.Diff(want, server.requests[route], protocmp.Transform()); diff != "" {
					t.Errorf("request (-want +got):\n%s", diff)
				}
			})
		}
	}
}

// entitlementImagesRecorder observes the generated HTTP binding for both image
// routes, keyed by route suffix; it implements no listing behavior.
type entitlementImagesRecorder struct {
	registry.UnimplementedEntitlementsServer
	requests map[string]*registry.EntitlementImagesFilter
}

func (s *entitlementImagesRecorder) record(route string, request *registry.EntitlementImagesFilter) (*registry.EntitlementImagesList, error) {
	if s.requests == nil {
		s.requests = make(map[string]*registry.EntitlementImagesFilter, 2)
	}
	s.requests[route] = request
	return &registry.EntitlementImagesList{}, nil
}

func (s *entitlementImagesRecorder) ListEntitlementCatalogImages(_ context.Context, request *registry.EntitlementImagesFilter) (*registry.EntitlementImagesList, error) {
	return s.record("catalog-images", request)
}

func (s *entitlementImagesRecorder) ListEntitlementImages(_ context.Context, request *registry.EntitlementImagesFilter) (*registry.EntitlementImagesList, error) {
	return s.record("images", request)
}
