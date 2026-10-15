/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package aeadstream is the chunked AES-256-GCM framing, in the STREAM
// construction, of the v2 upload stream in uploads, factored out so that
// other encrypted formats in the sdk can share it. It frames chunks and
// nothing else; each format owns its header, writes it before the chunks,
// and passes it in as the additional data every chunk authenticates.
//
// Chunk i is sealed under the nonce noncePrefix(7) || uint32 BE i ||
// final(1), where final is 1 only on the last chunk. The index stops a
// chunk moving within a stream and the flag stops a stream being cut
// short. Index and flag repeat in every stream, so under a DEK that seals
// more than one stream only the prefix keeps nonces apart: a caller must
// draw a fresh random prefix per stream (uploads' package doc gives the
// bound that buys).
//
// The 4-byte index caps a stream at 2^32 chunks, enforced on both the
// sealing and the opening side.
package aeadstream
