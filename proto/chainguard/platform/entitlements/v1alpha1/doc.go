/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/entitlements/v1alpha1/features.proto

// Package v1alpha1 contains the entitlements platform services: Features,
// which serves an organization's FEATURE entitlements (product features such
// as sandbox Checks and Workspaces, with their per-organization limits).
package v1alpha1
