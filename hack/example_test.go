/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package tools_test

// Example demonstrates the purpose of the tools package.
//
// The tools package pins tool dependencies for the public/sdk module so they
// appear in go.mod and go.sum. Build with the tools tag to include them:
//
//	go build -tags tools ./hack/...
func Example() {
	// This package has no callable API. It uses blank imports under the
	// //go:build tools constraint to track tool dependencies in the module.
}
