/*
Copyright 2022 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package uidp

import "testing"

func TestReparent(t *testing.T) {
	parent := NewUIDP("")       // root UID (40 hex chars)
	source := parent.NewChild() // parent/suid1
	child := source.NewChild()  // parent/suid1/suid2

	newParent := NewUIDP("") // different root

	reparented, err := child.Reparent(newParent)
	if err != nil {
		t.Fatalf("Reparent: unexpected error: %v", err)
	}

	if !Valid(string(reparented)) {
		t.Errorf("Reparent produced invalid UIDP: %s", reparented)
	}

	// The result should be a direct child of newParent (not deeper).
	if !IsAncestor(string(newParent), string(reparented)) {
		t.Errorf("Reparent result %s is not a child of %s", reparented, newParent)
	}
	if IsAncestor(string(reparented), string(newParent)) {
		t.Errorf("Reparent result %s is deeper than one level under %s", reparented, newParent)
	}

	// Reparent should be deterministic.
	reparented2, err := child.Reparent(newParent)
	if err != nil {
		t.Fatalf("Reparent (second call): unexpected error: %v", err)
	}
	if reparented != reparented2 {
		t.Errorf("Reparent not deterministic: %s != %s", reparented, reparented2)
	}

	// Different sources produce different results.
	other := source.NewChild()
	reparented3, err := other.Reparent(newParent)
	if err != nil {
		t.Fatalf("Reparent (other source): unexpected error: %v", err)
	}
	if reparented == reparented3 {
		t.Errorf("Reparent produced same result for different sources: %s", reparented)
	}

	// The SUID is preserved: the basename of the reparented path equals
	// the basename of the original.
	wantSUID := string(child)[len(string(source))+1:]
	gotSUID := string(reparented)[len(string(newParent))+1:]
	if gotSUID != wantSUID {
		t.Errorf("Reparent did not preserve SUID: want %s, got %s", wantSUID, gotSUID)
	}

	// Reparent on a root UIDP must return an error.
	t.Run("errors on root UIDP", func(t *testing.T) {
		root := NewUIDP("")
		_, err := root.Reparent(newParent)
		if err == nil {
			t.Error("expected error when calling Reparent on a root UIDP, got nil")
		}
	})
}

func TestUIDPValidity(t *testing.T) {
	var uidp UIDP

	for i := 0; i < 5; i++ {
		uidp = uidp.NewChild()

		if !Valid(string(uidp)) {
			t.Errorf("Invalid UIDP: %s", uidp)
		}
	}
}
