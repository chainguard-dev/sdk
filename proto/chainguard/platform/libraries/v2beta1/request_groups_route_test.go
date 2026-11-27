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

	v2beta1 "chainguard.dev/sdk/proto/chainguard/platform/libraries/v2beta1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// routeRecorder answers every RPC on the service by recording which one the
// gateway dispatched to, so a test can assert the routing rather than the
// behaviour.
type routeRecorder struct {
	v2beta1.UnimplementedRequestGroupsServiceServer

	called string
}

func (r *routeRecorder) ListRequestedLibraries(context.Context, *v2beta1.ListRequestedLibrariesRequest) (*v2beta1.ListRequestedLibrariesResponse, error) {
	r.called = "ListRequestedLibraries"
	return &v2beta1.ListRequestedLibrariesResponse{}, nil
}

func (r *routeRecorder) ListRequestedLibraryVersions(context.Context, *v2beta1.ListRequestedLibraryVersionsRequest) (*v2beta1.ListRequestedLibraryVersionsResponse, error) {
	r.called = "ListRequestedLibraryVersions"
	return &v2beta1.ListRequestedLibraryVersionsResponse{}, nil
}

func (r *routeRecorder) ListRequestedLibraryVersionCves(context.Context, *v2beta1.ListRequestedLibraryVersionCvesRequest) (*v2beta1.ListRequestedLibraryVersionCvesResponse, error) {
	r.called = "ListRequestedLibraryVersionCves"
	return &v2beta1.ListRequestedLibraryVersionCvesResponse{}, nil
}

// Test_RequestGroups_HTTPRoutesReachTheirRPC pins which RPC answers each HTTP
// route, which the http_rule annotations decide together rather than one at a
// time: runtime.ServeMux.Handle prepends, so a route registered later is
// matched earlier, and the three routes share one collection path that only
// the verb suffix tells apart.
//
// It also pins that the old /requestGroups routes match nothing, so a client
// still calling them gets a 404 from the gateway rather than reaching some
// other method.
func Test_RequestGroups_HTTPRoutesReachTheirRPC(t *testing.T) {
	const (
		org  = "f7d31f768c62310afac57b60fb89324ef4410ce1"
		uid  = "5baca976-0cf6-415e-9b25-a061b2f330ce"
		base = "/libraries/v2beta1/requestGroups"
		libs = "/libraries/v2beta1/requestedLibraries"
	)

	for _, tt := range []struct {
		name   string
		method string
		path   string
		// want is the RPC that must answer, or "" when no route may match.
		want string
	}{{
		name:   "the requested-libraries collection",
		method: http.MethodGet, path: libs, want: "ListRequestedLibraries",
	}, {
		name:   "requested-library versions",
		method: http.MethodGet, path: libs + ":listRequestedLibraryVersions",
		want: "ListRequestedLibraryVersions",
	}, {
		name:   "requested-library version CVEs",
		method: http.MethodGet, path: libs + ":listRequestedLibraryVersionCves",
		want: "ListRequestedLibraryVersionCves",
	}, {
		name:   "removed: list groups",
		method: http.MethodGet, path: base, want: "",
	}, {
		name:   "removed: get a group",
		method: http.MethodGet, path: base + "/" + uid, want: "",
	}, {
		name:   "removed: list a group's items",
		method: http.MethodGet, path: base + "/" + uid + "/items", want: "",
	}, {
		name:   "removed: create",
		method: http.MethodPost, path: base + "/" + org, want: "",
	}, {
		name:   "removed: submit",
		method: http.MethodPost, path: base + "/" + uid + ":submitRequestGroup", want: "",
	}, {
		name:   "removed: update",
		method: http.MethodPatch, path: base + "/" + uid, want: "",
	}, {
		name:   "removed: removeItems",
		method: http.MethodPost, path: base + "/" + uid + ":removeItems", want: "",
	}, {
		name:   "removed: restoreItems",
		method: http.MethodPost, path: base + "/" + uid + ":restoreItems", want: "",
	}, {
		name:   "removed: delete a draft",
		method: http.MethodDelete, path: base + "/" + uid, want: "",
	}, {
		name:   "removed: delete a submitted group",
		method: http.MethodDelete, path: base + "/" + uid + ":deleteSubmittedRequestGroup", want: "",
	}, {
		name:   "removed: refreshCoverage",
		method: http.MethodPost, path: base + "/" + uid + ":refreshCoverage", want: "",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := &routeRecorder{}
			mux := runtime.NewServeMux()
			if err := v2beta1.RegisterRequestGroupsServiceHandlerServer(t.Context(), mux, recorder); err != nil {
				t.Fatalf("registering the handlers: %v", err)
			}

			// An empty JSON object, so a route taking a body would get as far as
			// the handler. The routes that take none ignore it.
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader("{}"))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)

			if recorder.called != tt.want {
				t.Errorf("%s %s reached %q, want %q (status %d, body %s)",
					tt.method, tt.path, recorder.called, tt.want, resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			if tt.want == "" && resp.Code != http.StatusNotFound {
				t.Errorf("%s %s answered %d, want 404 for a removed route", tt.method, tt.path, resp.Code)
			}
		})
	}
}
