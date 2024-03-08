# SDK

This repository contains the public [gRPC](https://grpc.io/) protos supporting
our services and packages to ease integration with the Chainguard platform.

## Updating `*.proto` files

After updating a `*.proto` you'll need to update the corresponding generated go
code.

```shell
./hack/update-codegen.sh
```

### Prerequisites

Install `protoc`: https://grpc.io/docs/protoc-installation/

Example for MacOS:

```shell
brew install protobuf
```

We currently require `protoc` v21.12.

Install `protoc` codegen dependencies:

```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.10.0
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.10.0
```
