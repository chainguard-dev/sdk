/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

/*
Package v2beta1 provides Go clients for the Chainguard Ping v2beta1 API.

# Overview

This package contains generated protobuf types and gRPC clients for the
Chainguard platform health check endpoint. The Ping service requires no
authentication and is used to verify API availability.

# Services

PingService provides a health check endpoint:
  - Ping: Check API availability (no authentication required)
*/
package v2beta1

//go:generate protoc -I ../../../.. -I ../../../../.. --go_out=../../../.. --go_opt=paths=source_relative --go-grpc_out=../../../.. --go-grpc_opt=paths=source_relative --grpc-gateway_out=../../../.. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. --openapiv2_opt use_allof_for_refs=true,preserve_rpc_order=true,openapi_naming_strategy=fqn,enable_rpc_deprecation=true chainguard/platform/ping/v2beta1/ping.proto
