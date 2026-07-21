/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package capabilities

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/bits-and-blooms/bitset"
	"github.com/chainguard-dev/clog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

var (
	// Map of stringified name to capability.
	nameCapabilityMap = make(map[string]Capability, len(Capability_value))
	ponce             sync.Once

	// Map of Capability to result of Bitify(). Set in initBitifyMap().
	bitifiedMap = make(map[Capability]uint32, len(Capability_value))
	// Sorted list of {bit, cap} so we can skip sorting in UnmarshalJSON.
	bitCaps    = make([]bitcap, 0, len(Capability_value))
	bitifyOnce sync.Once
)

type bitcap struct {
	bit uint32
	cap Capability
}

// We can't do this in init() because init() ordering is hard.
func initBitifyMap() {
	for i := range Capability_name {
		capability := Capability(i) //nolint: revive
		if capability == Capability_UNKNOWN {
			continue
		}
		bit, err := bitify(capability)
		if err != nil {
			// This should never happen!
			continue
		}

		bitifiedMap[capability] = bit
		bitCaps = append(bitCaps, bitcap{bit, capability})
	}

	sort.Slice(bitCaps, func(i int, j int) bool {
		return bitCaps[i].cap < bitCaps[j].cap
	})
}

// Names returns a slice of all capabilities Stringify'd, sans UNKNOWN.
func Names() []string {
	all := make([]string, 0, len(Capability_name)-1) // One less, we don't want UNKNOWN
	for n := range Capability_name {
		if Capability(n) == Capability_UNKNOWN {
			continue
		}
		sc, err := Stringify(Capability(n))
		if err != nil {
			// This should never happen!
			continue
		}
		all = append(all, sc)
	}
	return all
}

// InternalOnly reports whether the capability is annotated as
// internal-only via the (internal_only) proto option. Capabilities
// that return true must never be granted by a role available to
// customer organisations.
func InternalOnly(capability Capability) bool {
	evd := capability.Descriptor().Values().ByNumber(capability.Number())
	if evd == nil {
		return false
	}

	opt := evd.Options()
	if opt == nil {
		return false
	}

	evo := opt.(*descriptorpb.EnumValueOptions)
	v, ok := proto.GetExtension(evo, E_InternalOnly).(bool)
	if !ok {
		return false
	}

	return v
}

func Deprecated(capability Capability) bool {
	evd := capability.Descriptor().Values().ByNumber(capability.Number())
	if evd == nil {
		return false
	}
	opt := evd.Options()
	if opt == nil {
		return false
	}
	evo := opt.(*descriptorpb.EnumValueOptions)
	return evo.GetDeprecated()
}

func Stringify(capability Capability) (string, error) {
	evd := capability.Descriptor().Values().ByNumber(capability.Number())
	if evd == nil {
		return "", status.Errorf(codes.Internal, "capability has no descriptor: %v", capability)
	}
	opt := evd.Options()
	if opt == nil {
		return "", status.Errorf(codes.Internal, "capability has no options: %v", capability)
	}
	evo := opt.(*descriptorpb.EnumValueOptions)
	name := proto.GetExtension(evo, E_Name)
	if name == nil {
		return "", status.Errorf(codes.Internal, "capability is missing the name option: %v", capability)
	}
	return name.(string), nil
}

// warnedUnknownCaps tracks capability values StringifyAllContext has already
// warned about, so a stale value carried on every request (a stored role or a
// JWT cap claim) logs once per process instead of once per call. Size-capped
// so unexpected inputs cannot grow it without bound; at the cap, dedup
// degrades to warning on every occurrence rather than going silent.
var warnedUnknownCaps = struct {
	sync.Mutex
	seen map[Capability]struct{}
}{seen: make(map[Capability]struct{})}

const warnedUnknownCapsLimit = 1024

func warnUnknownCapability(ctx context.Context, capability Capability) {
	warnedUnknownCaps.Lock()
	_, seen := warnedUnknownCaps.seen[capability]
	atLimit := false
	if !seen && len(warnedUnknownCaps.seen) < warnedUnknownCapsLimit {
		warnedUnknownCaps.seen[capability] = struct{}{}
		atLimit = len(warnedUnknownCaps.seen) == warnedUnknownCapsLimit
	}
	warnedUnknownCaps.Unlock()
	if seen {
		return
	}
	clog.WarnContextf(ctx, "skipping capability %d not in this binary's enum (deleted capability, or one newer than this build)", capability)
	if atLimit {
		clog.WarnContextf(ctx, "capability warn-dedup limit (%d) reached; novel unknown-capability values will now log on every occurrence (already-seen values remain deduplicated)", warnedUnknownCapsLimit)
	}
}

// StringifyAllContext converts capabilities to their string names. A value
// with no descriptor is not in this binary's enum: either deleted from the
// proto (a stale value on a stored role, CUS-843) or added after this binary
// was built (deploy skew). Either way it grants nothing here, so it is skipped
// with a once-per-process warning rather than failing the whole list —
// otherwise a single stale capability on a stored role makes that role
// impossible to read or list, which takes down the Console IDP settings page.
// Other Stringify errors (a present enum value missing its options or name
// extension) indicate a proto authoring bug and are surfaced to the caller.
func StringifyAllContext(ctx context.Context, caps []Capability) ([]string, error) {
	scs := make([]string, 0, len(caps))
	for _, capability := range caps {
		if capability.Descriptor().Values().ByNumber(capability.Number()) == nil {
			warnUnknownCapability(ctx, capability)
			continue
		}
		sc, err := Stringify(capability)
		if err != nil {
			return nil, err
		}
		scs = append(scs, sc)
	}
	return scs, nil
}

