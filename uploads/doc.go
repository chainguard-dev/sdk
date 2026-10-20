/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package uploads provides helpers for client-side-encrypted upload
// storage. Both formats it implements wrap a fresh AES-256 data
// encryption key with RSA-OAEP-SHA256 under the org's public key.
//
// SealEnvelope / ParseEnvelope / OpenEnvelope are the v1 JSON envelope,
// which mirrors the browser's encryption.ts: one AES-GCM ciphertext
// base64-encoded into a JSON object, so the whole payload is in memory
// at once.
//
// SealStream / OpenStream are the v2 binary stream envelope, for
// payloads too large to buffer. A header carries the format version, key
// version, wrapped key and chunk size, followed by AES-GCM chunks in the
// STREAM construction: each chunk's nonce derives from its index and a
// final-chunk flag, so a truncated, reordered or spliced stream fails
// to open rather than decrypting to a plausible prefix. NewCRC32CWriter
// gives a caller the ciphertext's CRC32C without buffering it.
//
// Both Open functions accept an unwrap callback so an org's RSA private key
// can stay in KMS. OpenStream passes the stored CryptoKeyVersion to that
// callback, allowing uploads to remain readable after key rotation.
//
// The stream nonce is a 7-byte random per-stream prefix, a 4-byte
// big-endian chunk index and a 1-byte final flag. The prefix defends
// against one thing, and it is not a live threat in this package:
// because SealStream generates a fresh DEK per stream, the index and
// the flag alone already make every nonce unique under that key, and
// the prefix adds nothing an attacker can reach. It is carried so that
// a caller who ever seals two streams under one DEK — a supplied key, a
// cached KMS key — gets a 56-bit birthday bound on nonce reuse instead
// of a guaranteed repeat on chunk 0. Read it as nothing else; in
// particular it is not what keeps chunks from moving between streams,
// which is the distinct DEK.
//
// The 4-byte index is what the prefix costs: 2^32 chunks, 4 PiB at the
// 1 MiB chunk size, enforced on both the sealing and the opening side.
// Spending all 11 remaining bytes on a counter would lift that cap, and
// a cap is the point — 2^32 chunks of 1 MiB is 2^48 AES blocks under
// one DEK, which still sits inside AES-GCM's bounds for deterministic
// nonces, while an uncapped counter does not.
package uploads
