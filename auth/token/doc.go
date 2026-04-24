/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package token manages the on-disk cache of Chainguard OIDC access and
// refresh tokens, organized by audience and optional alias.
package token //nolint:revive // redefines-builtin-id: collides with go/token, but renaming would break API
