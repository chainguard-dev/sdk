/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"
	"testing"

	delegate "chainguard.dev/go-grpc-kit/pkg/options"
	iam "chainguard.dev/sdk/proto/platform/iam/v1"
	iamtest "chainguard.dev/sdk/proto/platform/iam/v1/test"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type MockPlatformServer struct {
	GroupServer iamtest.MockGroupsServer
}

// StartPlatformServer starts the mock server and returns the connection string.
func (m *MockPlatformServer) StartPlatformServer(_ context.Context, t *testing.T) string {
	t.Helper()
	lis := bufconn.Listen(1024 * 1024)
	testScheme := delegate.RegisterListenerForTest(lis)
	t.Cleanup(func() {
		delegate.UnregisterTestListener(testScheme)
	})
	s := grpc.NewServer()
	iam.RegisterGroupsServer(s, m.GroupServer)
	go func() {
		t.Helper()
		if err := s.Serve(lis); err != nil {
			panic(fmt.Sprintf("Server exited with error: %v", err))
		}
	}()
	t.Cleanup(func() {
		s.GracefulStop()
	})

	return fmt.Sprintf("%s://%s", testScheme, lis.Addr().String())
}
