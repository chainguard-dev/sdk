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
	perror            error

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

func StringifyAll(caps []Capability) ([]string, error) {
	scs := make([]string, 0, len(caps))
	for _, capability := range caps {
		sc, err := Stringify(capability)
		if err != nil {
			return nil, err
		}
		scs = append(scs, sc)
	}
	return scs, nil
}

func Parse(name string) (Capability, error) {
	ponce.Do(func() {
		// Populate nameCapabilityMap
		for capability := range Capability_name {
			scap, perror := Stringify(Capability(capability))
			if perror == nil {
				nameCapabilityMap[scap] = Capability(capability)
			} else {
				clog.FromContext(context.Background()).Errorf("Failed to stringify capability %d, error: %v",
					capability, perror)
			}
		}
	})

	if perror != nil {
		return Capability_UNKNOWN, perror
	}
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

func (s Set) String() string {
	caps := make([]string, 0, len(s))
	for _, c := range s {
		st, err := Stringify(c)
		if err != nil {
			st = fmt.Sprintf("[ERROR(cap=%d)", c)
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
