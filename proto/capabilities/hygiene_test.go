/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	_ "embed"
	"regexp"
	"sort"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

// The proto source is embedded so the counter check compares the same file
// this binary's descriptor was generated from, independent of the test's
// working directory (module cache, exported chainguard.dev/sdk repo).
//
//go:embed capabilities.proto
var capabilitiesProtoSrc []byte

// retiredBit is a retirement-ledger entry: why the bit is frozen and, for a
// deleted capability, the enum number it carried (zero for bits that were
// never assigned). Recording the number makes `reserved` hygiene
// self-enforcing — TestReservedEnumNumbers requires every recorded number to
// stay reserved in the proto.
type retiredBit struct {
	number protoreflect.EnumNumber
	reason string
}

// neverAssignBits are bit indices below the current maximum that no live
// capability may use. Bits are the wire format of the JWT "cap" claim (Set
// MarshalJSON); once a bit has shipped in a token it must never be reassigned,
// or tokens minted before the change silently carry the new meaning. Proto
// `reserved` covers enum numbers only — this map is the reservation ledger for
// bits. When deleting a capability, add its bit here with its enum number and
// a reference.
var neverAssignBits = map[uint32]retiredBit{
	0:   {reason: "never assigned; keep unused"},
	40:  {number: 1601, reason: "registry.pull, deleted in mono#23247"},
	41:  {number: 1602, reason: "registry.push, deleted in mono#23247"},
	135: {reason: "never assigned (allocation skip); frozen for safety"},
	136: {reason: "never assigned (allocation skip); frozen for safety"},
	143: {reason: "never assigned (allocation skip); frozen for safety"},
}

// capabilityBits enumerates every non-UNKNOWN enum value's (bit) option via
// the compiled descriptor — the same source of truth bitify() uses.
func capabilityBits(t *testing.T) map[uint32][]protoreflect.Name {
	t.Helper()
	vals := Capability(0).Descriptor().Values()
	bits := make(map[uint32][]protoreflect.Name, vals.Len())
	for i := range vals.Len() {
		vd := vals.Get(i)
		if vd.Number() == 0 { // UNKNOWN carries no bit
			continue
		}
		opts, ok := vd.Options().(*descriptorpb.EnumValueOptions)
		if !ok || !proto.HasExtension(opts, E_Bit) {
			// Fatal: a truncated map would corrupt maxBit/holes math in the
			// callers and bury this root cause under secondary diagnostics.
			t.Fatalf("%s (=%d) has no (bit) option", vd.Name(), vd.Number())
		}
		bit := proto.GetExtension(opts, E_Bit).(uint32)
		bits[bit] = append(bits[bit], vd.Name())
	}
	return bits
}

func maxAssignedBit(bits map[uint32][]protoreflect.Name) uint32 {
	var maxBit uint32
	for bit := range bits {
		if bit > maxBit {
			maxBit = bit
		}
	}
	return maxBit
}

// highestKnownBit is the bit space's high-water mark: the maximum across live
// capability bits and the retirement ledger. Allocation must resume above it —
// retiring the current maximum bit does not free its slot.
func highestKnownBit(bits map[uint32][]protoreflect.Name) uint32 {
	maxBit := maxAssignedBit(bits)
	for b := range neverAssignBits {
		if b > maxBit {
			maxBit = b
		}
	}
	return maxBit
}

func TestCapabilityBitsUnique(t *testing.T) {
	for bit, names := range capabilityBits(t) {
		if len(names) > 1 {
			t.Errorf("bit %d is assigned to %d capabilities %v; every capability needs a distinct bit", bit, len(names), names)
		}
	}
}

func TestRetiredBitsNotReassigned(t *testing.T) {
	bits := capabilityBits(t)
	for bit, rb := range neverAssignBits {
		if names, ok := bits[bit]; ok {
			t.Errorf("bit %d is assigned to %v but is retired (%s); allocate a fresh bit instead", bit, names, rb.reason)
		}
	}
}

// TestBitLedgerComplete forces every gap in the bit space to be accounted for:
// deleting a capability creates a new hole, which fails this test until the
// bit is recorded in neverAssignBits with a reference. The scan extends to the
// highest LEDGER bit, not just the highest live bit, so retiring the current
// maximum (which lowers maxBit below its own ledger entry) still balances.
func TestBitLedgerComplete(t *testing.T) {
	bits := capabilityBits(t)
	scanMax := highestKnownBit(bits)

	var holes []int
	for b := range scanMax + 1 {
		if _, ok := bits[b]; !ok {
			holes = append(holes, int(b))
		}
	}
	ledger := make([]int, 0, len(neverAssignBits))
	for b := range neverAssignBits {
		ledger = append(ledger, int(b))
	}
	sort.Ints(holes)
	sort.Ints(ledger)

	if diff := cmp.Diff(ledger, holes); diff != "" {
		t.Errorf("unassigned bits up to %d must exactly match the neverAssignBits ledger (-ledger +holes):\n%s\n+N (unledgered hole): you deleted a capability — add its bit to neverAssignBits with a reference.\n-N (ledgered but live): a capability reuses a retired bit — allocate a fresh bit instead (see TestRetiredBitsNotReassigned).", scanMax, diff)
	}
}

// TestReservedEnumNumbers requires deleted-capability enum numbers to stay
// `reserved` in the proto — numbers are the durable DB format (role capability
// arrays), the parallel obligation to bits. The expected set derives from the
// retirement ledger, so a deletion recorded there cannot pass CI without its
// number reserved too.
func TestReservedEnumNumbers(t *testing.T) {
	rr := Capability(0).Descriptor().ReservedRanges()
	// 1 and 1873 are reserved without a ledgered bit: 1873's bit (168) was
	// recycled in 2026-06, before this ledger existed.
	want := []protoreflect.EnumNumber{1, 1873}
	for _, rb := range neverAssignBits {
		if rb.number != 0 {
			want = append(want, rb.number)
		}
	}
	for _, n := range want {
		if !rr.Has(n) {
			t.Errorf("enum number %d must stay reserved in capabilities.proto", n)
		}
	}
}

// TestNextBitComment cross-checks the "// next bit: N" allocation ledger in
// capabilities.proto against the compiled descriptor, so the comment authors
// follow when adding capabilities can't go stale. This pairing assumes
// capabilities.pb.go is regenerated whenever capabilities.proto changes
// (enforced by the verify workflow); a simultaneously stale comment AND
// descriptor would not be caught here.
func TestNextBitComment(t *testing.T) {
	ms := regexp.MustCompile(`(?m)^\s*// next bit: (\d+)\s*$`).FindAllSubmatch(capabilitiesProtoSrc, -1)
	if len(ms) != 1 {
		near := regexp.MustCompile(`(?mi)^.*next bit.*$`).FindAll(capabilitiesProtoSrc, -1)
		t.Fatalf("capabilities.proto must contain exactly one '// next bit: N' ledger comment, found %d (lines mentioning it: %q)", len(ms), near)
	}
	n, err := strconv.Atoi(string(ms[0][1]))
	if err != nil {
		t.Fatalf("parsing next-bit value: %v", err)
	}

	// One past the highest live-or-retired bit: retiring the current maximum
	// must not walk the counter back onto the freshly retired bit.
	maxBit := highestKnownBit(capabilityBits(t))
	if want := int(maxBit) + 1; n != want {
		t.Errorf("// next bit comment: got %d, want %d (highest live-or-retired bit is %d)", n, want, maxBit)
	}
}
