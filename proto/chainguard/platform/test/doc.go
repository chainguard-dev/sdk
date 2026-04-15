/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Generate the proto definitions
//go:generate protoc -I . --go_out=. --go_opt=paths=source_relative exemplar.proto

// Package test contains test-only proto definitions for SDK unit tests.
package test
