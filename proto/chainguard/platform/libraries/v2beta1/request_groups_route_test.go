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
	"google.golang.org/protobuf/types/known/emptypb"
)

// routeRecorder answers every RPC on the service by recording which one the
// gateway dispatched to, so a test can assert the routing rather than the
// behaviour.
type routeRecorder struct {
	v2beta1.UnimplementedRequestGroupsServiceServer

	called string
}

func (r *routeRecorder) CreateRequestGroup(context.Context, *v2beta1.CreateRequestGroupRequest) (*v2beta1.RequestGroup, error) {
	r.called = "CreateRequestGroup"
	return &v2beta1.RequestGroup{}, nil
}

func (r *routeRecorder) ListRequestGroups(context.Context, *v2beta1.ListRequestGroupsRequest) (*v2beta1.ListRequestGroupsResponse, error) {
	r.called = "ListRequestGroups"
	return &v2beta1.ListRequestGroupsResponse{}, nil
}

func (r *routeRecorder) GetRequestGroup(context.Context, *v2beta1.GetRequestGroupRequest) (*v2beta1.RequestGroup, error) {
	r.called = "GetRequestGroup"
	return &v2beta1.RequestGroup{}, nil
}

func (r *routeRecorder) ListRequestGroupItems(context.Context, *v2beta1.ListRequestGroupItemsRequest) (*v2beta1.ListRequestGroupItemsResponse, error) {
	r.called = "ListRequestGroupItems"
	return &v2beta1.ListRequestGroupItemsResponse{}, nil
}

func (r *routeRecorder) SubmitRequestGroup(context.Context, *v2beta1.SubmitRequestGroupRequest) (*v2beta1.RequestGroup, error) {
	r.called = "SubmitRequestGroup"
	return &v2beta1.RequestGroup{}, nil
}

func (r *routeRecorder) UpdateRequestGroup(context.Context, *v2beta1.UpdateRequestGroupRequest) (*v2beta1.RequestGroup, error) {
	r.called = "UpdateRequestGroup"
	return &v2beta1.RequestGroup{}, nil
}

func (r *routeRecorder) RemoveItems(context.Context, *v2beta1.RemoveItemsRequest) (*v2beta1.RemoveItemsResponse, error) {
	r.called = "RemoveItems"
	return &v2beta1.RemoveItemsResponse{}, nil
}

func (r *routeRecorder) RestoreItems(context.Context, *v2beta1.RestoreItemsRequest) (*v2beta1.RestoreItemsResponse, error) {
	r.called = "RestoreItems"
	return &v2beta1.RestoreItemsResponse{}, nil
}

func (r *routeRecorder) DeleteRequestGroup(context.Context, *v2beta1.DeleteRequestGroupRequest) (*emptypb.Empty, error) {
	r.called = "DeleteRequestGroup"
	return &emptypb.Empty{}, nil
}

func (r *routeRecorder) DeleteSubmittedRequestGroup(context.Context, *v2beta1.DeleteSubmittedRequestGroupRequest) (*emptypb.Empty, error) {
	r.called = "DeleteSubmittedRequestGroup"
	return &emptypb.Empty{}, nil
}

func (r *routeRecorder) RefreshCoverage(context.Context, *v2beta1.RefreshCoverageRequest) (*v2beta1.RequestGroup, error) {
	r.called = "RefreshCoverage"
	return &v2beta1.RequestGroup{}, nil
}

func (r *routeRecorder) ListRequestedLibraries(context.Context, *v2beta1.ListRequestedLibrariesRequest) (*v2beta1.ListRequestedLibrariesResponse, error) {
	r.called = "ListRequestedLibraries"
	return &v2beta1.ListRequestedLibrariesResponse{}, nil
}

func (r *routeRecorder) ListRequestedLibraryVersions(context.Context, *v2beta1.ListRequestedLibraryVersionsRequest) (*v2beta1.ListRequestedLibraryVersionsResponse, error) {
	r.called = "ListRequestedLibraryVersions"
	return &v2beta1.ListRequestedLibraryVersionsResponse{}, nil
}

