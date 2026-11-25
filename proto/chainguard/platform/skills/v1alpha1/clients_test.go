/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1alpha1_test

import (
	"context"
	"net"
	"testing"

	"cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/testing/protocmp"

	skills "chainguard.dev/sdk/proto/chainguard/platform/skills/v1alpha1"
)

type hardenServer struct {
	skills.UnimplementedSkillsHardenServer
}

// GetHardenOperation echoes the requested name, so the test checks the request path too.
func (hardenServer) GetHardenOperation(_ context.Context, req *skills.GetHardenOperationRequest) (*longrunningpb.Operation, error) {
	return &longrunningpb.Operation{Name: req.GetName(), Done: true}, nil
}

// TestClientsHarden proves Harden() is wired to the connection, not just non-nil.
func TestClientsHarden(t *testing.T) {
	want := &longrunningpb.Operation{Name: "harden/720909c9f5279097d847ad02a2f24ba8f59de36a/0123abcd", Done: true}

	lis := bufconn.Listen(1 << 20)
	srv := grpc.NewServer()
	skills.RegisterSkillsHardenServer(srv, hardenServer{})
	go func() { _ = srv.Serve(lis) }()
	t.Cleanup(srv.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) { return lis.DialContext(ctx) }),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	got, err := skills.NewClientsFromConnection(conn).Harden().GetHardenOperation(t.Context(), &skills.GetHardenOperationRequest{Name: want.GetName()})
	if err != nil {
		t.Fatalf("GetHardenOperation: %v", err)
	}
	if diff := cmp.Diff(want, got, protocmp.Transform()); diff != "" {
		t.Errorf("GetHardenOperation (-want, +got):\n%s", diff)
	}
}