// StringifyAll is StringifyAllContext with context.Background(). Prefer
// StringifyAllContext where a context is available so the skip warning carries
// request logging metadata; some callers (package-variable initialization,
// descriptor walks) have none.
func StringifyAll(caps []Capability) ([]string, error) {
	return StringifyAllContext(context.Background(), caps)
}

// Parse resolves a capability's string name (e.g. "groups.list") to its enum
// value. Parse never returns a non-nil error; the error result is retained
// for API compatibility. Unknown names return (Capability_UNKNOWN, nil), so
// callers must check for Capability_UNKNOWN explicitly; UNKNOWN is never
// granted, so unresolvable input fails closed at the caller.
//
// The name map is built once from this binary's own compiled-in enum table.
// An enum value whose (name) option cannot be read is a generated-code
// inconsistency, not runtime data — TestRoundTrip rejects it in CI. If such a
// binary ships anyway, the value is logged at error level and omitted from
// the map, so its name resolves to UNKNOWN and resolution of every other
// name is unaffected — the same per-value scoping as StringifyAllContext and
// Bitify (CUS-843), rather than failing every Parse call in the process.
func Parse(name string) (Capability, error) {
	ponce.Do(func() {
		// Populate nameCapabilityMap
		for capability := range Capability_name {
			scap, err := Stringify(Capability(capability))
			if err != nil {
				clog.ErrorContextf(context.Background(),
					"capability %d omitted from Parse name map (generated code missing its (name) option?): %v",
					capability, err)
				continue
			}
			nameCapabilityMap[scap] = Capability(capability)
		}
	})

	return nameCapabilityMap[name], nil
}

func Bitify(capability Capability) (uint32, error) {
	bitifyOnce.Do(initBitifyMap)

	bit, ok := bitifiedMap[capability]
	if !ok {
		// If it's missing in bitifiedMap, we ignored it because bitify() returned an error.
		// Just call bitify() again here to get whatever the error was.
		// This should almost never happen, so duplicating the work is fine.
		return bitify(capability)
	}

	return bit, nil
}

func bitify(capability Capability) (uint32, error) {
	evd := capability.Descriptor().Values().ByNumber(capability.Number())
	if evd == nil {
		return 0, status.Errorf(codes.Internal, "capability has no descriptor: %v", capability)
	}
	opt := evd.Options()
	if opt == nil {
		return 0, status.Errorf(codes.Internal, "capability has no options: %v", capability)
	}
	evo := opt.(*descriptorpb.EnumValueOptions)
	name := proto.GetExtension(evo, E_Bit)
	if name == nil {
		return 0, status.Errorf(codes.Internal, "capability is missing the bit option: %v", capability)
	}
	return name.(uint32), nil
}

// Set performs efficient encoding of a list of capabilities.
type Set []Capability

// String renders the set for diagnostics. Values missing from this binary's
// enum (deleted, or newer than this build) render as "unknown(cap=N)" and any
// other Stringify failure as "error(cap=N)", mirroring StringifyAllContext's
// two cases while keeping malformed values visible instead of silently
// dropping them.
func (s Set) String() string {
	caps := make([]string, 0, len(s))
	for _, c := range s {
		st, err := Stringify(c)
		if err != nil {
			if c.Descriptor().Values().ByNumber(c.Number()) == nil {
				st = fmt.Sprintf("unknown(cap=%d)", c)
			} else {
				st = fmt.Sprintf("error(cap=%d)", c)
			}
		}
		caps = append(caps, st)
	}
	sort.Strings(caps)
	return strings.Join(caps, ",")
}

// MarshalJSON implements json.Marshaler
func (s Set) MarshalJSON() ([]byte, error) {
	bs := bitset.New(50)
	for _, capability := range s {
		b, err := Bitify(capability)
		if err != nil {
			return nil, err
		}
		bs.Set(uint(b))
	}
	return bs.MarshalJSON()
}

// UnmarshalJSON implements json.Unmarshaler
func (s *Set) UnmarshalJSON(b []byte) error {
	switch {
	case len(b) == 0:
		return io.EOF

	case b[0] == '[':
		// Legacy decoding!
		var caps []Capability
		if err := json.Unmarshal(b, &caps); err != nil {
			return err
		}
		for _, capability := range caps {
			*s = append(*s, capability)
		}
		return nil

	default:
		bitifyOnce.Do(initBitifyMap)

		// Compact encoding
		var bs bitset.BitSet
		if err := json.Unmarshal(b, &bs); err != nil {
			return err
		}

		*s = make([]Capability, 0, bs.Count())

		for _, bitcap := range bitCaps {
			if bs.Test(uint(bitcap.bit)) {
				*s = append(*s, bitcap.cap)
				// This ensures that our unit testing checks that no two
				// enumeration values are assigned the same bit.
				bs.Clear(uint(bitcap.bit))
			}
		}

		return nil
	}
}
