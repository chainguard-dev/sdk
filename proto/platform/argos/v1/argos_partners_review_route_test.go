/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1_test

import (
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	argosv1 "chainguard.dev/sdk/proto/platform/argos/v1"
)

// The staff review methods must stay gRPC-only. A google.api.http option on
// either method, or generating the gateway with generate_unbound_methods,
// would route it and fail this test.
func TestArgosPartnersStaffMethodsHaveNoHTTPRoute(t *testing.T) {
	sd := argosv1.File_argos_partners_platform_proto.Services().ByName("ArgosPartners")
	mux := runtime.NewServeMux()
	if err := argosv1.RegisterArgosPartnersHandlerServer(t.Context(), mux, argosv1.UnimplementedArgosPartnersServer{}); err != nil {
		t.Fatalf("RegisterArgosPartnersHandlerServer: %v", err)
	}

	for _, tc := range []struct {
		name      string
		method    string
		path      string
		wantRoute bool
	}{{
		name:      "bound partner method is routed",
		method:    http.MethodGet,
		path:      "/argos/v1/mitigation-methods",
		wantRoute: true,
	}, {
		name:   "ListPartnerReviews has no route",
		method: http.MethodPost,
		path:   "/chainguard.platform.argos.ArgosPartners/ListPartnerReviews",
	}, {
		name:   "DecidePartnerCapability has no route",
		method: http.MethodPost,
		path:   "/chainguard.platform.argos.ArgosPartners/DecidePartnerCapability",
	}} {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, strings.NewReader("{}")))
			if got := rec.Code != http.StatusNotFound; got != tc.wantRoute {
				t.Errorf("%s %s: got status %d, want routed = %v", tc.method, tc.path, rec.Code, tc.wantRoute)
			}
			if tc.wantRoute {
				return
			}
			// The 404 covers only the unbound route; an option binding any
			// other path is caught here.
			rpc := path.Base(tc.path)
			md := sd.Methods().ByName(protoreflect.Name(rpc))
			if md == nil {
				t.Fatalf("method %s not found", rpc)
			}
			if proto.HasExtension(md.Options(), annotations.E_Http) {
				t.Errorf("%s: got a google.api.http option, want none", rpc)
			}
		})
	}
}
