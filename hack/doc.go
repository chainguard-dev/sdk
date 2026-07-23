/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package tools pins tool dependencies for the public/sdk module.
//
// This package uses blank imports with a build tag to ensure that tool
// dependencies are tracked in go.mod and go.sum without being included
// in the final binary. To use the tools, build with the tools tag:
//
//	go build -tags tools ./hack/...
package tools
