/*
Copyright 2026 Chainguard, Inc.
SPDX-License-Identifier: Apache-2.0
*/

package migration_test

import (
	"fmt"

	"chainguard.dev/sdk/guardener/migration"
	guardpb "chainguard.dev/sdk/proto/chainguard/platform/guardener/v1alpha1"
)

func ExampleAll() {
	for _, definition := range migration.All() {
		fmt.Println(definition.Name)
	}
	// Output:
	// Actions
	// Images
}

func ExampleLookup() {
	definition, ok := migration.Lookup(guardpb.FeatureType_FEATURE_TYPE_IMAGES)
	fmt.Println(definition.Name, ok)
	// Output: Images true
}

func ExampleDefinition_Key() {
	definition, _ := migration.Lookup(guardpb.FeatureType_FEATURE_TYPE_ACTIONS)
	fmt.Println(definition.Key("https://github.com/acme/demo"))
	// Output: github.com/acme/demo/tree/HEAD/.github/workflows
}