// Test_RequestGroups_HTTPRoutesReachTheirRPC pins which RPC answers each HTTP
// route, which the http_rule annotations decide together rather than one at a
// time.
//
// Two of grpc-gateway's rules make a route able to swallow its neighbour, and
// neither is visible in the annotation being read:
//
//   - A "**" deep wildcard matches zero segments as happily as many, so
//     "/requestGroups/{uid=**}" also matches the bare collection path.
//   - runtime.ServeMux.Handle prepends, so a route registered later is matched
//     earlier. Registration follows the order the RPCs appear in the service,
//     which makes every route a candidate to shadow the ones declared above it.
//
// Together those hid ListRequestGroups entirely: GetRequestGroup is declared
// after it, so it matched first, and its wildcard accepted a path with no uid
// in it at all. Every REST caller listing groups got "uid is required".
//
// A test of one handler cannot see this. It is a property of the whole table,
// so the whole table is what this asserts.
func Test_RequestGroups_HTTPRoutesReachTheirRPC(t *testing.T) {
	const (
		org  = "f7d31f768c62310afac57b60fb89324ef4410ce1"
		uid  = "5baca976-0cf6-415e-9b25-a061b2f330ce"
		base = "/libraries/v2beta1/requestGroups"
	)

	for _, tt := range []struct {
		name   string
		method string
		path   string
		// want is the RPC that must answer, or "" when no route may match.
		want string
	}{{
		name: "the collection path lists, rather than being eaten by the uid route",
		// The regression. With a "**" uid this reached GetRequestGroup.
		method: http.MethodGet, path: base, want: "ListRequestGroups",
	}, {
		name:   "a uid fetches one group",
		method: http.MethodGet, path: base + "/" + uid, want: "GetRequestGroup",
	}, {
		name:   "items hang off the uid",
		method: http.MethodGet, path: base + "/" + uid + "/items", want: "ListRequestGroupItems",
	}, {
		// uid is a flat UUID on every one of these messages, so a slash in it is
		// not a deeper resource, it is a path that names nothing. It has to miss
		// rather than bind uid to "a/b", which is what the wildcard did.
		name:   "a multi-segment uid matches no route",
		method: http.MethodGet, path: base + "/a/b", want: "",
	}, {
		// parent keeps its deep wildcard, and this is why: a Chainguard UIDP
		// nests, so a sub-organization's create path carries slashes.
		name:   "create takes a nested parent",
		method: http.MethodPost, path: base + "/" + org + "/nested/child", want: "CreateRequestGroup",
	}, {
		name:   "create takes a flat parent too",
		method: http.MethodPost, path: base + "/" + org, want: "CreateRequestGroup",
	}, {
		name:   "submit",
		method: http.MethodPost, path: base + "/" + uid + ":submitRequestGroup", want: "SubmitRequestGroup",
	}, {
		name:   "update",
		method: http.MethodPatch, path: base + "/" + uid, want: "UpdateRequestGroup",
	}, {
		name:   "removeItems",
		method: http.MethodPost, path: base + "/" + uid + ":removeItems", want: "RemoveItems",
	}, {
		name:   "restoreItems",
		method: http.MethodPost, path: base + "/" + uid + ":restoreItems", want: "RestoreItems",
	}, {
		name:   "delete a draft",
		method: http.MethodDelete, path: base + "/" + uid, want: "DeleteRequestGroup",
	}, {
		// The verb is what separates this from the draft delete above, on the
		// same method and the same path.
		name:   "delete a submitted group",
		method: http.MethodDelete, path: base + "/" + uid + ":deleteSubmittedRequestGroup",
		want: "DeleteSubmittedRequestGroup",
	}, {
		name:   "refreshCoverage",
		method: http.MethodPost, path: base + "/" + uid + ":refreshCoverage", want: "RefreshCoverage",
	}, {
		name:   "the requested-libraries collection",
		method: http.MethodGet, path: "/libraries/v2beta1/requestedLibraries", want: "ListRequestedLibraries",
	}, {
		name:   "requested-library versions",
		method: http.MethodGet, path: "/libraries/v2beta1/requestedLibraries:listRequestedLibraryVersions",
		want: "ListRequestedLibraryVersions",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := &routeRecorder{}
			mux := runtime.NewServeMux()
			if err := v2beta1.RegisterRequestGroupsServiceHandlerServer(t.Context(), mux, recorder); err != nil {
				t.Fatalf("registering the handlers: %v", err)
			}

			// An empty JSON object, because the routes taking a body need one to
			// get as far as the handler. The routes that take none ignore it.
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader("{}"))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)

			if recorder.called != tt.want {
				t.Errorf("%s %s reached %q, want %q (status %d, body %s)",
					tt.method, tt.path, recorder.called, tt.want, resp.Code, strings.TrimSpace(resp.Body.String()))
			}
		})
	}
}
