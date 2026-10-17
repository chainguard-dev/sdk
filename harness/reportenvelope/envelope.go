/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package reportenvelope

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"slices"

	"chainguard.dev/sdk/internal/aeadstream"
)

const (
	// FormatVersion is the format version string in every envelope's
	// immutable prefix.
	FormatVersion = "harness-report-envelope-v1alpha1"

	// ChunkSize is the body's plaintext chunk size. A reader refuses any
	// other, so the size is part of the format rather than a parameter.
	ChunkSize = 1 << 20

	// MaxWrappedKeySize bounds a slot's wrapped DEK, so a hostile key region
	// cannot name an absurd allocation. 2048 bytes is an RSA-16384 modulus.
	MaxWrappedKeySize = 2048
)

const (
	// magic differs from the v2 upload stream's "CGUS", so a reader handed
	// an upload refuses it before reading anything else.
	magic = "CGRE"

	aeadAES256GCMStream = 1

	noncePrefixAt = len(magic) + 1 + len(FormatVersion) + 1
	chunkSizeAt   = noncePrefixAt + aeadstream.NoncePrefixSize
	prefixSize    = chunkSizeAt + 4

	maxSlots = 2
)

// SlotKind names a key slot's recipient.
type SlotKind uint8

const (
	OrgKEK            SlotKind = 1 // the "org-kek" slot
	CustomerSubmitKey SlotKind = 2 // the "customer-submit-key" slot
)

// slotOrder is the kind each slot position must hold.
var slotOrder = [maxSlots]SlotKind{OrgKEK, CustomerSubmitKey}

// String returns k's slot name, or "SlotKind(n)" for an unknown kind.
func (k SlotKind) String() string {
	switch k {
	case OrgKEK:
		return "org-kek"
	case CustomerSubmitKey:
		return "customer-submit-key"
	default:
		return fmt.Sprintf("SlotKind(%d)", uint8(k))
	}
}

// Slot is the report DEK wrapped to one recipient.
type Slot struct {
	Kind SlotKind
	// KeyVersion is the CryptoKeyVersion of the harness org key that wraps
	// an org-kek slot. It is zero in a customer-submit-key slot.
	KeyVersion uint64
	Wrapped    []byte
}

// KeyRegion is the mutable part of the header: the hunt's report_key_epoch
// and the key slots, the org-kek slot first.
type KeyRegion struct {
	Epoch uint64
	Slots []Slot
}

// Slot returns the slot of the given kind, and whether there is one.
func (k KeyRegion) Slot(kind SlotKind) (Slot, bool) {
	i := slices.IndexFunc(k.Slots, func(s Slot) bool { return s.Kind == kind })
	if i < 0 {
		return Slot{}, false
	}
	return k.Slots[i], true
}

func (k KeyRegion) validate() error {
	if n := len(k.Slots); n < 1 || n > maxSlots {
		return fmt.Errorf("key region has %d slots, want 1..%d", n, maxSlots)
	}
	for i, s := range k.Slots {
		if want := slotOrder[i]; s.Kind != want {
			return fmt.Errorf("slot %d is %s, want %s", i, s.Kind, want)
		}
		if s.Kind == OrgKEK && s.KeyVersion == 0 {
			return fmt.Errorf("slot %d (%s) has key version zero", i, s.Kind)
		}
		if s.Kind == CustomerSubmitKey && s.KeyVersion != 0 {
			return fmt.Errorf("slot %d (%s) has key version %d, want 0", i, s.Kind, s.KeyVersion)
		}
		if l := len(s.Wrapped); l == 0 || l > MaxWrappedKeySize {
			return fmt.Errorf("slot %d (%s) holds a %d-byte wrapped key, want 1..%d", i, s.Kind, l, MaxWrappedKeySize)
		}
	}
	return nil
}

// appendKeyRegion appends k's encoding to b. k must be valid.
func appendKeyRegion(b []byte, k KeyRegion) []byte {
	b = binary.BigEndian.AppendUint64(b, k.Epoch)
	b = append(b, byte(len(k.Slots))) //nolint:gosec // validate bounds the count to 1..maxSlots
	for _, s := range k.Slots {
		b = append(b, byte(s.Kind))
		b = binary.BigEndian.AppendUint64(b, s.KeyVersion)
		b = binary.BigEndian.AppendUint16(b, uint16(len(s.Wrapped))) //nolint:gosec // validate bounds it to MaxWrappedKeySize
		b = append(b, s.Wrapped...)
	}
	return b
}

