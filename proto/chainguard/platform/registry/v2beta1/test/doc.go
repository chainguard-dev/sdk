/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package test provides mock implementations of the registry v2beta1
// service clients for use in unit tests. Each mock struct embeds the
// corresponding gRPC client interface and uses [test.On] / [test.Match]
// for request matching, allowing tests to define expected request/response
// pairs without a running server.
//
// Use [MockClients] as a drop-in replacement for [v2beta1.Clients] in
// tests that depend on the registry client interface.
package test
