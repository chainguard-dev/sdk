/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package test

import (
	"context"
	"fmt"

	longrunningpb "cloud.google.com/go/longrunning/autogen/longrunningpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/testing/protocmp"

	guardener "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1"
)

var _ guardener.GuardenerClient = (*MockGuardenerClient)(nil)

type MockGuardenerClient struct {
	OnGetEntitlement          []GuardenerOnGetEntitlement
	OnUpdateEntitlement       []GuardenerOnUpdateEntitlement
	OnMigrateRepository       []GuardenerOnMigrateRepository
	OnGetMigrationOperation   []GuardenerOnGetMigrationOperation
	OnListMigrationOperations []GuardenerOnListMigrationOperations
	OnListScans               []GuardenerOnListScans
	OnGetScan                 []GuardenerOnGetScan
}

type GuardenerOnGetEntitlement struct {
	Given       *guardener.GetEntitlementRequest
	Entitlement *guardener.Entitlement
	Error       error
}

type GuardenerOnUpdateEntitlement struct {
	Given       *guardener.UpdateEntitlementRequest
	Entitlement *guardener.Entitlement
	Error       error
}

type GuardenerOnMigrateRepository struct {
	Given     *guardener.MigrateRepositoryRequest
	Operation *longrunningpb.Operation
	Error     error
}

type GuardenerOnGetMigrationOperation struct {
	Given     *guardener.GetMigrationOperationRequest
	Operation *longrunningpb.Operation
	Error     error
}

type GuardenerOnListMigrationOperations struct {
	Given *guardener.ListMigrationOperationsRequest
	List  *guardener.ListMigrationOperationsResponse
	Error error
}

type GuardenerOnListScans struct {
	Given *guardener.ListScansRequest
	List  *guardener.ListScansResponse
	Error error
}

type GuardenerOnGetScan struct {
	Given *guardener.GetScanRequest
	Scan  *guardener.Scan
	Error error
}

func (m MockGuardenerClient) GetEntitlement(_ context.Context, given *guardener.GetEntitlementRequest, _ ...grpc.CallOption) (*guardener.Entitlement, error) {
	for _, o := range m.OnGetEntitlement {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Entitlement, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGuardenerClient) UpdateEntitlement(_ context.Context, given *guardener.UpdateEntitlementRequest, _ ...grpc.CallOption) (*guardener.Entitlement, error) {
	for _, o := range m.OnUpdateEntitlement {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Entitlement, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGuardenerClient) MigrateRepository(_ context.Context, given *guardener.MigrateRepositoryRequest, _ ...grpc.CallOption) (*longrunningpb.Operation, error) {
	for _, o := range m.OnMigrateRepository {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Operation, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGuardenerClient) GetMigrationOperation(_ context.Context, given *guardener.GetMigrationOperationRequest, _ ...grpc.CallOption) (*longrunningpb.Operation, error) {
	for _, o := range m.OnGetMigrationOperation {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Operation, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGuardenerClient) ListMigrationOperations(_ context.Context, given *guardener.ListMigrationOperationsRequest, _ ...grpc.CallOption) (*guardener.ListMigrationOperationsResponse, error) {
	for _, o := range m.OnListMigrationOperations {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGuardenerClient) ListScans(_ context.Context, given *guardener.ListScansRequest, _ ...grpc.CallOption) (*guardener.ListScansResponse, error) {
	for _, o := range m.OnListScans {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.List, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}

func (m MockGuardenerClient) GetScan(_ context.Context, given *guardener.GetScanRequest, _ ...grpc.CallOption) (*guardener.Scan, error) {
	for _, o := range m.OnGetScan {
		if cmp.Equal(o.Given, given, protocmp.Transform()) {
			return o.Scan, o.Error
		}
	}
	return nil, fmt.Errorf("mock not found for %v", given)
}
