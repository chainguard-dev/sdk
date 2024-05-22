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
	// Map of stringified name to capability. Initialized with initNameCapabilityMap()
	nameCapabilityMap = make(map[string]Capability, len(Capability_value))
	ponce             sync.Once
	perror            error
)

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

func Stringify(cap Capability) (string, error) {
	evd := cap.Descriptor().Values().ByNumber(cap.Number())
	if evd == nil {
		return "", status.Errorf(codes.Internal, "capability has no descriptor: %v", cap)
	}
	opt := evd.Options()
	if opt == nil {
		return "", status.Errorf(codes.Internal, "capability has no options: %v", cap)
	}
	evo := opt.(*descriptorpb.EnumValueOptions)
	name := proto.GetExtension(evo, E_Name)
	if name == nil {
		return "", status.Errorf(codes.Internal, "capability is missing the name option: %v", cap)
	}
	return name.(string), nil
}

func StringifyAll(caps []Capability) ([]string, error) {
	scs := make([]string, 0, len(caps))
	for _, cap := range caps {
		sc, err := Stringify(cap)
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
		for cap := range Capability_name {
			scap, perror := Stringify(Capability(cap))
			if perror == nil {
				nameCapabilityMap[scap] = Capability(cap)
			} else {
				clog.FromContext(context.Background()).Errorf("Failed to stringify capability %d, error: %v",
					cap, perror)
			}
		}
	})

	if perror != nil {
		return Capability_UNKNOWN, perror
	}
	return nameCapabilityMap[name], nil
}

func Bitify(cap Capability) (uint32, error) {
	evd := cap.Descriptor().Values().ByNumber(cap.Number())
	if evd == nil {
		return 0, status.Errorf(codes.Internal, "capability has no descriptor: %v", cap)
	}
	opt := evd.Options()
	if opt == nil {
		return 0, status.Errorf(codes.Internal, "capability has no options: %v", cap)
	}
	evo := opt.(*descriptorpb.EnumValueOptions)
	name := proto.GetExtension(evo, E_Bit)
	if name == nil {
		return 0, status.Errorf(codes.Internal, "capability is missing the bit option: %v", cap)
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
	for _, cap := range s {
		b, err := Bitify(cap)
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
		for _, cap := range caps {
			*s = append(*s, cap)
		}
		return nil

	default:
		// Compact encoding
		var bs bitset.BitSet
		if err := json.Unmarshal(b, &bs); err != nil {
			return err
		}
		for i := range Capability_name {
			cap := Capability(i) //nolint: revive
			if cap == Capability_UNKNOWN {
				continue
			}
			bit, err := Bitify(cap)
			if err != nil {
				return err
			}
			if bs.Test(uint(bit)) {
				*s = append(*s, cap)
				// This ensures that our unit testing checks that no two
				// enumeration values are assigned the same bit.
				bs.Clear(uint(bit))
			}
		}
		sort.Slice(*s, func(i int, j int) bool {
			return (*s)[i] < (*s)[j]
		})
		return nil
	}
}
