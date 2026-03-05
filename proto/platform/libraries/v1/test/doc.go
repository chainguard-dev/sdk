/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package test provides mock implementations of the Libraries API clients
// for testing purposes.
//
// # Overview
//
// This package contains mock clients that implement the same interfaces as
// the production Libraries API clients, allowing you to write tests without
// making actual API calls. The mocks use a simple pattern where you configure
// expected inputs and corresponding outputs.
//
// # Features
//
//   - MockLibrariesClients: Mock implementation of the Clients interface
//   - MockArtifactsClient: Mock for the Artifacts service client
//   - MockEntitlementsClient: Mock for the Entitlements service client
//   - MockNpmPackagesClient: Mock for the NPM Packages service client
//   - Configurable responses based on input matching
//   - Support for error simulation
//
// # Usage
//
// The mock clients use a slice-based configuration pattern. For each method,
// you provide a slice of expected inputs and their corresponding outputs.
// When the method is called, the mock searches for a matching input and
// returns the configured response.
//
// Example:
//
//	mock := test.MockLibrariesClients{
//	    EntitlementsClient: test.MockEntitlementsClient{
//	        OnList: []test.EntitlementsOnList{{
//	            Given: &v1.EntitlementFilter{Parent: "org/abc"},
//	            List:  &v1.EntitlementList{Items: []*v1.Entitlement{...}},
//	        }},
//	    },
//	}
//
//	// Use the mock in your tests
//	list, err := mock.Entitlements().List(ctx, &v1.EntitlementFilter{Parent: "org/abc"})
//
// # Integration Patterns
//
// The mock clients are designed to be drop-in replacements for the real
// clients. Any code that accepts the Clients interface can use these mocks
// for testing:
//
//	func processLibraries(clients v1.Clients) error {
//	    // This function works with both real and mock clients
//	    artifacts := clients.Artifacts()
//	    // ...
//	}
//
//	func TestProcessLibraries(t *testing.T) {
//	    mock := test.MockLibrariesClients{...}
//	    err := processLibraries(&mock)
//	    // Assert on err
//	}
//
// # Error Simulation
//
// Each mock configuration struct has an Error field that allows you to
// simulate error conditions:
//
//	mock := test.MockEntitlementsClient{
//	    OnCreate: []test.EntitlementsOnCreate{{
//	        Given: &v1.CreateEntitlementRequest{...},
//	        Error: errors.New("permission denied"),
//	    }},
//	}
//
// # Input Matching
//
// The mocks use google/go-cmp with protocmp.Transform() to compare inputs.
// This provides deep equality checking for protocol buffer messages. If no
// matching input is found, the mock returns an error indicating the mock
// was not configured for that input.
package test
