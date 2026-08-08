/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package test provides mock implementations of the ping v2
// service clients for use in unit tests. Use [MockPingServiceClients]
// as a drop-in replacement for [v2.Clients] in tests that depend
// on the ping client interface.
package test
