/*
Copyright 2023 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package v1 contains the v1 GRPC client and server definitions
// for implementing APK interactions.
//
//go:generate protoc -I . -I ../.. -I ../../.. -I ../../tenant/v1 --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt logtostderr=true --grpc-gateway_opt paths=source_relative --grpc-gateway_opt generate_unbound_methods=true --grpc-gateway_opt omit_package_doc=true --openapiv2_out=. apk.platform.proto
package v1
