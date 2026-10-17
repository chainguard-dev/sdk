/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package reportenvelope is harness-report-envelope-v1alpha1, the sealed
// object a hunt report is stored and delivered as. The release candidate and
// the delivery object share it, so a re-key can rewrite a key slot without
// touching the ciphertext.
//
// The body is the v2 upload stream's chunk framing unchanged (see
// chainguard.dev/sdk/uploads): AES-256-GCM over 1 MiB plaintext chunks in
// the STREAM construction, with a 7-byte random nonce prefix, a big-endian
// chunk index and a final-chunk flag. What differs from v2 is the header,
// which is in two parts, and how much of it the chunks authenticate.
//
// # Layout
//
// Integers are big-endian.
//
//	immutable prefix, 49 bytes, the additional data of every chunk
//	  magic         4  "CGRE"
//	  version_len   1  32
//	  version      32  [FormatVersion]
//	  aead          1  1, AES-256-GCM in the STREAM construction
//	  nonce_prefix  7  random per object
//	  chunk_size    4  [ChunkSize]
//	mutable key region
//	  epoch         8  the hunt's report_key_epoch
//	  slot_count    1  1 or 2
//	  slot_count times:
//	    kind        1  [OrgKEK] or [CustomerSubmitKey]
//	    key_version 8  org-kek: the CryptoKeyVersion, positive
//	                   customer-submit-key: 0
//	    wrapped_len 2  1..[MaxWrappedKeySize]
//	    wrapped     n  the wrapped report DEK
//	body
//	  AES-GCM chunks to the end of the object
//
// Slot 0 is always the org-kek slot. Slot 1, present only when the hunt has a
// customer key, is the customer-submit-key slot. Each wrapped value is the
// report DEK, prefixed with a tag naming its slot and with the report's
// artifact UUID, encrypted to the slot's recipient. This package treats it as
// opaque.
//
// # Why the key region is outside the AAD
//
// The prefix is the only additional data, so a re-key can rewrite the key
// region and leave the body byte-identical and still opening under the same
// DEK. Every prefix field stays bound because each one changes how the body is
// read: a changed byte fails every chunk.
//
// Leaving the key slots unauthenticated yields only failures to open. A
// rewritten slot wraps either the same DEK or one under which chunk 0 fails. A
// slot moved to another position fails the reader's position check, and
// wrapped values swapped between slots fail their recipient's unwrap or the
// check of the unwrapped tag. A region that grows or shrinks without a matching
// slot count moves where the body is read from.
//
// A slot removed along with its count is noticed by neither the AEAD nor the
// body digest, so [ReadHeader] requires the org-kek slot in every envelope:
// without it, Chainguard could never re-key the report. A customer-submit-key
// slot removed that way leaves an envelope its customer cannot open until a
// re-key restores the slot.
//
// The epoch is unauthenticated too, and opens normally when forged. It is a
// coordination field between the hunt row and the object, and only a
// principal that can write the report bucket can forge it, who could write
// any bytes there anyway.
//
// # The body digest
//
// The release record's digest is [BodyDigest], which covers the body and not
// the header, so a re-key during or after review leaves the record valid.
package reportenvelope
