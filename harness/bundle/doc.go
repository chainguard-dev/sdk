/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

// Package bundle implements harness-source-v1, the single artifact a
// "chainctl harness scan" submission becomes whatever source kind it was built
// from: a zstd-compressed tar of normalized entries with a manifest first.
//
// chainctl packs a bundle with [Writer] and the pipeline streams it back
// through [Validate]. Both ends link this package so the rules below have one
// implementation; the server validates every bundle regardless of what the
// client checked.
//
// # The tar
//
// The manifest is the first entry, at [ManifestPath]; that path is reserved and
// a source tree containing it is rejected. Every later entry is sorted after it
// by [Entry.Path] in byte order, and each path appears once. The ordering key is
// the entry's path, not the name it takes in the tar: a directory is stored with
// a trailing slash, and "cmd/" sorts after "cmd.go" where "cmd" sorts before it.
// Sorting on the tar name would put a directory and its siblings in a different
// order and reject a valid bundle.
//
// Paths are relative, slash-separated and clean: no leading slash, no ".." or
// "." segment, no empty segment, no NUL, no backslash, and valid UTF-8. A
// directory is stored under a name with one trailing slash, which is not part
// of its [Entry] path.
//
// Entry paths and link targets reject Windows drive prefixes (including
// drive-relative forms such as "C:outside"), backslashes, and components ending
// in a space or dot. Exact "." and ".." components remain valid in relative
// link targets. These rules apply on every host: Windows drive resolution,
// backslash separators and trailing-character aliases can otherwise change
// the tree the validator checks. Intake drops these unrepresentable entries.
//
// Only three entry types exist. A regular file has mode 0644, or 0755 if any
// execute bit was set on the source; a directory has mode 0755; a symlink has
// mode 0777 and a relative target that resolves inside the tree. A hardlink, a
// device, a FIFO, a socket and an escaping symlink are rejected here — intake
// drops them with a warning before they reach the writer, by the same rule:
// [Normalize] is what decides, for one entry, whether the format can carry it.
// The writer and reader also resolve symlink chains across the complete tree,
// following each link before applying a subsequent "..". Resolution follows at
// most 40 links; cycles and longer chains are rejected. No entry may have a
// regular file or symlink as an ancestor, so entries cannot alias each other
// through a link. Paths are compared case- and normalization-insensitively
// here, because Windows and macOS resolve names that way, and a path the
// bundle uses as a directory cannot also name a file or a link under that
// comparison. Nor may one directory appear under two spellings, whether an
// entry declares it or a descendant implies it: such a host keeps one of them,
// so the paths that land would not be the paths tree_sha256 pins.
// Streaming extraction must use os.Root as described by
// [Reader.Read], because forward links cannot be checked until the end.
//
// Metadata is zeroed: mtime is the Unix epoch, uid and gid are 0, uname and
// gname are empty, and no PAX record other than path and linkpath is written or
// accepted.
//
// After the end-of-archive marker a reader reads on to EOF and requires every
// remaining byte to be zero, within [Limits.MaxTrailingBytes]. Refusing padding
// outright would reject a conforming writer — GNU tar pads its output to a whole
// 20-block record by default — and leaving it unread would let a bundle carry
// bytes that no entry declares and tree_sha256 does not cover. What is left of
// that channel is the padding's length, a few bits that carry nothing on their
// own; the bundle's own SHA-256, declared at intake and checked after
// decryption, covers them along with every other byte.
//
// # Determinism
//
// The tar layer is reproducible: one tree, packed by any caller in any
// enumeration order, is the same tar bytes, because entry order, modes, mtimes,
// ownership and the PAX records written are all pinned here. That holds for a
// given Go toolchain — the exact framing comes from archive/tar, which is free
// to change when it writes an extended header, so a future toolchain may frame
// the same tree differently.
//
// The compressed bundle is not pinned. zstd's output depends on the encoder's
// version and its level, and the format pins the codec rather than the ratio
// (see [Writer.Write]), so two clients on different builds can send different
// bytes for one tree and both be correct. Nothing in the format compares
// compressed bundles.
//
// What pins exact bytes end to end is the bundle's own SHA-256, declared at
// intake and checked after decryption; what pins the tree independently of how
// it was compressed is tree_sha256.
//
// # tree_sha256
//
// tree_sha256 lets a reader confirm the tar it streamed is the tree the
// manifest describes, independently of how it was compressed. Both ends and any
// reimplementation must agree on it byte for byte, so it is defined exactly:
//
//	tree_sha256 = lowercase hex of SHA-256 over
//
//	    field("harness-source-v1-tree") || record(e₁) || record(e₂) || …
//
//	over every entry of the bundle except the manifest, in the one order the
//	format allows them in the tar — ascending by path, comparing paths as
//	unsigned bytes — where
//
//	    record(e) = field(path) || field(mode) || field(type) || field(payload)
//	    field(s)  = decimal(len(s)) || ":" || s || "\n"
//
//	    decimal the byte length of s in ASCII digits, shortest form: no
//	            leading zero, no sign, no padding ("0" for an empty s)
//	    ":"     one ASCII colon (0x3A), "\n" one ASCII newline (0x0A)
//	    s       the field's bytes, verbatim and unescaped, including any
//	            colon or newline of its own
//
//	    path    the entry's path, with no leading or trailing slash (a
//	            directory's trailing "/" in the tar is not part of it)
//	    mode    the normalized mode as four octal digits, one of exactly
//	            "0644" | "0755" | "0777"
//	    type    one of exactly "file" | "dir" | "symlink"
//	    payload a file's content as 64 lowercase hex digits of SHA-256;
//	            a directory's empty string;
//	            a symlink's target as stored in the tar linkpath, byte for
//	            byte, with no cleaning and no trailing slash
//
//	Nothing else is hashed: not the entry's size, not the manifest's own
//	fields, not any tar framing. A bundle with no entries hashes the domain
//	field alone.
//
// Every field is length-prefixed, as in chainguard.dev/sdk/skills.HardenJobID,
// and every record is exactly four fields, so the byte stream parses back to
// one entry list and no field or record boundary is ambiguous — colons and
// newlines inside a path or a link target are ordinary bytes the length already
// accounted for. Two distinct trees therefore cannot produce the same stream,
// and no payload of one type can be read as another: type is itself a field, so
// a symlink whose target is 64 hex digits and a file whose content hashes to
// them differ in the record before the payload. The leading domain tag pins the
// scheme, so a future change to the record layout yields distinct hashes rather
// than colliding with a v1 bundle's.
//
// # Limits
//
// [DefaultLimits] is the one policy both ends apply; no call site picks its own.
// Every limit is enforced before the thing it bounds exists: the reader checks
// each one against the stream so far, on every read, so a decompression bomb is
// refused while it inflates rather than once it has; the writer refuses a byte
// that would carry the tar stream or the compressed bundle past its limit
// instead of counting it afterwards. The one allowance is
// [Limits.RatioFloorBytes], the decoded size below which the ratio is not
// judged; that floor is all a bomb gets for free, and no more, because
// [Limits.MaxDecodedBytes] still bounds the stream it has to pay for.
// The writer also streams its compressed output through [Reader.Read] and
// waits for validation before returning a manifest, so it enforces the same
// prefix ratio policy without buffering the complete bundle.
//
// A rejection is an [Error] carrying a stable [Reason] the caller can branch on
// and surface.
package bundle
