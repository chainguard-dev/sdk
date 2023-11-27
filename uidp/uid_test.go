/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uidp

import "testing"

func TestUIDPValidity(t *testing.T) {
	var uidp UIDP

	for i := 0; i < 5; i++ {
		uidp = uidp.NewChild()

		if !Valid(string(uidp)) {
			t.Errorf("Invalid UIDP: %s", uidp)
		}
	}
}
