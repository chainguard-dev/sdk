/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package v1 contains the v1 Argos API.

//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true argos_documents.platform.proto
//go:generate protoc -I . -I ../.. -I ../../.. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true argos_osv.platform.proto
package v1