// Seal writes an envelope's header to w and returns a writer that seals
// everything written to it under dek as the body. dek is the 32-byte report
// DEK, and keys holds it wrapped to each recipient. Seal does not check that
// the slots wrap dek; a slot that does not cannot open the body.
//
// The caller must Close the returned writer to emit the final chunk; an
// envelope whose final chunk is missing does not open. Close does not close
// w. If Seal, a Write or Close fails, w holds a partial envelope that must be
// discarded rather than committed. The returned writer is not safe for
// concurrent use.
func Seal(w io.Writer, dek []byte, keys KeyRegion) (io.WriteCloser, error) {
	noncePrefix := make([]byte, aeadstream.NoncePrefixSize)
	if _, err := io.ReadFull(rand.Reader, noncePrefix); err != nil {
		return nil, fmt.Errorf("generate nonce prefix: %w", err)
	}
	return seal(w, dek, keys, noncePrefix)
}

func seal(w io.Writer, dek []byte, keys KeyRegion, noncePrefix []byte) (io.WriteCloser, error) {
	if w == nil {
		return nil, errors.New("nil writer")
	}
	if err := keys.validate(); err != nil {
		return nil, err
	}
	prefix := make([]byte, 0, prefixSize)
	prefix = append(prefix, magic...)
	prefix = append(prefix, byte(len(FormatVersion)))
	prefix = append(prefix, FormatVersion...)
	prefix = append(prefix, aeadAES256GCMStream)
	prefix = append(prefix, noncePrefix...)
	prefix = binary.BigEndian.AppendUint32(prefix, ChunkSize)

	// The writer retains prefix as its additional data, so the header is
	// built in a copy.
	sw, err := aeadstream.NewWriter(w, dek, prefix, noncePrefix, ChunkSize)
	if err != nil {
		return nil, err
	}
	if _, err := w.Write(appendKeyRegion(slices.Clone(prefix), keys)); err != nil {
		return nil, fmt.Errorf("write header: %w", err)
	}
	return sw, nil
}

// Header is an envelope's parsed header.
type Header struct {
	Keys KeyRegion

	prefix []byte
}

// ReadHeader reads and checks an envelope's header from r, leaving r at the
// first byte of the body. It refuses an unknown magic, version, AEAD or chunk
// size, a slot count outside 1..2, a slot out of position or malformed, and a
// region that ends before its declared slots do.
//
// Nothing in Keys is authenticated, before or after the body opens. A caller
// may use it to choose which slot to unwrap and to compare the epoch with the
// hunt row's; a forged value is the bucket-writer denial of service the
// design accepts, not a disclosure (see the package doc).
func ReadHeader(r io.Reader) (*Header, error) {
	if r == nil {
		return nil, errors.New("nil reader")
	}
	prefix := make([]byte, prefixSize)
	if _, err := io.ReadFull(r, prefix[:len(magic)]); err != nil {
		return nil, fmt.Errorf("read magic: %w", err)
	}
	if string(prefix[:len(magic)]) != magic {
		return nil, errors.New("input is not a report envelope")
	}
	if _, err := io.ReadFull(r, prefix[len(magic):]); err != nil {
		return nil, fmt.Errorf("read prefix: %w", err)
	}
	if l := int(prefix[len(magic)]); l != len(FormatVersion) {
		return nil, fmt.Errorf("envelope version is %d bytes, want %d", l, len(FormatVersion))
	}
	if v := string(prefix[len(magic)+1 : noncePrefixAt-1]); v != FormatVersion {
		return nil, fmt.Errorf("unsupported envelope version %q, want %q", v, FormatVersion)
	}
	if a := prefix[noncePrefixAt-1]; a != aeadAES256GCMStream {
		return nil, fmt.Errorf("unsupported AEAD %d, want %d", a, aeadAES256GCMStream)
	}
	if c := binary.BigEndian.Uint32(prefix[chunkSizeAt:]); c != ChunkSize {
		return nil, fmt.Errorf("envelope declares chunk size %d, want %d", c, ChunkSize)
	}
	keys, err := readKeyRegion(r)
	if err != nil {
		return nil, err
	}
	return &Header{Keys: keys, prefix: prefix}, nil
}

