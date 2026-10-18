/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/harness/v1alpha1/hunts.proto chainguard/platform/harness/v1alpha1/uploads.proto

// Package v1alpha1 contains the harness platform services: the Hunts intake
// API behind `chainctl harness hunt`, which records a customer's source bundle,
// hands back a signed URL to upload the sealed bundle to, and tracks the hunt
// the pipeline runs over it; and the Uploads API that replaces it, which
// records sealed bundles that an operator later starts hunts over.
package v1alpha1
