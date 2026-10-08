/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package migration

import (
	"testing"

	guardpb "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1"
)

func TestRegistryCoversFeatures(t *testing.T) {
	definitions := All()
	seen := make(map[guardpb.FeatureType]struct{}, len(definitions))
	keys := make(map[string]struct{}, len(definitions))
	for _, definition := range definitions {
		if _, ok := seen[definition.Feature]; ok {
			t.Fatalf("duplicate feature: %v", definition.Feature)
		}
		seen[definition.Feature] = struct{}{}
		if definition.Capability == 0 || definition.Name == "" || definition.QueueSuffix == "" {
			t.Fatalf("incomplete definition: %+v", definition)
		}
		key := definition.Key("https://github.com/owner/repo/")
		if _, ok := keys[key]; ok {
			t.Fatalf("duplicate queue key: %q", key)
		}
		keys[key] = struct{}{}
		if got := definition.Key("github.com/owner/repo"); got != key {
			t.Errorf("scheme-relative key = %q, want %q", got, key)
		}
	}
	for value := range guardpb.FeatureType_name {
		feature := guardpb.FeatureType(value)
		if feature == guardpb.FeatureType_FEATURE_TYPE_UNSPECIFIED {
			continue
		}
		if _, ok := seen[feature]; !ok {
			t.Errorf("feature %v has no definition", feature)
		}
	}
	definitions[0].Name = "changed"
	if got := All()[0].Name; got == "changed" {
		t.Fatal("All exposes mutable registry storage")
	}
	if _, ok := Lookup(guardpb.FeatureType(1000)); ok {
		t.Fatal("unknown feature was accepted")
	}
}
