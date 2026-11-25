/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"
	longrunningpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/testing/protocmp"
)

// hardenRouteRecorder records which RPC the gateway dispatched to and the
// request it bound, so a test can assert the routing table as a whole.
type hardenRouteRecorder struct {
	skills.UnimplementedSkillsHardenServer

	called string
	req    proto.Message
}

func (r *hardenRouteRecorder) HardenSkill(_ context.Context, req *skills.HardenSkillRequest) (*longrunningpb.Operation, error) {
	r.called, r.req = "HardenSkill", req
	return &longrunningpb.Operation{}, nil
}

func (r *hardenRouteRecorder) GetHardenOperation(_ context.Context, req *skills.GetHardenOperationRequest) (*longrunningpb.Operation, error) {
	r.called, r.req = "GetHardenOperation", req
	return &longrunningpb.Operation{}, nil
}

func (r *hardenRouteRecorder) ListHardenOperations(_ context.Context, req *skills.ListHardenOperationsRequest) (*skills.ListHardenOperationsResponse, error) {
	r.called, r.req = "ListHardenOperations", req
	return &skills.ListHardenOperationsResponse{}, nil
}

func (r *hardenRouteRecorder) CancelHardenOperation(_ context.Context, req *skills.CancelHardenOperationRequest) (*longrunningpb.Operation, error) {
	r.called, r.req = "CancelHardenOperation", req
	return &longrunningpb.Operation{}, nil
}

// TestSkillsHarden_HTTPRoutesReachTheirRPC pins the SkillsHarden routing
// table: the list collection must not be swallowed by Get's harden/** wildcard,
// and Get must still take a nested group's operation name.
func TestSkillsHarden_HTTPRoutesReachTheirRPC(t *testing.T) {
	const (
		org    = "f7d31f768c62310afac57b60fb89324ef4410ce1"
		folder = org + "/9a1e5c2b7d4f6e8a"
		job    = "3b8f0c"
	)
	listQuery := url.Values{
		"group":      {folder},
		"filter":     {`skill_name = "pdf" AND done = true`},
		"page_size":  {"7"},
		"page_token": {"tok"},
	}.Encode()

	for _, tt := range []struct {
		name   string
		method string
		path   string
		want   string
		// wantReq, when set, is the request the RPC must receive.
		wantReq proto.Message
	}{{
		name:   "submit",
		method: http.MethodPost, path: "/skills/v1alpha1/harden", want: "HardenSkill",
	}, {
		name:   "get an operation of an org",
		method: http.MethodGet, path: "/skills/v1alpha1/harden/" + org + "/" + job + "?group=" + org,
		want:    "GetHardenOperation",
		wantReq: &skills.GetHardenOperationRequest{Group: org, Name: "harden/" + org + "/" + job},
	}, {
		name:   "get an operation of a nested folder",
		method: http.MethodGet, path: "/skills/v1alpha1/harden/" + folder + "/" + job + "?group=" + folder,
		want:    "GetHardenOperation",
		wantReq: &skills.GetHardenOperationRequest{Group: folder, Name: "harden/" + folder + "/" + job},
	}, {
		name:   "list binds every query parameter",
		method: http.MethodGet, path: "/skills/v1alpha1/hardenOperations?" + listQuery,
		want: "ListHardenOperations",
		wantReq: &skills.ListHardenOperationsRequest{
			Group:     folder,
			Filter:    `skill_name = "pdf" AND done = true`,
			PageSize:  7,
			PageToken: "tok",
		},
	}, {
		name:   "cancel",
		method: http.MethodPost, path: "/skills/v1alpha1/harden/" + job + ":cancel", want: "CancelHardenOperation",
	}} {
		t.Run(tt.name, func(t *testing.T) {
			recorder := &hardenRouteRecorder{}
			mux := runtime.NewServeMux()
			if err := skills.RegisterSkillsHardenHandlerServer(t.Context(), mux, recorder); err != nil {
				t.Fatalf("registering the handlers: %v", err)
			}

			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader("{}"))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			mux.ServeHTTP(resp, req)

			if recorder.called != tt.want {
				t.Fatalf("%s %s reached %q, want %q (status %d, body %s)",
					tt.method, tt.path, recorder.called, tt.want, resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			if tt.wantReq != nil {
				if diff := cmp.Diff(tt.wantReq, recorder.req, protocmp.Transform()); diff != "" {
					t.Errorf("%s %s bound the wrong request (-want +got):\n%s", tt.method, tt.path, diff)
				}
			}
		})
	}
}
