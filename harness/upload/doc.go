/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package upload seals a harness bundle to an organization's harness key and
// uploads it through the harness Uploads API: GetOrgPublicKey, CreateUpload,
// one PUT of the sealed bundle to the signed URL CreateUpload returns, then
// CompleteUpload. chainctl and Chainguard's operator tooling both call it, so
// every client seals, sizes and checksums an upload the same way.
//
// The sealed bundle is never held whole. Its exact length comes from
// [uploads.SealedSize] before sealing starts, it is sealed while it streams
// to storage, and its CRC32C is taken over the same bytes as they go, so the
// checksum CompleteUpload compares is of exactly what storage received.
//
// The signed URL admits one PUT of a pinned length that cannot overwrite, so
// there is no resume: a failed upload starts again from CreateUpload, and the
// abandoned upload is left pending for the server to reclaim.
package upload
