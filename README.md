# SDK

This repository contains the public [gRPC](https://grpc.io/) protos supporting
our services and packages to ease integration with the Chainguard platform.

## Updating `*.proto` files

After updating a `*.proto` you'll need to update the corresponding generated go
code.

### Using Docker (recommended)

Run codegen in an ephemeral container with all dependencies pre-installed:

```shell
./hack/update-codegen-docker.sh
```

This requires Docker but no local installation of `protoc` or its plugins.

### Running locally

```shell
./hack/update-codegen.sh
```

#### Prerequisites

Install `protoc` [v34.1](https://github.com/protocolbuffers/protobuf/releases/tag/v34.1): https://grpc.io/docs/protoc-installation/

We currently require `protoc` v34.1.

Install `protoc` codegen dependencies:

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.11
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.6.1
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.22.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.22.0
```

## Breaking Changes

While we make every effort to maintain backward compatibility and avoid breaking changes, we cannot guarantee that future updates to this SDK will be entirely non-breaking. As our platform evolves and new features are added, some modifications to the API surface may be necessary. We recommend pinning to specific versions in production environments and thoroughly testing updates before deployment.
