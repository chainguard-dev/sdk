/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	apkotypes "chainguard.dev/apko/pkg/build/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	fuzz "github.com/google/gofuzz"
	"google.golang.org/protobuf/testing/protocmp"
)

// TestRoundTrip tests that converting an ImageConfiguration to protobuf and back
// yields the original ImageConfiguration, and vice versa, modulo ignored fields.
func TestRoundTrip(t *testing.T) {
	if err := quick.Check(func(apko apkotypes.ImageConfiguration) bool {
		pb := ToApkoProto(apko)
		apko2 := ToApkoNative(pb)
		// Include was deprecated in the proto.
		// BaseImage is not supported in the registry proto.
		// We ignore them here to avoid the diff.
		// https://github.com/chainguard-dev/apko/blob/main/pkg/build/types/types.go#L185-L186
		if d := cmp.Diff(apko, apko2,
			cmpopts.IgnoreFields(apkotypes.ImageConfiguration{}, "Include"),
			cmpopts.IgnoreFields(apkotypes.ImageContents{}, "BaseImage")); d != "" {
			t.Errorf("apko diff(-want,+got): %s", d)
			return false
		}

		pb2 := ToApkoProto(apko2)
		if d := cmp.Diff(pb, pb2, protocmp.Transform()); d != "" {
			t.Errorf("proto diff(-want,+got): %s", d)
		}
		return true
	}, &quick.Config{
		Values: func(vals []reflect.Value, r *rand.Rand) {
			// Use gofuzz to generate random ImageConfiguration values
			fz := fuzz.NewWithSeed(r.Int63())
			fz.NilChance(0.2)
			fz.NumElements(1, 3)

			var apko apkotypes.ImageConfiguration
			fz.Fuzz(&apko)
			vals[0] = reflect.ValueOf(apko)
		},
	}); err != nil {
		t.Error(err)
	}
}
