# SDK

This repository contains the public [gRPC](https://grpc.io/) protos supporting
our services, plus packages to ease integration with the Chainguard platform.

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

## Argos Platform Services

The `proto/platform/argos/v1` package exposes two gRPC services over the
platform endpoint.

### ArgosDocuments

`ArgosDocuments` manages organization-scoped, client-side-encrypted documents.
Callers require one of `CAP_ARGOS_DOCUMENTS_CREATE`, `CAP_ARGOS_DOCUMENTS_LIST`,
or `CAP_ARGOS_DOCUMENTS_DELETE` depending on the operation. Documents are
encrypted client-side before upload; the server never holds plaintext.

### ArgosOSV

`ArgosOSV` provides OSV-format vulnerability query access to the Private OSV
shared corpus. All Argos members holding `CAP_ARGOS_OSV_READ` see the same
record set — there is no per-organisation filtering. The service exposes three
RPCs:

| RPC | HTTP | Description |
|-----|------|-------------|
| `Query` | `POST /argos/v1/osv/query` | Query vulnerabilities affecting a package at a version |
| `QueryBatch` | `POST /argos/v1/osv/querybatch` | Batch form of Query |
| `GetVuln` | `GET /argos/v1/osv/vulns/{id}` | Fetch a single OSV record by ID |

Records are returned in [OSV schema](https://ossf.github.io/osv-schema/) format.

## Breaking Changes

While we make every effort to maintain backward compatibility and avoid breaking changes, we cannot guarantee that future updates to this SDK will be entirely non-breaking. As our platform evolves and new features are added, some modifications to the API surface may be necessary. We recommend pinning to specific versions in production environments and thoroughly testing updates before deployment.
