/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package test provides mock implementations of Policies gRPC clients for testing.
//
// # Overview
//
// This package contains mock implementations of the Policies service clients,
// enabling unit testing of code that depends on the Policies API without
// requiring a live gRPC connection. The mocks support configurable responses
// for all Policies operations including policies and bindings management.
//
// # Features
//
//   - Mock implementations of policies.Clients, policies.PoliciesClient,
//     and policies.BindingsClient interfaces
//   - Configurable responses based on request matching using protocmp
//   - Support for error simulation in test scenarios
//   - Zero external dependencies beyond the Policies proto definitions
//
// # Usage
//
// To use the mocks in tests, create a MockPoliciesClients instance and
// configure the expected requests and responses:
//
//	mock := &test.MockPoliciesClients{
//		PoliciesOnClient: test.MockPoliciesClient{
//			OnListPolicies: []test.OnListPolicies{{
//				Given: &policies.PolicyFilter{},
//				List: &policies.PolicyList{
//					Items: []*policies.Policy{
//						{Id: "policy-1", Name: "test-policy"},
//					},
//				},
//			}},
//		},
//	}
//
//	// Use the mock in your code under test
//	policies := mock.Policies()
//	list, err := policies.ListPolicies(ctx, &policies.PolicyFilter{})
//
// # Integration Patterns
//
// The mocks are designed to be used in table-driven tests where different
// scenarios require different mock configurations. Each mock method matches
// requests using protocmp.Transform() for accurate protobuf comparison.
//
// When a request doesn't match any configured mock, the methods return an
// error indicating the mock was not found, helping identify test configuration
// issues quickly.
//
// # Thread Safety
//
// The mock types are not thread-safe. Each test should create its own mock
// instance to avoid concurrent access issues.
package test
