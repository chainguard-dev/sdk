/*
Copyright 2024 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package v1

import (
	"testing"
	"testing/quick"

	apkotypes "chainguard.dev/apko/pkg/build/types"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestRoundTrip(t *testing.T) {
	if err := quick.Check(func(apko apkotypes.ImageConfiguration) bool {
		pb := ToApkoProto(apko)
		apko2 := ToApkoNative(pb)
		if d := cmp.Diff(apko, apko2); d != "" {
			t.Errorf("apko diff(-want,+got): %s", d)
			return false
		}

		pb2 := ToApkoProto(apko2)
		if d := cmp.Diff(pb, pb2, protocmp.Transform()); d != "" {
			t.Errorf("proto diff(-want,+got): %s", d)
		}
		return true
	}, &quick.Config{}); err != nil {
		t.Error(err)
	}
}