func readKeyRegion(r io.Reader) (KeyRegion, error) {
	var fixed [8 + 1]byte
	if _, err := io.ReadFull(r, fixed[:]); err != nil {
		return KeyRegion{}, fmt.Errorf("read key region: %w", err)
	}
	count := int(fixed[8])
	if count < 1 || count > maxSlots {
		return KeyRegion{}, fmt.Errorf("envelope declares %d key slots, want 1..%d", count, maxSlots)
	}
	keys := KeyRegion{
		Epoch: binary.BigEndian.Uint64(fixed[:8]),
		Slots: make([]Slot, 0, count),
	}
	for i := range count {
		var slot [1 + 8 + 2]byte
		if _, err := io.ReadFull(r, slot[:]); err != nil {
			return KeyRegion{}, fmt.Errorf("read slot %d of %d: %w", i, count, err)
		}
		l := int(binary.BigEndian.Uint16(slot[9:]))
		if l == 0 || l > MaxWrappedKeySize {
			return KeyRegion{}, fmt.Errorf("slot %d declares a %d-byte wrapped key, want 1..%d", i, l, MaxWrappedKeySize)
		}
		wrapped := make([]byte, l)
		if _, err := io.ReadFull(r, wrapped); err != nil {
			return KeyRegion{}, fmt.Errorf("read slot %d of %d: %w", i, count, err)
		}
		keys.Slots = append(keys.Slots, Slot{
			Kind:       SlotKind(slot[0]),
			KeyVersion: binary.BigEndian.Uint64(slot[1:9]),
			Wrapped:    wrapped,
		})
	}
	if err := keys.validate(); err != nil {
		return KeyRegion{}, err
	}
	return keys, nil
}

// Open returns a reader of the plaintext of body, the rest of the reader h
// was read from, under dek.
//
// Read the returned reader to EOF to verify the complete body. Each chunk
// authenticates before its plaintext is returned, and modification,
// deletion, reordering, splicing, a missing final chunk, a changed prefix or
// the wrong dek causes a read error. Open may be called more than once on the
// same Header. The returned reader is not safe for concurrent use.
func (h *Header) Open(body io.Reader, dek []byte) (io.Reader, error) {
	if len(h.prefix) != prefixSize {
		return nil, errors.New("header was not read by ReadHeader")
	}
	return aeadstream.NewReader(body, dek, h.prefix, h.prefix[noncePrefixAt:chunkSizeAt], ChunkSize)
}

// Rewrite copies the envelope in src to dst with its key region replaced by
// keys. It is [ReadHeader] followed by [Header.Rewrite].
func Rewrite(dst io.Writer, src io.Reader, keys KeyRegion) error {
	h, err := ReadHeader(src)
	if err != nil {
		return err
	}
	return h.Rewrite(dst, src, keys)
}

// Rewrite writes to dst the envelope h was read from, with its key region
// replaced by keys. body is the rest of the reader h was read from. The
// prefix and the body are copied byte for byte, so the body still opens under
// the same DEK and [BodyDigest] is unchanged. Rewrite does not open the body,
// and does not check that keys wrap the DEK the body is sealed under.
//
// If Rewrite fails, dst holds a partial envelope that must be discarded
// rather than committed. A re-key writes the new generation before deleting
// the old, so committing a partial one would leave a report that no key opens.
func (h *Header) Rewrite(dst io.Writer, body io.Reader, keys KeyRegion) error {
	if dst == nil {
		return errors.New("nil writer")
	}
	if err := keys.validate(); err != nil {
		return err
	}
	if len(h.prefix) != prefixSize {
		return errors.New("header was not read by ReadHeader")
	}
	if _, err := dst.Write(appendKeyRegion(slices.Clone(h.prefix), keys)); err != nil {
		return fmt.Errorf("write header: %w", err)
	}
	if _, err := io.Copy(dst, body); err != nil {
		return fmt.Errorf("copy body: %w", err)
	}
	return nil
}

// BodyDigest reads the envelope in r and returns the SHA-256 of its body,
// every byte after the key region, which is where [ReadHeader] leaves its
// reader. It is the release record's digest, and a rewritten key region
// leaves it unchanged. BodyDigest checks the header but does not open the
// body.
func BodyDigest(r io.Reader) ([sha256.Size]byte, error) {
	if _, err := ReadHeader(r); err != nil {
		return [sha256.Size]byte{}, err
	}
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return [sha256.Size]byte{}, fmt.Errorf("read body: %w", err)
	}
	return [sha256.Size]byte(h.Sum(nil)), nil
}
