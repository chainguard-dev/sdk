/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package token manages the on-disk cache of Chainguard OIDC access and
// refresh tokens, organized by audience and optional alias. Refresh tokens
// read and written WithIssuer are kept in refresh-token families, one per
// login, and are exchanged through Refresh, which serializes every process
// sharing the cache.
package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API
